# PR Submission Plan for crewjam/saml

Date: 2026-02-09 (Revised)
Target: https://github.com/crewjam/saml

## Strategy Overview

Submit 22 security/correctness fixes incrementally to build trust with maintainer(s).

**Key Insight:** Lead with high-value critical fixes, not easy trivial ones. The first PR
sets the tone for how the maintainer perceives you. A well-tested critical security fix
signals "serious contributor worth engaging with." A typo fix signals "low-value noise."

### Why Lead with Critical Issues (Not Easy Ones)

**The flawed "easy first" approach:**
- ❌ PR #1: Fix typo in example code (< 5 min review)
- Maintainer thinks: "This person sends trivial stuff, review when I have spare time"
- Gets categorized as low-priority noise
- Even when you send critical fixes later, you're already in the "trivial contributor" bucket

**The correct "value first" approach:**
- ✅ PR #1: Fix signature validation bypass (15 min review, critical security)
- Maintainer thinks: "This person found a real vulnerability, I need to review this NOW"
- Gets categorized as serious security contributor
- Follow-up PRs get benefit of the doubt and faster review

**Maintainer priority queue:**
1. **Severity** - Critical security > nice-to-have improvements
2. **Trust** - Known contributor > unknown
3. **Review burden** - Well-tested complex > vague simple

A well-tested critical fix with excellent description is **easier to review** than a
vague one-liner because:
- Tests prove the vulnerability exists
- Tests prove the fix works
- Description provides full context
- Diff is focused (one issue, one fix)

**Counter-intuitive truth:** Complex critical fixes get reviewed **faster** than trivial
ones because busy maintainers prioritize by impact, not difficulty.

## Submission Waves

### Wave 1: Critical Security Fixes (Submit First)

**Goal:** Establish yourself as a serious security contributor who finds real vulnerabilities

**Why lead with critical fixes:**
- Maintainers prioritize critical security issues (get reviewed immediately)
- Sets expectations: "This person finds serious issues, not trivial ones"
- Proves expertise: deep understanding of SAML, XML signatures, attack vectors
- Well-tested complex fixes are easier to review than vague simple ones
- After these merge, smaller PRs get benefit of the doubt

| Order | Issue | Branch | Why First | Severity | Breaking? |
|-------|-------|--------|-----------|----------|-----------|
| 1 | #03 | fix/issue-03 | Authentication bypass - InResponseTo with AllowIDPInitiated | **High** | No |
| 2 | #12 | fix/issue-12 | Replay attack - Destination not required on unsigned responses | **High** | Yes (behavioral) |

**Timing:** Submit with 3-5 day gaps. Wait for merge or substantive feedback before next.

**Signal sent:**
- "I found multiple real vulnerabilities in core authentication flows"
- "I understand SAML protocol deeply (not just running scanners)"
- "I write thorough tests that prove vulnerabilities and fixes"
- "I understand codebase quirks (goxmldsig, validation hooks)"
- "I make good security trade-offs (willing to propose breaking changes when justified)"

---

### Wave 2: Additional Security Fixes

**Goal:** Continue demonstrating security expertise with high-value fixes

**Prerequisites:** At least 1-2 Wave 1 PRs merged or positive engagement

| Order | Issue | Branch | Description | Severity | Breaking? |
|-------|-------|--------|-------------|----------|-----------|
| 4 | #07 | fix/issue-07 | Open redirect via RelayState in IDP-initiated SSO | **High** | No |
| 5 | #21 | fix/issue-21 | Timing oracle - signature errors returned after content processing | **Medium** | Yes (error order) |
| 6 | #23 | fix/issue-23 | SubjectConfirmation Method not validated (holder-of-key bypass) | **High** | No |
| 7 | #02 | fix/issue-02 | Padding oracle in CBC decryption (constant-time validation) | **Critical** | No |

**Timing:** After Wave 1 acceptance, submit 2-3 at a time with week-long gaps

