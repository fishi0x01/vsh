load ../../util/common
load ../../util/concurrency-setup
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} concurrency 'mv'" {
  #######################################
  echo "==== case: move large directory tree ===="
  run ${APP_BIN} -c "mv /KV2/src/ /KV2/dest/"
  assert_success

  echo "ensure at least one file got moved to destination"
  run get_vault_value "value" "/KV2/dest/a/a/50"
  assert_success
  assert_output "1"

  echo "ensure at least one src file got removed"
  run get_vault_value "value" "/KV2/src/a/a/50"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"
}
