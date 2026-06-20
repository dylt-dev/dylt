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

### pushing code changes

All code changes will be done on new branches, with short names -- 2 or 3 words or terms separated by hyphens. Ask me to approve branch names. After committing and pushing a branch, create a PR for the change, where the body of the PR contains information similar or identical to the plan markdown. Create an issue and link it to the PR. Ask what the issue should be labelled - Bug, Task, or Feature.

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
