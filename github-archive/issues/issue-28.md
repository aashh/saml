# Issue #28: AuthnContext ClassRef not validated against RequestedAuthnContext

**State:** OPEN  
**Created:** 2026-02-08T05:50:10Z  
**Updated:** 2026-02-08T21:32:46Z  
**URL:** https://github.com/aashh/saml/issues/28

## Description

The SP can request a specific authentication strength via `RequestedAuthnContext` (e.g., MFA), but `validateAssertion` never checks whether the IdP actually honored the request. If the IdP returns a password-only session when MFA was requested, the SP accepts it without complaint.

The fix checks `AuthnContextClassRef` in the assertion against `sp.RequestedAuthnContext` when the latter is configured. No check when it's nil (backward compatible). Only `exact` Comparison (the SAML default) is supported; other Comparison values (`minimum`, `maximum`, `better`) are rejected with a clear error rather than silently treated as `exact`.

Fix: [`fix/issue-28`](../tree/fix/issue-28)

