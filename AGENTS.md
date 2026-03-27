# AGENTS.md

This file provides guidance for AI agents and contributors working in this repository.

## Project overview

`vsh` is an interactive [HashiCorp Vault](https://www.vaultproject.io/) shell and CLI tool written in Go.
It treats Vault paths like filesystem directories and files, exposing familiar shell-like commands (`cp`, `mv`, `rm`, `ls`, `cat`, `grep`, `replace`, `add`, `append`, `cd`).

Key design goals:
- Transparent KV1/KV2 interoperability — secrets can be moved freely between both backends.
- Recursive operations (`cp`, `mv`, `rm`, `append`) with concurrent execution via a goroutine worker pool.
- Both interactive (REPL with tab-completion) and non-interactive (`-c "<cmd>"`) modes.

## Repository layout

```
main.go              # Entry point; flag parsing, REPL loop
tokenhelper.go       # Vault token-helper integration (build-tag guarded)
cli/                 # One file per command (cp.go, mv.go, …); shared args/confirmation helpers
client/              # Vault API client abstraction (read, write, delete, list, traverse)
completer/           # Path auto-completion logic for interactive mode
log/                 # Logging helpers
doc/commands/        # Markdown docs for each command
test/                # Integration tests (bats framework)
build/               # Compiled binaries (git-ignored)
vendor/              # Vendored Go dependencies
```

## Build

Requirements: `golang >= 1.24`, `docker`, `make`.

A [`.mise.toml`](https://mise.jdx.dev/getting-started.html) is provided to set up the Go toolchain quickly.

```sh
make compile                    # build for the current platform into build/
make compile-releases           # cross-compile for linux/darwin × amd64/arm64
```

## Testing

Integration tests use the [bats](https://bats-core.readthedocs.io/) framework and run against a real Vault container.

```sh
make get-bats                   # download bats and plugins into test/bin/
make integration-tests          # run all test suites

# Run a single suite (useful during development):
make single-test KV_BACKEND=KV2 VAULT_VERSION=1.20.2 TEST_SUITE=commands/cp
```

Test suites live under `test/suites/`. Every test runs against both KV1 and KV2 and against multiple Vault versions (currently `1.13.4` and `1.20.2`).

Start a local Vault instance for manual testing:

```sh
make local-vault-standard-test-instance      # standard secrets setup
make local-vault-concurrency-test-instance   # large dataset for concurrency tests
```

## Code style and linting

```sh
make lint       # golangci-lint
make format     # gofmt + gci + golines + goimports
```

All PRs must pass the CI lint check. Run `make lint` before submitting.

## Adding a new command

1. Create `cli/<command>.go` implementing the `Command` interface defined in `cli/command.go`.
2. Register the command in `main.go`.
3. Add a doc page at `doc/commands/<command>.md` following the existing format.
4. Add integration tests under `test/suites/commands/<command>/`.
5. Update `README.md` (supported commands list) and `CHANGELOG.md` (unreleased section).

## Environment variables

| Variable            | Description                                      |
|---------------------|--------------------------------------------------|
| `VAULT_ADDR`        | Vault server address                             |
| `VAULT_TOKEN`       | Vault token (also supports `~/.vault-token` and token-helper) |
| `VAULT_CACERT`      | Path to PEM CA certificate for TLS               |
| `VAULT_PATH`        | Initial working path (optional)                  |
| `VAULT_KV1_MOUNTS`  | Comma-delimited list of KV1 secret mounts        |
| `VAULT_KV2_MOUNTS`  | Comma-delimited list of KV2 secret mounts        |

## Important notes for agents

- **Never modify vendored code** in `vendor/`. Run `make vendor` after changing `go.mod`/`go.sum`.
- **Always vendor changes**: the release builds use `-mod vendor`; tests will fail if `vendor/` is out of sync.
- **Integration tests are the source of truth** — unit tests are minimal; correctness is verified by the bats suites against a real Vault.
- **KV1 and KV2 parity**: any change to secret read/write/delete logic must be verified on both backends.
- **Concurrency**: recursive operations use a goroutine pool (see `--worker-count` flag). Changes to `cp`, `mv`, or `rm` must account for concurrent execution.
