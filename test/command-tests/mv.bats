load ../util/util
load ../bin/plugins/bats-support/load
load ../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'mv'" {
  #######################################
  echo "==== case: move single file ===="
  run ${APP_BIN} -c "mv ${KV_BACKEND}/src/prod/all ${KV_BACKEND}/dest/prod/all"
  assert_success

  echo "ensure the file got moved to destination"
  run get_vault_value "value" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "all"
  run get_vault_value "example" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "test"

  echo "ensure the src file got removed"
  run get_vault_value "value" "${KV_BACKEND}/src/prod/all"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  #######################################
  echo "==== case: move single directory without trailing '/' ===="
  run ${APP_BIN} -c "mv ${KV_BACKEND}/src/dev ${KV_BACKEND}/dest/dev"
  assert_success

  echo "ensure the directory got moved to destination"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev/1"
  assert_success
  assert_output "1"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev/2"
  assert_success
  assert_output "2"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev/3"
  assert_success
  assert_output "3"

  echo "ensure the src directory got removed"
  run get_vault_value "value" "${KV_BACKEND}/src/dev/1"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"
  run get_vault_value "value" "${KV_BACKEND}/src/dev/2"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"
  run get_vault_value "value" "${KV_BACKEND}/src/dev/3"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  #######################################
  echo "==== case: move single directory with trailing '/' ===="
  run ${APP_BIN} -c "mv ${KV_BACKEND}/dest/dev/ ${KV_BACKEND}/dest/dev.copy"
  assert_success

  echo "ensure the directory got moved to destination"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev.copy/1"
  assert_success
  assert_output "1"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev.copy/2"
  assert_success
  assert_output "2"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev.copy/3"
  assert_success
  assert_output "3"

  echo "ensure the src directory got removed"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev/1"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev/2"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev/3"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  #######################################
  echo "==== TODO case: move ambigious file ===="

  #######################################
  echo "==== TODO case: move ambigious directory ===="
}
