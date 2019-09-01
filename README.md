### Status
[![CircleCI](https://circleci.com/gh/fishi0x01/vsh.svg?style=svg)](https://circleci.com/gh/fishi0x01/vsh)
[![Go Report Card](https://goreportcard.com/badge/github.com/fishi0x01/vsh)](https://goreportcard.com/report/github.com/fishi0x01/vsh)
[![Code Climate](https://codeclimate.com/github/fishi0x01/vsh/badges/gpa.svg)](https://codeclimate.com/github/fishi0x01/vsh)

# vsh

vsh is a simple interactive HashiCorp Vault shell which treats vault secret paths like directories. That way you can do recursive operations on paths. Vault KV v1 and v2 are both supported.
Commands can also be executed in a non-interacitve way.

## Local Development

Requirements:
- `golang` v1.12.7
- `docker` for integration testing
- `make` for simplified commands