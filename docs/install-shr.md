# How To Install a Self-Hosted GitHub Actions Runner on a VM

This doc walks through installing a self-hosted runner (SHR) for a repo,
using a GitHub user access token obtained via device-code flow.

The recommended approach uses **ephemeral mode** — the runner processes
one job, then self-deregisters and exits. No permanent credentials are
left on the VM. This eliminates the biggest security risk with
self-hosted runners.

For a deeper discussion of the security model, see
[runner-security-hardening.md](runner-security-hardening.md).

---

## Prerequisites
# Approaches to Deploying a Service on a VM via GitHub

Three approaches, ordered by integration level.

| | Approach | Credential | Complexity |
|--|----------|------------|------------|
| Best | Permanent SHR | `.credentials` on disk | Medium (svc.sh) |
| Simplest | DIY + cron | Short-lived or PAT | Low (bash + crontab) |
| Safest | DIY long-poll | Short-lived, refreshed | Medium (agent loop) |

---

## Approach 1: Permanent Self-Hosted Runner (Best)

A permanent SHR stays online as a systemd service, long-polls GitHub for
workflow jobs, and needs no open inbound ports. It's the same model
GitHub uses for its own hosted runners, just on your hardware.

### Prerequisites

- A GitHub repo you own (or have admin access to)
- SSH or console access to a Linux VM (x86_64 or arm64)
- `jq` for JSON parsing
- A browser (for the device-code authentication step)

---

## Step 1: Get a User Access Token via Device-Code Flow

A GitHub user access token lets you act on behalf of your GitHub user,
scoped to whatever repos the GitHub App used for authentication is
installed on.

We'll use the **`shrboy`** GitHub App, which has
`organization_self_hosted_runners: write` permission and is installed
on the repos that need it.

### Via daylight.sh (if available)

```bash
source /opt/bin/daylight.sh
github-create-user-access-token MY_TOKEN shrboy
# Opens browser for device-code flow. After completion, $MY_TOKEN is set.
```

### Manual curl

```bash
# 1. Get the client ID for the shrboy app
CLIENT_ID=$(curl -s https://api.github.com/apps/shrboy | jq -r '.client_id')
echo "Client ID: $CLIENT_ID"

# 2. Start device-code flow
### Step 1: Get a User Access Token via Device-Code Flow

A GitHub user access token lets you act on behalf of your GitHub user,
scoped to whatever repos the GitHub App used for authentication is
installed on. We'll use **`shrboy`**, which has
`organization_self_hosted_runners: write` permission.

**Via daylight.sh:**

```bash
/opt/bin/daylight.sh github-create-user-access-token ACCESS_TOKEN shrboy
```
`ACCESS_TOKEN` (no dollar-sign) is the name of an envvar that will receive the value of the token. Most `daylight.sh` function don't work this way, instead emitting data which can be assigned to a variable via VAR=$(func foo bar). But that requires the function does not write anything to `stdout`, and `github-create-user-access-token` needs to write to `stdout` to prompt the user for input.

**Manual curl:**

```bash
CLIENT_ID=$(curl -s https://api.github.com/apps/shrboy | jq -r '.client_id')

DATA=$(curl -s -X POST "https://github.com/login/device/code?client_id=$CLIENT_ID")
DEVICE_CODE=$(echo "$DATA" | jq -r '.device_code')
USER_CODE=$(echo "$DATA" | jq -r '.user_code')
VERIFICATION_URI=$(echo "$DATA" | jq -r '.verification_uri')

echo
echo "Go to $VERIFICATION_URI in your browser"
echo "Enter the code: $USER_CODE"
echo
read -r -p "Press Enter after completing the browser step ..."

# 3. Poll for the access token
echo "Go to $VERIFICATION_URI and enter: $USER_CODE"
read -r -p "Press Enter after completing the browser step ..."

GRANT_TYPE="urn:ietf:params:oauth:grant-type:device_code"
ACCESS_TOKEN=$(curl -s -X POST \
  "https://github.com/login/oauth/access_token?client_id=$CLIENT_ID&device_code=$DEVICE_CODE&grant_type=$GRANT_TYPE" \
  | jq -r '.access_token')

echo "Token acquired: ${ACCESS_TOKEN:0:20}..."
```

Keep this token in a variable — you'll use it in the next step. The
token is short-lived (8 hours by default for device-code tokens) and
should not be stored on disk.

---

## Step 2: Download the Runner

### Via daylight.sh (if available)
```

Keep `$ACCESS_TOKEN` in a shell variable — it's short-lived (~8 hours)
and should not be stored on disk.

### Step 2: Download the Runner

**Via daylight.sh:**

```bash
download-shr-tarball /opt/actions-runner
```

### Manual curl
**Manual:**

```bash
mkdir -p /opt/actions-runner
curl -sL https://api.github.com/repos/actions/runner/releases/latest \
  | jq -r '.assets[] | select(.name | test("actions-runner-linux-x64-.*\\.tar\\.gz")) | .browser_download_url' \
  | xargs curl -sL \
  | tar xz -C /opt/actions-runner
```

---

## Step 3: Register an Ephemeral Runner (Recommended)

An ephemeral runner processes one job, then self-deregisters and
deletes its credentials. No `.credentials` file persists on disk.
### Step 3: Register and Start

```bash
ORG=your-org
REPO=your-repo
RUNNER_NAME="ephemeral-$(hostname)-$(date +%s)"

# Get a registration token using the user access token

REG_TOKEN=$(curl -sL \
  -X POST \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Accept: application/vnd.github+json" \
  "https://api.github.com/repos/$ORG/$REPO/actions/runners/registration-token" \
  | jq -r '.token')

# Register as ephemeral
/opt/actions-runner/config.sh \
  --unattended \
  --url "https://github.com/$ORG/$REPO" \
  --token "$REG_TOKEN" \
  --name "$RUNNER_NAME" \
  --labels "ephemeral,linux,$(uname -m)" \
  --ephemeral \
  --replace

# Run once — the runner processes one job and exits
cd /opt/actions-runner && ./run.sh
```

After the job completes, the runner deregisters. No cleanup needed.

### Running as a one-shot service

To run the ephemeral runner on every workflow dispatch, wrap it in a
script and trigger it via cron or a webhook:

```bash
#!/usr/bin/env bash
# /usr/local/bin/run-ephemeral-shr.sh
set -euo pipefail
ORG=your-org
REPO=your-repo
ACCESS_TOKEN=$(cat /run/secrets/github-token)  # injected securely
REG_TOKEN=$(curl -sL -X POST \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Accept: application/vnd.github+json" \
  "https://api.github.com/repos/$ORG/$REPO/actions/runners/registration-token" \
  | jq -r '.token')
/opt/actions-runner/config.sh --unattended \
  --url "https://github.com/$ORG/$REPO" \
  --token "$REG_TOKEN" \
  --name "ephemeral-$(hostname)-$(date +%s)" \
  --labels "ephemeral,linux" \
  --ephemeral --replace
cd /opt/actions-runner && exec ./run.sh
```

---

## Step 4: Alternative — Permanent Runner (Legacy)

If ephemeral mode doesn't fit your workflow (e.g., you need the runner
to stay online for multiple jobs), you can register a permanent runner.

