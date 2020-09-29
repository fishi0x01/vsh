load ../../util/util
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'append'" {
  #######################################
  echo "==== case: append value to non existing destination ===="
  run ${APP_BIN} -c "append ${KV_BACKEND}/src/prod/all ${KV_BACKEND}/dest/prod/all"
  assert_success

  echo "ensure the file got appended to destination"
  run get_vault_value "value" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "all"
  run get_vault_value "example" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "test"

  echo "ensure the src file still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/prod/all"
  assert_success
  assert_output "all"
  run get_vault_value "example" "${KV_BACKEND}/src/prod/all"
  assert_success
  assert_output "test"

  #######################################
  echo "==== case: append to non-existing file ===="
  run ${APP_BIN} -c "append ${KV_BACKEND}/src/does/not/exist ${KV_BACKEND}/src/aa"
  assert_success

  echo "ensure proper error message"
  assert_line --partial "Not a valid path for operation: /${KV_BACKEND}/src/does/not/exist"

  #######################################
  echo "==== case: append value to existing destination with conflicting keys (default merge strategy) ===="
  run ${APP_BIN} -c "append ${KV_BACKEND}/src/dev/1 ${KV_BACKEND}/dest/prod/all"
  assert_success

  echo "ensure the file got appended to destination"
  run get_vault_value "value" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "all"
  run get_vault_value "example" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "test"
  run get_vault_value "fruit" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "apple"

  #######################################
  echo "==== case: append value to existing destination with conflicting keys (skip strategy) ===="
  run ${APP_BIN} -c "append ${KV_BACKEND}/src/tooling/v1 ${KV_BACKEND}/dest/prod/all --skip"
  assert_success

  echo "ensure the file got appended to destination"
  run get_vault_value "value" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "all"
  run get_vault_value "example" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "test"
  run get_vault_value "fruit" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "apple"
  run get_vault_value "drink" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "juice"
  run get_vault_value "key" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "B"

  #######################################
  echo "==== case: append value to existing destination with conflicting keys (overwrite strategy) ===="
  run ${APP_BIN} -c "append ${KV_BACKEND}/src/tooling/v2 ${KV_BACKEND}/dest/prod/all --force"
  assert_success

  echo "ensure the file got appended to destination"
  run get_vault_value "value" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "v2"
  run get_vault_value "example" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "test"
  run get_vault_value "fruit" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "apple"
  run get_vault_value "drink" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "water"
  run get_vault_value "key" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "C"

  #######################################
  echo "==== case: append value to existing destination with conflicting keys (rename strategy) ===="
  run ${APP_BIN} -c "append ${KV_BACKEND}/src/tooling/v1 ${KV_BACKEND}/dest/prod/all --rename"
  assert_success

  echo "ensure the file got appended to destination"
  run get_vault_value "value" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "v2"
  run get_vault_value "example" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "test"
  run get_vault_value "fruit" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "apple"
  run get_vault_value "drink" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "water"
  run get_vault_value "key" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "C"
  run get_vault_value "drink_1" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "juice"
  run get_vault_value "key_1" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "B"
  run get_vault_value "value_1" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "v1"
}
