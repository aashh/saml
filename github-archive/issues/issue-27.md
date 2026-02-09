# Issue #27: Logout redirect binding validates embedded XML signature instead of detached URL signature

**State:** OPEN  
**Created:** 2026-02-08T05:50:09Z  
**Updated:** 2026-02-08T05:50:09Z  
**URL:** https://github.com/aashh/saml/issues/27

## Description

`ValidateLogoutResponseRedirect` parses the XML from the URL parameter and calls `validateSignature` on it -- but the HTTP-Redirect binding doesn't use embedded XML signatures. The signature is detached: it signs the raw query string (`SAMLResponse=...&RelayState=...&SigAlg=...`) and is sent as separate URL parameters.

This means legitimate logout responses from compliant IdPs fail validation (no embedded signature to find), while an attacker can forge a logout response with a self-signed embedded XML signature and it'll pass.

Fixing this requires implementing detached URL query parameter signature verification, which is a different verification path from what exists today. Closely related to #19 (the middleware doesn't even route to the SLO handler). Both should be addressed together.

Deferred for now.