```bash
ORG=your-org
REPO=your-repo

# Get a registration token
REG_TOKEN=$(curl -sL \
  -X POST \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Accept: application/vnd.github+json" \
  "https://api.github.com/repos/$ORG/$REPO/actions/runners/registration-token" \
  | jq -r '.token')

# Register (permanent)
/opt/actions-runner/config.sh \
  --unattended \
  --url "https://github.com/$ORG/$REPO" \
  --token "$REG_TOKEN" \
  --name "$(hostname)" \
  --labels "linux,$(uname -m)" \
  --replace

# Install as a system service
cd /opt/actions-runner
sudo ./svc.sh install
sudo ./svc.sh start
```

**Security caveat:** The runner's `.credentials` file is permanent.
Anyone who gains access to this VM can reuse those credentials to
receive workflow jobs and exfiltrate secrets. Prefer ephemeral mode
when possible.

### Via daylight.sh (legacy, permanent only)

```bash
install-shr-token $ORG $REPO my-shr $ACCESS_TOKEN "linux"
```

Note: `install-shr-token` does not support `--ephemeral`. It registers
a permanent runner as a systemd service.

---

## Step 5: Verify

Trigger a test workflow that targets your runner's labels and confirm
it picks up the job.

```bash
gh workflow run <workflow-name> --repo $ORG/$REPO --ref main
```

Or via the GitHub web UI: Actions → workflow → Run workflow → select
your runner labels.

Check that the job runs on your VM and completes successfully.

---

## Step 6: Cleanup

### Remove a permanent runner

```bash
# Get a fresh token (device-code flow again, or use the same one if still valid)
# Then deregister:
  --name "$(hostname)" \
  --labels "linux,$(uname -m)" \
  --replace

cd /opt/actions-runner
sudo ./svc.sh install
sudo ./svc.sh start
```

**Via daylight.sh:**

```bash
install-shr-token $ORG $REPO my-shr $ACCESS_TOKEN "linux"
```

### Security Note

The runner's `.credentials` file is a long-lived credential. Mitigations:

