#!/bin/sh

BINARY="simd"

BASE_DIR=~/.simapp
CHAIN_DIR="${BASE_DIR}/simapp0"
BINARY_DIR="$(which ${BINARY})"

export DAEMON_NAME="${BINARY}"
export DAEMON_HOME="${CHAIN_DIR}"

if [ -z ${BINARY_DIR} ]; then
  echo "Failed to get ${BINARY_DIR}. Aborting..."
  exit 1
fi

if [ ! -d ${CHAIN_DIR} ]; then
  echo "${CHAIN_DIR} is not exist. Aborting..."
  exit 1
fi

BIN_DIR="${CHAIN_DIR}/cosmovisor/genesis/bin"
if ! mkdir -p ${BIN_DIR}; then
  echo "Failed to create cosmovisor/genesis/bin folder(${CHAIN_DIR}). Aborting..."
  exit 1
fi

if ! cp ${BINARY_DIR} ${BIN_DIR}; then
  echo "Failed to copy ${BINARY_DIR} to ${BIN_DIR}. Aborting..."
  exit 1
fi

echo "cosmovisor version: $(cosmovisor version)"
