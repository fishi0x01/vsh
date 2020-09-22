load ../util/util
load ../bin/plugins/bats-support/load
load ../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'grep'" {
  #######################################
  echo "==== case: grep term on '${KV_BACKEND}/' ===="
  run ${APP_BIN} -c "grep value ${KV_BACKEND}/"
  assert_success
  assert_line --partial "/${KV_BACKEND}/src/dev/1"
  assert_line --partial "/${KV_BACKEND}/src/dev/2"
  assert_line --partial "/${KV_BACKEND}/src/dev/3"
  assert_line --partial "/${KV_BACKEND}/src/prod/all"

  #######################################
  echo "==== case: grep non-existing file ===="
  run ${APP_BIN} -c "grep test ${KV_BACKEND}/src/does/not/exist"
  assert_success

  echo "ensure proper error message"
  assert_line --partial "Not a valid path for operation: /${KV_BACKEND}/src/does/not/exist"

  #######################################
  echo "==== case: grep term on ambigious directory ===="
  run ${APP_BIN} -c "grep juice ${KV_BACKEND}/src/tooling/"
  assert_line --partial "/${KV_BACKEND}/src/tooling/v1"

  #######################################
  echo "==== case: grep term on ambigious file ===="
  run ${APP_BIN} -c "grep beer ${KV_BACKEND}/src/tooling"
  assert_line --partial "/${KV_BACKEND}/src/tooling"

  #######################################
  echo "==== TODO case: grep term on directory with reduced permissions ===="

  #######################################
  echo "==== TODO case: grep term on '/' ===="
}
