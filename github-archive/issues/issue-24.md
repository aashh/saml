# Issue #24: OneTimeUse and ProxyRestriction conditions silently ignored

**State:** OPEN  
**Created:** 2026-02-08T05:50:06Z  
**Updated:** 2026-02-08T05:50:06Z  
**URL:** https://github.com/aashh/saml/issues/24

## Description

The `Conditions` struct parses `<OneTimeUse>` and `<ProxyRestriction>`, but `validateAssertion` never looks at them.

For `OneTimeUse`: without a replay cache (#1), we can't actually enforce single-use semantics. But silently accepting an assertion the IdP explicitly marked as one-time-use is worse than rejecting it honestly. The fix rejects assertions carrying `OneTimeUse` unless the caller explicitly opts out via `IgnoreOneTimeUse`.

For `ProxyRestriction`: the fix rejects assertions carrying this condition, since the SP doesn't implement proxy counting.

Fix: [`fix/issue-24`](../tree/fix/issue-24)

