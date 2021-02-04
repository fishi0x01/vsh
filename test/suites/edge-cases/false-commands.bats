load ../../util/util
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} - false commands" {
  #######################################
  echo "==== case: non-existing flag ===="
  run ${APP_BIN} -x not
  assert_line --partial "unknown argument -x"
  assert_failure 255

  echo "==== case: non-existing command ===="
  run ${APP_BIN} -c "nono xD"
  assert_failure 1
}
