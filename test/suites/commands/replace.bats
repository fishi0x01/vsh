load ../../util/util
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'replace'" {
  #######################################
  echo "==== case: replace nonexistant string ===="
  run ${APP_BIN} -c "replace 'foobarbaz' 'pie' ${KV_BACKEND}/src/dev/1 -y"
  assert_success
  assert_line "No matches found to replace."

  #######################################
  echo "==== case: replace in dry-run ===="
  run ${APP_BIN} -c "replace 'fruit' 'pie' ${KV_BACKEND}/src/dev/1 -y -n"
  assert_success
  assert_line "Skipping write."
  run get_vault_value "fruit" "${KV_BACKEND}/src/dev/1"
  assert_line apple

  #######################################
  echo "==== case: replace value with empty string ===="
  run ${APP_BIN} -c "replace 'beer' '' ${KV_BACKEND}/src/tooling -y"
  assert_success
  assert_line "Writing!"
  run get_vault_value "drink" "${KV_BACKEND}/src/tooling"
  refute_line beer

  #######################################
  echo "==== case: replace key in single path without scope ===="
  run ${APP_BIN} -c "replace 'fruit' 'pie' ${KV_BACKEND}/src/dev/1 -y"
  assert_success
  assert_line "Writing!"
  run get_vault_value "pie" "${KV_BACKEND}/src/dev/1"
  assert_line apple

  #######################################
  echo "==== case: replace value in single path without scope ===="
  run ${APP_BIN} -c "replace 'banana' 'something' ${KV_BACKEND}/src/dev/2 -y"
  assert_success
  assert_line "Writing!"
  run get_vault_value "fruit" "${KV_BACKEND}/src/dev/2"
  assert_line something

  #######################################
  echo "==== case: replace key in single path with scope ===="
  run ${APP_BIN} -c "replace 'tree' 'flora' ${KV_BACKEND}/src/staging/all -k -y"
  assert_success
  assert_line "Writing!"
  run get_vault_value "flora" "${KV_BACKEND}/src/staging/all"
  assert_line palm

  #######################################
  echo "==== case: replace value in single path with scope ===="
  run ${APP_BIN} -c "replace 'test' 'exhibit' ${KV_BACKEND}/src/prod/all -v -y"
  assert_success
  assert_line "Writing!"
  run get_vault_value "example" "${KV_BACKEND}/src/prod/all"
  assert_line exhibit
}

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'replace' regexp" {
  #######################################
  echo "==== case: replace nonexistant string ===="
  run ${APP_BIN} -c "replace '^ruit' 'pie' ${KV_BACKEND}/src/dev/1 -y -e"
  assert_success
  assert_line "No matches found to replace."

  #######################################
  echo "==== case: replace key in single path without scope ===="
  run ${APP_BIN} -c "replace '^fru.*' 'pie' ${KV_BACKEND}/src/dev/1 -y -e"
  assert_success
  assert_line "Writing!"
  run get_vault_value "pie" "${KV_BACKEND}/src/dev/1"
  assert_line apple

  #######################################
  echo "==== case: replace value in single path without scope ===="
  run ${APP_BIN} -c "replace '[ba]+nana' 'something' ${KV_BACKEND}/src/dev/2 -y -e"
  assert_success
  assert_line "Writing!"
  run get_vault_value "fruit" "${KV_BACKEND}/src/dev/2"
  assert_line something

  #######################################
  echo "==== case: replace key in single path with scope ===="
  run ${APP_BIN} -c "replace 'tre{2}' 'flora' ${KV_BACKEND}/src/staging/all -k -y -e"
  assert_success
  assert_line "Writing!"
  run get_vault_value "flora" "${KV_BACKEND}/src/staging/all"
  assert_line palm

  #######################################
  echo "==== case: replace value in single path with scope ===="
  run ${APP_BIN} -c "replace '(test)' '\${1}exhibit' ${KV_BACKEND}/src/prod/all -v -y -e"
  assert_success
  assert_line "Writing!"
  run get_vault_value "example" "${KV_BACKEND}/src/prod/all"
  assert_line testexhibit
}
