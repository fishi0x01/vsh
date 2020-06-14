load ../util/util
load ../bin/plugins/bats-support/load
load ../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'cat'" {
  #######################################
  echo "==== case: cat file ===="
  run ${APP_BIN} -c "cat ${KV_BACKEND}/src/dev/1"
  assert_success
  assert_line "value = 1"

  #######################################
  echo "==== case: cat directory ===="
  run ${APP_BIN} -c "cat ${KV_BACKEND}/src/dev"
  assert_success
  assert_output --partial "is not a file"

  #######################################
  echo "==== case: cat ambigious file ===="
  run ${APP_BIN} -c "cat ${KV_BACKEND}/src/tooling"
  assert_success
  assert_output "value = tooling"
}
