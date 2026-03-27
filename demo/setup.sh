#!/usr/bin/env bash
# Seeds a local Vault instance with demo data for the vsh recording.
# Usage: bash demo/setup.sh [--teardown]
#
# Requires: docker, vault CLI (for VAULT_ADDR/VAULT_TOKEN env vars)

set -euo pipefail

VAULT_CONTAINER="vsh-demo-vault"
VAULT_VERSION="${VAULT_VERSION:-1.16.2}"
VAULT_PORT="${VAULT_PORT:-8200}"

if [[ "${1:-}" == "--teardown" ]]; then
    echo "Removing demo vault container..."
    docker rm -f "$VAULT_CONTAINER" 2>/dev/null || true
    exit 0
fi

# Remove any leftover container from a previous run
docker rm -f "$VAULT_CONTAINER" 2>/dev/null || true

echo "Starting Vault ${VAULT_VERSION} on port ${VAULT_PORT}..."
docker run -d \
    --name="$VAULT_CONTAINER" \
    -p "${VAULT_PORT}:8200" \
    --cap-add=IPC_LOCK \
    -e VAULT_DEV_ROOT_TOKEN_ID=root \
    -e VAULT_DEV_LISTEN_ADDRESS="0.0.0.0:8200" \
    "hashicorp/vault:${VAULT_VERSION}" &>/dev/null

echo "Waiting for Vault to be ready..."
for i in $(seq 1 30); do
    docker exec -e VAULT_ADDR=http://127.0.0.1:8200 "$VAULT_CONTAINER" vault status &>/dev/null && break
    sleep 1
done
docker exec -e VAULT_ADDR=http://127.0.0.1:8200 "$VAULT_CONTAINER" vault status &>/dev/null \
    || { echo "Vault did not become ready in time"; exit 1; }

exec_vault() {
    docker exec -e VAULT_ADDR=http://127.0.0.1:8200 -e VAULT_TOKEN=root "$VAULT_CONTAINER" /bin/sh -c "$1" &>/dev/null
}

echo "Seeding demo secrets..."

# prod secrets
exec_vault "vault kv put secret/prod/database \
    host=db.prod.example.com \
    port=5432 \
    username=admin \
    password=s3cr3t-prod"

exec_vault "vault kv put secret/prod/cache \
    host=cache.prod.example.com \
    port=6379 \
    password=cache-s3cr3t-prod"

exec_vault "vault kv put secret/prod/api \
    endpoint=https://api.prod.example.com \
    key=prod-api-key-abc123 \
    timeout=30"

# staging secrets
exec_vault "vault kv put secret/staging/database \
    host=db.staging.example.com \
    port=5432 \
    username=admin \
    password=s3cr3t-staging"

exec_vault "vault kv put secret/staging/cache \
    host=cache.staging.example.com \
    port=6379 \
    password=cache-s3cr3t-staging"

exec_vault "vault kv put secret/staging/api \
    endpoint=https://api.staging.example.com \
    key=staging-api-key-xyz789 \
    timeout=30"

echo "Done. Vault is ready at http://localhost:${VAULT_PORT} (token: root)"

# Create a plain 'vsh' symlink in build/ for easy use in the demo tape
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)  ARCH=amd64 ;;
    aarch64|armv8*) ARCH=arm64 ;;
esac
BINARY="build/vsh_linux_${ARCH}"
if [[ -f "$BINARY" ]]; then
    ln -sf "vsh_linux_${ARCH}" build/vsh
    echo "Symlinked build/vsh -> vsh_linux_${ARCH}"
else
    echo "Binary $BINARY not found — run 'make compile' first"
fi
