# Issue #4: Double XML parsing via etree serialization round-trip

**State:** CLOSED  
**Created:** 2026-02-08T05:49:46Z  
**Updated:** 2026-02-08T05:54:28Z  
**URL:** https://github.com/aashh/saml/issues/4

## Description

The codebase bridges `etree` (DOM/signatures) and `encoding/xml` (struct marshaling) by serializing elements to bytes and re-parsing them. This is a performance cost and a potential consistency risk if the two serializers ever diverge.

That said, this pattern actually *helps* with XSW mitigation -- the struct gets populated from the exact verified element, not from a second parse of the original document. Fixing it would mean committing to one XML library and rewriting everything that uses the other. That's a full architectural rewrite with no clear security upside and a real security downside.

Not planning to change this.

