load ../../util/common
load ../../util/standard-setup
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} KV migrations with 'mv'" {
  #######################################
  echo "==== case: move single file from KV1 to KV2 ===="
  run ${APP_BIN} -c "mv KV1/src/prod/all KV2/dest/prod/all"
  assert_success

  echo "ensure the file got moved to destination"
  run get_vault_value "value" "KV2/dest/prod/all"
  assert_success
  assert_output "all"
  run get_vault_value "example" "KV2/dest/prod/all"
  assert_success
  assert_output "test"

  echo "ensure the src file got removed"
  run get_vault_value "value" "KV1/src/prod/all"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  #######################################
  echo "==== case: move single directory from KV1 to KV2 ===="
  run ${APP_BIN} -c "mv KV1/src/dev KV2/dest/dev"
  assert_success

  echo "ensure the directory got moved to destination"
  run get_vault_value "value" "KV2/dest/dev/1"
  assert_success
  assert_output "1"

  run get_vault_value "value" "KV2/dest/dev/2"
  assert_success
  assert_output "2"

  run get_vault_value "value" "KV2/dest/dev/3"
  assert_success
  assert_output "3"

  echo "ensure the src directory got removed"
  run get_vault_value "value" "KV1/src/dev/1"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  run get_vault_value "value" "KV1/src/dev/2"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  run get_vault_value "value" "KV1/src/dev/3"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  #######################################
  echo "==== case: move single file from KV2 to KV1 ===="
  run ${APP_BIN} -c "mv KV2/src/prod/all KV1/dest/prod/all"
  assert_success

  echo "ensure the file got moved to destination"
  run get_vault_value "value" "KV1/dest/prod/all"
  assert_success
  assert_output "all"
  run get_vault_value "example" "KV1/dest/prod/all"
  assert_success
  assert_output "test"

  echo "ensure the src file got removed"
  run get_vault_value "value" "KV2/src/prod/all"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  #######################################
  echo "==== case: move single directory from KV1 to KV2 ===="
  run ${APP_BIN} -c "mv KV2/src/dev KV1/dest/dev"
  assert_success

  echo "ensure the directory got moved to destination"
  run get_vault_value "value" "KV1/dest/dev/1"
  assert_success
  assert_output "1"

  run get_vault_value "value" "KV1/dest/dev/2"
  assert_success
  assert_output "2"

  run get_vault_value "value" "KV1/dest/dev/3"
  assert_success
  assert_output "3"

  echo "ensure the src directory got removed"
  run get_vault_value "value" "KV2/src/dev/1"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  run get_vault_value "value" "KV2/src/dev/2"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  run get_vault_value "value" "KV2/src/dev/3"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"
}
