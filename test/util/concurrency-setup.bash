data_for_path() {
    result=""
    for i in $(seq 1 400); do
        result="${result} vault kv put ${1}/${i} value=1;";
    done
    echo $result
}

setup() {
    rm -f vsh_trace.log
    docker run -d \
        --name=${VAULT_CONTAINER_NAME} \
        -p "${VAULT_HOST_PORT}:8200" \
        --cap-add=IPC_LOCK \
        -e "VAULT_ADDR=http://127.0.0.1:8200" \
        -e "VAULT_TOKEN=root" \
        -e "VAULT_DEV_ROOT_TOKEN_ID=root" \
        -e "VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200" \
        "hashicorp/vault:${VAULT_VERSION}" &> /dev/null
    docker cp "$DIR/policy-no-root.hcl" ${VAULT_CONTAINER_NAME}:.
    docker cp "$DIR/policy-delete-only.hcl" ${VAULT_CONTAINER_NAME}:.
    # need some time for GH Actions CI
    sleep 3
    vault_exec "vault secrets disable secret;
                vault policy write no-root policy-no-root.hcl;
                vault token create -id=no-root -policy=no-root;
                vault policy write delete-only policy-delete-only.hcl;
                vault token create -id=delete-only -policy=delete-only;
                vault secrets enable -version=1 -path=KV1 kv;
                vault secrets enable -version=2 -path=KV2 kv"

    dirs="src src/a src/b src/1 src/a/a src/a/a/a src/b/a src/b/b src/b/b/a"
    for dir in $dirs; do
        vault_exec "$(data_for_path "/KV2/${dir}")" &
    done
    wait
}
