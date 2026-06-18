<img src="./images/vclt_banner.png" alt="vclt logo" />

# vclt

A stripped-down, opinionated Hashicorp Vault client for KV v2 secret management and basic server administration.

---

## Table of Contents

- [Overview](#overview)
- [Global Flags](#global-flags)
- [Authentication & Server Address Resolution](#authentication--server-address-resolution)
- [Commands](#commands)
  - [version](#version)
  - [completion](#completion)
  - [secrets list](#secrets-list)
  - [secrets read](#secrets-read)
  - [secrets write](#secrets-write)
  - [secrets rm](#secrets-rm)
  - [secrets destroy](#secrets-destroy)
  - [admin setrootkeys](#admin-setrootkeys)
  - [admin seal](#admin-seal)
  - [admin unseal](#admin-unseal)
- [Vault Policy Requirements](#vault-policy-requirements)
- [Building from Source](#building-from-source)
- [Binary Package Building](#binary-package-building)

---

## Overview

`vclt` wraps the [`vaultLib/kv`](https://github.com/jeanfrancoisgratton/vaultLib) and [`vaultLib/admin`](https://github.com/jeanfrancoisgratton/vaultLib) libraries to expose a focused CLI covering the most common day-to-day Vault operations: reading, writing, listing, deleting, and destroying KV v2 secrets, plus sealing and unsealing the server.

Per-user configuration lives under `$HOME/.config/JFG/vclt/`, which is created automatically on first run.

---

## Global Flags

These flags are available on every command and subcommand.

| Flag | Short | Default | Description |
|---|---|---|---|
| `--token` | `-t` | — | Vault auth token. Overrides `VAULT_TOKEN` and `~/.vault-token`. |
| `--address` | `-a` | — | Vault server URL (e.g. `https://vault.example.com:8200`). Overrides `VAULT_ADDR`. |
| `--output` | `-o` | `text` | Output format: `text` or `json`. |
| `--quiet` | `-q` | `false` | Suppress human-readable output (useful in scripts). |
| `--debug` | `-d` | `false` | Enable debug mode. |

---

## Authentication & Server Address Resolution

Every command that talks to Vault resolves the token and server address through the same lookup chain before making any API call.

**Token resolution** (`--token` / `-t`):
1. `-t` flag value, if provided.
2. `VAULT_TOKEN` environment variable.
3. Contents of `~/.vault-token` (whitespace-trimmed).
4. If none of the above yield a value, the command exits with `ERR_VAULTTOKENMISSING`.

**Server address resolution** (`--address` / `-a`):
1. `-a` flag value, if provided.
2. `VAULT_ADDR` environment variable.
3. If neither is set, the command exits with `ERR_VAULTADDRESSMISSING`.

> **Note:** `admin unseal` does not require a token — it only needs the server address, because an unsealing operation is performed against a sealed (i.e. unauthenticated) server.

---

## Commands

### version

```
vclt version
```

Prints the build version, release date, and the Go runtime version used to compile the binary.

**Token required:** No  
**Policy required:** None

---

### completion

```
vclt completion bash
vclt completion zsh
```

Generates shell completion scripts via Cobra. Output is written to stdout and can be redirected to the appropriate completion directory.

**Bash (current session):**
```sh
source <(vclt completion bash)
```

**Bash (persistent):**
```sh
vclt completion bash | sudo tee /etc/bash_completion.d/vclt > /dev/null
```

**Zsh:**
```sh
vclt completion zsh > ~/.zsh/vclt
echo 'fpath=($HOME/.zsh $fpath)' >> ~/.zshrc
echo 'autoload -Uz compinit && compinit' >> ~/.zshrc
```

**Token required:** No  
**Policy required:** None

---

### secrets list

```
vclt secrets list <KV_ENGINE>
vclt secrets ls   <KV_ENGINE>
vclt secrets show <KV_ENGINE>
```

Lists all secret paths available under the given KV v2 engine mount. Output is rendered as a sorted table. The `--extended` flag fetches and displays the latest version number for each secret.

| Flag | Short | Default | Description |
|---|---|---|---|
| `--extended` | `-x` | `false` | Include version numbers in the listing. |

**Example:**
```sh
vclt secrets list mysecrets
vclt secrets list mysecrets -x
```

**Token required:** Yes  
**Policy required:**
```hcl
path "<KV_ENGINE>/metadata/*" {
  capabilities = ["list"]
}
# With --extended, metadata reads are also performed:
path "<KV_ENGINE>/metadata/*" {
  capabilities = ["list", "read"]
}
```

---

### secrets read

```
vclt secrets read <KV_ENGINE> <SECRET_PATH>
vclt secrets get  <KV_ENGINE> <SECRET_PATH>
```

Reads a secret from the KV v2 engine. Without `--field`, all key/value pairs in the secret are printed. With `--field`, only the value of the specified field is printed — useful for scripting.

| Flag | Short | Default | Description |
|---|---|---|---|
| `--field` | `-f` | — | Print only the value of the named field. |
| `--version` | `-v` | `0` | Read a specific version. `0` resolves to the latest available non-destroyed version. |

**Examples:**
```sh
vclt secrets read mysecrets db/credentials
vclt secrets read mysecrets db/credentials -f password
vclt secrets read mysecrets db/credentials -v 3
vclt -o json secrets read mysecrets db/credentials
```

**Token required:** Yes  
**Policy required:**
```hcl
path "<KV_ENGINE>/data/<SECRET_PATH>" {
  capabilities = ["read"]
}
```

---

### secrets write

```
vclt secrets write <KV_ENGINE> <SECRET_PATH> <KEY> <VALUE>
vclt secrets put   <KV_ENGINE> <SECRET_PATH> <KEY> <VALUE>
```

Writes a single key/value field to the secret at the given path. If the secret already exists, a new version is created (KV v2 versioning). If the path does not yet exist, it is created.

**Example:**
```sh
vclt secrets write mysecrets db/credentials password s3cr3t
```

**Token required:** Yes  
**Policy required:**
```hcl
path "<KV_ENGINE>/data/<SECRET_PATH>" {
  capabilities = ["create", "update"]
}
```

---

### secrets rm

```
vclt secrets rm     <KV_ENGINE> <SECRET_PATH>
vclt secrets delete <KV_ENGINE> <SECRET_PATH>
```

Soft-deletes a secret or a single field within a secret. A soft delete marks the version as deleted but retains the data; it can be undeleted if needed (KV v2 behaviour). When `--field` is specified, only that key is removed from the secret's data map, and a new version is written.

| Flag | Short | Default | Description |
|---|---|---|---|
| `--field` | `-f` | — | Delete only the named field from the secret. |
| `--version` | `-v` | `0` | Target a specific version for deletion. `0` targets the latest. |

**Examples:**
```sh
# Soft-delete the entire secret (latest version)
vclt secrets rm mysecrets db/credentials

# Soft-delete a specific version
vclt secrets rm mysecrets db/credentials -v 2

# Remove a single field
vclt secrets rm mysecrets db/credentials -f password
```

**Token required:** Yes  
**Policy required:**
```hcl
# Whole-secret soft delete
path "<KV_ENGINE>/data/<SECRET_PATH>" {
  capabilities = ["delete"]
}
# Field-level delete (reads current data, rewrites without the field)
path "<KV_ENGINE>/data/<SECRET_PATH>" {
  capabilities = ["read", "update", "delete"]
}
```

---

### secrets destroy

```
vclt secrets destroy <KV_ENGINE> <SECRET_PATH>
```

Permanently destroys a secret version. Unlike `rm`, a destroy is irreversible — the secret data is permanently removed from Vault storage. The `--version` flag targets a specific version; without it, version `0` is passed to the underlying library (resolves to latest).

| Flag | Short | Default | Description |
|---|---|---|---|
| `--version` | `-v` | `0` | Permanently destroy the specified version. |

**Example:**
```sh
vclt secrets destroy mysecrets db/old-credentials
vclt secrets destroy mysecrets db/old-credentials -v 1
```

**Token required:** Yes  
**Policy required:**
```hcl
path "<KV_ENGINE>/destroy/<SECRET_PATH>" {
  capabilities = ["update"]
}
```

---

### admin setrootkeys

```
vclt admin setrootkeys [filename]
```

Interactively collects the unseal key shards for a Vault initialized with Shamir secret sharing and saves them to an encrypted JSON file under `$HOME/.config/JFG/vclt/`. If no filename is provided, `rootkeys.json` is used. A `.json` extension is appended automatically if omitted.

The command queries the Vault server for its current seal threshold (`minimumRequired`) and will refuse to save a file that contains fewer key shards than that threshold.

Key shards are stored encoded (via `helperFunctions/v5` `EncodeString`) — they are not stored in plaintext.

**Example:**
```sh
vclt admin setrootkeys
vclt admin setrootkeys prod-keys
```

**Token required:** Yes (to query seal status)  
**Policy required:**
```hcl
path "sys/seal-status" {
  capabilities = ["read"]
}
```

---

### admin seal

```
vclt admin seal
```

Seals the Vault server. Once sealed, all secrets become inaccessible until the server is unsealed again. This operation requires a token with `sys/seal` permission.

**Example:**
```sh
vclt admin seal
vclt -a https://vault.example.com:8200 -t hvs.xxxx admin seal
```

**Token required:** Yes  
**Policy required:**
```hcl
path "sys/seal" {
  capabilities = ["update"]
}
```

---

### admin unseal

```
vclt admin unseal [filename]
```

Unseals a sealed Vault using the key shards previously saved by `admin setrootkeys`. If no filename is provided, `$HOME/.config/JFG/vclt/rootkeys.json` is used. The command decodes the stored key shards and submits exactly `minimumRequired` of them to the Vault unseal endpoint.

**Example:**
```sh
vclt admin unseal
vclt admin unseal prod-keys.json
```

**Token required:** No (unsealing operates against a sealed, unauthenticated server)  
**Policy required:** None (unauthenticated endpoint — `sys/unseal` is always accessible on a sealed node)

---

## Vault Policy Requirements

Summary table for quick reference. All secret operations assume a KV v2 engine.

| Command | Vault path | Capabilities |
|---|---|---|
| `secrets list` | `<engine>/metadata/*` | `list` |
| `secrets list -x` | `<engine>/metadata/*` | `list`, `read` |
| `secrets read` | `<engine>/data/<path>` | `read` |
| `secrets write` | `<engine>/data/<path>` | `create`, `update` |
| `secrets rm` (whole) | `<engine>/data/<path>` | `delete` |
| `secrets rm` (field) | `<engine>/data/<path>` | `read`, `update`, `delete` |
| `secrets destroy` | `<engine>/destroy/<path>` | `update` |
| `admin setrootkeys` | `sys/seal-status` | `read` |
| `admin seal` | `sys/seal` | `update` |
| `admin unseal` | `sys/unseal` | — (unauthenticated) |

---

## Building from Source

### Requirements

- Go **1.26.4** or later (see `go.version`)
- Network access to the Go module proxy (or a pre-populated module cache)

### Dependencies

Direct dependencies are managed via `go.mod`:

| Module | Version |
|---|---|
| `github.com/jeanfrancoisgratton/customError/v3` | v3.0.0 |
| `github.com/jeanfrancoisgratton/helperFunctions/v5` | v5.2.2 |
| `github.com/jeanfrancoisgratton/vaultLib` | v1.5.0 |
| `github.com/jedib0t/go-pretty/v6` | v6.8.1 |
| `github.com/spf13/cobra` | v1.10.2 |

### Quick build

```sh
cd src
go mod download
go build -o vclt .
```

### Using `build.sh`

The `build.sh` script at the repo root adds branch-aware output naming and optional permission checks on the target directory. The default output path is `/opt/bin`.

```sh
# Build to the default output path (/opt/bin)
./build.sh

# Build to a custom output path
./build.sh ~/bin

# Build with output directory permission check (expects group 'devops', mode 775)
./build.sh --checkperms
./build.sh ~/bin --checkperms
```

When the current Git branch is `master`, `main`, or `develop`, the binary is named `vclt`. On any other branch it is named `vclt-<branch>`, which avoids overwriting a stable binary during feature development.

### Static build (CGO disabled)

All packaging targets build with CGO disabled for maximum portability:

```sh
cd src
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -buildid=" -o vclt .
```

---

## Binary Package Building

Packaging stubs for Alpine (APK), Arch Linux (PKGBUILD), Debian (DEB), and Red Hat/RHEL/Fedora (RPM) are provided under `__alpine/`, `__archlinux/`, `__debian/`, and `__redhat/` respectively.

All package formats are built inside dedicated Docker containers. The following container images are required:

| Container | Format |
|---|---|
| `apkbuilder` | Alpine APK |
| `archbuilder` | Arch Linux PKGBUILD |
| `debbuilder` | Debian/Ubuntu DEB |
| `rpmbuilder` | Red Hat / Fedora RPM |

> **Note:** The Docker build contexts for these containers are not provided in this repository and are not intended for distribution outside of the author's own environment. The containers are assumed to be available locally.

### RPM

The `__redhat/` directory contains a `Makefile` that drives the full RPM workflow.

```sh
cd __redhat

# Create the source tarball (git archive of HEAD)
make tarball

# Build the RPM (includes tarball step)
make rpm

# Build the RPM and update the %changelog from git log
make rpmcl

# Upload the built RPM to the configured Nexus repository
make upload

# Commit the updated spec file (after changelog update)
make commitcl
```

The changelog can also be updated standalone via `__redhat/updateChangelog.sh`.

### Debian

```sh
cd __debian
./1.install-build-deps.sh   # install Go and build toolchain inside the debbuilder container
./2.build_binary.sh         # compile and package the .deb
./3.restore_repo.sh         # restore the local apt repository
```

### Alpine

The `__alpine/APKBUILD` follows the standard `abuild` workflow. It copies the `src/` tree into the build directory, disables CGO, and builds a statically linked binary installed to `/opt/bin/vclt`. Post-install hooks register bash and zsh completions automatically.

### Arch Linux

```sh
cd __archlinux
./1.install-build-deps.sh   # install Go inside the archbuilder container
./2.build-package.sh        # run makepkg and produce the .pkg.tar.zst
```
