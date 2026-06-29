# Self-Hosted Runner Security Hardening

A tour of escalating security measures for GitHub Actions self-hosted
runners, from naive to paranoid. Each stage demonstrates a weakness in
the previous approach, then eliminates it.

## Prerequisites

- A GitHub repo you own (for testing)
- `gh` CLI authenticated for API calls
- Root or sudo on the target machine (for installing the runner)
- `jq` for JSON parsing in demo commands

---

## Stage 1: The Naive Approach

A long-lived PAT stored on the runner, a permanent runner registration,
manual setup.

### Setup

```bash
# Set the token once. It lives here forever.
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx

# Download the runner
curl -sL https://api.github.com/repos/actions/runner/releases/latest \
  | jq -r '.assets[] | select(.name | test("actions-runner-linux-x64-.*\\.tar\\.gz")) | .browser_download_url' \
  | xargs curl -sL \
  | tar xz -C /opt/actions-runner

# Get a registration token
REG_TOKEN=$(curl -sL \
  -X POST \
  -H "Authorization: Bearer $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github+json" \
  "https://api.github.com/repos/your-org/your-repo/actions/runners/registration-token" \
  | jq -r '.token')

# Register the runner
/opt/actions-runner/config.sh \
  --unattended \
  --url "https://github.com/your-org/your-repo" \
  --token "$REG_TOKEN" \
  --name "naive-runner" \
  --labels "naive" \
  --replace

# Install as a service
/opt/actions-runner/svc.sh install
/opt/actions-runner/svc.sh start
```

### What's wrong?

**Problem 1:** The PAT sits in `GITHUB_TOKEN` (env, bashrc, or a file).
Any process on the machine can read it. If this runner processes
untrusted code (e.g., PRs from forks), that code sees the token.

**Problem 2:** The PAT is scoped to your user or org — it can read repos,
create issues, trigger workflows, etc. It has far more power than the
runner needs.

**Problem 3:** The runner's `.credentials` file is permanent. If someone
gets shell access to this machine, they can reuse those credentials
indefinitely to receive jobs and exfiltrate workflow secrets.

### The exploit

```bash
# On the compromised runner — steal the PAT
echo "$GITHUB_TOKEN"

# Or if it's not in the env, check common places
cat ~/.bashrc
cat /etc/environment
systemctl show actions.runner.* 2>/dev/null | grep -i token

# Use it from anywhere, forever
curl -H "Authorization: Bearer $GITHUB_TOKEN" \
  https://api.github.com/repos/owner/repo

# Steal the runner's permanent credentials instead
cat /opt/actions-runner/.credentials
# Returns something like:
# {"clientId":"abc123","token":"eyJhbGciOiJSUzI1NiIs..."}

# Those creds let you impersonate this runner from any machine.
# No PAT needed — just copy .credentials and .runner to a new box.
```

---

## Stage 2: Ephemeral Mode

GitHub Actions supports `--ephemeral` on `config.sh`. The runner
processes exactly one job, then self-deregisters and exits.

### Setup

Same download as before, but registration changes:

```bash
# Get a registration token (still needs some auth)
REG_TOKEN=$(curl -sL \
  -X POST \
  -H "Authorization: Bearer $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github+json" \
  "https://api.github.com/repos/your-org/your-repo/actions/runners/registration-token" \
  | jq -r '.token')

# Register as ephemeral
/opt/actions-runner/config.sh \
  --unattended \
  --url "https://github.com/your-org/your-repo" \
  --token "$REG_TOKEN" \
  --name "ephemeral-runner" \
  --labels "ephemeral" \
  --ephemeral \
  --replace

# Run once — the runner processes one job and exits
./run.sh

# .credentials is GONE. The runner no longer exists in GitHub.
```

### What's still wrong?

**Problem fixed:** Runner creds are one-shot. No `.credentials` to
persist. An attacker who shells this runner during a job can only exfiltrate
that job's secrets (*), not receive future jobs.

**Problem remaining:** We still have a PAT somewhere. The orchestrator
(whatever calls `config.sh` for the next job) needs auth to generate
registration tokens. That PAT is still long-lived, still over-scoped,
still lives on disk.

(*) This is serious enough that you should use different runners for
untrusted (fork PR) workflows vs. trusted ones. GitHub's [security
hardening guide][gh-runner-security] covers this.

[gh-runner-security]: https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions

### The exploit (Stage 2)

```bash
# On the orchestrator machine — steal the PAT
cat /etc/sunbeam/env  # or wherever it's stored
# Still game over for all repos that PAT can access.
```

---

## Stage 3: GitHub App Instead of PAT

A GitHub App has its own identity. It can be installed on specific repos
or orgs with narrowly scoped permissions. Authentication uses a private
key — more secure than a PAT because:

- The key's only job is signing JWTs for installation tokens.
- Installation tokens expire in 1 hour.
- The app's permissions are visible, auditable, and revocable without
  affecting anything else.

### Setup

**One-time: Register the app**

1. Go to `https://github.com/settings/apps/new`
2. Name it (e.g., "runner-orchestrator")
3. Set permissions: at minimum `administration: write` (for runner tokens)
   and `metadata: read`
