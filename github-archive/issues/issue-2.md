# Issue #2: Padding oracle in CBC decryption

**State:** OPEN  
**Created:** 2026-02-08T05:49:44Z  
**Updated:** 2026-02-08T05:49:44Z  
**URL:** https://github.com/aashh/saml/issues/2

## Description

`stripPadding()` in `xmlenc/cbc.go` uses early returns with distinct error messages for different padding failures. If an SP surfaces these errors (via HTTP status codes, response bodies, or timing), an attacker can decrypt encrypted assertions byte-by-byte.

The fix rewrites `stripPadding` to validate padding in constant time using `crypto/subtle`, and unifies all decryption error messages so callers can't distinguish padding errors from other failures.

Fix: [`fix/issue-02`](../tree/fix/issue-02)

