#!/usr/bin/env sh
set -ex

PASSWORD="1234567890"
# initialize
rm -rf ~/.linkd ~/.linkcli

# Configure your CLI to eliminate need for chain-id flag
linkcli config chain-id link
linkcli config output json
linkcli config indent true
linkcli config trust-node true

# Initialize configuration files and genesis file
# moniker is the name of your node
linkd init solo --chain-id link


# Copy the `Address` output here and save it for later use
# [optional] add "--ledger" at the end to use a Ledger Nano S
linkcli keys add jack  <<< ${PASSWORD} <<< ${PASSWORD}

# Copy the `Address` output here and save it for later use
linkcli keys add alice  <<< ${PASSWORD} <<< ${PASSWORD}

# Add both accounts, with coins to the genesis file
linkd add-genesis-account $(linkcli keys show jack -a) 1000link,100000000stake
linkd add-genesis-account $(linkcli keys show alice -a) 1000link,100000000stake 


linkd gentx --name jack <<< ${PASSWORD}

linkd collect-gentxs

linkd validate-genesis

# linkd start
linkd start --log_level *:debug
