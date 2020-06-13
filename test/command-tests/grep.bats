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
  echo "==== TODO case: grep term on directory with reduced permissions ===="

  #######################################
  echo "==== TODO case: grep term on ambigious directory ===="

  #######################################
  echo "==== TODO case: grep term on '/' ===="
}
