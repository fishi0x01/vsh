#!/bin/bash

VAULT_VERSIONS=("1.0.0" "1.4.2")
KV_BACKENDS=("KV1" "KV2")

export DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

for vault_version in "${VAULT_VERSIONS[@]}"
do
    VAULT_VERSION=${vault_version} ${DIR}/bin/core/bin/bats ${DIR}/special-tests

    for kv_backend in "${KV_BACKENDS[@]}"
    do
        VAULT_VERSION=${vault_version} KV_BACKEND="${kv_backend}" ${DIR}/bin/core/bin/bats ${DIR}/command-tests
    done
done
