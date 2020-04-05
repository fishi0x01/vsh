### Status
[![CircleCI](https://circleci.com/gh/fishi0x01/vsh.svg?style=svg)](https://circleci.com/gh/fishi0x01/vsh)
[![Go Report Card](https://goreportcard.com/badge/github.com/fishi0x01/vsh)](https://goreportcard.com/report/github.com/fishi0x01/vsh)
[![Code Climate](https://codeclimate.com/github/fishi0x01/vsh/badges/gpa.svg)](https://codeclimate.com/github/fishi0x01/vsh)

# vsh

![vsh usage](https://user-images.githubusercontent.com/10799507/66355982-9872a980-e969-11e9-8ca4-6a2ff215f835.gif)

`vsh` is an interactive HashiCorp Vault shell which treats vault secret paths like directories. 
That way you can do recursive operations on the paths. 
Both, vault KV v1 and v2 are supported. 
Further, copying/moving secrets between both KV versions is supported.

`vsh` also supports a non-interactive mode (similar to `bash -c "<cmd>"`), which 
makes it easier to integrate with automation.

Integration tests are running against vault `1.3.4`.

## Supported commands

```
mv <from-path> <to-path>
cp <from-path> <to-path>
rm <dir-path or filel-path>
ls <dir-path // optional>
grep <search-term> <path>
cd <dir-path>
cat <file-path>
```

`cp`, `rm` and `grep` command always have the `-r/-R` flag implied, i.e., every operation works recursively on the paths.

### grep

`grep` recursively searches the given term in key and value pairs and does not support regex. 
[SSOT](https://en.wikipedia.org/wiki/Single_source_of_truth) is very much desired, however, in praxis it is probably not always applied consistently.
 If you are looking for copies or just trying to find the path to a certain term, this command might come in handy. 

## Setting the vault token

In order to get a valid token, `vsh` uses vault's TokenHelper mechanism (`github.com/hashicorp/vault/command/config`).
That means `vsh` supports setting vault tokens via `~/.vault-token`, `VAULT_TOKEN` and external `token_helper`.

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

## Non-interactive mode

```
export VAULT_ADDR=<addr>
export VAULT_TOKEN=<token>
./vsh -c "rm secret/dir/to/remove/"
```

## Misc

`vsh` attempts to reliably discover all available backends. 
Ideally, the vault token used by `vsh` has `list` permissions on `sys/mount`. 
If this is not the case, then `vsh` does not know the available backends beforehand. 
That means initially there won't be path auto-completion on the top level. 
However, `vsh` will try with best effort, to reliably determine the kv version of every entered path. 

## Local Development

Requirements:
- `golang` v1.12.7
- `docker` for integration testing
- `make` for simplified commands

```
make compile
make integration-test
```

## TODOs

- `tree` command
- currently `mv` and `cp` behave a little different from UNIX. `mv /secret/source/a /secret/target/` should yield `/secret/target/a`
