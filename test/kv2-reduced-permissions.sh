#!/bin/bash

source $(dirname ${0})/util.sh

export APP_BIN="./build/vsh_linux_amd64"
export VAULT_PORT=8888
export VAULT_TOKEN="root"
export VAULT_VERSION="1.2.3"
export VAULT_ADDR="http://localhost:${VAULT_PORT}"
export VAULT_CONTAINER_NAME="vault-kv1-reduced-permissions"
export VAULT_TEST_VALUE="test"

{ # Try

## Setup v2 KV
start_vault ${VAULT_VERSION} ${VAULT_CONTAINER_NAME} ${VAULT_PORT}

docker cp test/reduced-policy.hcl ${VAULT_CONTAINER_NAME}:.
vault_exec ${VAULT_CONTAINER_NAME} "vault policy write reduced-access reduced-policy.hcl"
vault_exec ${VAULT_CONTAINER_NAME} "vault token create -id=reduced -policy=reduced-access"

vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/a value=${VAULT_TEST_VALUE}"
vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/b value=${VAULT_TEST_VALUE}"
vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/x value=${VAULT_TEST_VALUE}"
vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/c/d value=${VAULT_TEST_VALUE}"
vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/c/e value=${VAULT_TEST_VALUE}"

vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/remove/x value=${VAULT_TEST_VALUE}"
vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/remove/y/z value=${VAULT_TEST_VALUE}"

## Run App
export VAULT_TOKEN="reduced"
${APP_BIN} -c "mv secret/source/x secret/target2/x"
${APP_BIN} -c "cp secret/source/ secret/target/"
${APP_BIN} -c "rm secret/remove"

## Verify result
vault_value_must_be ${VAULT_CONTAINER_NAME} "secret/target2/x" ${VAULT_TEST_VALUE}
vault_value_must_be ${VAULT_CONTAINER_NAME} "secret/target/a" ${VAULT_TEST_VALUE}
} || { # Catch
  echo "Error running Tests"
  stop_vault ${VAULT_CONTAINER_NAME}
  exit 1
}

# Finally - Cleanup
stop_vault ${VAULT_CONTAINER_NAME}
