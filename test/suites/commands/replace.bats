load ../../util/util
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'replace'" {
  #######################################
  echo "==== case: replace content in single file ===="
  run ${APP_BIN} -c "replace --value 'apple' 'something' ${KV_BACKEND}/src/dev/1"
  assert_success
  
  echo "ensure value(s) got replaced"
  run get_vault_value "fruit" "${KV_BACKEND}/src/dev/1"
  assert_line "fruit = something"
}
