#!/bin/bash

export VAULT_VERSION=${VAULT_VERSION:-"1.6.1"}
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
    rm -f vsh_trace.log
    docker run -d \
        --name=${VAULT_CONTAINER_NAME} \
        -p "${VAULT_HOST_PORT}:8200" \
        --cap-add=IPC_LOCK \
        -e "VAULT_ADDR=http://127.0.0.1:8200" \
        -e "VAULT_TOKEN=root" \
        -e "VAULT_DEV_ROOT_TOKEN_ID=root" \
        -e "VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200" \
        "vault:${VAULT_VERSION}" &> /dev/null
    docker cp "$DIR/policy-no-root.hcl" ${VAULT_CONTAINER_NAME}:.
    docker cp "$DIR/policy-delete-only.hcl" ${VAULT_CONTAINER_NAME}:.
    # need some time for GH Actions CI
    sleep 3
    vault_exec "vault secrets disable secret"
    vault_exec "vault policy write no-root policy-no-root.hcl"
    vault_exec "vault token create -id=no-root -policy=no-root"
    vault_exec "vault policy write delete-only policy-delete-only.hcl"
    vault_exec "vault token create -id=delete-only -policy=delete-only"

    KV_BACKENDS=("KV1" "KV2")
    vault_exec "vault secrets enable -version=1 -path=KV1 kv"
    vault_exec "vault secrets enable -version=2 -path=KV2 kv"
    for kv_backend in "${KV_BACKENDS[@]}"
    do
        vault_exec "vault kv put ${kv_backend}/src/data/1 data=1"
        vault_exec "vault kv put ${kv_backend}/src/data/2 value=1 data=2"
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
        vault_exec "vault kv put ${kv_backend}/src/ambivalence/1 value=1 fruit=apple"
        vault_exec "vault kv put ${kv_backend}/src/ambivalence/1/a value=2 fruit=banana"
        vault_exec "vault kv put ${kv_backend}/src/a/foo value=1"
        vault_exec "vault kv put ${kv_backend}/src/a/foo/bar value=2"
        vault_exec "vault kv put ${kv_backend}/src/b/foo value=1"
        vault_exec "vault kv put ${kv_backend}/src/b/foo/bar value=2"
        vault_exec "echo -n \"a spaced value\" | vault kv put ${kv_backend}/src/spaces/foo bar=-"
        vault_exec "vault kv put ${kv_backend}/src/apostrophe/foo bar=steve\'s"
        vault_exec "echo -n 'a \"quoted\" value' | vault kv put ${kv_backend}/src/quoted/foo bar=-"
    done
}

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
