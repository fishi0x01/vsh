load ../../util/util
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} - print client version" {
  #######################################
  export BIN_VERSION=$(git describe --tags --always --dirty)
  echo "==== case: print client version ===="
  run ${APP_BIN} -version
  assert_success
  assert_line "${BIN_VERSION}"
}
