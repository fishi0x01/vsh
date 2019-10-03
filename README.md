### Status
[![CircleCI](https://circleci.com/gh/fishi0x01/vsh.svg?style=svg)](https://circleci.com/gh/fishi0x01/vsh)
[![Go Report Card](https://goreportcard.com/badge/github.com/fishi0x01/vsh)](https://goreportcard.com/report/github.com/fishi0x01/vsh)
[![Code Climate](https://codeclimate.com/github/fishi0x01/vsh/badges/gpa.svg)](https://codeclimate.com/github/fishi0x01/vsh)

# vsh

**The project is still in an alpha stage. Expect bugs.**

vsh is a simple interactive HashiCorp Vault shell which treats vault secret paths like directories. That way you can do recursive operations on paths. Vault KV v1 and v2 are both supported.
Commands can also be executed in a non-interacitve way.

## Interactive mode

```
export VAULT_ADDR=http://localhost:8080
export VAULT_TOKEN=root
./vsh
http://localhost:8888 > cd secret/
http://localhost:8888 secret/>
```

**Note: in order to query the root `/` the `VAULT_TOKEN` should have permissions to list the available secret backends (`sys/mounts/`).**

**Note: the given token is also used for auto-completion feature, i.e., quite some `List()` queries are done with that token, even if you do not `rm` or `mv` anything.**

## Non-interactive mode

```
export VAULT_ADDR=<addr>
export VAULT_TOKEN=<token>
./vsh -c "rm secret/dir/to/remove"
```

## Local Development

Requirements:
- `golang` v1.12.7
- `docker` for integration testing
- `make` for simplified commands
