load ../../util/common
load ../../util/standard-setup
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} - test issue https://github.com/fishi0x01/vsh/issues/114" {
  #######################################
  echo "==== case: https://github.com/fishi0x01/vsh/issues/114 ===="
  run ${APP_BIN} -c "cp KV1/src/data KV1/dest/data"
  assert_success

  echo "ensure the file got copied to destination"
  run get_vault_value "data" "KV1/dest/data"
  assert_success
  assert_output "2"
}
