# GitHub Communication Style Guide

How we write issues, PR descriptions, and comments for this project.

## Principles

**Say what's wrong, say how to fix it, stop talking.** Every sentence should
either establish the problem, explain why it matters, or describe the solution.
If a sentence doesn't do one of those three things, cut it.

**Assume the reader is technically competent.** Don't explain what a padding
oracle is. Don't explain what `encoding/xml` does. Do explain *this specific*
padding oracle -- which function, which early return, why it's exploitable here.
The reader's time is the scarcest resource.

**Be direct about severity without being inflammatory.** Describe what's broken
precisely: "has multiple bugs preventing decryption" rather than "completely broken"
or "entirely non-functional". Avoid language that reads as dismissive or judgmental
toward the original authors. The goal is to communicate the technical reality clearly,
not to editorialize.

**Acknowledge trade-offs and context honestly.** If a fix is a behavioral change,
say so. If an issue is in test/demo code, acknowledge the context. If something
is technically a vulnerability but not worth fixing, say that too and close the
issue. Don't hedge -- decide and state the decision.

## Structure

### Issue titles

Short, lowercase-ish, describes the defect rather than the fix. Reads like a
commit message subject line or a one-line bug description you'd give a
colleague. Use neutral factual language.

Good:
- "Padding oracle in CBC decryption"
- "Open redirect via RelayState"
- "AES-GCM Encrypt has multiple bugs preventing decryption"
- "SubjectConfirmation Method not checked"
- "InResponseTo not validated when AllowIDPInitiated is set"

Bad:
- "Fix: Implement constant-time padding validation for CBC mode" (describes fix, not problem)
- "Security: Critical AES-GCM Vulnerability in xmlenc Package" (bureaucratic)
- "Ensure SubjectConfirmation Method attribute is properly validated" (prescriptive, wordy)
- "AES-GCM Encrypt is completely broken" (inflammatory language)
- "InResponseTo skipped entirely when flag set" (implies negligence)

### Issue body

**Paragraph 1: What's broken.** State the bug or vulnerability. Include the
file/function location inline (not as a separate "Location" section). If a code
snippet makes the problem immediately obvious, include it. If not, don't.

**Paragraph 2 (if needed): Why it matters.** Threat model, attack scenario, or
spec violation. Only when the "so what" isn't obvious from paragraph 1. Don't
write out full attack chains for issues where the impact is self-evident.

**Paragraph 3 (if needed): The fix.** What was done or what should be done.
Mention backward-compat considerations, the nil-means-default pattern, or
behavioral changes here. If there's a branch, link it.

**Last line (for fixed issues):**
```
Fix: [`fix/issue-NN`](../tree/fix/issue-NN)
```

**For deferred issues**, end with a brief statement of why and the word
"Deferred": `Deferred for now.`

**For won't-fix issues**, state the reasoning plainly and close the issue.
Don't apologize for the decision.

### Things to avoid

- **Section headers in issues.** No `## Severity`, `## Location`,
  `## Analysis`, `## Recommendation`. These are bureaucratic filler. Weave
  the information into prose.

- **Severity labels in the text.** The issue speaks for itself. "AES-GCM is
  completely broken" doesn't need `## Severity: Critical` above it.

- **"Ease of fix" ratings.** These are internal planning artifacts, not useful
  to readers. If the fix is complex, explain why in the fix description.

- **Filler transitions.** "It is worth noting that", "Additionally", "It
  should be mentioned that". Just say the thing.

- **Passive voice when active is clearer.** "The fix strips KeyInfo" not
  "KeyInfo is stripped by the fix."

- **Performative uncertainty.** Don't write "this could potentially allow an
  attacker to possibly..." if you mean "an attacker can." Reserve hedging for
  situations where you're genuinely uncertain, and then be specific about what
  you're uncertain about.

## Tone

**Technically precise, professionally direct.** Write like you're explaining
a bug to a maintainer you respect -- not writing a formal report, not writing a
tweet. No jargon for its own sake, but no dumbing down either. Avoid language
that could be read as dismissive or judgmental.

**Respectful of context.** When discussing test/demo code, acknowledge its intended
scope: "While `samlidp` is intended for testing, adding optional authentication
improves security for non-isolated deployments." Avoid framing test code issues as
if they're production security failures.

**Factual, not judgmental.** Prefer "validation is not performed" over "check is
skipped entirely". The former describes what is; the latter implies what should
have been. Similarly, avoid "trivial", "obvious", "whatsoever", "zero friction" --
these editorialize rather than inform.

**Opinionated where warranted.** "This allows assertion replay attacks" is better
than "this may potentially enable possible replay scenarios." State your view,
support it, move on. But distinguish technical claims from value judgments.

