# Security fixes for crewjam/saml

We maintain 22 fix branches for upstream submission to `github.com/crewjam/saml`.
Our fork is at `github.com/aashh/saml`. Issues tracking each fix live on the fork.

## Git setup

Two remotes:
- `origin` — `aashh/saml` (fork), configured with PAT in `.git/config`
- `upstream` — `crewjam/saml` (upstream target)

PAT is at `~/saml-github-pat.txt`. Use it for fork operations:
```
GH_TOKEN=$(cat ~/saml-github-pat.txt) gh ...
```
For upstream PRs, use `gh auth login` (different GitHub account).

## Branches

Each `fix/issue-NN` branch must be exactly 1 commit ahead of `upstream/main`.
No `.gitignore` or other local changes. Verify with:
```
git fetch upstream main
for b in $(git branch --list 'fix/*' | tr -d '* '); do
  echo "$b: +$(git rev-list --count upstream/main..$b)"
done
```

All should show `+1`. If upstream/main moves, rebase fix branches onto it.

## Opening PRs

PRs go to `crewjam/saml` from `aashh/saml` fix branches.

Follow `github-archive/style.md` for PR descriptions:
- No emojis, no bold, no arrows, no `##` headers in the body
- Use plain text labels: `Fix:`, `Changes:`, `Testing:`, `Breaking Change:`
- Backticks for all code elements
- Long lines, let GitHub wrap
- Use `--body-file` with heredoc to avoid shell backtick interpretation

Submission order and pacing: `github-archive/pr-submission-plan.md`

## Local files

`CLAUDE.md` and `github-archive/` live in the worktree but are excluded from all
branches via `.git/info/exclude`. They are version-controlled on the orphan `meta`
branch (no shared history with main). To update them:

```
git stash        # if needed
git checkout meta
# edit files
git add -A && git commit -m "update meta"
git checkout main
git stash pop    # if needed
```

After switching back, the files stay in the worktree because `.git/info/exclude`
tells git not to touch them.

Contents of `github-archive/`:
- `pr-submission-plan.md` — wave ordering, PR status, pacing rules
- `style.md` — writing style for issues, PRs, and comments
- `issues/` — backup of all 29 fork issues (JSON + markdown)

## Testing

```
go test ./... -count=1
```

Golden files: `go test ./... -test.update-golden`

Upstream CI runs golangci-lint v2 and tests on Go 1.22/1.23/1.24.
