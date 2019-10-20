### Status
[![CircleCI](https://circleci.com/gh/fishi0x01/vsh.svg?style=svg)](https://circleci.com/gh/fishi0x01/vsh)
[![Go Report Card](https://goreportcard.com/badge/github.com/fishi0x01/vsh)](https://goreportcard.com/report/github.com/fishi0x01/vsh)
[![Code Climate](https://codeclimate.com/github/fishi0x01/vsh/badges/gpa.svg)](https://codeclimate.com/github/fishi0x01/vsh)

# vsh

![vsh usage](https://user-images.githubusercontent.com/10799507/66355982-9872a980-e969-11e9-8ca4-6a2ff215f835.gif)

vsh is an interactive HashiCorp Vault shell which treats vault secret paths like directories. 
That way you can do recursive operations on the paths. 
Both, vault KV v1 and v2 are supported. 
Further, copying/moving secrets between both versions is supported.

vsh can also act as an executor in a non-interactive way (similar to `bash -c "<cmd>"`).

Integration tests are running against vault `1.2.2`.

## Supported commands

```
mv <from-path> <to-path>
cp <from-path> <to-path>
rm <dir-path or filel-path>
ls <dir-path // optional>
cd <dir-path>
cat <file-path>
```

Unlike unix, `cp` and `rm` always have the `-r` flag implied, i.e., every operation works recursively on the paths.

## Interactive mode

```
export VAULT_ADDR=http://localhost:8080
export VAULT_TOKEN=root
export VAULT_PATH=secret/  # VAULT_PATH is optional
./vsh
http://localhost:8080 /secret/> 
```

**Note:** in order to query the root `/` the `VAULT_TOKEN` should have permissions to list the available secret backends (`sys/mounts/`).

**Note:** the given token is used for auto-completion, i.e., quite some `List()` queries are done with that token, even if you do not `rm` or `mv` anything.
If your token has a limited number of uses, then consider using the non-interactive mode to avoid auto-completion queries.

## Non-interactive mode

```
export VAULT_ADDR=<addr>
export VAULT_TOKEN=<token>
./vsh -c "rm secret/dir/to/remove/"
```

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

- sys/mounts/ permission needed at the moment for auto-completion --> disable auto-completion on top level if permission not given
- `tree` command
- currently `mv` behaves a little different from UNIX. `mv /secret/source/a /secret/target/` should yield `/secret/target/a`
- caching `List()` queries to reduce IO / token usage (?)
- more integration tests!