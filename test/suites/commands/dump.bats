load ../../util/util
load ../../bin/plugins/bats-support/load
load ../../bin/plugins/bats-assert/load

@test "vault-${VAULT_VERSION} ${KV_BACKEND} 'dump'" {
  #######################################
  echo "==== case: dump directory ===="
  run ${APP_BIN} -c "dump ${KV_BACKEND}/src/"
  assert_success
  export DUMP_DIR=$(ls -l | grep vsh-dump | awk '{print $9}')
  sed -i "s/${KV_BACKEND}\/src/${KV_BACKEND}\/restore/g" ${DUMP_DIR}/restore.sh
  docker cp ${DUMP_DIR} vsh-integration-test-vault:/dump-data
  rm -rf ${DUMP_DIR}

  echo "==== case: restore directory dump ===="
  run docker exec vsh-integration-test-vault sh -c "cd /dump-data && ./restore.sh"
  assert_success
  run get_vault_value "value" "${KV_BACKEND}/restore/a/foo"
  assert_success
  assert_output "1"


  #######################################
  echo "==== case: dump file ===="
  run ${APP_BIN} -c "dump ${KV_BACKEND}/src/a/foo"
  assert_success
  export DUMP_DIR=$(ls -l | grep vsh-dump | awk '{print $9}')
  sed -i "s/${KV_BACKEND}\/src/${KV_BACKEND}\/restore2/g" ${DUMP_DIR}/restore.sh
  docker cp ${DUMP_DIR} vsh-integration-test-vault:/dump-data2
  rm -rf ${DUMP_DIR}

  echo "==== case: restore file dump ===="
  run docker exec vsh-integration-test-vault sh -c "cd /dump-data2 && ./restore.sh"
  assert_success
  run get_vault_value "value" "${KV_BACKEND}/restore2/a/foo"
  assert_success
  assert_output "1"

  #######################################
  echo "==== case: dump non-existing path ===="
  run ${APP_BIN} -c "dump ${KV_BACKEND}/notexist"
  assert_failure 1
}
