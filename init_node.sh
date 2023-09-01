#!/bin/sh

display_usage() {
  echo "\nMissing $1 parameter. Please check if all parameters were specified."
  echo "\nUsage: ./run_node [CHAIN_ID] [N, default=1]"
  exit 1
}

if [ -z "$1" ]; then
  display_usage "[CHAIN_ID]"
fi

KEYRING="--keyring-backend=test"
SILENT=1

redirect() {
  if [ "${SILENT}" -eq 1 ]; then
    "$@" > /dev/null 2>&1
  else
    "$@"
  fi
}

BINARY=simd
BASE_DIR=~/.l2simapp
CHAIN_DIR_PREFIX="${BASE_DIR}/simapp"
GENTXS_DIR="${BASE_DIR}/gentxs"
CHAIN_ID=$1
MONIKER_PREFIX="node"

# Control N node count
N=1
if [ -n "$2" ]; then
  N=$2

  if [ "${N}" -le 0 ]; then
    echo "N must be positive int. Aborting..."
    exit 1
  elif [ "${N}" -gt 10 ]; then
    echo "N must be smaller than 10. Aborting..."
    echo "If you need to create node more than 10, modify init_node.sh line 38"
    exit 1
  fi
fi

VALIDATOR_PREFIX="validator"
COINS="100000000000stake,100000000000ukrw"
DELEGATE="10000000000stake"

# Create base dir and gentxs dir
if ! mkdir -p ${GENTXS_DIR} 2> /dev/null; then
  echo "Failed to create chain folder(${GENTXS_DIR}). Aborting..."
  exit 1
fi

# Please do not use the TEST_MNEMONIC for production purpose
TEST_MNEMONIC="mind flame tobacco sense move hammer drift crime ring globe art gaze cinnamon helmet cruise special produce notable negative wait path scrap recall have"

# Initialize config files and genesis file
# Create genesis account and gentx
for ((i = 0; i < N; i++))
  do
    # ~/.simapp0, ~/.simapp1, ...
    CHAIN_DIR="${CHAIN_DIR_PREFIX}${i}"

    # Add dir for chain, exit if error
    if ! mkdir -p ${CHAIN_DIR} 2>/dev/null; then
      echo "Failed to create chain folder(${CHAIN_DIR}). Aborting..."
      exit 1
    fi

    # Initialize configuration files and genesis file
    # moniker is the name of your node
    MONIKER="${MONIKER_PREFIX}${i}"
    redirect ${BINARY} init ${MONIKER} --home ${CHAIN_DIR} --chain-id ${CHAIN_ID}

    # Create N genesis account(with mnemonic), so N-th chain's N account is same with M-th chain's
    for ((j = 0; j < N; j++))
      do
        VALIDATOR="${VALIDATOR_PREFIX}${j}"
        ${BINARY} keys add ${VALIDATOR} ${KEYRING} --home ${CHAIN_DIR} --recover --account ${j} --output json <<< ${TEST_MNEMONIC} >> ${CHAIN_DIR}/validator_seed.json 2> /dev/null
        redirect ${BINARY} add-genesis-account $(${BINARY} --home ${CHAIN_DIR} keys ${KEYRING} show ${VALIDATOR} -a) ${COINS} --home ${CHAIN_DIR}
      done

    # Make gentx file and move it to GENTXS folder
    VALIDATOR="${VALIDATOR_PREFIX}${i}"
    redirect ${BINARY} gentx ${VALIDATOR} ${DELEGATE} ${KEYRING} --home ${CHAIN_DIR} --chain-id ${CHAIN_ID}
    mv "${CHAIN_DIR}/config/gentx/$(ls ${CHAIN_DIR}/config/gentx | grep .json)" "${GENTXS_DIR}/${MONIKER}.json"
    rm -r "${CHAIN_DIR}/config/gentx"
  done

SRC_GENESIS_FIlE="${CHAIN_DIR_PREFIX}0/config/genesis.json"
RPC_PORT=26657
P2P_PORT=26656
PROF_PORT=6060
GRPC_PORT=9090
GRPC_WEB_PORT=9091

