# Issue #12: Empty Destination accepted on unsigned responses

**State:** OPEN  
**Created:** 2026-02-08T05:49:54Z  
**Updated:** 2026-02-08T05:49:54Z  
**URL:** https://github.com/aashh/saml/issues/12

## Description

The Destination check in `parseResponse` is gated on `responseHasSignature || response.Destination != ""`. When the response envelope is unsigned (common -- many IdPs only sign the assertion) and the `Destination` attribute is absent, validation is not performed.

For unsigned responses, `Destination` is the only non-cryptographic binding to the intended SP. Without this binding, a captured response could be replayed to a different SP.

The fix requires `Destination` to always be present. The SAML spec says MUST for signed responses and SHOULD otherwise; the fix enforces MUST in both cases for defense in depth.

Fix: [`fix/issue-12`](../tree/fix/issue-12)

