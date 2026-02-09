# Issue #10: Signature KeyInfo not always stripped before validation

**State:** OPEN  
**Created:** 2026-02-08T05:49:52Z  
**Updated:** 2026-02-08T05:49:52Z  
**URL:** https://github.com/aashh/saml/issues/10

## Description

`validateSignature` only removes `KeyInfo` when `X509Certificate` is absent:

```go
if el.FindElement("./Signature/KeyInfo/X509Data/X509Certificate") == nil {
    // remove KeyInfo
}
```

When `X509Certificate` is present, `KeyInfo` is retained. While `goxmldsig` validates embedded certificates against the trusted store (preventing direct bypass), stripping `KeyInfo` unconditionally is cleaner — it ensures validation relies solely on metadata-provided certificates regardless of `goxmldsig` behavior.

Separately, `goxmldsig` only auto-selects a certificate from the store when there is exactly one root. When an IdP has multiple signing certificates in its metadata, validation fails even for valid signatures.

The fix always strips `KeyInfo` and restructures the validation loop to try each trusted certificate individually.

Fix: [`fix/issue-10`](../tree/fix/issue-10)