# Set genesis file and config(port, peer, ...)
CHAIN_0_DIR="${CHAIN_DIR_PREFIX}0"
for ((i = 0; i < N; i++))
  do
    CHAIN_DIR="${CHAIN_DIR_PREFIX}${i}"

    # Set genesis file of 0-th chain dir and copy to other chains
    # If we call collect-gentxs at each chains, genesis_time values can be different.
    if [ ${i} -eq 0 ]; then
      redirect ${BINARY} collect-gentxs --home ${CHAIN_DIR} --gentx-dir ${GENTXS_DIR}
    else
      cp ${SRC_GENESIS_FIlE} "${CHAIN_DIR}/config"
    fi

    MONIKER="${MONIKER_PREFIX}${i}"
    MEMO=`sed 's/"//g' <<< \`cat ${GENTXS_DIR}/${MONIKER}.json | jq '.body.memo'\``
    MEMO_SPLIT=(`echo ${MEMO} | tr ":" "\n"`)

    # Set proper defaults and change ports (use a different sed for Mac or Linux)
    if [ "`uname`" = "Linux" ]; then
      sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"${RPC_PORT}"'"#g' ${CHAIN_DIR}/config/config.toml
      sed -i 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"${P2P_PORT}"'"#g' ${CHAIN_DIR}/config/config.toml
      sed -i 's#"localhost:6060"#"localhost:'"${PROF_PORT}"'"#g' ${CHAIN_DIR}/config/config.toml
      sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/config/config.toml
      sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/config/config.toml
      sed -i 's/addr_book_strict = true/addr_book_strict = false/g' ${CHAIN_DIR}/config/config.toml  # for local test
      sed -i 's/allow_duplicate_ip = false/allow_duplicate_ip = true/g' ${CHAIN_DIR}/config/config.toml  # allow duplicated ip

      sed -i 's#'"${MEMO}"'#'"${MEMO_SPLIT[1]}"':'"${P2P_PORT}"'#g' ${CHAIN_0_DIR}/config/config.toml  # change port of persistent_peers

      sed -i 's/pruning = "default"/pruning = "nothing"/g' ${CHAIN_DIR}/config/app.toml
      sed -i 's#"0.0.0.0:9091"#"0.0.0.0:'"${GRPC_WEB_PORT}"'"#g' ${CHAIN_DIR}/config/app.toml
      sed -i 's#"0.0.0.0:9090"#"0.0.0.0:'"${GRPC_PORT}"'"#g' ${CHAIN_DIR}/config/app.toml
    else
      sed -i '' 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"${RPC_PORT}"'"#g' ${CHAIN_DIR}/config/config.toml
      sed -i '' 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"${P2P_PORT}"'"#g' ${CHAIN_DIR}/config/config.toml
      sed -i '' 's#"localhost:6060"#"localhost:'"${PROF_PORT}"'"#g' ${CHAIN_DIR}/config/config.toml
      sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/config/config.toml
      sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/config/config.toml
      sed -i '' 's/addr_book_strict = true/addr_book_strict = false/g' ${CHAIN_DIR}/config/config.toml  # for local test
      sed -i '' 's/allow_duplicate_ip = false/allow_duplicate_ip = true/g' ${CHAIN_DIR}/config/config.toml  # allow duplicated ip
      
      sed -i '' 's#'"${MEMO}"'#'"${MEMO_SPLIT[1]}"':'"${P2P_PORT}"'#g' ${CHAIN_0_DIR}/config/config.toml  # change port of persistent_peers
      
      sed -i '' 's/pruning = "default"/pruning = "nothing"/g' ${CHAIN_DIR}/config/app.toml
      sed -i '' 's#"0.0.0.0:9091"#"0.0.0.0:'"${GRPC_WEB_PORT}"'"#g' ${CHAIN_DIR}/config/app.toml
      sed -i '' 's#"0.0.0.0:9090"#"0.0.0.0:'"${GRPC_PORT}"'"#g' ${CHAIN_DIR}/config/app.toml
    fi

    echo "${BINARY} instance: home ${CHAIN_DIR} | chain-id ${CHAIN_ID} | p2p=:${P2P_PORT} | rpc=:${RPC_PORT} | profiling=:${PROF_PORT} | grpc=:${GRPC_PORT} | grpc-web=:${GRPC_WEB_PORT}"
    RPC_PORT=`expr ${RPC_PORT} + 2`
    P2P_PORT=`expr ${P2P_PORT} + 2`
    PROF_PORT=`expr ${PROF_PORT} + 1`
    GRPC_PORT=`expr ${GRPC_PORT} + 2`
    GRPC_WEB_PORT=`expr ${GRPC_WEB_PORT} + 2`
  done
  
