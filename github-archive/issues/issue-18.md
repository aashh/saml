# Issue #18: AudienceRestriction only holds one Audience

**State:** OPEN  
**Created:** 2026-02-08T05:50:00Z  
**Updated:** 2026-02-08T05:50:00Z  
**URL:** https://github.com/aashh/saml/issues/18

## Description

The SAML schema defines `<Audience>` with `maxOccurs="unbounded"`, but the Go struct has:

```go
type AudienceRestriction struct {
    Audience Audience // singular
}
```

When an assertion contains multiple `<Audience>` elements, `encoding/xml` maps only the last one. If the SP's audience URI isn't last, validation fails (false rejection). Conversely, an attacker who controls element ordering could hide the intended audience from the validator.

The fix changes `Audience Audience` to `Audiences []Audience` and updates `validateAudienceRestriction` to check whether any entry matches. This is a breaking API change for anyone accessing the field directly, but it's necessary for spec compliance.

Fix: [`fix/issue-18`](../tree/fix/issue-18)

