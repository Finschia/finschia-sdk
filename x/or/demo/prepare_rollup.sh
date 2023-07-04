#!/bin/bash

#
DENOM="stake"
ROLLUPNAME="test-rollup"
DEPOSIT="20"
WITHDRAW="10"

# Get address which want to register for sequencer
SEQUENCER=$(simd keys list --keyring-backend=test --home ~/.simapp/simapp0 --output json | jq -r '.[0]'.name)
SEQUENCERADDRESS=$(simd keys list --keyring-backend=test --home ~/.simapp/simapp0 --output json | jq -r '.[0]'.address)
SEQUENCERPUBKEY=$(simd keys list --keyring-backend=test --home ~/.simapp/simapp0 --output json | jq -r '.[0]'.pubkey)
echo "Sequencer Info"
echo $SEQUENCER
echo $SEQUENCERADDRESS
echo $SEQUENCERPUBKEY

echo "# Check init balance"
INITBALANCE=$(simd query bank balances $SEQUENCERADDRESS --home ~/.simapp/simapp0 --output json | jq -r '.balances[0].amount')
echo $INITBALANCE

echo "# Create rollup"
simd tx rollup create-rollup $ROLLUPNAME --from $SEQUENCER --keyring-backend=test --home ~/.simapp/simapp0 --chain-id sim -y

sleep 2

echo "# Check created rollup"
simd query rollup show-rollup $ROLLUPNAME

echo "# Check rollup list"
simd query rollup list

echo "# Register sequencer"
simd tx rollup register-sequencer test-rollup ${SEQUENCERPUBKEY} --from $SEQUENCER --amount $DEPOSIT$DENOM --keyring-backend test --home ~/.simapp/simapp0 --chain-id sim -y

sleep 2

echo "# Check balance after registered sequencer"
BALANCEAFTERREGISTER=$(simd query bank balances $SEQUENCERADDRESS --home ~/.simapp/simapp0 --output json | jq -r '.balances[0].amount')
echo $BALANCEAFTERREGISTER

if [ $((${INITBALANCE}-${DEPOSIT})) -ne ${BALANCEAFTERREGISTER} ]; then
    echo "The balance after registering the sequencer does not match."
    exit 1
else
    echo "Register sequencer done"
fi

echo "# Check sequencer by rollup name"
simd query rollup show-sequencers-by-rollup $ROLLUPNAME --output json

echo "# Check sequencer"
simd query rollup show-sequencer $SEQUENCERADDRESS --home ~/.simapp/simapp0 --output json

echo "# Withdraw deposit"
simd tx rollup withdraw-by-sequencer $ROLLUPNAME --from $SEQUENCER --amount $WITHDRAW$DENOM --keyring-backend test --home ~/.simapp/simapp0 --chain-id sim -y

sleep 2

echo "# Check balance after withdraw"
BALANCEAFTERWITHDRAW=$(simd query bank balances $SEQUENCERADDRESS --home ~/.simapp/simapp0 --output json | jq -r '.balances[0].amount')

if [ $((${INITBALANCE}-${DEPOSIT}+${WITHDRAW})) -ne ${BALANCEAFTERWITHDRAW} ]; then
    echo "The balance after withdraw does not match."
    exit 1
else
    echo "Withdraw done"
fi

echo "# Deposit by sequencer again"
simd tx rollup deposit-by-sequencer test-rollup --from $SEQUENCER --amount $DEPOSIT$DENOM --keyring-backend=test --home ~/.simapp/simapp0 --chain-id sim -y

sleep 2

echo "# Check balance after deposit again"
BALANCEAFTERDEPOSITAGAIN=$(simd query bank balances $SEQUENCERADDRESS --home ~/.simapp/simapp0 --output json | jq -r '.balances[0].amount')

if [ $((${INITBALANCE}-2*${DEPOSIT}+${WITHDRAW})) -ne ${BALANCEAFTERDEPOSITAGAIN} ]; then
    echo "The balance after deposit again does not match."
    exit 1
else
    echo "Deposit again done"
fi

echo "Prepare rollup all done."
