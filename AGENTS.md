### sunbeam.sh

sunbeam.sh is a bash 4+ script that contains functions to help a user develop daylight apps on their local workstation, a VPS, wherever they want to work. sunbeam.sh might be installed on Windows, OSX, or Mac, and sunbeam functions will be run as the local user, who may or may not have sudo. This makes sunbeam.sh very different from daylight.sh, from https://github.com/daylight-public/daylight. But it is still desirable for sunbeam.sh to be consistent with daylight.sh. Someone familiar with reading and using one script should feel comfortable reading and using the other.

sunbeam.sh functions are alphabetized, except for a main function at the end.

#### comments
sunbeam.sh functions should begin with comments that look like this.

```
#-------------------------------------------------------------------------------
#
# download-shr-tarball()
#
# Download the tarball for the latest GitHub Actions Self-Hosted Runner release
#
download-shr-tarball ()
{
```

`download-shr-tarball` is the function name.
"Download" etc is a one or two line description. In rare cases a description will be longer because the function will warrant a longer description. Do not shorten or edit existing descriptions.

Functions that appear in the `main()` case statement should NOT use the `@internal` tag in their comment block — it signifies the function is a helper that lives outside the dispatch table.

```
#-------------------------------------------------------------------------------
#
# sb()
#
# @internal
# Shorthand description of the function goes here on the next line
#
```

The `@internal` tag sits alone on its own line between the function name and the description. This makes it trivially grepable (`grep -B3 '^# @internal$'`) while keeping the description readable and adjacent.

### tools/

New helper scripts go in `./tools/` by default, not the repo root. Tools follow the same conventions as sunbeam.sh — comment blocks, alphabetized functions, case dispatch, and `@internal` for helpers. Each tool script is self-contained and calls sunbeam.sh functions via case dispatch.

Test scripts follow the same conventions and are manual verification tools, not CI. Daylight repo has test scripts for shared infrastructure (`github-curl-parse-args`, `--token` flag handling) that can be referenced by dylt tests.

Test scripts must be executable (`chmod +x`). Name them
`tools/test-<branchname>.sh` where the branch name includes the issue number
(e.g. `tools/test-138-trigger-func-changes.sh`).

### pushing code changes (issue-driven workflow)

1. Propose an issue title, a short branch name (without issue number), and an issue body (can include markdown)
2. User confirms or edits each
3. Create the issue via `gh issue create`
4. Prepend the issue number to the branch name (e.g. `42-fix-thing`)
5. Commit, push, create PR with `Closes #N` in the body
6. Do the work; user merges the PR when ready

Exception: Meta changes to AGENTS.md itself use the `update-agents-md` persistent branch with no issue (see below).

### download-dylt flags

`download-dylt` uses a custom flag parser (not `github-parse-args`) for branch/release selection.

| Flag | Value | Behavior |
|---|---|---|
| (none) | — | Branch mode, defaults to `main` |
| `--branch` | (no value) | Branch mode, defaults to `main` |
| `--branch <name>` | branch name | Branch mode, specific branch |
| `--release` | (no value) | Release mode, latest release |
| `--release <tag>` | tag name | Release mode, specific tag |
| `--release --latest` | — | Same as `--release` with no value |
| `--latest` alone | — | Error: requires `--release` |
| `--token <value>` | token | GitHub API token for release mode |
| `--branch` + `--release` | — | Error: incompatible |

Rules:
- Flags are parsed in order. The optional value after `--branch` or `--release` is consumed only if it doesn't start with `--`.
- `--latest` must follow `--release` (either immediately or as a later flag).
- The destination folder is the first positional argument after all flags; an optional platform string is the second positional.
- Exactly one of branch mode or release mode must be active.

### download-daylight

`download-daylight` is a wrapper around `download-daylight-batch` that adds a
`--gen-bash-completions [<path>]` flag and an interactive prompt to generate
bash completions, delegating to the downloaded `daylight.sh gen-completion-script`.

### reminders

- Cleanup dylt release yamls
- Investigate why dylt's release matrix always hits every platform
- Explore externalizing label creation into a separate function
- Create custom git function for setting URL
- Use a GitHub App for auth instead of PATs
- Sort out all the github functions starting with flags/args

### AGENTS.md changes

AGENTS.md is meta — it holds conventions and reminders. Changes to it don't
need issues, labels, or approval. Use the `update-agents-md` persistent branch:

- Check it out from `main`, push commits to it over time
- Open a PR against `main` when there's a batch ready (no issue link needed)
- Self-merge, then rebase `update-agents-md` onto fresh `main`
