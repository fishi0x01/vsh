load ../../util/util
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'cat'" {
  #######################################
  echo "==== case: cat file ===="
  run ${APP_BIN} -c "cat ${KV_BACKEND}/src/dev/1"
  assert_success
  assert_line "value = 1"

  #######################################
  echo "==== case: cat non-existing file ===="
  run ${APP_BIN} -c "cat ${KV_BACKEND}/src/does/not/exist"
  assert_failure 1

  echo "ensure proper error message"
  assert_line --partial "Not a valid path for operation: /${KV_BACKEND}/src/does/not/exist"

  #######################################
  echo "==== case: cat directory ===="
  run ${APP_BIN} -c "cat ${KV_BACKEND}/src/dev"
  assert_failure 1
  assert_line --partial "Not a valid path for operation: /${KV_BACKEND}/src/dev"

  #######################################
  echo "==== case: cat ambiguous file ===="
  run ${APP_BIN} -c "cat ${KV_BACKEND}/src/tooling"
  assert_success
  assert_line "value = tooling"
  assert_line "drink = beer"
  assert_line "key = A"

  #######################################
  echo "==== case: cat ambiguous directory ===="
  run ${APP_BIN} -c "cat ${KV_BACKEND}/src/tooling/"
  assert_failure 1
  assert_line --partial "Not a valid path for operation: /${KV_BACKEND}/src/tooling/"
}
