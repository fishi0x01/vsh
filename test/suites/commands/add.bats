load ../../util/common
load ../../util/standard-setup
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'add'" {
  #######################################
  echo "==== case: add value to non existing path ===="
  run ${APP_BIN} -c "add --confirm test value ${KV_BACKEND}/fake/path"
  assert_failure

  #######################################
  echo "==== case: add key to existing path ===="
  run ${APP_BIN} -c "add --confirm test value ${KV_BACKEND}/src/a/foo"
  assert_success

  echo "ensure the key was written to destination"
  run get_vault_value "test" "${KV_BACKEND}/src/a/foo"
  assert_success
  assert_output "value"

  #######################################
  echo "==== case: add existing key to existing path ===="
  run ${APP_BIN} -c "add --confirm value another ${KV_BACKEND}/src/a/foo"
  assert_failure
  assert_output --partial "Key already exists at path: ${KV_BACKEND}/src/a/foo"

  #######################################
  echo "==== case: overwrite existing key to existing path ===="
  run ${APP_BIN} -c "add --confirm -f value another ${KV_BACKEND}/src/a/foo"
  assert_success

  echo "ensure the key was written to destination"
  run get_vault_value "value" "${KV_BACKEND}/src/a/foo"
  assert_success
  assert_output "another"

  #######################################
  echo "==== case: add dryrun, ensure key not added ===="
  run ${APP_BIN} -c "add --dry-run -f proposedvalue willnotexist ${KV_BACKEND}/src/a/foo"
  assert_success
  assert_line "Skipping write."

  echo "ensure the key was NOT written to destination"
  run get_vault_value "proposedvalue" "${KV_BACKEND}/src/a/foo"
  assert_output --partial "not present in secret"

}
