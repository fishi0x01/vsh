### Status
[![Latest release](https://img.shields.io/github/release/fishi0x01/vsh.svg)](https://github.com/fishi0x01/vsh/releases/latest)
![CI](https://github.com/fishi0x01/vsh/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fishi0x01/vsh)](https://goreportcard.com/report/github.com/fishi0x01/vsh)
[![Code Climate](https://codeclimate.com/github/fishi0x01/vsh/badges/gpa.svg)](https://codeclimate.com/github/fishi0x01/vsh)

# vsh

![vsh usage](https://user-images.githubusercontent.com/10799507/66355982-9872a980-e969-11e9-8ca4-6a2ff215f835.gif)

`vsh` is an interactive [HashiCorp Vault](https://www.vaultproject.io/) shell which treats paths and keys like directories and files.
Key features are:

- recursive operations on paths with `cp`, `mv` or `rm`
- term search with `grep`
- transparency towards differences between KV1 and KV2, i.e., you can freely move/copy secrets between both
- non-interactive mode for automation (`vsh -c "<cmd>"`)
- merging keys with different strategies through `append`

## Supported commands

```text
mv <from-path> <to-path>
cp <from-path> <to-path>
append <from-secret> <to-secret> [flag]
rm <dir-path or filel-path>
ls <dir-path // optional>
grep <search-term> <path>
cd <dir-path>
cat <file-path>
```

`cp`, `rm` and `grep` command always have the `-r/-R` flag implied, i.e., every operation works recursively.

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

`grep` recursively searches the given term in key and value pairs. It does not support regex.
 If you are looking for copies or just trying to find the path to a certain term, this command might come in handy.

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
Most tests run on KV1 and KV2 and every test runs against vault `1.0.0` and `1.5.4`, i.e., versions in between should also be compatible.

## Local Development

Requirements:
- `golang` (compiled and tested with `v1.13.12`)
- `docker` for integration testing
- `make` for simplified commands

```
make compile
make get-bats
make integration-tests
```

## Debugging

`-v DEBUG` sets debug log level, which also creates a `vsh_trace.log` file to log any error object from the vault API.
