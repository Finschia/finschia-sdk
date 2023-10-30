#!/bin/bash

L2BINARYNAME=rollupd
ROLLUP_NAME=test-rollup
L1_CHAIN_ID=sim
L2_CHAIN_ID=sim2
L1_KEYRING_DIR=~/.simapp
L2_KEYRING_DIR=~/.l2simapp
NAMESPACE_ID=$(openssl rand -hex 8)
RPC_URI=http://localhost:26659
TEST_SEQUENCER_ADDRESS=link1twsfmuj28ndph54k4nw8crwu8h9c8mh3rtx705
DA_BLOCK_HEIGH=1
SEQUENCER_DIR=simapp0
TEST_MNEMONIC="mind flame tobacco sense move hammer drift crime ring globe art gaze cinnamon helmet cruise special produce notable negative wait path scrap recall have"

# Reset
rm -rf $L2_KEYRING_DIR
pkill $L2BINARYNAME
rm $GOPATH/bin/$L2BINARYNAME

# Run L1 chain
./run_chain.sh

sleep 10

# Prepare rollup info
./prepare_rollup.sh

cd ..

# Build & rename
make build && cp -r build/simd $GOPATH/bin/$L2BINARYNAME
$L2BINARYNAME version

# Init sequencer

$L2BINARYNAME init rollupdemo --home $L2_KEYRING_DIR/$SEQUENCER_DIR --chain-id $L2_CHAIN_ID > /dev/null 2>&1
$L2BINARYNAME keys add validator --keyring-backend=test --home $L2_KEYRING_DIR/$SEQUENCER_DIR --recover --account=0 <<< ${TEST_MNEMONIC} > /dev/null 2>&1
$L2BINARYNAME add-genesis-account $($L2BINARYNAME --home $L2_KEYRING_DIR/$SEQUENCER_DIR keys show validator -a --keyring-backend=test) 100000000000stake,100000000000tcony --home $L2_KEYRING_DIR/$SEQUENCER_DIR > /dev/null 2>&1
$L2BINARYNAME gentx validator 10000000000stake --keyring-backend=test --home $L2_KEYRING_DIR/$SEQUENCER_DIR --chain-id=$L2_CHAIN_ID > /dev/null 2>&1
$L2BINARYNAME collect-gentxs --home $L2_KEYRING_DIR/$SEQUENCER_DIR > /dev/null 2>&1

# Run L2 sequencer
$L2BINARYNAME start --home $L2_KEYRING_DIR/$SEQUENCER_DIR --p2p.laddr "tcp://0.0.0.0:26556" --grpc.address "0.0.0.0:9190" --grpc-web.address "0.0.0.0:9191" --rollkit.sequencer "true" --rollkit.da_layer finschia --rollkit.da_config='{"rpc_uri":"'$RPC_URI'","chain_id":"'$L1_CHAIN_ID'","keyring_dir":"'$L1_KEYRING_DIR'","from":"'$TEST_SEQUENCER_ADDRESS'", "rollup_name":"'$ROLLUP_NAME'"}' --rollkit.namespace_id $NAMESPACE_ID  --rollkit.da_start_height $DA_BLOCK_HEIGH > $L2_KEYRING_DIR/$L2_CHAIN_ID.log 2>&1 &
sleep 10

# Send test
$L2BINARYNAME keys add alice --home $L2_KEYRING_DIR/$SEQUENCER_DIR --keyring-backend=test
sleep 1
ALICE_ADDR=$($L2BINARYNAME keys show alice -a --home $L2_KEYRING_DIR/$SEQUENCER_DIR --keyring-backend=test)
echo $ALICE_ADDR
$L2BINARYNAME tx bank send validator $ALICE_ADDR 100stake --home $L2_KEYRING_DIR/$SEQUENCER_DIR --keyring-backend=test --chain-id $L2_CHAIN_ID -y
sleep 30

# Check alice's balance
BALANCE=$($L2BINARYNAME query bank balances $ALICE_ADDR --home $L2_KEYRING_DIR/$SEQUENCER_DIR --output json | jq -r '.balances[0].amount')
if [ 100 -eq ${BALANCE} ]; then
    echo "send success!"
else
    echo "send failed..."
fi

# NOTE: There is a bug in the current Ramus that only validators can execute tx.
$L2BINARYNAME keys add bob --home $L2_KEYRING_DIR/$SEQUENCER_DIR --keyring-backend=test
sleep 1
BOB_ADDR=$($L2BINARYNAME keys show bob -a --home $L2_KEYRING_DIR/$SEQUENCER_DIR --keyring-backend=test)
echo $ALICE_ADDR
$L2BINARYNAME tx bank send $ALICE_ADDR $BOB_ADDR 10stake --home $L2_KEYRING_DIR/$SEQUENCER_DIR --keyring-backend=test --chain-id $L2_CHAIN_ID -y
sleep 30

# Check alice's balance
BALANCE=$($L2BINARYNAME query bank balances $BOB_ADDR --home $L2_KEYRING_DIR/$SEQUENCER_DIR --output json | jq -r '.balances[0].amount')
if [ 10 -eq ${BALANCE} ]; then
    echo "send success!"
else
    echo "send failed..."
fi
