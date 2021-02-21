load ../../util/common
load ../../util/standard-setup
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} - print client version" {
  #######################################
  echo "==== case: print client version ===="
  export BIN_VERSION=$(git describe --tags --always --dirty)
  run ${APP_BIN} --version
  assert_success
  assert_line "${BIN_VERSION}"
}
