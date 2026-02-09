# Issue #16: IdP session cookie Secure flag unreliable, SameSite missing

**State:** OPEN  
**Created:** 2026-02-08T05:49:58Z  
**Updated:** 2026-02-08T05:49:58Z  
**URL:** https://github.com/aashh/saml/issues/16

## Description

In `samlidp/session.go`:

```go
Secure: r.URL.Scheme == "https",
```

`r.URL.Scheme` is usually empty in Go's `http.Server` unless explicitly set by middleware. This means `Secure` evaluates to `false` and the session cookie gets sent over plaintext HTTP. The `SameSite` attribute is also unset, which legacy browsers interpret as `None`.

The fix checks `r.TLS != nil` instead (reliable without reverse-proxy cooperation) and sets `SameSite: Lax`.

Fix: [`fix/issue-16`](../tree/fix/issue-16)

