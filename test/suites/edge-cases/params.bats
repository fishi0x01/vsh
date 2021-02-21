load ../../util/common
load ../../util/standard-setup
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} whitespaces between parameters" {
  #######################################
  echo "==== case: copy with multiple whitespaces ===="
  run ${APP_BIN} -v DEBUG -c "cp    /KV2/src/prod/all      /KV2/dest/prod/all"
  assert_success
  assert_output --partial "Copied /KV2/src/prod/all to /KV2/dest/prod/all"

  echo "==== case: copy with tabs ===="
  run ${APP_BIN} -v DEBUG -c "cp     /KV2/src/prod/all      /KV2/dest/prod/all       "
  assert_success
  assert_output --partial "Copied /KV2/src/prod/all to /KV2/dest/prod/all"

  echo "==== case: append with multiple whitespaces ===="
  run ${APP_BIN} -v DEBUG -c "append    --rename    /KV2/src/prod/all      /KV2/dest/prod/all"
  assert_success
  assert_output --partial "Appended values from /KV2/src/prod/all to /KV2/dest/prod/all"
}
