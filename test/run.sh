#!/bin/bash
set -e # required to fail test suite when a single test fails

VAULT_VERSIONS=("1.5.3" "1.0.0")
KV_BACKENDS=("KV1" "KV2")

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
export DIR
BATS="${DIR}/bin/core/bin/bats"

for vault_version in "${VAULT_VERSIONS[@]}"
do
    VAULT_VERSION=${vault_version} ${BATS} "${DIR}/special-tests/"

    for kv_backend in "${KV_BACKENDS[@]}"
    do
        VAULT_VERSION=${vault_version} KV_BACKEND="${kv_backend}" ${BATS} "${DIR}/command-tests/"
    done
done
