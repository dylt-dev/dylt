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

### TL;DR

One command end-to-end:

```bash
sudo daylight.sh github-shr-install your-org your-repo
# ... device-code prompt (go to URL, enter code) ...
# ... runner downloads, registers, service installs and starts ...
# ... verification workflow triggers, tailing runner logs ...
# Runner is online and verified for your-org/your-repo
```

No separate download, token creation, or svc start steps.

### Prerequisites

- A GitHub repo you own (or have admin access to)
- SSH or console access to a Linux VM (x86_64 or arm64)
- `jq` for JSON parsing
- A browser (for the device-code authentication step)

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

echo "Go to $VERIFICATION_URI and enter: $USER_CODE"
read -r -p "Press Enter after completing the browser step ..."

GRANT_TYPE="urn:ietf:params:oauth:grant-type:device_code"
ACCESS_TOKEN=$(curl -s -X POST \
  "https://github.com/login/oauth/access_token?client_id=$CLIENT_ID&device_code=$DEVICE_CODE&grant_type=$GRANT_TYPE" \
  | jq -r '.access_token')
```

Keep `$ACCESS_TOKEN` in a shell variable — it's short-lived (~8 hours)
and should not be stored on disk.

### Step 2: Download the Runner

**Via daylight.sh:**

```bash
download-shr-tarball /opt/actions-runner
```

**Manual:**

```bash
mkdir -p /opt/actions-runner
curl -sL https://api.github.com/repos/actions/runner/releases/latest \
  | jq -r '.assets[] | select(.name | test("actions-runner-linux-x64-.*\\.tar\\.gz")) | .browser_download_url' \
  | xargs curl -sL \
  | tar xz -C /opt/actions-runner
```

### Step 3: Register and Start (One Command)

```bash
sudo daylight.sh github-shr-install $ORG $REPO "linux"
```

This single command runs the entire flow:

1. Device-code flow to get a user access token for the `shrboy` GitHub App
2. Downloads and extracts the runner tarball into `/opt/actions-runner/shr-$ORG-$REPO`
3. Exchanges the UAT for a registration token and runs `config.sh`
4. Installs and starts the runner as a systemd service
5. Optionally creates and triggers a test workflow

The svc name is derived automatically as `shr-$ORG-$REPO`. The runner directory
is `/opt/actions-runner/shr-$ORG-$REPO`. The UAT is saved to `.uat` in that
directory for later cleanup operations.

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
sudo ./config.sh remove --token "$REG_TOKEN"
rm -rf /opt/actions-runner
```

---

## daylight.sh Function Reference

| Function | Use for |
|----------|---------|
| `github-create-user-access-token tokenvar appslug` | Device-code flow → token (if not using github-shr-install) |
| `detect-runner-platform` | Detect `linux-x64`, `linux-arm64`, etc. |
| `download-shr-tarball targetFolder [$version]` | Download + extract runner release (optionally pinned) |
| `github-shr-install org repoName [$labels] [$version]` | **End-to-end** — device-code auth, download, register, svc start, optional test. Requires `sudo`. |
| `github-shr-install-runner org repoName svcName uatPath labels [$version]` | Low-level: register runner with existing UAT (no svc start). |
| `github-shr-start org repoName` | Install svc, start service, run verification workflow. Requires `sudo`. |
| `github-shr-test org repoName svcName uatPath` | Create/trigger a test workflow and tail runner logs |
| `github-shr-clean org repoName` | Stop svc, deregister runner, delete directory. Requires `sudo`. |

All functions are in `daylight.sh`: `source /opt/bin/daylight.sh`.