**Note:** #02 is critical but intricate (constant-time crypto). Submit after initial trust established.

---

### Wave 3: Spec Compliance & Validation Gaps

**Goal:** Address remaining validation and spec compliance issues

**Prerequisites:** Multiple Wave 1-2 PRs merged, trust established

| Order | Issue | Branch | Description | Severity | Breaking? |
|-------|-------|--------|-------------|----------|-----------|
| 8 | #28 | fix/issue-28 | AuthnContext ClassRef not validated against RequestedAuthnContext | **Medium** | No |
| 9 | #29 | fix/issue-29 | SubjectConfirmationData NotBefore not checked | **Low-Med** | No |
| 10 | #13 | fix/issue-13 | AES-GCM Encrypt bugs (4 interlocking bugs) | **Critical** | No |
| 11 | #06 | fix/issue-06 | Session JWTs signed with SAML SP key (optional SessionKey) | **Medium** | No |
| 12 | #05 | fix/issue-05 | TripleDES registered by default (Sweet32) | **Medium** | No |
| 13 | #10 | fix/issue-10 | Defense-in-depth KeyInfo stripping + multi-cert IdP support | **Medium** | No |

**Timing:** After Wave 2 acceptance, submit 2-3 at a time

**Note:** #13 is critical but GCM may not be widely used. Better after trust established.

---

### Wave 4: Breaking Changes & API Modifications

**Goal:** Propose necessary breaking changes with good justification

**Prerequisites:** High trust established (10+ PRs merged), maintainer trusts your judgment

| Order | Issue | Branch | Description | Breaking? | Notes |
|-------|-------|--------|-------------|-----------|-------|
| 14 | #09 | fix/issue-09 | RSA-PKCS1v1.5 not registered by default (Bleichenbacher) | No (adds opt-in) | Similar to #05 |
| 15 | #24 | fix/issue-24 | Reject OneTimeUse/ProxyRestriction conditions | Yes - adds IgnoreOneTimeUse | Spec strictness |
| 16 | #25 | fix/issue-25 | Host-only cookie scope by default | Yes - behavioral | Security improvement |
| 17 | #18 | fix/issue-18 | AudienceRestriction multiple Audiences | **YES - API** (.Audience → .Audiences[]) | Major breaking |

**Timing:** Submit individually with 1-2 week gaps. Get buy-in for each before next.

**Note:** #18 is the most disruptive API change - save for last, require strong trust.

---

### Wave 5: Enhancements & Lower-Priority Fixes

**Goal:** Quality-of-life improvements and test code enhancements

**Prerequisites:** Fully trusted contributor (15+ PRs merged)

| Order | Issue | Branch | Description | Breaking? | Priority |
|-------|-------|--------|-------------|-----------|----------|
| 18 | #20 | fix/issue-20 | IdP Encryptor field + AES-256-CBC default | Yes - defaults | Medium |
| 19 | #26 | fix/issue-26 | OAEP SHA-1 for rsa-oaep-mgf1p spec compliance | **YES - OAEP() defaults** | Medium |
| 20 | #08 | fix/issue-08 | MemoryStore size cap with FIFO eviction (samlidp) | No | Low (test code) |
| 21 | #22 | fix/issue-22 | AdminMiddleware for samlidp endpoints (samlidp) | No | Low (test code) |
| 22 | #16 | fix/issue-16 | Cookie Secure flag + SameSite (samlidp) | No | Low (test code) |
| 23 | #15 | fix/issue-15 | Use crypto/rand in example code | No | Low (example) |

**Timing:** Can batch submit 2-3 at a time. Lower priority = less urgency.

**Note:** #20 and #26 interact (both affect OAEP defaults). Submit together or in sequence.
Samlidp and example code fixes saved for last - valuable but not critical.

---

## Wave 1 PRs

### PR #659: fix/issue-03 (InResponseTo Validation)

Open on crewjam/saml. All CI passing.

---

### PR #660: fix/issue-12 (Destination Required)

