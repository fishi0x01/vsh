load ../util/util
load ../bin/plugins/bats-support/load
load ../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} KV migrations with 'append'" {
  #######################################
  echo "==== case: append file from KV1 to KV2 ===="
  run ${APP_BIN} -c "append KV1/src/prod/all KV2/src/staging/all/v1"
  assert_success

  echo "ensure the file got appended to destination"
  run get_vault_value "value" "KV2/src/staging/all/v1"
  assert_success
  assert_output "v1"
  run get_vault_value "example" "KV2/src/staging/all/v1"
  assert_success
  assert_output "test"
  run get_vault_value "tree" "KV2/src/staging/all/v1"
  assert_success
  assert_output "oak"

  echo "ensure the src file still exists"
  run get_vault_value "value" "KV1/src/prod/all"
  assert_success
  assert_output "all"
  run get_vault_value "example" "KV1/src/prod/all"
  assert_success
  assert_output "test"

  #######################################
  echo "==== case: append file from KV1 to KV2 ===="
  run ${APP_BIN} -c "append KV1/src/staging/all/v1 KV2/src/dev/2"
  assert_success

  echo "ensure the file got appended to destination"
  run get_vault_value "value" "KV2/src/dev/2"
  assert_success
  assert_output "2"
  run get_vault_value "tree" "KV2/src/dev/2"
  assert_success
  assert_output "oak"
  run get_vault_value "fruit" "KV2/src/dev/2"
  assert_success
  assert_output "banana"

  echo "ensure the src file still exists"
  run get_vault_value "value" "KV1/src/dev/2"
  assert_success
  assert_output "2"
  run get_vault_value "fruit" "KV1/src/dev/2"
  assert_success
  assert_output "banana"
}
