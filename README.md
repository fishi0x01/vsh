# vsh

[![Latest release](https://img.shields.io/github/release/fishi0x01/vsh.svg)](https://github.com/fishi0x01/vsh/releases/latest)
![CI](https://github.com/fishi0x01/vsh/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fishi0x01/vsh)](https://goreportcard.com/report/github.com/fishi0x01/vsh)
[![Code Climate](https://codeclimate.com/github/fishi0x01/vsh/badges/gpa.svg)](https://codeclimate.com/github/fishi0x01/vsh)

![vsh usage](https://user-images.githubusercontent.com/10799507/66355982-9872a980-e969-11e9-8ca4-6a2ff215f835.gif)

`vsh` is an interactive [HashiCorp Vault](https://www.vaultproject.io/) shell which treats paths and keys like directories and files.
Core features are:

- recursive operations on paths with `cp`, `mv` or `rm`
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

```text
append <from-secret> <to-secret> [flag]
cat <file-path>
cd <dir-path>
cp <from-path> <to-path>
grep <search> <path> [-e|--regexp] [-k|--keys] [-v|--values]
ls <dir-path // optional>
mv <from-path> <to-path>
replace <search> <replacement> <path> [-e|--regexp] [-k|--keys] [-v|--values] [-y|--confirm] [-n|--dry-run]
rm <dir-path or file-path>
```

`cp`, `grep`, `replace` and `rm` command always have the `-r/-R` flag implied, i.e., every operation works recursively.

### append

Append operation reads secrets from `<from-secret>` and merges it to `<to-secret>`.
The `<to-secret>` will be created with a placeholder value if it does not exists.
Both `<from-secret>` and `<to-secret>` must be leaves (path cannot end with `/`).

By default, `append` does not overwrite secrets if the `<to-secret>` already contains a key.
The default behavior can be explicitly set using flag: `-s` or `--skip`. Example:

```bash
> cat /secret/from

fruit=apple
vegetable=tomato

> cat /secret/to

fruit=pear
tree=oak

> append /secret/from /secret/to -s

> cat /secret/to

fruit=pear
vegetable=tomato
tree=oak
```

Setting flag `-f` or `--force` will cause the conflicting keys from the `<to-secret>` to be overwritten with keys from the `<from-secret`>. Example:

```bash
> cat /secret/from

fruit=apple
vegetable=tomato

> cat /secret/to

fruit=pear
tree=oak

> append /secret/from /secret/to -f

> cat /secret/to

fruit=apple
vegetable=tomato
tree=oak
```

Setting flag `-r` or `--rename` will cause the conflicting keys from the `<to-secret>` to be kept as they are. Instead the keys from the `<from-secret`> will be stored under a renamed key. Example:

```bash
> cat /secret/from

fruit=apple
vegetable=tomato

> cat /secret/to

fruit=pear
tree=oak

> append /secret/from /secret/to -r

> cat /secret/to

fruit=pear
fruit_1=apple
vegetable=tomato
tree=oak
```

### grep

`grep` recursively searches the given substring in key and value pairs. To treat the search string as a regular-expression, add `-e` or `--regexp` to the end of the command. By default, both keys and values will be searched. If you would like to limit the search, you may add `-k` or `--keys` to the end of the command to search only a path's keys, or `-v` or `--values` to search only a path's values.
 If you are looking for copies or just trying to find the path to a certain string, this command might come in handy.

### replace

`replace` works similarly to `grep` above, but has the ability to mutate data inside Vault. By default, confirmation is required before writing data. You may skip confirmation by using the `-y`/`--confirm` flags. Conversely, you may use the `-n`/`--dry-run` flags to skip both confirmation and any writes. Changes that would be made are presented in red (delete) and green (add) coloring.

## Setting the vault token

In order to get a valid token, `vsh` uses vault's TokenHelper mechanism (`github.com/hashicorp/vault/command/config`).
That means `vsh` supports setting vault tokens via `~/.vault-token`, `VAULT_TOKEN` and external `token_helper`.

## Secret Backend Discovery

`vsh` attempts to reliably discover all available backends.
Ideally, the vault token used by `vsh` has `list` permissions on `sys/mount`.
If this is not the case, then `vsh` does not know the available backends beforehand.
That means initially there won't be path auto-completion on the top (backend) level.
However, `vsh` will try with best-effort strategy to reliably determine the kv version of every entered path.

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
If your token has a limited number of uses, then consider using the non-interactive mode to avoid auto-completion queries.

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

## Permission requirements

`vsh` requires `List` permission on the operated paths.
This is necessary to determine if a path points to a node or leaf in the path tree.
Further, it is needed to gather auto-completion data.

For operations like `cp` or `mv`, `vsh` additionally requires `Read` and `Write` permissions on the operated paths.

## Quality

Working on vault secrets can be critical, making quality and correct behavior a first class citizen for `vsh`.
That being said, `vsh` is still a small open source project, meaning we cannot make any guarantees.
However, we put strong emphasis on [TDD](https://en.wikipedia.org/wiki/Test-driven_development).
Every PR is tested with an extensive [suite](test/suites) of integration tests.
Most tests run on KV1 and KV2 and every test runs against vault `1.0.0` and `1.6.1`, i.e., versions in between should also be compatible.

## Local Development

Requirements:

- `golang` (compiled and tested with `v1.15.3`)
- `docker` for integration testing
- `make` for simplified commands

```
make compile
make get-bats
make integration-tests
```

## Debugging

`-v DEBUG` sets debug log level, which also creates a `vsh_trace.log` file to log any error object from the vault API.
