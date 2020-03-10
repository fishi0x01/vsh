#!/bin/bash

source $(dirname ${0})/util.sh

export APP_BIN="./build/vsh_linux_amd64"
export VAULT_PORT=8888
export VAULT_VERSION="1.3.3"
export VAULT_ADDR="http://localhost:${VAULT_PORT}"
# Paths are relative to Makefile
export VAULT_CONFIG_PATH="$(pwd)/test/vault-config"
export VAULT_CONTAINER_NAME="vault-kv2-auth-test"
export VAULT_TEST_VALUE="test"

{ # Try

## Setup v2 KV
start_vault ${VAULT_VERSION} ${VAULT_CONTAINER_NAME} ${VAULT_PORT}

vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/a value=${VAULT_TEST_VALUE}"

## Run App
${APP_BIN} -c "mv secret/source/a secret/target2/a"

## Verify result
vault_value_must_be ${VAULT_CONTAINER_NAME} "secret/target2/a" ${VAULT_TEST_VALUE}
} || { # Catch
  echo "Error running Tests"
  stop_vault ${VAULT_CONTAINER_NAME}
  exit 1
}

# Finally - Cleanup
stop_vault ${VAULT_CONTAINER_NAME}
