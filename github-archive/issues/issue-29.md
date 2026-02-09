# Issue #29: SubjectConfirmationData NotBefore not checked

**State:** OPEN  
**Created:** 2026-02-08T05:50:11Z  
**Updated:** 2026-02-08T05:50:11Z  
**URL:** https://github.com/aashh/saml/issues/29

## Description

`validateAssertion` checks `NotOnOrAfter` on `SubjectConfirmationData` but not `NotBefore`. The `Conditions` element has its own time window check, but the SAML Core spec (Section 2.4.1.2) gives `SubjectConfirmationData` its own `NotBefore` attribute that should be independently validated.

The fix mirrors the existing `Conditions.NotBefore` pattern, with a zero-value check for backward compatibility with assertions that don't set it.

Fix: [`fix/issue-29`](../tree/fix/issue-29)

