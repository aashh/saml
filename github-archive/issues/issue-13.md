# Issue #13: AES-GCM Encrypt has multiple bugs preventing decryption

**State:** OPEN
**Created:** 2026-02-08T05:49:55Z
**Updated:** 2026-02-08T21:31:48Z
**URL:** https://github.com/aashh/saml/issues/13

## Description

Four bugs in `xmlenc/gcm.go` `Encrypt()` prevent successful round-trip encryption/decryption:

1. `aesgcm.Seal(nil, nonce, ciphertext, nil)` encrypts a zero-filled buffer instead of `plaintext`
2. The `nonce` variable is shadowed inside an `if` block (`:=` vs `=`), causing generated nonces to be discarded
3. The nonce is not prepended to the output, but `Decrypt` expects `nonce || ciphertext || tag`
4. Calls `appendPadding`, which is unnecessary for GCM (a stream cipher) and prevents decryption

The existing test only verified encryption without attempting decryption, so these bugs were not caught.

Fix: [`fix/issue-13`](../tree/fix/issue-13)

