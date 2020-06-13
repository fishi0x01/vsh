load ../util/util
load ../bin/plugins/bats-support/load
load ../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'rm'" {
  #######################################
  echo "==== case: remove single file ===="
  run get_vault_value "value" "${KV_BACKEND}/src/prod/all"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  run ${APP_BIN} -c "rm ${KV_BACKEND}/src/prod/all"
  assert_success

  echo "ensure the file got removed"
  run get_vault_value "value" "${KV_BACKEND}/src/prod/all"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  #######################################
  echo "==== case: remove single directory ===="
  run get_vault_value "value" "${KV_BACKEND}/src/dev/1"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  run get_vault_value "value" "${KV_BACKEND}/src/dev/2"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  run get_vault_value "value" "${KV_BACKEND}/src/dev/3"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  run ${APP_BIN} -c "rm ${KV_BACKEND}/src/dev"
  assert_success

  echo "ensure the directory got removed"
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
  echo "==== case: remove ambigious directory ===="
  run get_vault_value "value" "${KV_BACKEND}/src/staging/all/v1"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  run get_vault_value "value" "${KV_BACKEND}/src/staging/all/v2"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  run ${APP_BIN} -c "rm ${KV_BACKEND}/src/staging/all/"
  assert_success

  echo "ensure the directory got removed"
  run get_vault_value "value" "${KV_BACKEND}/src/staging/all/v1"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  run get_vault_value "value" "${KV_BACKEND}/src/staging/all/v2"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  echo "ensure the ambigious file still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/staging/all"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  #######################################
  echo "==== case: remove ambigious file ===="
  run get_vault_value "value" "${KV_BACKEND}/src/tooling"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  run ${APP_BIN} -c "rm ${KV_BACKEND}/src/tooling"
  assert_success

  echo "ensure the file got removed"
  run get_vault_value "value" "${KV_BACKEND}/src/tooling"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  echo "ensure the ambigious directory still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/tooling/v1"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"

  run get_vault_value "value" "${KV_BACKEND}/src/tooling/v2"
  assert_success
  assert_output "${VAULT_TEST_VALUE}"  
}
