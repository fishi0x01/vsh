#!/bin/bash
set -e # required to fail test suite when a single test fails

VAULT_VERSION=${VAULT_VERSION:-"1.8.3"}
KV_BACKEND=${KV_BACKEND:-"KV2"}
TEST_SUITE=${TEST_SUITE:-"commands/cp"}

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
export DIR
BATS="${DIR}/bin/core/bin/bats"

VAULT_VERSION=${VAULT_VERSION} KV_BACKEND="${KV_BACKEND}" ${BATS} "${DIR}/suites/${TEST_SUITE}.bats"
