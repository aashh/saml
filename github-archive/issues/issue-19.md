# Issue #19: SLO endpoints advertised in metadata but not handled

**State:** OPEN  
**Created:** 2026-02-08T05:50:01Z  
**Updated:** 2026-02-08T05:50:01Z  
**URL:** https://github.com/aashh/saml/issues/19

## Description

The SP generates `SingleLogoutService` endpoints in its metadata (based on `SloURL`), but the middleware never routes requests there -- they 404. When an IdP tries to propagate a global logout, this SP does not process the request and the user's session stays active.

Fixing this properly means implementing LogoutRequest/LogoutResponse parsing, session invalidation, and response generation. It's closely related to #27 (redirect binding signature verification requires fixes). Both should be addressed together as part of a broader SLO implementation.

Deferred for now.

