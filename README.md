# vsh

[![Latest release](https://img.shields.io/github/release/fishi0x01/vsh.svg)](https://github.com/fishi0x01/vsh/releases/latest)
![CI](https://github.com/fishi0x01/vsh/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fishi0x01/vsh)](https://goreportcard.com/report/github.com/fishi0x01/vsh)
[![Code Climate](https://codeclimate.com/github/fishi0x01/vsh/badges/gpa.svg)](https://codeclimate.com/github/fishi0x01/vsh)

![vsh usage](https://user-images.githubusercontent.com/10799507/66355982-9872a980-e969-11e9-8ca4-6a2ff215f835.gif)

`vsh` is an interactive [HashiCorp Vault](https://www.vaultproject.io/) shell and cli tool. It comes with multiple common operations and treats paths like directories and files.
Core features are:

- recursive operations on paths for many operations, e.g., `cp`, `rm`, `mv`
- search with `grep` (substring or regular-expression)
- substitute patterns in keys and/or values (substring or regular-expression) with `replace`
- transparency towards differences between KV1 and KV2, i.e., you can freely move/copy secrets between both
- non-interactive mode for automation (`vsh -c "<cmd>"`)
- merging keys with different strategies through `append`

## Installation

### Homebrew

```sh
brew install vsh
```

### Nix

```sh
nix-env -i vsh
```

### Static binaries for Linux / MacOS

Download latest static binaries from [release page](https://github.com/fishi0x01/vsh/releases).

## Supported commands

- [append](doc/commands/append.md) merges secrets with different strategies (allows recursive operation on paths)
- [cat](doc/commands/cat.md) shows the key/value pairs of a path
- [cd](doc/commands/cd.md) allows interactive navigation through the paths
- [cp](doc/commands/cp.md) copies secrets from one location to another (allows recursive operation on paths)
- [grep](doc/commands/grep.md) searches for substrings or regular expressions (allows recursive operation on paths)
- [ls](doc/commands/ls.md) shows the subpaths of a given path
- [mv](doc/commands/mv.md) moves secrets from one location to another (allows recursive operation on paths)
- [replace](doc/commands/replace.md) substrings or regular expressions (allows recursive operation on paths)
- [rm](doc/commands/rm.md) removes secret(s) (allows recursive operation on paths)

## Setting the vault token

In order to get a valid token, `vsh` uses vault's TokenHelper mechanism.
That means `vsh` supports setting vault tokens via `~/.vault-token`, `VAULT_TOKEN` and external [token-helper](https://www.vaultproject.io/docs/commands/token-helper).

## Token permission requirements

`vsh` requires `List` permission on the operated paths.
This is necessary to determine if a path points to a node or leaf in the path tree.
Further, it is needed to gather auto-completion data.

Commands which alter the data like `cp` or `mv`, additionally require `Read` and `Write` permissions on the operated paths.

In order to reliably discover all available backends, ideally the vault token used by `vsh` has `List` permission on `sys/mount`. However, this is not a hard requirement.
If the token doesn't have `List` permission on `sys/mount`, then `vsh` does not know the available backends beforehand.
That means initially there won't be path auto-completion on the top (backend) level.
Regardless, `vsh` will try with best-effort strategy to reliably determine the kv version of every entered path.

## Interactive mode

```
export VAULT_ADDR=http://localhost:8080
export VAULT_TOKEN=root
export VAULT_PATH=secret/  # VAULT_PATH is optional
./vsh
http://localhost:8080 /secret/>
```

**Note:** the given token is used for auto-completion, i.e., `List()` queries are done with that token, even if you do not `rm` or `mv` anything.
`vsh` caches `List()` results to reduce the amount of queries. However, after execution of each command the cache is cleared
in order to do accurate tab-completion.
If your token has a limited number of uses, then consider using the non-interactive mode or toggle auto-completion off, to avoid `List()` queries.

### Toggle auto-completion

To reduce the number of queries against vault, you can disable path auto-completion in 2 ways:

1. Disable at start time:

```
./vsh --disable-auto-completion
```

2. Toggle inside interactive mode:

```
./vsh
http://localhost:8080 /secret/> toggle-auto-completion
Use path auto-completion: false
http://localhost:8080 /secret/> toggle-auto-completion
Use path auto-completion: true
```

## Non-interactive mode

```
export VAULT_ADDR=<addr>
export VAULT_TOKEN=<token>
./vsh -c "rm secret/dir/to/remove/"
```

## Some words about the quality

Working on vault secrets can be critical, making quality and correct behavior a first class citizen for `vsh`.
That being said, `vsh` is still a small open source project, meaning we cannot give any guarantees.
However, we put strong emphasis on test-driven development.
Every PR is tested with an extensive [suite](test/suites) of integration tests.
Vast majority of tests run on KV1 and KV2 and every test runs against vault `1.0.0` and `1.6.2`, i.e., vault versions in between should also be compatible.

## Contributions

Contributions in any form are always welcome! Without contributions from the community, `vsh` wouldn't be the tool it is today.

### Local Development

Requirements:

- `golang` (compiled and tested with `v1.15.7`)
- `docker` for integration testing
- `make` for simplified commands

```
make compile
make get-bats
make integration-tests
```

### Debugging

`-v DEBUG` sets debug log level, which also creates a `vsh_trace.log` file to log any error object from the vault API.
