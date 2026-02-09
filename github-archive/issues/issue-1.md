# Issue #1: No replay protection for assertion IDs

**State:** OPEN  
**Created:** 2026-02-08T05:49:43Z  
**Updated:** 2026-02-08T05:49:43Z  
**URL:** https://github.com/aashh/saml/issues/1

## Description

The SP doesn't track seen assertion IDs, so a valid SAMLResponse can be replayed within its validity window (~4.5 minutes given `MaxIssueDelay` + `MaxClockSkew`). The SAML Core spec requires SPs to ensure an assertion is used at most once.

The right fix is probably an optional `ReplayCache` interface on `ServiceProvider` -- nil means no change in behavior (preserving backward compat), non-nil means assertion IDs get checked. An in-memory implementation covers single-instance deployments; anyone running horizontally can implement the interface against Redis or a database.

This is a real gap but the fix needs to be designed carefully to avoid forcing a storage dependency on everyone. Planned, not yet implemented.

Worth noting: #24 already rejects assertions carrying the `OneTimeUse` condition, which is the honest thing to do until replay detection exists.

