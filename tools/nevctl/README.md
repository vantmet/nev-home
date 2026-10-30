# nevctl

Command and Control CLI for OpenTofu with Linux System Keyring secret management.

`nevctl` securely stores credentials in the Linux system keyring (via the DBus Secret Service API / libsecret) and dynamically hydrates configuration templates (`*.tmpl`) before executing OpenTofu commands. Sensitive files marked as ephemeral (`terraform.tfvars`) are guaranteed to be cleaned up upon exit, ensuring real secrets are never committed to git.

---

## Installation & Build

Build the binary directly with Go:

```bash
cd tools/nevctl
go build -o ../../nevctl ./cmd/nevctl
```

You can also install it to your user path (e.g. `~/.local/bin` or `~/go/bin`):

```bash
go install ./cmd/nevctl
```

---

## Usage

### 1. Managing Secrets

Store, view, and inspect secrets in the system keyring:

```bash
# Set a secret (prompts securely for masked input if value is omitted)
nevctl secret set routeros_password
nevctl secret set influxdb_token "my-token-value"

# Retrieve a secret value
nevctl secret get routeros_password

# List all stored secrets for the current service
nevctl secret list

# Delete a secret
nevctl secret delete old_secret
```

### 2. Running OpenTofu

`nevctl` wraps standard OpenTofu commands:

```bash
nevctl plan
nevctl apply
nevctl init
nevctl destroy
```

Any extra flags are forwarded directly to OpenTofu:

```bash
nevctl plan -out=tfplan
nevctl apply tfplan
```

You can also run arbitrary OpenTofu commands:

```bash
nevctl tofu state list
nevctl tofu show
```

#### Lifecycle Behavior:
1. **Hydrate:** Evaluates all templates declared in `nevctl.yaml` with secrets from the keyring.
2. **Execute:** Runs `tofu` while forwarding process signals (such as `SIGINT` / Ctrl+C).
3. **Clean:** Ephemeral files (like `terraform.tfvars`) are automatically deleted upon completion (even if OpenTofu fails or is interrupted). Persistent files (like container configs) remain intact for running Docker/Podman services.

### 3. Rendering & Cleaning Manually

```bash
# Render all templates without running OpenTofu
nevctl render

# Clean all rendered files
nevctl clean

# Clean only ephemeral files
nevctl clean --ephemeral-only
```

---

## Configuration (`nevctl.yaml`)

`nevctl` automatically discovers `nevctl.yaml` in the current working directory or any parent directory.

```yaml
service: nev-home     # Keyring service/namespace
tofu_binary: tofu     # Can be overridden by TOFU_BIN environment variable

templates:
  - source: terraform.tfvars.tmpl
    dest: terraform.tfvars
    ephemeral: true   # Automatically wiped after tofu runs

  - source: telegraf/telegraf.conf.tmpl
    dest: telegraf/telegraf.conf
    ephemeral: false  # Kept on disk for container mounts
```

### Template Functions
Templates use standard Go `text/template` syntax and support the following functions:
- `{{ secret "key_name" }}`: Fetches secret from keyring (fails fast if missing)
- `{{ secret "key_name" | quote }}`: Escapes and double-quotes value
- `{{ secret "key_name" | squote }}`: Single-quotes value
- `{{ defaultVal "fallback" (secret "key_name") }}`: Fallback value
- `{{ env "ENV_VAR" }}`: Reads environment variable
- `{{ b64enc "string" }}`: Base64 encodes
- `{{ b64dec "encoded" }}`: Base64 decodes
