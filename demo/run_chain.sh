#!/bin/bash

# Prepare chain
TEST_MNEMONIC="mind flame tobacco sense move hammer drift crime ring globe art gaze cinnamon helmet cruise special produce notable negative wait path scrap recall have"
BASE_DIR=~/.simapp
CHAIN_ID=sim

./L1/simd-darwin-arm64 init rollupdemo --home $BASE_DIR --chain-id $CHAIN_ID || exit 1
./L1/simd-darwin-arm64 keys add validator --keyring-backend=test --home $BASE_DIR --recover --account=0 <<< ${TEST_MNEMONIC} || exit 1
./L1/simd-darwin-arm64 keys add sequencer --keyring-backend=test --home $BASE_DIR --recover --account=1 <<< ${TEST_MNEMONIC} || exit 1
./L1/simd-darwin-arm64 keys add challenger --keyring-backend=test --home $BASE_DIR --recover --account=2 <<< ${TEST_MNEMONIC} || exit 1
./L1/simd-darwin-arm64 add-genesis-account $(./L1/simd-darwin-arm64 --home $BASE_DIR keys show validator -a --keyring-backend=test) 100000000000stake,100000000000tcony --home $BASE_DIR || exit 1
./L1/simd-darwin-arm64 add-genesis-account $(./L1/simd-darwin-arm64 --home $BASE_DIR keys show sequencer -a --keyring-backend=test) 100000000000stake,100000000000tcony --home $BASE_DIR || exit 1
./L1/simd-darwin-arm64 add-genesis-account $(./L1/simd-darwin-arm64 --home $BASE_DIR keys show challenger -a --keyring-backend=test) 100000000000stake,100000000000tcony --home $BASE_DIR || exit 1
./L1/simd-darwin-arm64 gentx validator 10000000000stake --keyring-backend=test --home $BASE_DIR --chain-id=$CHAIN_ID || exit 1
./L1/simd-darwin-arm64 collect-gentxs --home $BASE_DIR || exit 1

# Run chain
./L1/simd-darwin-arm64 start --rpc.laddr "tcp://127.0.0.1:26659" --home $BASE_DIR > $BASE_DIR/$CHAIN_ID.log 2>&1 &
