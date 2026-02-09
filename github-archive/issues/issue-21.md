# Issue #21: Signature error returned after content processing (information leak)

**State:** OPEN  
**Created:** 2026-02-08T05:50:03Z  
**Updated:** 2026-02-08T05:50:03Z  
**URL:** https://github.com/aashh/saml/issues/21

## Description

`parseResponse` computes the signature validation result early but defers returning it until after parsing response attributes (Issuer, Destination, timestamps). The code even has a TODO acknowledging this:

```go
// TODO(ross): adjust the test cases so that we can abort here
// if the Response signature is invalid.
```

This lets an attacker probe whether their XML structure is correct: "Issuer mismatch" means parsing succeeded; "Invalid Signature" means the Issuer was right. Useful for iterating toward a working XSW payload.

The fix returns signature errors immediately, before processing any content. Test cases updated accordingly.

Fix: [`fix/issue-21`](../tree/fix/issue-21)

