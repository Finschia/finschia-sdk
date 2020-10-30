#!/usr/bin/env bash
set -ex

LINKCLI=${LINKCLI:-linkwasmcli}
LINKD=${LINKD:-linkwasmd}

# initialize
rm -rf ~/.linkwasmd ~/.linkwasmcli

# Configure your CLI to eliminate need for chain-id flag
${LINKCLI} config chain-id linkwasm
${LINKCLI} config output json
${LINKCLI} config indent true
${LINKCLI} config trust-node true
${LINKCLI} config keyring-backend test

# Initialize configuration files and genesis file
# moniker is the name of your node
${LINKD} init solo --chain-id linkwasm

${LINKCLI} keys add jack
${LINKCLI} keys add alice
${LINKCLI} keys add bob
${LINKCLI} keys add rinah
${LINKCLI} keys add sam
${LINKCLI} keys add evelyn

# Add both accounts, with coins to the genesis file
${LINKD} add-genesis-account $(${LINKCLI} keys show jack -a) 1000link,100000000stake
${LINKD} add-genesis-account $(${LINKCLI} keys show alice -a) 1000link,100000000stake
${LINKD} add-genesis-account $(${LINKCLI} keys show bob -a) 1000link,100000000stake
${LINKD} add-genesis-account $(${LINKCLI} keys show rinah -a) 1000link,100000000stake
${LINKD} add-genesis-account $(${LINKCLI} keys show sam -a) 1000link,100000000stake
${LINKD} add-genesis-account $(${LINKCLI} keys show evelyn -a) 1000link,100000000stake

${LINKD} --keyring-backend=test gentx --name jack

${LINKD} collect-gentxs

${LINKD} validate-genesis

