#!/bin/bash

# Remove data
pkill ./L1/simd
rm -rf ~/.simapp

if ! [ -x "$(which simd)" ]; then
  echo "Error: simd is not installed. Try running 'cd ../../../ && make install'" >&2
  exit 1
fi

# Prepare chain
TEST_MNEMONIC="mind flame tobacco sense move hammer drift crime ring globe art gaze cinnamon helmet cruise special produce notable negative wait path scrap recall have"
BASE_DIR=~/.simapp
CHAIN_ID=sim

./L1/simd init rollupdemo --home $BASE_DIR --chain-id $CHAIN_ID > /dev/null 2>&1
./L1/simd keys add validator --keyring-backend=test --home $BASE_DIR --recover --account=0 <<< ${TEST_MNEMONIC} > /dev/null 2>&1
./L1/simd keys add sequencer --keyring-backend=test --home $BASE_DIR --recover --account=1 <<< ${TEST_MNEMONIC} > /dev/null 2>&1
./L1/simd keys add challenger --keyring-backend=test --home $BASE_DIR --recover --account=2 <<< ${TEST_MNEMONIC} > /dev/null 2>&1
./L1/simd add-genesis-account $(./L1/simd --home $BASE_DIR keys show validator -a --keyring-backend=test) 100000000000stake,100000000000tcony --home $BASE_DIR > /dev/null 2>&1
./L1/simd add-genesis-account $(./L1/simd --home $BASE_DIR keys show sequencer -a --keyring-backend=test) 100000000000stake,100000000000tcony --home $BASE_DIR > /dev/null 2>&1
./L1/simd add-genesis-account $(./L1/simd --home $BASE_DIR keys show challenger -a --keyring-backend=test) 100000000000stake,100000000000tcony --home $BASE_DIR > /dev/null 2>&1
./L1/simd gentx validator 10000000000stake --keyring-backend=test --home $BASE_DIR --chain-id=$CHAIN_ID > /dev/null 2>&1
./L1/simd collect-gentxs --home $BASE_DIR > /dev/null 2>&1

# Run chain
./L1/simd start --rpc.laddr "tcp://127.0.0.1:26659" --home $BASE_DIR > $BASE_DIR/$CHAIN_ID.log 2>&1 &
