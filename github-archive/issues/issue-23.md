# Issue #23: SubjectConfirmation Method not checked

**State:** OPEN  
**Created:** 2026-02-08T05:50:05Z  
**Updated:** 2026-02-08T05:50:05Z  
**URL:** https://github.com/aashh/saml/issues/23

## Description

`validateAssertion` iterates over `SubjectConfirmations` but never checks the `Method` attribute. The SAML Web SSO profile requires `urn:oasis:names:tc:SAML:2.0:cm:bearer`. Other methods like Holder-of-Key require the presenter to prove possession of a cryptographic key.

By ignoring `Method`, the SP treats everything as a bearer token. If an attacker intercepts a Holder-of-Key assertion, this SP accepts it without the key proof.

The fix skips non-bearer subject confirmations and rejects the assertion if no bearer confirmation is found.

Fix: [`fix/issue-23`](../tree/fix/issue-23)