- Restrict what workflows target this runner via labels + environments
- Use a dedicated VM with minimal other services
- Rotate periodically by deregistering and re-registering

---

## Approach 2: DIY + Cron (Simplest)

Forget GitHub Actions runners entirely. A cron job runs a deploy script
on a schedule — every hour, every 5 minutes, whatever suits you.

### Setup

```bash
# /usr/local/bin/deploy.sh
set -euo pipefail
cd /opt/svc/my-service
git pull
make update
sudo systemctl restart my-service
```

```bash
# crontab — runs every hour
0 * * * * /usr/local/bin/deploy.sh
```

### Advantages

- **Simple.** One bash script, one crontab line.
- **No permanent credentials.** The service only needs `git pull` access
  (deploy key or PAT scoped to `contents: read`). No `.credentials` file.
- **Manual override.** If you don't want to wait for the next cron tick:
  `ssh vm && sudo systemctl restart my-service`.

### Tradeoffs

- Worst-case delay = the cron interval. Set it to match your tolerance.
- No GitHub-side trigger. You can't start a deployment from the Actions
  tab. It runs on its own schedule.

---

## Approach 3: DIY Long-Poll Agent (Safest)

A purpose-built agent that loops on the VM, polls the GitHub API for
new releases or commits, and deploys when something changes. Like a
SHR, but scoped to exactly one job — deploy the service.

### How It Works

1. Agent starts, does a device-code flow to get a user access token
   (short-lived, ~8 hours)
2. Enters a loop: every N seconds, check the GitHub API for new releases
   or commits on a target branch
3. When a change is detected, pull the latest code, run the deploy
   script, report status back via the GitHub API
4. When the token is close to expiry, the agent re-runs the device-code
   flow or exits with a prompt

### Advantages

- **No permanent credential on disk.** The token lives in memory and is
  refreshed regularly.
- **No inbound ports.** The agent polls outbound, same as a SHR.
- **Full control over the credential scope.** The `shrboy` app's
  user-access token is scoped to exactly the repos the app is installed
  on — nothing else.

### Tradeoffs

- **More complex** than cron — you're writing and maintaining an agent.
- **Still needs human interaction** for the initial device-code flow
  (or a GitHub App private key for fully automated bootstrapping).

---

## Cleanup

Remove a permanent SHR:

```bash
REG_TOKEN=$(curl -sL -X POST \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Accept: application/vnd.github+json" \
  "https://api.github.com/repos/$ORG/$REPO/actions/runners/registration-token" \
  | jq -r '.token')

cd /opt/actions-runner
sudo ./svc.sh stop
sudo ./svc.sh uninstall
./config.sh remove --token "$REG_TOKEN"
rm -rf /opt/actions-runner
```

### Ephemeral runners

No cleanup needed — the runner self-deregisters after its one job. You
can delete the runner directory if desired:

```bash
rm -rf /opt/actions-runner
```

---

## Appendix: daylight.sh Function Reference

| Function | When to use | Notes |
|----------|-------------|-------|
| `github-create-user-access-token tokenvar appslug` | Step 1 | Device-code flow; returns token in a nameref |
| `detect-runner-platform` | Step 2 | Prints `linux-x64`, `linux-arm64`, etc. |
| `download-shr-tarball targetFolder` | Step 2 | Downloads + extracts latest runner release |
| `install-shr-token org repo svcName token labels` | Step 4 | Permanent runner only; no `--ephemeral` support |

All functions are in `daylight.sh` source: `source /opt/bin/daylight.sh`.

---

## Token Lifecycle

| Token type | Acquired via | Lifetime | Where it lives |
|------------|-------------|----------|----------------|
| User access token | Device-code flow | ~8 hours | Shell variable only |
| Registration token | API (POST) | ~1 hour | Used immediately |
| Runner credentials | `config.sh` | Permanent (Stage 1) or one-shot (Stage 2) | `.credentials` file or gone |

**Best practice:** Acquire the user access token, use it immediately to
register the runner, and don't save it to disk. For automated setups,
use a GitHub App private key (Stage 3 in
[runner-security-hardening.md](runner-security-hardening.md)) instead
of a user access token.
---

## daylight.sh Function Reference

| Function | Use for |
|----------|---------|
| `github-create-user-access-token tokenvar appslug` | Device-code flow → token |
| `detect-runner-platform` | Detect `linux-x64`, `linux-arm64`, etc. |
| `download-shr-tarball targetFolder` | Download + extract latest runner release |
| `install-shr-token org repo svcName token labels` | Register permanent SHR + install as service |

All functions are in `daylight.sh`: `source /opt/bin/daylight.sh`.
