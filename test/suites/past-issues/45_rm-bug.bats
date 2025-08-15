load ../../util/common
load ../../util/standard-setup
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} - test issue https://github.com/fishi0x01/vsh/issues/45" {
  #######################################
  echo "==== case: https://github.com/fishi0x01/vsh/issues/45#issuecomment-700178099 ===="
  run bash -c "VAULT_TOKEN=delete-only ${APP_BIN} -c 'ls /KV2/src/a'"
  assert_success
  assert_line "foo"
  assert_line "foo/"

  run bash -c "VAULT_TOKEN=delete-only ${APP_BIN} -c 'ls /KV2/src/a/foo'"
  assert_failure 1
  assert_line --partial "not a valid path for operation: /KV2/src/a/foo"

  run bash -c "VAULT_TOKEN=delete-only ${APP_BIN} -c 'ls /KV2/src/a/foo/'"
  assert_success
  assert_line "bar"

  run bash -c "VAULT_TOKEN=delete-only ${APP_BIN} -c 'rm /KV2/src/a/foo'"
  assert_success
  run get_vault_value "value" "/KV2/src/a/foo/bar"
  assert_success
  assert_output "2"
  run get_vault_value "value" "/KV2/src/a/foo"
  assert_success
  assert_output --partial "${NO_VALUE_FOUND}"

  run bash -c "VAULT_TOKEN=delete-only ${APP_BIN} -c 'ls /KV2/src/a/foo'"
  assert_success
  assert_line "bar"

  run bash -c "VAULT_TOKEN=delete-only ${APP_BIN} -c 'ls /KV2/src/a/'"
  assert_success
  assert_output "foo/"
}
