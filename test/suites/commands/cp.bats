load ../../util/util
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'cp'" {
  #######################################
  echo "==== case: copy single file ===="
  run ${APP_BIN} -c "cp ${KV_BACKEND}/src/prod/all ${KV_BACKEND}/dest/prod/all"
  assert_success

  echo "ensure the file got copied to destination"
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
  echo "==== case: copy non-existing file ===="
  run ${APP_BIN} -c "cp ${KV_BACKEND}/src/does/not/exist ${KV_BACKEND}/dest/a"
  assert_failure 1

  echo "ensure proper error message"
  assert_line --partial "Not a valid path for operation: /${KV_BACKEND}/src/does/not/exist"

  #######################################
  echo "==== case: copy single directory without trailing '/' ===="
  run ${APP_BIN} -c "cp '${KV_BACKEND}/src/dev' '${KV_BACKEND}/dest/dev'"
  assert_success

  echo "ensure the directory got copied to destination"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev/1"
  assert_success
  assert_output "1"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev/2"
  assert_success
  assert_output "2"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev/3"
  assert_success
  assert_output "3"

  echo "ensure the src directory still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/dev/1"
  assert_success
  assert_output "1"
  run get_vault_value "value" "${KV_BACKEND}/src/dev/2"
  assert_success
  assert_output "2"
  run get_vault_value "value" "${KV_BACKEND}/src/dev/3"
  assert_success
  assert_output "3"

  #######################################
  echo "==== case: copy single directory with trailing '/' ===="
  run ${APP_BIN} -c "cp \"${KV_BACKEND}/src/dev/\" \"${KV_BACKEND}/dest/dev.copy\""
  assert_success

  echo "ensure the directory got copied to destination"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev.copy/1"
  assert_success
  assert_output "1"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev.copy/2"
  assert_success
  assert_output "2"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev.copy/3"
  assert_success
  assert_output "3"

  echo "ensure the src directory still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/dev/1"
  assert_success
  assert_output "1"
  run get_vault_value "value" "${KV_BACKEND}/src/dev/2"
  assert_success
  assert_output "2"
  run get_vault_value "value" "${KV_BACKEND}/src/dev/3"
  assert_success
  assert_output "3"

  ######################################
  echo "==== case: copy single directory with dest trailing '/' ===="
  run ${APP_BIN} -c "cp ${KV_BACKEND}/src/dev ${KV_BACKEND}/dest/dev.copy2/"
  assert_success

  echo "ensure the directory got copied to destination"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev.copy2/1"
  assert_success
  assert_output "1"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev.copy2/2"
  assert_success
  assert_output "2"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev.copy2/3"
  assert_success
  assert_output "3"

  echo "ensure the src directory still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/dev/1"
  assert_success
  assert_output "1"
  run get_vault_value "value" "${KV_BACKEND}/src/dev/2"
  assert_success
  assert_output "2"
  run get_vault_value "value" "${KV_BACKEND}/src/dev/3"
  assert_success
  assert_output "3"

  ######################################
  echo "==== case: copy single directory with src and dest trailing '/' ===="
  run ${APP_BIN} -c "cp ${KV_BACKEND}/src/dev/ ${KV_BACKEND}/dest/dev.copy3/"
  assert_success

  echo "ensure the directory got copied to destination"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev.copy3/1"
  assert_success
  assert_output "1"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev.copy3/2"
  assert_success
  assert_output "2"
  run get_vault_value "value" "${KV_BACKEND}/dest/dev.copy3/3"
  assert_success
  assert_output "3"

  echo "ensure the src directory still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/dev/1"
  assert_success
  assert_output "1"
  run get_vault_value "value" "${KV_BACKEND}/src/dev/2"
  assert_success
  assert_output "2"
  run get_vault_value "value" "${KV_BACKEND}/src/dev/3"
  assert_success
  assert_output "3"

  #######################################
  echo "==== case: copy ambiguous directory ===="
  run ${APP_BIN} -c "cp ${KV_BACKEND}/src/staging/all/ ${KV_BACKEND}/dest/staging/all/"
  assert_success

  echo "ensure the directory got copied"
  run get_vault_value "value" "${KV_BACKEND}/dest/staging/all/v1"
  assert_success
  assert_output "v1"
  run get_vault_value "value" "${KV_BACKEND}/dest/staging/all/v2"
  assert_success
  assert_output "v2"

  echo "ensure the source directory still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/staging/all/v1"
  assert_success
  assert_output "v1"
  run get_vault_value "value" "${KV_BACKEND}/src/staging/all/v2"
  assert_success
  assert_output "v2"

  echo "ensure the ambiguous file still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/staging/all"
  assert_success
  assert_output "all"

  #######################################
  echo "==== case: copy ambiguous file ===="
  run ${APP_BIN} -c "cp ${KV_BACKEND}/src/tooling ${KV_BACKEND}/dest/tooling"
  assert_success

  echo "ensure the file got copied"
  run get_vault_value "value" "${KV_BACKEND}/dest/tooling"
  assert_success
  assert_output "tooling"

  echo "ensure the source file still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/tooling"
  assert_success
  assert_output "tooling"

  echo "ensure the ambiguous directory still exists"
  run get_vault_value "value" "${KV_BACKEND}/src/tooling/v1"
  assert_success
  assert_output "v1"

  run get_vault_value "value" "${KV_BACKEND}/src/tooling/v2"
  assert_success
  assert_output "v2"

  #######################################
  echo "==== case: copy ambiguous directory parent ===="
  run ${APP_BIN} -c "cp ${KV_BACKEND}/src/staging/ ${KV_BACKEND}/dest/staging5/"
  assert_success

  echo "ensure ambiguous file got copied"
  run get_vault_value "value" "${KV_BACKEND}/dest/staging5/all"
  assert_success
  assert_output "all"
  run get_vault_value "value" "${KV_BACKEND}/dest/staging/all/v2"
  assert_success
  assert_output "v2"
}
