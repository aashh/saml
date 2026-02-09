# Issue #9: RSA-PKCS1v1.5 decrypter registered by default

**State:** OPEN  
**Created:** 2026-02-08T05:49:51Z  
**Updated:** 2026-02-08T05:49:51Z  
**URL:** https://github.com/aashh/saml/issues/9

## Description

`init()` in `xmlenc/pubkey.go` registers both OAEP and PKCS1v1.5. Since the library accepts whatever `EncryptionMethod` appears in the incoming XML, an attacker who can intercept and modify encrypted assertions can downgrade the key transport to PKCS1v1.5 and mount a Bleichenbacher attack to recover the plaintext.

The fix removes `RegisterDecrypter(PKCS1v15())` from `init()`. The constructor remains available for consumers who explicitly need legacy compatibility.

Fix: [`fix/issue-09`](../tree/fix/issue-09)

