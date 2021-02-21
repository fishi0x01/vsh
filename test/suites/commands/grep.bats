load ../../util/common
load ../../util/standard-setup
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'grep'" {
  #######################################
  echo "==== case: grep term on '${KV_BACKEND}/' ===="
  run ${APP_BIN} -c "grep value ${KV_BACKEND}/"
  assert_success
  assert_line --partial "/${KV_BACKEND}/src/dev/1"
  assert_line --partial "/${KV_BACKEND}/src/dev/2"
  assert_line --partial "/${KV_BACKEND}/src/dev/3"
  assert_line --partial "/${KV_BACKEND}/src/prod/all"

  #######################################
  echo "==== case: grep non-existing file ===="
  run ${APP_BIN} -c "grep test ${KV_BACKEND}/src/does/not/exist"
  assert_failure 1

  echo "ensure proper error message"
  assert_line --partial "Not a valid path for operation: /${KV_BACKEND}/src/does/not/exist"

  #######################################
  echo "==== case: grep term on ambigious directory ===="
  run ${APP_BIN} -c "grep juice ${KV_BACKEND}/src/tooling/"
  assert_line --partial "/${KV_BACKEND}/src/tooling/v1"

  #######################################
  echo "==== case: grep term on ambigious file ===="
  run ${APP_BIN} -c "grep beer ${KV_BACKEND}/src/tooling"
  assert_line --partial "/${KV_BACKEND}/src/tooling"

  #######################################
  echo "==== case: grep value with quotes ===="
  run ${APP_BIN} -c "grep \\\"quoted\\\" ${KV_BACKEND}/src/quoted/foo"
  assert_line --partial "/${KV_BACKEND}/src/quoted/foo"

  #######################################
  echo "==== case: regexp pattern ===="
  run ${APP_BIN} -c "grep app.* ${KV_BACKEND}/src -e"
  assert_line --partial "/${KV_BACKEND}/src/dev/1"
  assert_line --partial "/${KV_BACKEND}/src/ambivalence/1"

  #######################################
  echo "==== case: fails on invalid regex pattern ===="
  run ${APP_BIN} -c "grep '][' ${KV_BACKEND}/src/dev -e"
  assert_line --partial "cannot parse regex"
  assert_failure 1

  #######################################
  echo "==== case: regex pattern on a long value ===="
  run ${APP_BIN} -c "grep -e 'value-for-testing' ${KV_BACKEND}/src/a/foo"
  assert_line --partial this-is-a-really-long-value-for-testing
  assert_success

  #######################################
  echo "==== case: pattern with spaces ===="
  run ${APP_BIN} -c "grep 'a spaced val' ${KV_BACKEND}/src/spaces"
  assert_line --partial "/${KV_BACKEND}/src/spaces/foo"

  #######################################
  echo "==== case: pattern with escaped spaces ===="
  run ${APP_BIN} -c "grep a\ spaced\ val ${KV_BACKEND}/src/spaces"
  assert_line --partial "/${KV_BACKEND}/src/spaces/foo"

  #######################################
  echo "==== case: pattern with apostrophe ===="
  run ${APP_BIN} -c "grep \"steve's\" ${KV_BACKEND}/src/apostrophe"
  assert_line --partial "/${KV_BACKEND}/src/apostrophe"

  #######################################
  echo "==== case: no match when only searching keys ===="
  run ${APP_BIN} -c "grep 'apple' ${KV_BACKEND}/src/dev/1 -k"
  refute_line --partial "/${KV_BACKEND}/src/dev/1"

  #######################################
  echo "==== case: no match when only searching values ===="
  run ${APP_BIN} -c "grep 'fruit' ${KV_BACKEND}/src/dev/1 -v"
  refute_line --partial "/${KV_BACKEND}/src/dev/1"

  #######################################
  echo "==== case: match when only searching keys ===="
  run ${APP_BIN} -c "grep 'fruit' ${KV_BACKEND}/src/dev -k"
  assert_line --partial "/${KV_BACKEND}/src/dev/1"
  assert_line --partial "/${KV_BACKEND}/src/dev/2"
  assert_line --partial "/${KV_BACKEND}/src/dev/3"

  #######################################
  echo "==== case: match when only searching values ===="
  run ${APP_BIN} -c "grep 'apple' ${KV_BACKEND}/src/dev -v"
  assert_line --partial "/${KV_BACKEND}/src/dev/1"

  #######################################
  echo "==== case: fails on invalid flag ===="
  run ${APP_BIN} -c "grep 'apple' ${KV_BACKEND}/src/dev --foo"
  assert_line --partial "unknown argument --foo"
  assert_failure 1

  #######################################
  echo "==== TODO case: grep term on directory with reduced permissions ===="

  #######################################
  echo "==== TODO case: grep term on '/' ===="
}
