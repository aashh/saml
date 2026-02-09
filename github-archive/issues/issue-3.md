# Issue #3: InResponseTo not validated when AllowIDPInitiated is set

**State:** OPEN
**Created:** 2026-02-08T05:49:45Z
**Updated:** 2026-02-08T05:49:45Z
**URL:** https://github.com/aashh/saml/issues/3

## Description

When `AllowIDPInitiated` is true, the `InResponseTo` validation is not performed -- even for SP-initiated flows that have an `InResponseTo` value. This means enabling IdP-initiated SSO also disables request binding validation for SP-initiated flows, allowing assertion replay attacks.

The fix: only skip validation when `InResponseTo` is empty AND `AllowIDPInitiated` is true (i.e., a genuine unsolicited response). When `InResponseTo` is present and non-empty, validate it against tracked request IDs regardless of the flag.

Some IdPs (notably Rippling) set `InResponseTo` to arbitrary values on IdP-initiated flows. The existing `ValidateRequestID` hook handles that case.

Fix: [`fix/issue-03`](../tree/fix/issue-03)

