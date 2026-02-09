# Issue #11: Library panics instead of returning errors

**State:** CLOSED  
**Created:** 2026-02-08T05:49:53Z  
**Updated:** 2026-02-08T05:54:28Z  
**URL:** https://github.com/aashh/saml/issues/11

## Description

Several utility functions use `panic` for error conditions: `randomBytes` panics on entropy failure, `HandleStartAuthFlow` panics if the request path matches the ACS URL, `defaultSigningMethodForKey` panics on unsupported key types.

After looking at this more carefully, these panics catch programming errors (misconfiguration, unsupported key types) early rather than letting them propagate silently. That's a reasonable trade-off and consistent with Go idioms for genuinely unrecoverable states. The entropy-exhaustion panic in particular is preferable to silently generating weak random values.

Not planning to change this.

