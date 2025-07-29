#!/usr/bin/env bats

load ../../util/common
load ../../util/standard-setup
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} - mount configuration via environment variables" {
  #######################################
  echo "==== case: no environment variables set - should use auto-discovered mounts ===="
  run ${APP_BIN} -v DEBUG -c "ls /"
  assert_success
  assert_line "KV1/"
  assert_line "KV2/"

  #######################################
  echo "==== case: VAULT_KV1_MOUNTS environment variable ===="
  run bash -c "VAULT_KV1_MOUNTS=custom-kv1,another-kv1 ${APP_BIN} -v DEBUG -c 'ls /'"
  assert_success
  assert_line "custom-kv1/"
  assert_line "another-kv1/"

  #######################################
  echo "==== case: VAULT_KV2_MOUNTS environment variable ===="
  run bash -c "VAULT_KV2_MOUNTS=custom-kv2,another-kv2 ${APP_BIN} -v DEBUG -c 'ls /'"
  assert_success
  assert_line "custom-kv2/"
  assert_line "another-kv2/"

  #######################################
  echo "==== case: both VAULT_KV1_MOUNTS and VAULT_KV2_MOUNTS set ===="
  run bash -c "VAULT_KV1_MOUNTS=custom-kv1 VAULT_KV2_MOUNTS=custom-kv2 ${APP_BIN} -v DEBUG -c 'ls /'"
  assert_success
  assert_line "custom-kv1/"
  assert_line "custom-kv2/"

  #######################################
  echo "==== case: mount paths with leading slash ===="
  run bash -c "VAULT_KV1_MOUNTS=/custom-kv1 VAULT_KV2_MOUNTS=/custom-kv2 ${APP_BIN} -v DEBUG -c 'ls /'"
  assert_success
  assert_line "custom-kv1/"
  assert_line "custom-kv2/"

  #######################################
  echo "==== case: mount paths without trailing slash ===="
  run bash -c "VAULT_KV1_MOUNTS=custom-kv1 VAULT_KV2_MOUNTS=custom-kv2 ${APP_BIN} -v DEBUG -c 'ls /'"
  assert_success
  assert_line "custom-kv1/"
  assert_line "custom-kv2/"

  #######################################
  echo "==== case: multiple mounts in single environment variable ===="
  run bash -c "VAULT_KV1_MOUNTS=kv1-a,kv1-b,kv1-c VAULT_KV2_MOUNTS=kv2-a,kv2-b ${APP_BIN} -v DEBUG -c 'ls /'"
  assert_success
  assert_line "kv1-a/"
  assert_line "kv1-b/"
  assert_line "kv1-c/"
  assert_line "kv2-a/"
  assert_line "kv2-b/"

  #######################################
  echo "==== case: empty environment variables should be ignored ===="
  run bash -c "VAULT_KV1_MOUNTS='' VAULT_KV2_MOUNTS='' ${APP_BIN} -v DEBUG -c 'ls /'"
  assert_success
  assert_line "KV1/"
  assert_line "KV2/"

  #######################################
  echo "==== case: environment variables with spaces should be handled ===="
  run bash -c "VAULT_KV1_MOUNTS=' kv1-a , kv1-b ' VAULT_KV2_MOUNTS=' kv2-a , kv2-b ' ${APP_BIN} -v DEBUG -c 'ls /'"
  assert_success
  assert_line " kv1-a /"
  assert_line " kv1-b /"
  assert_line " kv2-a /"
  assert_line " kv2-b /"

  #######################################
  echo "==== case: using configured mounts for operations ===="
  # First, create some test data in the custom mounts
  vault_exec "vault secrets enable -version=1 -path=custom-kv1 kv"
  vault_exec "vault secrets enable -version=2 -path=custom-kv2 kv"
  vault_exec "vault kv put custom-kv1/test key1=value1"
  vault_exec "vault kv put custom-kv2/test key2=value2"

  # Test listing with custom mounts
  run bash -c "VAULT_KV1_MOUNTS=custom-kv1 VAULT_KV2_MOUNTS=custom-kv2 ${APP_BIN} -c 'ls /custom-kv1/'"
  assert_success
  assert_line "test"

  run bash -c "VAULT_KV1_MOUNTS=custom-kv1 VAULT_KV2_MOUNTS=custom-kv2 ${APP_BIN} -c 'ls /custom-kv2/'"
  assert_success
  assert_line "test"

  # Test reading from custom mounts
  run bash -c "VAULT_KV1_MOUNTS=custom-kv1 VAULT_KV2_MOUNTS=custom-kv2 ${APP_BIN} -c 'cat /custom-kv1/test'"
  assert_success
  assert_line "key1 = value1"

  run bash -c "VAULT_KV1_MOUNTS=custom-kv1 VAULT_KV2_MOUNTS=custom-kv2 ${APP_BIN} -c 'cat /custom-kv2/test'"
  assert_success
  assert_line "key2 = value2"

  #######################################
  echo "==== case: default mount when no environment variables and no auto-discovery ===="
  # Test with a token that doesn't have list permissions
  run bash -c "VAULT_TOKEN=no-root ${APP_BIN} -v DEBUG -c 'ls /'"
  assert_success
  assert_output --partial "Cannot auto-discover mount backends"
  assert_output --partial "No KV mounts found or specified, adding default KV version 2 mount at /secrets"
  assert_line "secrets/"

  #######################################
  echo "==== case: environment variables are added to auto-discovery ===="
  # Environment variables are added alongside auto-discovered mounts
  run bash -c "VAULT_KV1_MOUNTS=custom-kv1 VAULT_KV2_MOUNTS=custom-kv2 ${APP_BIN} -v DEBUG -c 'ls /'"
  assert_success
  assert_line "custom-kv1/"
  assert_line "custom-kv2/"
  # Should also include the auto-discovered mounts
  assert_line "KV1/"
  assert_line "KV2/"
}
