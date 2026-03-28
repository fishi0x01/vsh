load ../../util/common
load ../../util/standard-setup
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} - test issue https://github.com/fishi0x01/vsh/issues/129" {
  #######################################
  echo "==== case: add key using relative path while cwd is set ===="
  run bash -c "VAULT_PATH=KV2/src/a ${APP_BIN} -c 'add --confirm newkey newvalue foo'"
  assert_success

  echo "ensure the key was written to the correct absolute path"
  run get_vault_value "newkey" "KV2/src/a/foo"
  assert_success
  assert_output "newvalue"
}
