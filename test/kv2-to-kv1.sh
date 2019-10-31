#!/bin/bash

source $(dirname ${0})/util.sh

export APP_BIN="./build/vsh_linux_amd64"
export VAULT_PORT=8888
export VAULT_TOKEN="root"
export VAULT_VERSION="1.2.3"
export VAULT_ADDR="http://localhost:${VAULT_PORT}"
export VAULT_CONTAINER_NAME="vault-kv2-to-kv1-test"
export VAULT_TEST_VALUE="test"

{ # Try
start_vault ${VAULT_VERSION} ${VAULT_CONTAINER_NAME} ${VAULT_PORT}

## Setup v2 KV
vault_exec ${VAULT_CONTAINER_NAME} "vault secrets enable -version=2 -path=secretkv2 kv"

vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secretkv2/source/a value=${VAULT_TEST_VALUE}"
vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secretkv2/source/b value=${VAULT_TEST_VALUE}"
vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secretkv2/source/x value=${VAULT_TEST_VALUE}"
vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secretkv2/source/c/d value=${VAULT_TEST_VALUE}"
vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secretkv2/source/c/e value=${VAULT_TEST_VALUE}"

vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secretkv2/remove/x value=${VAULT_TEST_VALUE}"
vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secretkv2/remove/y/z value=${VAULT_TEST_VALUE}"

## Setup v1 KV
vault_exec ${VAULT_CONTAINER_NAME} "vault secrets disable secret"
vault_exec ${VAULT_CONTAINER_NAME} "vault secrets enable -version=1 -path=secretkv1 kv"

## Run App
${APP_BIN} -c "cp secretkv2/source/x secretkv1/target2/x"
${APP_BIN} -c "mv secretkv2/source/ secretkv1/target/"

## Verify result
vault_value_must_be ${VAULT_CONTAINER_NAME} "secretkv1/target2/x" ${VAULT_TEST_VALUE}
vault_value_must_be ${VAULT_CONTAINER_NAME} "secretkv1/target/a" ${VAULT_TEST_VALUE}
vault_value_must_be ${VAULT_CONTAINER_NAME} "secretkv1/target/b" ${VAULT_TEST_VALUE}
} || { # Catch
  echo "Error running Tests"
  stop_vault ${VAULT_CONTAINER_NAME}
  exit 1
}

# Finally - Cleanup
stop_vault ${VAULT_CONTAINER_NAME}
