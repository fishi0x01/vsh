#!/bin/bash

export VAULT_VERSION=${VAULT_VERSION:-"1.16.2"}
export VAULT_CONTAINER_NAME="vsh-integration-test-vault"
export VAULT_HOST_PORT=${VAULT_HOST_PORT:-"8888"}

export VAULT_TOKEN="root"
export VAULT_ADDR="http://localhost:${VAULT_HOST_PORT}"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
export DIR
UNAME=$(uname | tr '[:upper:]' '[:lower:]')
case "$(uname -m)" in
  x86_64)
    ARCH=amd64 ;;
  arm64|aarch64|armv8b|armv8l)
    ARCH=arm64 ;;
  arm*)
    ARCH=arm ;;
  i386|i686)
    ARCH=386 ;;
  *)
    ARCH=$(uname -m) ;;
esac
export ARCH
export APP_BIN="${DIR}/../../build/vsh_${UNAME}_${ARCH}"
export NO_VALUE_FOUND="No value found at"

teardown() {
    docker rm -f ${VAULT_CONTAINER_NAME} &> /dev/null
}

vault_exec() {
    vault_exec_output "$@" &> /dev/null
}

vault_exec_output() {
    docker exec ${VAULT_CONTAINER_NAME} /bin/sh -c "$1"
}

get_vault_value() {
    vault_exec_output "vault kv get -field=\"${1}\" \"${2}\" || true"
}