Open on crewjam/saml. All CI passing.

---

### PR #2: fix/issue-03 (InResponseTo Validation)

**Files:** `service_provider.go`, `service_provider_test.go`
**Change:** Only skip InResponseTo validation when empty AND AllowIDPInitiated=true
**Tests:** Both response-level and assertion-level validation covered

**PR Description Template:**
```markdown
# Fix authentication bypass with AllowIDPInitiated

## Issue

When `AllowIDPInitiated` is true, InResponseTo validation is skipped at both the
Response and Assertion levels - even for SP-initiated flows that have an
InResponseTo value. This allows replay attacks against SPs that enable
IdP-initiated SSO.

## Root Cause

The validation logic treated AllowIDPInitiated as "skip all InResponseTo checks"
rather than "allow empty InResponseTo for unsolicited responses."

## Fix

Only skip InResponseTo validation when both conditions are met:
1. InResponseTo is empty (genuinely unsolicited)
2. AllowIDPInitiated is true

When InResponseTo is present and non-empty, validate it against tracked request
IDs regardless of the AllowIDPInitiated flag.

## Edge Case

Some IdPs (notably Rippling) set InResponseTo to arbitrary values in IdP-initiated
flows. Users can handle this with the existing ValidateRequestID hook.

## Changes
- service_provider.go: Fix validation logic in validateAssertion and parseResponse
- service_provider_test.go: Add TestSPValidateRequestIDWithAllowIDPInitiated

## Testing
- AllowIDPInitiated=true, InResponseTo="": accepted (IdP-initiated)
- AllowIDPInitiated=true, InResponseTo="wrong-id": rejected (replay attack blocked)
- AllowIDPInitiated=true, InResponseTo="valid-id": accepted (SP-initiated works)
- AllowIDPInitiated=false: existing behavior preserved

Fixes #3
```

**Why second:**
- **Authentication bypass** - Real security issue in common configuration
- **Subtle but serious** - Shows deep understanding of SAML flows
- **Well-tested** - Covers both validation layers
- **Handles edge case** - Mentions Rippling quirk and ValidateRequestID hook
- **Demonstrates:**
  - SP-initiated vs IdP-initiated flow understanding
  - Complex flag interaction knowledge
  - Real-world deployment awareness

**Maintainer thinks:** "They understand SAML protocol flows deeply and found a
subtle validation gap in a common configuration."

---

### PR #3: fix/issue-12 (Destination Required)

**Files:** `service_provider.go`, `service_provider_test.go`
**Change:** Require Destination on all responses (remove conditional check)
**Tests:** Unsigned response without Destination rejected

**PR Description Template:**
```markdown
# Require Destination on all SAML responses

## Issue

The Destination check in `parseResponse` is gated on:
```go
responseHasSignature || response.Destination != ""
```

When the response envelope is unsigned (common - many IdPs only sign the
assertion) and Destination is absent, validation is not performed. This allows
captured responses to be replayed to different SPs.

## Security Impact

For unsigned responses, Destination is the only non-cryptographic binding to
the intended SP. Without this check, an attacker who captures a response can
replay it to any SP that trusts the same IdP.

## Fix

Require Destination unconditionally. The SAML spec says Destination MUST be
present for signed responses and SHOULD be present otherwise. This PR enforces
MUST in both cases for defense-in-depth.

## Breaking Change

This is technically breaking: unsigned responses without Destination that were
previously accepted will now be rejected. However:
- This configuration is spec non-compliant
- It represents a real security gap
- Most IdPs include Destination by default

## Changes
- service_provider.go: Remove conditional, always require Destination
- service_provider_test.go: Update tests, add unsigned-without-Destination test

## Testing
- Signed response with Destination: accepted
- Unsigned response with Destination: accepted
- Unsigned response without Destination: now rejected (breaking)
- All existing tests pass (all had Destination)

Fixes #12
```

