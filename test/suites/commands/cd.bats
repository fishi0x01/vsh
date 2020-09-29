load ../../util/util
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'cd'" {
  #######################################
  echo "==== case: cd to sub-sub-dir ===="
  run ${APP_BIN} -c "cd ${KV_BACKEND}/src/dev"
  assert_success

  #######################################
  echo "==== case: cd to KV ===="
  run ${APP_BIN} -c "cd ${KV_BACKEND}/"
  assert_success

  #######################################
  echo "==== case: cd to KV first level sub-dir ===="
  run ${APP_BIN} -c "cd ${KV_BACKEND}/src/"
  assert_success

  #######################################
  echo "==== case: cd to non-existing dir ===="
  run ${APP_BIN} -c "cd ${KV_BACKEND}/src/does/not/exist"
  assert_success
}
