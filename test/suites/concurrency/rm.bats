load ../../util/common
load ../../util/concurrency-setup
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} concurrency 'rm'" {
  #######################################
  echo "==== case: remove large directory tree ===="
  run ${APP_BIN} -c "rm -r /KV2/src/"
  assert_success

  echo "ensure at least one src file got removed"
  run get_vault_value "value" "/KV2/src/a/a/50"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"
}
