# Issue #20: IdP hardcodes AES-128-CBC and SHA-1 for assertion encryption

**State:** OPEN  
**Created:** 2026-02-08T05:50:02Z  
**Updated:** 2026-02-08T05:50:02Z  
**URL:** https://github.com/aashh/saml/issues/20

## Description

`MakeAssertionEl()` in `identity_provider.go` hardcodes:

```go
encryptor := xmlenc.OAEP()
encryptor.BlockCipher = xmlenc.AES128CBC
encryptor.DigestMethod = &xmlenc.SHA1
```

No way to configure this. Every SP in the federation gets AES-128-CBC + SHA-1 regardless of what it supports.

The fix adds a configurable `Encryptor *xmlenc.RSA` field to `IdentityProvider` with a nil-means-default pattern. The default is upgraded to AES-256-CBC. Existing code that doesn't set the field gets the new default; anyone who needs the old behavior can set it explicitly.

Fix: [`fix/issue-20`](../tree/fix/issue-20)

