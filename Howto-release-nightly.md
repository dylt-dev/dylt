# How to Release Nightly

## Background

Nightly releases are managed by a GitHub Actions workflow (`nightly-release.yml`)
in each repo, triggered via `workflow_dispatch`. The local systemd timer calls
`trigger-nightly-release` which dispatches the workflow — the tagging and release
creation happens entirely on GitHub.

## Automated (timer)

The systemd indexed service runs daily at 03:00:

```bash
# install the service for a repo
install-build-nightly-svc dylt-dev/dylt

# check timer status
systemctl status nightly-release@dylt-dev-dylt.timer
```

Admin must create `/opt/svc/nightly-release/<instance>/env` with:

```
GITHUB_PAT=ghp_xxxxxxxx
```

## Manual (ad-hoc label)

```bash
# via sunbeam.sh
source sunbeam.sh
trigger-nightly-release "dylt-dev/dylt"

# via GitHub CLI (with a label)
gh workflow run nightly-release.yml -f label="fix oauth timeout"

# via curl (same as what trigger-nightly-release does)
GITHUB_PAT=xxx \
  curl -X POST .../actions/workflows/nightly-release.yml/dispatches \
  -H "Authorization: Bearer $GITHUB_PAT" \
  -d '{"ref": "main"}'
```

## Tag format

```
# clean (first run of the day)
v1.0.7-nightly-20240618

# clean (second run same day — timestamp fallback)
v1.0.7-nightly-20240618-093022

# labelled
v1.0.7-nightly-20240618-fix-oauth-timeout-153102
```

The tag push triggers `release.yml` (goreleaser) which builds cross-platform
binaries and publishes a release.

## Removed

The following functions were removed from `sunbeam.sh` in favor of the GHA approach:

- `gen-nightly-tagname`
- `gen-nightly-timestamp`
- `git-tag-nightly`
- `git-do-nightly-release`

Use `trigger-nightly-release` instead.
