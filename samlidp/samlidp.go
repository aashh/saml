// Package samlidp a rudimentary SAML identity provider suitable for
// testing or as a starting point for a more complex service.
package samlidp

import (
	"crypto"
	"crypto/x509"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/logger"
)

// Options represent the parameters to New() for creating a new IDP server
type Options struct {
	URL               url.URL
	Key               crypto.PrivateKey
	Signer            crypto.Signer
	Logger            logger.Interface
	Certificate       *x509.Certificate
	Store             Store
	LoginFormTemplate *template.Template

	// AdminMiddleware, if set, is applied to all admin API endpoints
	// (users, services, sessions, shortcuts). Use this to require
	// authentication on management endpoints. If nil, admin endpoints
	// are served without any access control — only appropriate for
	// development and testing.
	AdminMiddleware func(http.Handler) http.Handler
}

// Server represents an IDP server. The server provides the following URLs:
//
//	/metadata     - the SAML metadata
//	/sso          - the SAML endpoint to initiate an authentication flow
//	/login        - prompt for a username and password if no session established
//	/login/:shortcut - kick off an IDP-initiated authentication flow
//	/services     - RESTful interface to Service objects
//	/users        - RESTful interface to User objects
//	/sessions     - RESTful interface to Session objects
//	/shortcuts    - RESTful interface to Shortcut objects
type Server struct {
	http.Handler
	idpConfigMu       sync.RWMutex // protects calls into the IDP
	logger            logger.Interface
	serviceProviders  map[string]*saml.EntityDescriptor
	IDP               saml.IdentityProvider // the underlying IDP
	Store             Store                 // the data store
	LoginFormTemplate *template.Template
	AdminMiddleware   func(http.Handler) http.Handler
}

// New returns a new Server
func New(opts Options) (*Server, error) {
	opts.URL.Path = strings.TrimSuffix(opts.URL.Path, "/")

	metadataURL := opts.URL
	metadataURL.Path += "/metadata"
	ssoURL := opts.URL
	ssoURL.Path += "/sso"
	loginURL := opts.URL
	loginURL.Path += "/login"
	logr := opts.Logger
	if logr == nil {
		logr = logger.DefaultLogger
	}

	s := &Server{
		serviceProviders: map[string]*saml.EntityDescriptor{},
		IDP: saml.IdentityProvider{
			Key:         opts.Key,
			Signer:      opts.Signer,
			Logger:      logr,
			Certificate: opts.Certificate,
			MetadataURL: metadataURL,
			SSOURL:      ssoURL,
			LoginURL:    loginURL,
		},
		logger:            logr,
		Store:             opts.Store,
		LoginFormTemplate: opts.LoginFormTemplate,
		AdminMiddleware:   opts.AdminMiddleware,
	}

	s.IDP.SessionProvider = s
	s.IDP.ServiceProviderProvider = s

	if err := s.initializeServices(); err != nil {
		return nil, err
	}
	s.InitializeHTTP()
	return s, nil
}

// InitializeHTTP sets up the HTTP handler for the server. (This function
// is called automatically for you by New, but you may need to call it
// yourself if you don't create the object using New.)
func (s *Server) InitializeHTTP() {
	mux := http.NewServeMux()
	s.Handler = mux

	mux.HandleFunc("GET /metadata", func(w http.ResponseWriter, r *http.Request) {
		s.idpConfigMu.RLock()
		defer s.idpConfigMu.RUnlock()
		s.IDP.ServeMetadata(w, r)
	})
	mux.HandleFunc("/sso", func(w http.ResponseWriter, r *http.Request) {
		s.IDP.ServeSSO(w, r)
	})

	mux.HandleFunc("/login", s.HandleLogin)
	mux.HandleFunc("/login/{shortcut}", s.HandleIDPInitiated)
	mux.HandleFunc("/login/{shortcut}/{suffix}", s.HandleIDPInitiated)

	admin := s.adminHandler

	mux.Handle("GET /services/", admin(http.HandlerFunc(s.HandleListServices)))
	mux.Handle("GET /services/{id}", admin(http.HandlerFunc(s.HandleGetService)))
	mux.Handle("PUT /services/{id}", admin(http.HandlerFunc(s.HandlePutService)))
	mux.Handle("POST /services/{id}", admin(http.HandlerFunc(s.HandlePutService)))
	mux.Handle("DELETE /services/{id}", admin(http.HandlerFunc(s.HandleDeleteService)))

	mux.Handle("GET /users/", admin(http.HandlerFunc(s.HandleListUsers)))
	mux.Handle("GET /users/{id}", admin(http.HandlerFunc(s.HandleGetUser)))
	mux.Handle("PUT /users/{id}", admin(http.HandlerFunc(s.HandlePutUser)))
	mux.Handle("DELETE /users/{id}", admin(http.HandlerFunc(s.HandleDeleteUser)))

	mux.Handle("GET /sessions/", admin(http.HandlerFunc(s.HandleListSessions)))
	mux.Handle("GET /sessions/{id}", admin(http.HandlerFunc(s.HandleGetSession)))
	mux.Handle("DELETE /sessions/{id}", admin(http.HandlerFunc(s.HandleDeleteSession)))

	mux.Handle("GET /shortcuts/", admin(http.HandlerFunc(s.HandleListShortcuts)))
	mux.Handle("GET /shortcuts/{id}", admin(http.HandlerFunc(s.HandleGetShortcut)))
	mux.Handle("PUT /shortcuts/{id}", admin(http.HandlerFunc(s.HandlePutShortcut)))
	mux.Handle("DELETE /shortcuts/{id}", admin(http.HandlerFunc(s.HandleDeleteShortcut)))
}

// adminHandler returns the AdminMiddleware if configured, otherwise a no-op passthrough.
func (s *Server) adminHandler(h http.Handler) http.Handler {
	if s.AdminMiddleware != nil {
		return s.AdminMiddleware(h)
	}
	return h
}
