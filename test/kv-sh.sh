#!/bin/bash

source $(dirname ${0})/util.sh

export APP_BIN="./build/vsh_linux_amd64"
export VAULT_PORT=8889
export VAULT_TOKEN="root"
export VAULT_VERSION="1.3.1"
export VAULT_ADDR="http://localhost:${VAULT_PORT}"
export VAULT_CONTAINER_NAME="vault"

{ # Try
start_vault ${VAULT_VERSION} ${VAULT_CONTAINER_NAME} ${VAULT_PORT}

## Setup v2 KV
vault_exec ${VAULT_CONTAINER_NAME} "vault secrets disable secret"
vault_exec ${VAULT_CONTAINER_NAME} "vault secrets enable -version=2 -path=secretkv2 kv"

vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secretkv2/departmentA/keyA value=secret-key-A"
vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secretkv2/departmentA/keyB value=secret-key-B"
vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secretkv2/departmentA/subkeys/keyC value=secret-key-C"
vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secretkv2/departmentA/subkeys/keyD value=secret-key-D"

## Setup v1 KV
vault_exec ${VAULT_CONTAINER_NAME} "vault secrets enable -version=1 -path=secretkv1 kv"

vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secretkv1/departmentB/secretA value=secretA"
vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secretkv1/departmentB/secretB value=secretB"
vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secretkv1/departmentB/subsecret/secretC value=secretC"
}
