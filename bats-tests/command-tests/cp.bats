load ../util/util
load ../bin/plugins/bats-support/load
load ../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'cp'" {
  #######################################
  echo "==== case: copy single file ===="
  run ${APP_BIN} -c "cp ${KV_BACKEND}/src/prod/all ${KV_BACKEND}/dest/prod/all"
  assert_success

  echo "ensure the file got copied to destination"
  run get_vault_value "value" "${KV_BACKEND}/dest/prod/all"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  echo "ensure the src file still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/prod/all"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  #######################################
  echo "==== case: copy single directory ===="
  run ${APP_BIN} -c "cp ${KV_BACKEND}/src/dev ${KV_BACKEND}/dest/dev"
  assert_success

  echo "ensure the directory got copied to destination"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev/1"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  run get_vault_value "value" "${KV_BACKEND}/dest/dev/2"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  run get_vault_value "value" "${KV_BACKEND}/dest/dev/3"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  echo "ensure the src directory still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/dev/1"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  run get_vault_value "value" "${KV_BACKEND}/src/dev/2"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  run get_vault_value "value" "${KV_BACKEND}/src/dev/3"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  #######################################
  echo "==== TODO case: copy ambigious file ===="

  #######################################
  echo "==== TODO case: copy ambigious directory ===="
}
