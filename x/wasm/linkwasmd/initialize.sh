#!/usr/bin/env bash
set -ex

LINKD=${LINKD:-linkwasmd}

# initialize
rm -rf ~/.linkwasmd

# Configure your CLI to eliminate need for chain-id flag
# ${LINKD} config chain-id linkwasm
# ${LINKD} config output json
# ${LINKD} config indent true
# ${LINKD} config trust-node true
# ${LINKD} config keyring-backend test

# Initialize configuration files and genesis file
# moniker is the name of your node
${LINKD} init solo --chain-id linkwasm

${LINKD} keys add jack --keyring-backend=test
${LINKD} keys add alice --keyring-backend=test
${LINKD} keys add bob --keyring-backend=test
${LINKD} keys add rinah --keyring-backend=test
${LINKD} keys add sam --keyring-backend=test
${LINKD} keys add evelyn --keyring-backend=test

# Add both accounts, with coins to the genesis file
${LINKD} add-genesis-account $(${LINKD} keys show jack -a --keyring-backend=test) 1000link,100000000stake
${LINKD} add-genesis-account $(${LINKD} keys show alice -a --keyring-backend=test) 1000link,100000000stake
${LINKD} add-genesis-account $(${LINKD} keys show bob -a --keyring-backend=test) 1000link,100000000stake
${LINKD} add-genesis-account $(${LINKD} keys show rinah -a --keyring-backend=test) 1000link,100000000stake
${LINKD} add-genesis-account $(${LINKD} keys show sam -a --keyring-backend=test) 1000link,100000000stake
${LINKD} add-genesis-account $(${LINKD} keys show evelyn -a --keyring-backend=test) 1000link,100000000stake

${LINKD} gentx jack 100000000stake --keyring-backend=test --chain-id=linkwasm

${LINKD} collect-gentxs

${LINKD} validate-genesis

