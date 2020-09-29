load ../../util/util
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'rm'" {
  #######################################
  echo "==== case: remove single file ===="
  run get_vault_value "value" "${KV_BACKEND}/src/prod/all"
  assert_success
  assert_output "all"

  run ${APP_BIN} -c "rm ${KV_BACKEND}/src/prod/all"
  assert_success

  echo "ensure the file got removed"
  run get_vault_value "value" "${KV_BACKEND}/src/prod/all"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  #######################################
  echo "==== case: remove non-existing file ===="
  run ${APP_BIN} -c "rm ${KV_BACKEND}/src/does/not/exist"
  assert_success

  echo "ensure proper error message"
  assert_line --partial "Not a valid path for operation: /${KV_BACKEND}/src/does/not/exist"

  #######################################
  echo "==== case: remove single directory ===="
  run get_vault_value "value" "${KV_BACKEND}/src/dev/1"
  assert_success
  assert_output "1"

  run get_vault_value "value" "${KV_BACKEND}/src/dev/2"
  assert_success
  assert_output "2"

  run get_vault_value "value" "${KV_BACKEND}/src/dev/3"
  assert_success
  assert_output "3"

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
  echo "==== case: remove ambiguous directory ===="
  run get_vault_value "value" "${KV_BACKEND}/src/staging/all/v1"
  assert_success
  assert_output "v1"

  run get_vault_value "value" "${KV_BACKEND}/src/staging/all/v2"
  assert_success
  assert_output "v2"

  run ${APP_BIN} -c "rm ${KV_BACKEND}/src/staging/all/"
  assert_success

  echo "ensure the directory got removed"
  run get_vault_value "value" "${KV_BACKEND}/src/staging/all/v1"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  run get_vault_value "value" "${KV_BACKEND}/src/staging/all/v2"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  echo "ensure the ambiguous file still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/staging/all"
  assert_success
  assert_output "all"

  #######################################
  echo "==== case: remove ambiguous file ===="
  run get_vault_value "value" "${KV_BACKEND}/src/tooling"
  assert_success
  assert_output "tooling"

  run ${APP_BIN} -c "rm ${KV_BACKEND}/src/tooling"
  assert_success

  echo "ensure the file got removed"
  run get_vault_value "value" "${KV_BACKEND}/src/tooling"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  echo "ensure the ambiguous directory still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/tooling/v1"
  assert_success
  assert_output "v1"

  run get_vault_value "value" "${KV_BACKEND}/src/tooling/v2"
  assert_success
  assert_output "v2"

  #######################################
  echo "==== case: remove ambiguous directory ===="
  run get_vault_value "value" "${KV_BACKEND}/src/ambivalence/1"
  assert_success
  assert_output "1"

  run ${APP_BIN} -c "rm ${KV_BACKEND}/src/ambivalence/1/"
  assert_success

  echo "ensure the directory got removed"
  run get_vault_value "value" "${KV_BACKEND}/src/ambivalence/1/a"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  echo "ensure the ambiguous file still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/ambivalence/1"
  assert_success
  assert_output "1"

  #######################################
  echo "==== case: remove ambiguous file without read permissions ===="
  run get_vault_value "value" "${KV_BACKEND}/src/a/foo"
  assert_success
  assert_output "1"
  run get_vault_value "value" "${KV_BACKEND}/src/a/foo/bar"
  assert_success
  assert_output "2"

  run bash -c "VAULT_TOKEN=delete-only ${APP_BIN} -c 'rm ${KV_BACKEND}/src/a/foo'"
  assert_success

  echo "ensure file deletion"
  run get_vault_value "value" "${KV_BACKEND}/src/a/foo"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  echo "ensure that the ambiguous directory still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/a/foo/bar"
  assert_success
  assert_output "2"
}
