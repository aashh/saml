# Issue #22: samlidp admin endpoints have no authentication

**State:** OPEN  
**Created:** 2026-02-08T05:50:04Z  
**Updated:** 2026-02-08T05:50:04Z  
**URL:** https://github.com/aashh/saml/issues/22

## Description

The `samlidp` package exposes admin endpoints (`PUT /users/{id}`, `DELETE /users/{id}`, `PUT /services/{id}`, etc.) without authentication.

While `samlidp` is documented as a testing/demo package, adding an optional authentication mechanism improves security for deployments in non-isolated environments.

The fix adds an `AdminMiddleware func(http.Handler) http.Handler` field. When set, admin endpoints are wrapped with the middleware. When nil, the default behavior is preserved for backward compatibility.

Fix: [`fix/issue-22`](../tree/fix/issue-22)

