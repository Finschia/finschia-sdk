#!/bin/bash

#
DENOM="stake"
ROLLUPNAME="test-rollup"
DEPOSIT="20"
WITHDRAW="10"

BASE_DIR=~/.simapp

# Get address which want to register for sequencer
SEQUENCER=$(./L1/simd keys list --keyring-backend=test --home $BASE_DIR --output json | jq -r '.[1]'.name)
SEQUENCERADDRESS=$(./L1/simd keys list --keyring-backend=test --home $BASE_DIR --output json | jq -r '.[1]'.address)
SEQUENCERPUBKEY=$(./L1/simd keys list --keyring-backend=test --home $BASE_DIR --output json | jq -r '.[1]'.pubkey)
echo "Sequencer Info"
echo $SEQUENCER
echo $SEQUENCERADDRESS
echo $SEQUENCERPUBKEY

echo "# Check init balance"
INITBALANCE=$(./L1/simd query bank balances $SEQUENCERADDRESS --node "tcp://localhost:26659" --home $BASE_DIR --output json | jq -r '.balances[0].amount')
echo $INITBALANCE

echo "# Create rollup"
./L1/simd tx rollup create-rollup $ROLLUPNAME 5 --from $SEQUENCER --keyring-backend=test --node "tcp://localhost:26659" --home $BASE_DIR --chain-id sim -y

sleep 5

echo "# Check created rollup"
./L1/simd query rollup show-rollup $ROLLUPNAME --node "tcp://localhost:26659"

echo "# Check rollup list"
./L1/simd query rollup list --node "tcp://localhost:26659"

echo "# Register sequencer"
./L1/simd tx rollup register-sequencer "test-rollup" ${SEQUENCERPUBKEY} $DEPOSIT$DENOM --from $SEQUENCER --keyring-backend test --node "tcp://localhost:26659" --home $BASE_DIR --chain-id sim -y

sleep 5

echo "# Check balance after registered sequencer"
BALANCEAFTERREGISTER=$(./L1/simd query bank balances $SEQUENCERADDRESS --node "tcp://localhost:26659" --home $BASE_DIR --output json | jq -r '.balances[0].amount')
echo $BALANCEAFTERREGISTER

if [ $((${INITBALANCE}-${DEPOSIT})) -ne ${BALANCEAFTERREGISTER} ]; then
    echo "The balance after registering the sequencer does not match."
    exit 1
else
    echo "Register sequencer done"
fi

echo "# Check sequencer by rollup name"
./L1/simd query rollup show-sequencers-by-rollup $ROLLUPNAME --node "tcp://localhost:26659" --output json

echo "# Check sequencer"
./L1/simd query rollup show-sequencer $SEQUENCERADDRESS --node "tcp://localhost:26659" --home $BASE_DIR --output json

echo "# Withdraw deposit"
./L1/simd tx rollup withdraw-by-sequencer $ROLLUPNAME $WITHDRAW$DENOM --from $SEQUENCER --keyring-backend test --node "tcp://localhost:26659" --home $BASE_DIR --chain-id sim -y

sleep 5

echo "# Check balance after withdraw"
BALANCEAFTERWITHDRAW=$(./L1/simd query bank balances $SEQUENCERADDRESS --node "tcp://localhost:26659" --home $BASE_DIR --output json | jq -r '.balances[0].amount')

if [ $((${INITBALANCE}-${DEPOSIT}+${WITHDRAW})) -ne ${BALANCEAFTERWITHDRAW} ]; then
    echo "The balance after withdraw does not match."
    exit 1
else
    echo "Withdraw done"
fi

echo "# Deposit by sequencer again"
./L1/simd tx rollup deposit-by-sequencer "test-rollup" $DEPOSIT$DENOM --from $SEQUENCER --keyring-backend=test --node "tcp://localhost:26659" --home $BASE_DIR --chain-id sim -y

sleep 5

echo "# Check balance after deposit again"
BALANCEAFTERDEPOSITAGAIN=$(./L1/simd query bank balances $SEQUENCERADDRESS --node "tcp://localhost:26659" --home $BASE_DIR --output json | jq -r '.balances[0].amount')

if [ $((${INITBALANCE}-2*${DEPOSIT}+${WITHDRAW})) -ne ${BALANCEAFTERDEPOSITAGAIN} ]; then
    echo "The balance after deposit again does not match."
    exit 1
else
    echo "Deposit again done"
fi

echo "Prepare rollup all done."
