load ../../util/common
load ../../util/standard-setup
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} - test issue https://github.com/fishi0x01/vsh/issues/130" {
  #######################################
  echo "==== case: add new key to KV2 path with existing keys ===="
  run ${APP_BIN} -c "add --confirm newkey newvalue KV2/src/a/foo"
  assert_success

  echo "ensure the new key was written"
  run get_vault_value "newkey" "KV2/src/a/foo"
  assert_success
  assert_output "newvalue"

  echo "ensure pre-existing keys were not destroyed by the add operation"
  run get_vault_value "value" "KV2/src/a/foo"
  assert_success
  assert_output "1"

  run get_vault_value "long" "KV2/src/a/foo"
  assert_success
  assert_output "this-is-a-really-long-value-for-testing"
}
