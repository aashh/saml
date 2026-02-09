# Issue #6: Session JWTs signed with the SAML SP key

**State:** OPEN  
**Created:** 2026-02-08T05:49:48Z  
**Updated:** 2026-02-08T05:49:48Z  
**URL:** https://github.com/aashh/saml/issues/6

## Description

`samlsp.New` reuses the SAML signing key (`opts.Key`) to sign session JWTs. If either the JWT signing path or the SAML signing path has a vulnerability, both are compromised simultaneously.

The fix adds an optional `SessionKey crypto.Signer` field to `Options`. When set, it's used for session JWTs; when nil, falls back to `opts.Key` for backward compatibility.

Fix: [`fix/issue-06`](../tree/fix/issue-06)

