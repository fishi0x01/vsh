#!/bin/bash

export VAULT_VERSION=${VAULT_VERSION:-"1.3.4"}
export KV_BACKEND=${KV_BACKEND:-"KV2"}
export VAULT_CONTAINER_NAME="vsh-integration-test-vault"
export VAULT_HOST_PORT=${VAULT_HOST_PORT:-"8888"}

export VAULT_TOKEN="root"
export VAULT_ADDR="http://localhost:${VAULT_HOST_PORT}"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
export DIR
UNAME=$(uname | tr '[:upper:]' '[:lower:]')
export APP_BIN="${DIR}/../../build/vsh_${UNAME}_amd64"
export NO_VALUE_FOUND="No value found at"

setup() {
    docker run -d \
        --name=${VAULT_CONTAINER_NAME} \
        -p "${VAULT_HOST_PORT}:8200" \
        --cap-add=IPC_LOCK \
        -e "VAULT_ADDR=http://127.0.0.1:8200" \
        -e "VAULT_TOKEN=root" \
        -e "VAULT_DEV_ROOT_TOKEN_ID=root" \
        -e "VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200" \
        "vault:${VAULT_VERSION}" &> /dev/null
    docker cp "$DIR/reduced-policy.hcl" ${VAULT_CONTAINER_NAME}:.
    # need some time for GH Actions CI
    sleep 3
    vault_exec "vault secrets disable secret"
    vault_exec "vault policy write reduced-access reduced-policy.hcl"
    vault_exec "vault token create -id=reduced -policy=reduced-access"

    KV_BACKENDS=("KV1" "KV2")
    vault_exec "vault secrets enable -version=1 -path=KV1 kv"
    vault_exec "vault secrets enable -version=2 -path=KV2 kv"
    for kv_backend in "${KV_BACKENDS[@]}"
    do
        vault_exec "vault kv put ${kv_backend}/src/dev/1 value=1 fruit=apple"
        vault_exec "vault kv put ${kv_backend}/src/dev/2 value=2 fruit=banana"
        vault_exec "vault kv put ${kv_backend}/src/dev/3 value=3 fruit=berry"
        vault_exec "vault kv put ${kv_backend}/src/staging/all value=all tree=palm"
        vault_exec "vault kv put ${kv_backend}/src/staging/all/v1 value=v1 tree=oak"
        vault_exec "vault kv put ${kv_backend}/src/staging/all/v2 value=v2 tree=bonsai"
        vault_exec "vault kv put ${kv_backend}/src/prod/all value=all example=test"
        vault_exec "vault kv put ${kv_backend}/src/tooling value=tooling drink=beer key=A"
        vault_exec "vault kv put ${kv_backend}/src/tooling/v1 value=v1 drink=juice key=B"
        vault_exec "vault kv put ${kv_backend}/src/tooling/v2 value=v2 drink=water key=C"
    done
}

teardown() {
    docker rm -f ${VAULT_CONTAINER_NAME} &> /dev/null
}

vault_exec() {
    docker exec ${VAULT_CONTAINER_NAME} ${1} &> /dev/null
}

get_vault_value() {
    docker exec ${VAULT_CONTAINER_NAME} vault kv get -field="${1}" "${2}" || true
}
