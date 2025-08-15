load ../../util/common
load ../../util/standard-setup
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load
load ../../bin/plugins/bats-file/load

@test "vault-${VAULT_VERSION} - logging" {
  #######################################
  echo "==== case: default verbosity ===="
  run bash -c "VAULT_TOKEN=delete-only ${APP_BIN} -c 'ls /KV2/src/a'"
  assert_success
  assert_line "foo"
  assert_line "foo/"
  assert_file_not_exist vsh_trace.log

  echo "==== case: print debug info for user ===="
  run bash -c "VAULT_TOKEN=delete-only ${APP_BIN} -v DEBUG -c 'ls /KV2/src/a'"
  assert_success
  assert_line "Cannot auto-discover mount backends: Token does not have list permission on sys/mounts"
  assert_line "foo"
  assert_line "foo/"
  assert_file_exist vsh_trace.log

  echo "==== case: invalid verbosity level ===="
  run bash -c "VAULT_TOKEN=delete-only ${APP_BIN} -v NOTEXIST -c 'ls /KV2/src/a'"
  assert_failure 2
  assert_line --partial "Not a valid verbosity level"

  echo "==== case: login with false token ===="
  run bash -c "VAULT_TOKEN=false ${APP_BIN} -v DEBUG -c 'ls /KV2/src/a'"
  assert_failure 1
  run cat vsh_trace.log
  assert_line --partial "* permission denied"
}
