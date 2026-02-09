# Issue #26: OAEP uses wrong hash for MGF1

**State:** OPEN  
**Created:** 2026-02-08T05:50:08Z  
**Updated:** 2026-02-08T21:32:18Z  
**URL:** https://github.com/aashh/saml/issues/26

## Description

The `rsa-oaep-mgf1p` algorithm URI (2001 spec) mandates MGF1 with SHA-1 as the mask generation function. Go's `rsa.EncryptOAEP` uses a single hash for both the label digest and MGF1, so when `OAEP()` defaults to SHA-256, the MGF is also SHA-256. Strict SAML implementations expecting MGF1-SHA1 will fail to decrypt.

The fix makes `OAEP()` use SHA-1 (matching the spec for `rsa-oaep-mgf1p`) and adds `OAEP_SHA256()` for the 2009 URI where SHA-256 for both is correct.

The IdP default encryption is also updated: the old code hardcoded AES-128-CBC, the new code uses `OAEP()` defaults which include AES-256-CBC. A new `Encryptor *xmlenc.RSA` field on `IdentityProvider` allows explicit configuration if the default needs to be overridden.

Fix: [`fix/issue-26`](../tree/fix/issue-26)

