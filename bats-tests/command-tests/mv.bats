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
  assert_output "${VAULT_TEST_VALUE}"

  echo "ensure the src file got removed"
  run get_vault_value "value" "${KV_BACKEND}/src/prod/all"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  #######################################
  echo "==== case: move single directory ===="
  run ${APP_BIN} -c "mv ${KV_BACKEND}/src/dev ${KV_BACKEND}/dest/dev"
  assert_success

  echo "ensure the directory got moved to destination"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev/1"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  run get_vault_value "value" "${KV_BACKEND}/dest/dev/2"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  run get_vault_value "value" "${KV_BACKEND}/dest/dev/3"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

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
  echo "==== TODO case: move ambigious file ===="

  #######################################
  echo "==== TODO case: move ambigious directory ===="
}
