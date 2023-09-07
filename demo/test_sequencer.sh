#!/bin/bash

L2BINARYNAME=rollupd
ROLLUP_NAME=test-rollup
GOPATH=$HOME/go
L1_CHAIN_ID=sim
L2_CHAIN_ID=sim2
L1_KEYRING_DIR=~/.simapp
L2_KEYRING_DIR=~/.l2simapp
RPC_URI=http://localhost:26659
TEST_SEQUENCER_ADDRESS=link1twsfmuj28ndph54k4nw8crwu8h9c8mh3rtx705
SEQUENCER_DIR=simapp0
TEST_MNEMONIC="mind flame tobacco sense move hammer drift crime ring globe art gaze cinnamon helmet cruise special produce notable negative wait path scrap recall have"

# Send test
${L2BINARYNAME} keys add alice --home $L2_KEYRING_DIR/$SEQUENCER_DIR --keyring-backend=test
sleep 1
ALICE_ADDR=$(${L2BINARYNAME} keys show alice -a --home $L2_KEYRING_DIR/$SEQUENCER_DIR --keyring-backend=test)
echo $ALICE_ADDR
${L2BINARYNAME} tx bank send validator $ALICE_ADDR 100stake --home $L2_KEYRING_DIR/$SEQUENCER_DIR --keyring-backend=test --chain-id $L2_CHAIN_ID -y
sleep 30

# Check alice's balance
BALANCE=$(${L2BINARYNAME} query bank balances $ALICE_ADDR --home $L2_KEYRING_DIR/$SEQUENCER_DIR --output json | jq -r '.balances[0].amount')
if [ 100 -eq ${BALANCE} ]; then
    echo "send success!"
else
    echo "send failed..."
fi

## NOTE: There is a bug in the current Ramus that only validators can execute tx.
# ${L2BINARYNAME} keys add bob --home $L2_KEYRING_DIR/$SEQUENCER_DIR --keyring-backend=test
# sleep 1
# BOB_ADDR=$(${L2BINARYNAME} keys show bob -a --home $L2_KEYRING_DIR/$SEQUENCER_DIR --keyring-backend=test)
# echo $ALICE_ADDR
# ${L2BINARYNAME} tx bank send $ALICE_ADDR $BOB_ADDR 10stake --home $L2_KEYRING_DIR/$SEQUENCER_DIR --keyring-backend=test --chain-id $L2_CHAIN_ID -y
# sleep 30

# # Check alice's balance
# BALANCE=$(${L2BINARYNAME} query bank balances $BOB_ADDR --home $L2_KEYRING_DIR/$SEQUENCER_DIR --output json | jq -r '.balances[0].amount')
# if [ 10 -eq ${BALANCE} ]; then
#     echo "send success!"
# else
#     echo "send failed..."
# fi