**Honest about what we don't know or haven't done.** "Deferred for now" is
fine. "Planned, not yet implemented" is fine. Don't pretend something is done
when it isn't, and don't pretend uncertainty when you've made a decision.

## Cross-references

Link related issues by number (`#24 already rejects assertions carrying the
OneTimeUse condition`). Don't over-explain the relationship -- the reader can
click through.

Link fix branches as relative URLs: `[fix/issue-NN](../tree/fix/issue-NN)`.

## Language to Avoid

Based on revisions made before submitting to high-prestige upstream project:

**Inflammatory absolutes:**
- ❌ "completely broken", "entirely non-functional"
- ✅ "has multiple bugs preventing...", "does not perform..."

**Dismissive intensifiers:**
- ❌ "trivial", "trivially crashable"
- ❌ "with no auth whatsoever", "zero friction"
- ❌ "anyone can see", "obviously"
- ✅ "without authentication", "accessible to..."

**Judgmental phrasing:**
- ❌ "skipped entirely" (implies negligence)
- ❌ "that's not a reason to ship..."
- ❌ "should never have..."
- ✅ "not validated", "not performed"
- ✅ "adding X improves..."

**Overly dramatic attack scenarios for test code:**
- ❌ "An attacker can trivially crash the server"
- ❌ "Anyone with network access can take over with zero friction"
- ✅ "Could lead to memory exhaustion under sustained load"
- ✅ "Without authentication, endpoints are accessible to any network client"

## Revisions Made

Examples of language changes before upstream submission:

| Issue | Before | After |
|-------|--------|-------|
| #3 title | "InResponseTo skipped entirely..." | "InResponseTo not validated..." |
| #8 body | "it shouldn't be trivially crashable" | "adding a size cap improves robustness" |
| #12 body | "the check is skipped entirely" | "validation is not performed" |
| #13 title | "AES-GCM Encrypt is completely broken" | "AES-GCM Encrypt has multiple bugs preventing decryption" |
| #13 body | "make AES-GCM entirely non-functional" | "prevent successful round-trip encryption/decryption" |
| #15 body | "entirely attacker-controlled" | "derived from request body data, making them predictable" |
| #22 body | "with no auth whatsoever... full takeover with zero friction... that's not a reason to ship" | "without authentication... improves security for non-isolated environments" |

| Commit | Before | After |
|--------|--------|-------|
| fix/issue-08 | "making it trivial for an attacker to exhaust" | "allowing memory exhaustion through" |
| fix/issue-13 | "fix: broken AES-GCM encryption" | "fix: AES-GCM encryption bugs" |

## PR Description Format