4. Generate a private key — download the `.pem` file
5. Install the app on your org/repo

**Now the orchestrator can get tokens programmatically:**

```bash
# These variables are the ONLY thing stored on the orchestrator.
# The private key is a PEM file, not a PAT.
APP_ID=123456
INSTALLATION_ID=789012
PRIVATE_KEY_PATH=/etc/sunbeam/runner-bot.pem

# Step 1: Sign a JWT with the private key
JWT=$(python3 -c "
import time, jwt

with open('$PRIVATE_KEY_PATH') as f:
    key = f.read()

now = int(time.time())
payload = {
    'iat': now - 60,        # issued 60s ago (clock skew)
    'exp': now + 600,       # expires in 10 minutes
    'iss': '$APP_ID'
}
print(jwt.encode(payload, key, algorithm='RS256'))
")

# Step 2: Exchange JWT for an installation token (1-hour expiry)
INSTALL_TOKEN=$(curl -sL \
  -X POST \
  -H "Authorization: Bearer $JWT" \
  -H "Accept: application/vnd.github+json" \
  "https://api.github.com/app/installations/$INSTALLATION_ID/access_tokens" \
  | jq -r '.token')

# Step 3: Use the installation token to create a registration token
REG_TOKEN=$(curl -sL \
  -X POST \
  -H "Authorization: Bearer $INSTALL_TOKEN" \
  -H "Accept: application/vnd.github+json" \
  "https://api.github.com/repos/your-org/your-repo/actions/runners/registration-token" \
  | jq -r '.token')
```

### What's still wrong?

**Problem fixed:** No human credential stored on the orchestrator. The
private key is narrowly scoped — it can only act as the app, which has
specific permissions. No PAT leak, no blast radius across unrelated
repos.

**Problem remaining:** The private key is still a long-lived secret on
disk. If the orchestrator is compromised, the attacker can generate
installation tokens for an hour at a time until the app is uninstalled
or the key is rotated.

### The exploit (Stage 3)

```bash
# On the compromised orchestrator
cp /etc/sunbeam/runner-bot.pem /tmp/
# Generate tokens from anywhere with that PEM + APP_ID + INSTALLATION_ID
```

---

## Stage 4: Private Key in a Secrets Manager

The private key moves out of the filesystem and into a secrets manager
with audit logging, access control, and short-lived sessions.

### Setup (1Password CLI example)

```bash
# The orchestrator has ONLY a short-lived 1Password session token.
# This could come from a CI OIDC claim, a service account, or manual auth.

# Fetch the key at startup, use it, drop it
PRIVATE_KEY=$(op read "op://infra/runner-bot-app/private-key")
APP_ID=$(op read "op://infra/runner-bot-app/app-id")
INSTALLATION_ID=$(op read "op://infra/runner-bot-app/installation-id")

# Generate JWT + installation token as before, but NEVER store the key
JWT=$(python3 -c "
import time, jwt
with open('/dev/stdin') as f: key = f.read()
now = int(time.time())
payload = {'iat': now - 60, 'exp': now + 600, 'iss': '$APP_ID'}
print(jwt.encode(payload, key, algorithm='RS256'))
" <<< "$PRIVATE_KEY")

# Wipe it from memory — the key exists only for the duration of the script
unset PRIVATE_KEY

# Use $JWT → $INSTALL_TOKEN → $REG_TOKEN as before
```

### What's solved

| Threat | Stage 1 | Stage 2 | Stage 3 | Stage 4 |
|---|---|---|---|---|
| Runner creds persist after job | Vulnerable | **Secure** | Secure | Secure |
| PAT on orchestrator disk | Vulnerable | Vulnerable | **Gone** | Gone |
| Private key on orchestrator disk | — | — | Vulnerable | **Gone** |
| Audit trail for who used what credential | None | None | **Partial** (app install) | **Full** (secrets manager logs) |

---

## Putting It All Together: dylt's Role

The `sunbeam.sh` orchestrator should ultimately support all four
approaches so you can pick your threat model:

| sunbeam function | Stage | Description |
|---|---|---|
| `github-shr-install` | 1–2 | Manual runner setup via UAT |
| *(planned)* `github-get-installation-token` | 3–4 | Generate JWT → installation token |
| *(planned)* `install-shr-ephemeral` | 2–4 | One-shot runner registration with any auth source |

The token resolution chain in `github-curl()` already supports
`--token` → `GITHUB_TOKEN` → `GH_TOKEN` → `gh auth token`. Adding
`GITHUB_APP_ID` + `GITHUB_APP_KEY` + `GITHUB_APP_INSTALLATION_ID` would
let it fall through to GitHub App auth when no PAT is available —
completing the path from Stage 1 to Stage 4 without breaking existing
usage.

---

## Key Takeaways

1. **Ephemeral mode is the biggest single win.** Use it first.
2. **PATs are the weakest link.** Replace them with a GitHub App when
   you can, a secrets manager when you must.
3. **Security is a chain, not a switch.** Each stage above removes one
   category of risk without overpromising. You can stop at any stage
   that matches your threat model.
4. **The runner's `.credentials` is permanent by default.** Treat it
   like a password. Or better, use `--ephemeral` and don't have one at
   all.
