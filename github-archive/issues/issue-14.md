# Issue #14: IdP does not verify AuthnRequest signatures

**State:** OPEN  
**Created:** 2026-02-08T05:49:56Z  
**Updated:** 2026-02-08T05:49:56Z  
**URL:** https://github.com/aashh/saml/issues/14

## Description

The IdP has a `TODO(ross): support signed authn requests` and either errors out (if metadata requires signing) or silently ignores the signature (if metadata is permissive). An attacker can intercept a signed AuthnRequest, strip or replace the signature, modify the `AssertionConsumerServiceURL`, and the IdP will happily issue an assertion to the attacker's endpoint.

Fixing this requires implementing two signature verification paths: embedded XML signatures for POST binding and detached query-string signatures for Redirect binding. That's a significant chunk of new code (~100+ lines), separate from the targeted fixes we're working on now.

Deferred for now, but it's a real gap.