Structure:
- Clear title (one line, describes the fix)
- Opening paragraph (what's broken) - no header, just describe the issue
- Fix: label followed by description of what you did
- Changes: label with bullet list of files modified
- Testing: label with bullet list of test coverage
- Optional sections: Breaking Change:, Backward Compatibility:

Formatting:
- Professional and technical
- Plain text, minimal formatting
- No markdown headers (##) within the body - use plain text labels like "Fix:", "Changes:", "Testing:"
- Use backticks liberally for code elements:
  - Field names: `AllowIDPInitiated`, `InResponseTo`
  - Function names: `validateAssertion`, `parseResponse`
  - File names: `service_provider.go`
  - Type names: `Response`, `Assertion`
  - Variable values: `InResponseTo=""`
- No emojis, no bold, no arrows
- Write natural paragraphs with long lines
- Let GitHub handle line wrapping

Example format:
```
The validation logic has a bug that allows X.

Fix: Changed the validation to check Y before Z.

Changes:
- service_provider.go: Updated validation logic
- service_provider_test.go: Added regression tests

Testing:
- Valid case: accepted
- Invalid case: rejected
- All existing tests pass
```

Avoid:
- Markdown headers: ## Issue, ## Fix, ## Testing
- Emojis: ✅ ❌ ⚠️
- Decorative formatting: **bold** for emphasis
- Arrows: → or =>
- Checkboxes: [x] or [ ]
- Manual line breaks within paragraphs

## Pull Request Submission Strategy

### The Psychology of Open Source Maintainers

Maintainers of mature security libraries operate under constraints:
- **Limited time** - usually volunteers with day jobs
- **High risk** - breaking changes affect many users
- **Pattern matching** - quickly categorize contributors as helpful/harmful
- **Defensive posture** - seen many low-quality PRs, bad security research

**The credibility curve:**
- First PR: High scrutiny, skeptical baseline
- Second PR: "Maybe this person knows what they're doing"
- Third+ PR: "This is a trusted contributor"

### Strategic Submission Order

**Don't dump all PRs at once.** This signals:
- ❌ "I'm a security researcher doing a drive-by"
- ❌ "I don't respect your review bandwidth"
- ❌ "I'm padding my CV, not here to help long-term"

**Instead: Build credibility incrementally.** Start with easy-to-review, high-value changes.

#### Wave 1: Trust Builders (First 1-3 PRs)

Submit **obviously correct**, **non-breaking**, **quick to review** fixes:

1. **#15 (example crypto/rand)** - One-line fix
   - Review time: < 5 min
   - Signal: "I'm here to help, not critique"

2. **#29 (SubjectConfirmationData NotBefore)** - Small, follows existing pattern
   - Review time: ~10 min
   - Signal: "I understand SAML and your codebase"

3. **#16 (cookie security)** - Practical improvement to test code
   - Review time: ~10 min
   - Signal: "I care about real-world usage, not just theory"

**Timing:** Submit with 2-3 day gaps. Wait for merge or substantive feedback before submitting next wave.

#### Wave 2: Credibility Established (PRs 4-6)

After 1-2 merges, submit **non-breaking** but **more substantial** fixes:

4. **#23 (SubjectConfirmation Method)**
5. **#28 (AuthnContext validation)**
6. **#06 (Optional SessionKey)**
7. **#05 (TripleDES not registered by default)**

#### Wave 3: Higher Impact (PRs 7-12)

After multiple merges, tackle **complex fixes** or **behavioral changes**:

8. **#10 (Always strip KeyInfo)** - Critical but complex
9. **#12 (Destination always required)** - Breaking change
10. **#07 (RelayState open redirect)**
11. **#03 (InResponseTo validation)**
12. **#21 (Return signature errors immediately)**

#### Wave 4: API Changes (PRs 13-18)

Only after established trust:

13. **#09 (RSA-PKCS1v1.5 not registered)**
14. **#24 (Reject OneTimeUse/ProxyRestriction)**
15. **#25 (Host-only cookie scope)** - Breaking change
16. **#18 (AudienceRestriction multiple Audiences)** - Breaking API
17. **#08 (MemoryStore eviction)**
18. **#22 (AdminMiddleware)**

#### Wave 5: Architectural Changes (PRs 19-22)

For established contributors only:

19. **#20 (IdP Encryptor field + AES-256)**
20. **#26 (OAEP SHA-1)** - Breaking change
21. **#02 (Constant-time CBC padding)** - Critical but intricate
22. **#13 (AES-GCM fixes)** - Four interlocking bugs

### PR Submission Tactics

**Pacing:**
- Don't submit more than 3 PRs before getting merge/feedback
- Match maintainer's response speed (if they're slow, you slow down)
- After first merge, wait 3-5 days before next batch

**First PR is critical:**
- Extra thorough description
- Include checklist: ✅ Tests added, ✅ Backward compatible, ✅ Documentation updated
- Show you've thought about edge cases

**Response to feedback:**
- Respond within 24 hours if possible
- Implement feedback eagerly: "Great point, fixed in [commit SHA]"
- If you disagree, be diplomatic: "I see your concern. Another approach would be..."
- Never be defensive about the original code

**Signals to send:**
- "I'm here long-term" (mention you have more fixes if helpful)
- "I respect your judgment" (don't argue about priority)
- "I care about this project" (engage with other issues/PRs)

### Anti-Patterns to Avoid

❌ **Don't submit all 22 PRs at once**
- Overwhelming, signals you're dumping work

❌ **Don't lead with breaking API changes**
- Signals you don't care about users

❌ **Don't lead with critical-but-complex security fixes**
- Hard to review = high barrier for first PR

❌ **Don't submit PRs faster than maintainer can review**
- Disrespectful of their time

❌ **Don't mention you have 22 fixes in first PR**
- Overwhelming

❌ **Don't say "I found X vulnerabilities"**
- Sounds like you're criticizing, not helping

### Reality Checks

**If maintainer ignores you:**
- Wait 2-3 weeks
- Open an issue asking: "Would you be open to a PR that fixes X?"
- If still no response after issue + easy PR, project may be inactive

**If they reject your approach:**
- Be flexible: "Would you prefer I do Y instead?"
- Be willing to walk away from contentious fixes
- Some issues may not align with maintainer's vision

**If they want to batch review:**
- Some maintainers prefer reviewing multiple PRs together
- If they say "submit the rest", do so
- Otherwise, stick to incremental approach

### Why This Approach Works

✅ **Minimizes maintainer burden** - Easy reviews first
✅ **Builds trust incrementally** - Prove competence before asking for trust
✅ **Shows long-term commitment** - Not a drive-by researcher
✅ **Demonstrates codebase understanding** - Start simple, show you learned
✅ **Respects their time** - Match their review pace
✅ **Reduces risk** - Start with non-breaking changes
✅ **Creates momentum** - Quick wins lead to more engagement

## PR Descriptions

(To be filled in when we open the first PRs. Should match issue voice but with additional detail on testing, backward compatibility, and merge considerations.)
