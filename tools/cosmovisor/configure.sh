#!/bin/sh

set -e

# cosmovisor config assertions
[ -n "$DAEMON_HOME" ]
[ -d $DAEMON_HOME ]
[ -n "$DAEMON_NAME" ]
[ -x "$(which $DAEMON_NAME)" ]

# app config assertions
[ -n "$CHAIN_ID" ]
keyring_backend=test
[ -n "$KEY_MNEMONIC" ]
[ -n "$KEY_INDEX" ]
[ $KEY_INDEX -ge 0 ]

# setup the chain
$DAEMON_NAME --home $DAEMON_HOME init validator$KEY_INDEX --chain-id $CHAIN_ID >/dev/null 2>&1

# modify genesis for the upgrade test
genesis_file=$DAEMON_HOME/config/genesis.json
genesis=$(cat $genesis_file)
echo $genesis | jq -f alter_genesis.jq >$genesis_file

# add validators
id=$KEY_INDEX
echo $KEY_MNEMONIC | $DAEMON_NAME --home $DAEMON_HOME keys --keyring-backend $keyring_backend add validator$id --recover --account $KEY_INDEX >/dev/null
$DAEMON_NAME --home $DAEMON_HOME add-genesis-account --keyring-backend $keyring_backend validator$id 1000000000stake
$DAEMON_NAME --home $DAEMON_HOME gentx --keyring-backend $keyring_backend validator$id 1000000stake --chain-id $CHAIN_ID
$DAEMON_NAME --home $DAEMON_HOME collect-gentxs
