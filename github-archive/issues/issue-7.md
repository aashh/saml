# Issue #7: Open redirect via RelayState

**State:** OPEN  
**Created:** 2026-02-08T05:49:49Z  
**Updated:** 2026-02-08T05:49:49Z  
**URL:** https://github.com/aashh/saml/issues/7

## Description

When `AllowIDPInitiated` is enabled, the `RelayState` POST parameter is used as the redirect target without validation:

```go
if uri := r.Form.Get("RelayState"); uri != "" {
    redirectURI = uri
}
http.Redirect(w, r, redirectURI, http.StatusFound)
```

An attacker can set `RelayState=https://evil.com` and the SP will redirect there after successful authentication. Classic open redirect -- useful for phishing because the URL starts at the legitimate SP domain.

The fix validates that the redirect URI is a relative path (starts with `/`, not `//`).

Fix: [`fix/issue-07`](../tree/fix/issue-07)