**Why third:**
- **Replay attack prevention** - Real security issue
- **Breaking change** - Shows you understand implications and make good trade-offs
- **Well-justified** - Clear security rationale
- **Spec compliance** - Upgrades SHOULD to MUST
- **Demonstrates:**
  - Defense-in-depth thinking
  - Willingness to propose necessary breaking changes
  - Good security judgment

**Maintainer thinks:** "This person makes thoughtful security trade-offs and
isn't afraid to propose justified breaking changes."

---

## Pacing Rules

1. **Don't submit faster than maintainer can review**
   - If maintainer takes 1 week to review, wait 1 week between submissions
   - Match their pace

2. **Wait for merge or substantive feedback**
   - Don't submit Wave 2 until at least 1 Wave 1 PR is merged
   - Substantive feedback (even if not merged) shows engagement

3. **Max 3 PRs without response**
   - Don't queue more than 3 PRs before getting merge/feedback
   - Shows respect for their review bandwidth

4. **Breaking changes submitted individually**
   - #12, #14, #15, #16, #18, #20, #26 should be solo submissions
   - Give maintainer space to consider impact

5. **Response timing**
   - Respond to feedback within 24 hours if possible
   - Implement requested changes eagerly
   - Never be defensive

## Signals to Send

✅ **Do:**
- "I'm here long-term" - Mention you have more fixes if helpful
- "I respect your judgment" - Don't argue about approach/priority
- "I care about this project" - Engage with other issues/PRs
- "I'm flexible" - Willing to adjust approach based on feedback

❌ **Don't:**
- Mention you have 22 fixes total (overwhelming)
- Say "I found 22 vulnerabilities" (sounds like criticism)
- Submit PRs faster than they review (disrespectful)
- Argue about severity or urgency (pushy)
- Be defensive about original code (rude)

## If Things Go Wrong

**Scenario: No response after 2-3 weeks on easy PR**
- Open an issue asking: "Would you be open to a PR that fixes X?"
- If still no response, project may be inactive/unmaintained
- Consider alternative: fork with your fixes, or find different project

**Scenario: Maintainer rejects approach**
- Be flexible: "Would you prefer I do Y instead?"
- Ask clarifying questions about their vision
- Be willing to walk away from contentious fixes
- Don't take it personally

**Scenario: Maintainer wants batch review**
- Some maintainers prefer reviewing multiple PRs at once
- If they say "submit the rest", do so
- Otherwise stick to incremental approach

**Scenario: Maintainer has strong opinions that conflict**
- Example: "We don't support OneTimeUse so rejecting it is wrong"
- Acknowledge their perspective: "That's a fair point about..."
- Offer compromise: "Would X approach work better?"
- If irreconcilable, close your PR gracefully

## Success Metrics

**After Wave 1 (3 PRs):**
- ✅ At least 1 merged
- ✅ Positive feedback tone
- ✅ Maintainer engaged (asks questions, gives feedback)

**After Wave 2 (7 PRs total):**
- ✅ 3-5 merged
- ✅ Maintainer references you by name
- ✅ Review time decreasing (trust building)

**After Wave 3 (12 PRs total):**
- ✅ 7-10 merged
- ✅ Maintainer accepts more complex PRs
- ✅ Less pushback on behavioral changes

**After Wave 4-5 (22 PRs total):**
- ✅ 15-20 merged
- ✅ Maintainer trusts your judgment on breaking changes
- ✅ You're a recognized contributor to the project

## Current Status

- **Fork:** https://github.com/aashh/saml
- **All 22 branches:** Pushed to fork
- **All 29 issues:** Created in fork with professional tone
- **Wave 1 PRs open:** #659 (fix/issue-03), #660 (fix/issue-12) — all CI passing
- **PR #658 (fix/issue-10):** Closed; security impact overestimated, demoted to Wave 3
- **Next action:** Wait for maintainer engagement on #659/#660

## References

- Style guide: `github-archive/style.md`
- Issue backups: `github-archive/issues/`
