# Issue #25: Session cookie scoped to all subdomains by default

**State:** OPEN  
**Created:** 2026-02-08T05:50:07Z  
**Updated:** 2026-02-08T05:50:07Z  
**URL:** https://github.com/aashh/saml/issues/25

## Description

`CookieSessionProvider` explicitly sets the `Domain` attribute on session cookies after stripping the port. This makes them wildcard cookies valid for `*.example.com` rather than host-only cookies valid for just `example.com`.

If an attacker compromises any subdomain (dev boxes, user-content hosts, etc.), they can read the main application's session cookies.

The fix only sets `Domain` when explicitly configured via `CookieDomain`. When omitted, browsers default to host-only scope. This is a behavioral change for deployments that rely on subdomain cookie sharing, but those deployments can set `CookieDomain` explicitly.

Fix: [`fix/issue-25`](../tree/fix/issue-25)

