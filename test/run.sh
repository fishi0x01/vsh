#!/bin/bash
set -e # required to fail test suite when a single test fails

VAULT_VERSIONS=("1.5.4" "1.0.0")
KV_BACKENDS=("KV1" "KV2")

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
export DIR
BATS="${DIR}/bin/core/bin/bats"

for vault_version in "${VAULT_VERSIONS[@]}"
do
    #VAULT_VERSION=${vault_version} ${BATS} "${DIR}/suites/past-issues/"
    #VAULT_VERSION=${vault_version} ${BATS} "${DIR}/suites/edge-cases/"

    for kv_backend in "${KV_BACKENDS[@]}"
    do
        VAULT_VERSION=${vault_version} KV_BACKEND="${kv_backend}" ${BATS} "${DIR}/suites/commands/dump.bats"
    done
done
