#!/usr/bin/env sh
set -ex

if [[ $1 == "docker" ]]
then
    LINKCLI="docker run -i --net=host -v ${HOME}/.linkd:/root/.linkd -v ${HOME}/.linkcli:/root/.linkcli line/link linkcli"
    LINKD="docker run -i -p 26656:26656 -p 26657:26657 -v ${HOME}/.linkd:/root/.linkd -v ${HOME}/.linkcli:/root/.linkcli line/link linkd"
fi

LINKCLI=${LINKCLI:-linkcli}
LINKD=${LINKD:-linkd}

# initialize
rm -rf ~/.linkd ~/.linkcli

# Configure your CLI to eliminate need for chain-id flag
${LINKCLI} config chain-id link
${LINKCLI} config output json
${LINKCLI} config indent true
${LINKCLI} config trust-node true
${LINKCLI} config keyring-backend test

# Initialize configuration files and genesis file
# moniker is the name of your node
${LINKD} init solo --chain-id link

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

# ${LINKD} start --log_level *:debug --rpc.laddr=tcp://0.0.0.0:26657 --p2p.laddr=tcp://0.0.0.0:26656
