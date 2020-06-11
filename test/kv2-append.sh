#!/bin/bash
set -x
set -e

# shellcheck source=test/util.sh
source "$(dirname ${0})/util.sh"
UNAME=$(uname | tr '[:upper:]' '[:lower:]')

export APP_BIN="./build/vsh_${UNAME}_amd64"
export VAULT_PORT=8888
export VAULT_TOKEN="root"
export VAULT_VERSION="1.3.4"
export VAULT_ADDR="http://localhost:${VAULT_PORT}"
export VAULT_CONTAINER_NAME="vault-kv2-test"

export VALUE1="A"
export VALUE2="B"

cleanup() {
  ${APP_BIN} -c "rm /secret/target" &&
  ${APP_BIN} -c "rm /secret/source"
}

## CASE: Append value to not existing destination
testCase1() {
  echo "=== TEST: Append value to not existing destination" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/a value1=A" &&
  ${APP_BIN} -c "append /secret/source/a /secret/target/a" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "value1" "/secret/target/a" "A" &&
  echo "PASSED" &&
  cleanup
}

## CASE: Append value to existing destination
testCase2() {
  echo "=== TEST: Append value to existing destination" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/a key1=A" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/target/dest key2=B" &&
  # Check initial-state
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key1" "/secret/source/a" "A" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key2" "/secret/target/dest" "B" &&
  # Do append
  ${APP_BIN} -c "append /secret/source/a /secret/target/dest" &&
  # Sources must remain!
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key1" "/secret/source/a" "A" &&
  # Merged secret must have correct values
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key1" "/secret/target/dest" "A" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key2" "/secret/target/dest" "B" &&
  echo "PASSED" &&
  cleanup
}

## CASE: Append multiple values to existing destination
testCase3() {
  echo "=== TEST: Append multiple values to existing destination" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/a key1=A" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/b key2=B" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/c key3=C" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/target/dest key0=O" &&
  # Check initial-state
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key1" "/secret/source/a" "A" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key2" "/secret/source/b" "B" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key3" "/secret/source/c" "C" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key0" "/secret/target/dest" "O" &&
  # Do append
  ${APP_BIN} -c "append /secret/source/a /secret/target/dest" &&
  ${APP_BIN} -c "append /secret/source/b /secret/target/dest" &&
  ${APP_BIN} -c "append /secret/source/c /secret/target/dest" &&
  # Sources must remain!
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key1" "/secret/source/a" "A" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key2" "/secret/source/b" "B" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key3" "/secret/source/c" "C" &&
  # Merged secret must have all values
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key0" "/secret/target/dest" "O" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key1" "/secret/target/dest" "A" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key2" "/secret/target/dest" "B" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key3" "/secret/target/dest" "C" &&
  echo "PASSED" &&
  cleanup
}

## CASE: Append conflicting keys - skip strategy (implicit flag)
testCase4a() {
  echo "=== TEST: Append conflicting keys - skip strategy (implicit flag)" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/a key=A" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/target/dest key=B" &&
  # Check initial-state
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/source/a" "A" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/target/dest" "B" &&
  # Do append
  ${APP_BIN} -c "append /secret/source/a /secret/target/dest" &&
  # Sources must remain!
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/source/a" "A" &&
  # Merged secret must have all values
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/target/dest" "B" &&
  echo "PASSED" &&
  cleanup
}

## CASE: Append conflicting keys - skip strategy (explicit flag short)
testCase4b() {
  echo "=== TEST: Append conflicting keys - skip strategy (explicit flag short)" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/a key=A" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/target/dest key=B" &&
  # Check initial-state
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/source/a" "A" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/target/dest" "B" &&
  # Do append
  ${APP_BIN} -c "append /secret/source/a /secret/target/dest -s" &&
  # Sources must remain!
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/source/a" "A" &&
  # Merged secret must have all values
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/target/dest" "B" &&
  echo "PASSED" &&
  cleanup
}

## CASE: Append conflicting keys - skip strategy (explicit flag long)
testCase4c() {
  echo "=== TEST: Append conflicting keys - skip strategy (explicit flag long)" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/a key=A" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/target/dest key=B" &&
  # Check initial-state
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/source/a" "A" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/target/dest" "B" &&
  # Do append
  ${APP_BIN} -c "append /secret/source/a /secret/target/dest --skip" &&
  # Sources must remain!
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/source/a" "A" &&
  # Merged secret must have all values
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/target/dest" "B" &&
  echo "PASSED" &&
  cleanup
}

## CASE: Append conflicting keys - overwrite strategy
testCase5() {
  echo "=== TEST: Append conflicting keys - overwrite strategy" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/a key=A" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/target/dest key=B" &&
  # Check initial-state
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/source/a" "A" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/target/dest" "B" &&
  # Do append
  ${APP_BIN} -c "append /secret/source/a /secret/target/dest --force" &&
  # Sources must remain!
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/source/a" "A" &&
  # Merged secret must have all values
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/target/dest" "A" &&
  echo "PASSED" &&
  cleanup
}

## CASE: Append conflicting keys - rename strategy
testCase6a() {
  echo "=== TEST: Append conflicting keys - rename strategy" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/a key=A" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/target/dest key=B" &&
  # Check initial-state
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/source/a" "A" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/target/dest" "B" &&
  # Do append
  ${APP_BIN} -c "append /secret/source/a /secret/target/dest --rename" &&
  # Sources must remain!
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/source/a" "A" &&
  # Merged secret must have all values
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/target/dest" "B" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key_1" "/secret/target/dest" "A" &&
  echo "PASSED" &&
  cleanup
}

## CASE: Append conflicting keys - rename strategy (deep rename)
testCase6b() {
  echo "=== TEST: Append conflicting keys - rename strategy (deep rename)" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/source/a key=A" &&
  vault_exec ${VAULT_CONTAINER_NAME} "vault kv put secret/target/dest key=B key_1=B1 key_2=B2" &&
  # Check initial-state
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/source/a" "A" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/target/dest" "B" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key_1" "/secret/target/dest" "B1" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key_2" "/secret/target/dest" "B2" &&
  # Do append
  ${APP_BIN} -c "append /secret/source/a /secret/target/dest --rename" &&
  # Sources must remain!
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/source/a" "A" &&
  # Merged secret must have all values
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key" "/secret/target/dest" "B" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key_1" "/secret/target/dest" "B1" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key_2" "/secret/target/dest" "B2" &&
  vault_filed_must_be ${VAULT_CONTAINER_NAME} "key_3" "/secret/target/dest" "A" &&
  echo "PASSED" &&
  cleanup
}

{ # Try

## Setup v2 KV
start_vault "${VAULT_VERSION}" "${VAULT_CONTAINER_NAME}" "${VAULT_PORT}"
testCase1 &&
testCase2 &&
testCase3 &&
testCase4a &&
testCase4b &&
testCase4c &&
testCase5 &&
testCase6a &&
testCase6b &&
true

} || { # Catch
  echo "Tests failed"
  stop_vault ${VAULT_CONTAINER_NAME}
  exit 1
}

# Finally - Cleanup
stop_vault ${VAULT_CONTAINER_NAME}
