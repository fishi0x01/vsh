#!/bin/bash
set -e
set -o pipefail
set -o nounset

UNAME=$(uname | tr '[:upper:]' '[:lower:]')
export APP_BIN="./build/vsh_${UNAME}_amd64"

# Input: <vault-version> <container-name> <vault-host-port>
start_vault() {
  docker run --name="${2}" -d -p "${3}:8200" --cap-add=IPC_LOCK -e "VAULT_ADDR=http://127.0.0.1:8200" -e "VAULT_TOKEN=root" -e "VAULT_DEV_ROOT_TOKEN_ID=root" -e "VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200" vault:${1} &> /dev/null
  # need at least 2s for GH Actions CI
  sleep 2
}

# Input: <container-name>
stop_vault() {
  docker rm -f ${1} &> /dev/null
}

# Input: <container-name> <command>
vault_exec() {
  docker exec ${1} ${2} &> /dev/null
}

# Input: <container-name> <path> <key> <value>
vault_field_must_be() {
  vault_val=$(docker exec "${1}" vault kv get -field="${2}" "${3}")
  if [ "$vault_val" = "$4" ]; then
    return 0
  else
    echo "Error: $vault_val (got) != $4 (expected)"
    return 1
  fi
}

# Input: <container-name> <path> <value>
vault_value_must_be() {
  vault_field_must_be "${1}" "value" "${2}" "${3}"
}

# Input: <given> <wanted>
value_must_be() {
  if [ "$1" = "$2" ]; then
    return 0
  else
    echo "Error: $1 != $2"
    return 1
  fi
}
