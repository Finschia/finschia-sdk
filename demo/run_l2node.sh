#!/bin/bash

L2NODE_BINARYNAME=rollupd2
ROLLUP_NAME=test-rollup
L1_CHAIN_ID=sim
L2_CHAIN_ID=sim2
L1_KEYRING_DIR=~/.simapp
L2_KEYRING_DIR=~/.l2simapp2
NAMESPACE_ID=$(openssl rand -hex 8)
RPC_URI=http://localhost:26659
TEST_SEQUENCER_P2P_ID=12D3KooWKUFu5UWCwuggvCR5sFFXRx8WxzxDJDX4dani1NQVAXWD
TEST_SEQUENCER_ADDRESS=link1twsfmuj28ndph54k4nw8crwu8h9c8mh3rtx705
DA_BLOCK_HEIGH=1
SEQUENCER_DIR=simapp0
TEST_MNEMONIC="mind flame tobacco sense move hammer drift crime ring globe art gaze cinnamon helmet cruise special produce notable negative wait path scrap recall have"

# Reset
rm -rf $L2_KEYRING_DIR
pkill $L2NODE_BINARYNAME
rm $GOPATH/bin/$L2NODE_BINARYNAME

cd ..

# Build & rename
make build && cp -r build/simd $GOPATH/bin/$L2NODE_BINARYNAME
${L2NODE_BINARYNAME} version

# Init sequencer

${L2NODE_BINARYNAME} init rollupdemo --home $L2_KEYRING_DIR/$SEQUENCER_DIR --chain-id $L2_CHAIN_ID > /dev/null 2>&1
${L2NODE_BINARYNAME} keys add validator --keyring-backend=test --home $L2_KEYRING_DIR/$SEQUENCER_DIR --recover --account=0 <<< ${TEST_MNEMONIC} > /dev/null 2>&1
${L2NODE_BINARYNAME} add-genesis-account $(${L2NODE_BINARYNAME} --home $L2_KEYRING_DIR/$SEQUENCER_DIR keys show validator -a --keyring-backend=test) 100000000000stake,100000000000tcony --home $L2_KEYRING_DIR/$SEQUENCER_DIR > /dev/null 2>&1
${L2NODE_BINARYNAME} gentx validator 10000000000stake --keyring-backend=test --home $L2_KEYRING_DIR/$SEQUENCER_DIR --chain-id=$L2_CHAIN_ID > /dev/null 2>&1
${L2NODE_BINARYNAME} collect-gentxs --home $L2_KEYRING_DIR/$SEQUENCER_DIR > /dev/null 2>&1

# Run L2 sequencer
${L2NODE_BINARYNAME} start --log_level "debug" --home $L2_KEYRING_DIR/$SEQUENCER_DIR --p2p.laddr "tcp://0.0.0.0:26555" --p2p.seeds="tcp://${TEST_SEQUENCER_P2P_ID}@127.0.0.1:26556" --rpc.laddr=tcp://0.0.0.0:26654 --grpc.address "0.0.0.0:9192" --grpc-web.address "0.0.0.0:9193" --rollkit.da_layer finschia --rollkit.da_config='{"rpc_uri":"'$RPC_URI'","chain_id":"'$L1_CHAIN_ID'","keyring_dir":"'$L1_KEYRING_DIR'","from":"'$TEST_SEQUENCER_ADDRESS'", "rollup_name":"'$ROLLUP_NAME'"}' --rollkit.namespace_id $NAMESPACE_ID  --rollkit.da_start_height $DA_BLOCK_HEIGH > $L2_KEYRING_DIR/$L2_CHAIN_ID.log 2>&1 &
