# Issue #5: TripleDES registered as a decrypter by default

**State:** OPEN  
**Created:** 2026-02-08T05:49:47Z  
**Updated:** 2026-02-08T05:49:47Z  
**URL:** https://github.com/aashh/saml/issues/5

## Description

TripleDES-CBC is registered in `init()` and will be used if an incoming assertion specifies it as the encryption method. The 64-bit block size makes it vulnerable to birthday attacks (Sweet32).

The fix removes TripleDES from default registration. Consumers who need it for legacy IdP compatibility can opt back in with `xmlenc.RegisterDecrypter(xmlenc.TripleDES)`. Also adds `xmlenc.UnregisterDecrypter()` for runtime algorithm management.

Fix: [`fix/issue-05`](../tree/fix/issue-05)

