load ../../util/common
load ../../util/standard-setup
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
  run get_vault_value "key" "${KV_BACKEND}/src/tooling"
  assert_line A

  #######################################
  echo "==== case: replace key in single path without scope ===="
  run ${APP_BIN} -c "replace 'fruit' 'pie' ${KV_BACKEND}/src/dev/1 -y"
  assert_success
  assert_line "Writing!"
  run get_vault_value "pie" "${KV_BACKEND}/src/dev/1"
  assert_line apple
  run get_vault_value "value" "${KV_BACKEND}/src/dev/1"
  assert_line 1

  #######################################
  echo "==== case: replace value in single path without scope ===="
  run ${APP_BIN} -c "replace 'banana' 'something' ${KV_BACKEND}/src/dev/2 -y"
  assert_success
  assert_line "Writing!"
  run get_vault_value "fruit" "${KV_BACKEND}/src/dev/2"
  assert_line something
  run get_vault_value "value" "${KV_BACKEND}/src/dev/2"
  assert_line 2

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
  run get_vault_value "value" "${KV_BACKEND}/src/prod/all"
  assert_line all

  #######################################
  echo "==== case: replace with invalid output format ===="
  run ${APP_BIN} -c "replace -s 'produce' 'apple' 'orange' ${KV_BACKEND}/src/selector/1 -o invalid"
  assert_failure
  assert_line --partial "invalid output format: invalid"

  #######################################
  echo "==== case: replace with diff output format ===="
  run ${APP_BIN} -c "replace -s 'produce' 'apple' 'orange' ${KV_BACKEND}/src/selector/1 -n -o diff"
  assert_success
  assert_line "- /${KV_BACKEND}/src/selector/1> produce = apple"
  assert_line "+ /${KV_BACKEND}/src/selector/1> produce = orange"

  #######################################
  echo "==== case: replace value in single path with selector ===="
  run ${APP_BIN} -c "replace -s 'produce' 'apple' 'orange' ${KV_BACKEND}/src/selector/1 -y"
  assert_success
  assert_line "Writing!"
  run get_vault_value "produce" "${KV_BACKEND}/src/selector/1"
  assert_line orange
  run get_vault_value "fruit" "${KV_BACKEND}/src/selector/1"
  assert_line apple
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
  run get_vault_value "value" "${KV_BACKEND}/src/dev/1"
  assert_line 1

  #######################################
  echo "==== case: replace value in single path without scope ===="
  run ${APP_BIN} -c "replace '[ba]+nana' 'something' ${KV_BACKEND}/src/dev/2 -y -e"
  assert_success
  assert_line "Writing!"
  run get_vault_value "fruit" "${KV_BACKEND}/src/dev/2"
  assert_line something
  run get_vault_value "value" "${KV_BACKEND}/src/dev/2"
  assert_line 2

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
  run get_vault_value "value" "${KV_BACKEND}/src/prod/all"
  assert_line all

  #######################################
  echo "==== case: replace value in single path with selector ===="
  run ${APP_BIN} -c "replace -e -s 'prod.*' '^apple' 'orange' ${KV_BACKEND}/src/selector/1 -y"
  assert_success
  assert_line "Writing!"
  run get_vault_value "produce" "${KV_BACKEND}/src/selector/1"
  assert_line orange
  run get_vault_value "fruit" "${KV_BACKEND}/src/selector/1"
  assert_line apple

  #######################################
  echo "==== case: replace fails with bad regex selector ===="
  run ${APP_BIN} -c "replace -e -s '][' '^apple' 'orange' ${KV_BACKEND}/src/selector/1 -y"
  assert_failure
  assert_line  --partial "key-selector: error parsing regexp"

  #######################################
  echo "==== case: shallow replace"
  run ${APP_BIN} -c "replace -e -S -k 'tree' '.*' 'maple' /${KV_BACKEND}/src/staging/all"
  assert_line --partial "/${KV_BACKEND}/src/staging/all"
  refute_line --partial "/${KV_BACKEND}/src/staging/all/v1"
  refute_line --partial "/${KV_BACKEND}/src/staging/all/v2"

  #######################################
  echo "==== case: replace in pwd ===="
  export VAULT_PATH=${KV_BACKEND}/src/a
  run ${APP_BIN} -c "replace 'long-value' 'interesting-value' -y"
  assert_success
  assert_line "Writing!"
  unset VAULT_PATH
  run get_vault_value "long" "${KV_BACKEND}/src/a/foo"
  assert_line this-is-a-really-interesting-value-for-testing
}
