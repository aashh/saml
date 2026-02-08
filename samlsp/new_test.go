package samlsp

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"gotest.tools/assert"
)

func TestNewCanAcceptCookieName(t *testing.T) {

	testCases := []struct {
		testName   string
		cookieName string
		expected   string
	}{
		{"Works with alt name", "altCookie", "altCookie"},
		{"Works with default", "", "token"},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			opts := Options{
				CookieName: tc.cookieName,
			}
			sp, err := New(opts)
			assert.Assert(t, err)
			cookieProvider := sp.Session.(CookieSessionProvider)
			assert.Equal(t, tc.expected, cookieProvider.Name)

		})
	}

}

func TestSessionKeyFallsBackToKey(t *testing.T) {
	samlKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Assert(t, err)

	opts := Options{
		Key: samlKey,
	}
	codec := DefaultSessionCodec(opts)
	assert.Equal(t, crypto.Signer(samlKey), codec.Key)

	tracker := DefaultTrackedRequestCodec(opts)
	assert.Equal(t, crypto.Signer(samlKey), tracker.Key)
}

func TestSessionKeyUsesSessionKeyWhenSet(t *testing.T) {
	samlKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Assert(t, err)
	sessKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Assert(t, err)

	opts := Options{
		Key:        samlKey,
		SessionKey: sessKey,
	}
	codec := DefaultSessionCodec(opts)
	assert.Equal(t, crypto.Signer(sessKey), codec.Key)

	tracker := DefaultTrackedRequestCodec(opts)
	assert.Equal(t, crypto.Signer(sessKey), tracker.Key)
}
