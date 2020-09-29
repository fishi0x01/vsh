load ../../util/util
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'ls'" {
  #######################################
  echo "==== case: list directory ===="
  run ${APP_BIN} -c "ls ${KV_BACKEND}/src/dev"
  assert_success
  assert_line --index 0 "1"
  assert_line --index 1 "2"
  assert_line --index 2 "3"

  #######################################
  echo "==== case: ls non-existing dir ===="
  run ${APP_BIN} -c "ls ${KV_BACKEND}/src/does/not/exist"
  assert_success

  echo "ensure proper error message"
  assert_line --partial "Not a valid path for operation: /${KV_BACKEND}/src/does/not/exist"

  #######################################
  echo "==== case: list backends ===="
  run ${APP_BIN} -c "ls /"
  assert_success
  assert_line "KV1/"
  assert_line "KV2/"

  #######################################
  echo "==== case: list backends with reduced permissions ===="
  run bash -c "VAULT_TOKEN=no-root ${APP_BIN} -v -c 'ls /'"
  assert_success
  assert_output --partial "Cannot auto-discover mount backends"

  #######################################
  echo "==== case: list directory with reduced permissions ===="
  run bash -c "VAULT_TOKEN=no-root ${APP_BIN} -c 'ls ${KV_BACKEND}/src/dev'"
  assert_success
  assert_line "1"
  assert_line "2"
  assert_line "3"
}
