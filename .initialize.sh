#!/usr/bin/env bash
set -ex

mode="mainnet"

if [[ $1 == "docker" ]]
then
    if [[ $2 == "testnet" ]]
    then
        mode="testnet"
    fi
    LINKCLI="docker run -i --net=host -v ${HOME}/.linkd:/root/.linkd -v ${HOME}/.linkcli:/root/.linkcli line/link linkcli"
    LINKD="docker run -i -p 26656:26656 -p 26657:26657 -v ${HOME}/.linkd:/root/.linkd -v ${HOME}/.linkcli:/root/.linkcli line/link linkd"
elif [[ $1 == "testnet" ]]
then
    mode="testnet"
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

# configure for testnet
if [[ ${mode} == "testnet" ]]
then
    if [[ $1 == "docker" ]]
    then
        docker run -i --net=host -v ${HOME}/.linkd:/root/.linkd -v ${HOME}/.linkcli:/root/.linkcli line/link sh -c "echo 'testnet = true' >> /root/.linkcli/config/config.toml"
        docker run -i -p 26656:26656 -p 26657:26657 -v ${HOME}/.linkd:/root/.linkd -v ${HOME}/.linkcli:/root/.linkcli line/link sh -c "echo 'testnet = true' >> /root/.linkd/config/app.toml"
    else
        echo "testnet = true" >> ~/.linkcli/config/config.toml
        echo "testnet = true" >> ~/.linkd/config/app.toml
    fi
fi

${LINKCLI} keys add jack
${LINKCLI} keys add alice
${LINKCLI} keys add bob
${LINKCLI} keys add rinah
${LINKCLI} keys add sam
${LINKCLI} keys add evelyn

if [[ ${mode} == "testnet" ]]
then
   ${LINKD} add-genesis-account tlink15la35q37j2dcg427kfy4el2l0r227xwhc2v3lg 9223372036854775807link,1stake
else
   ${LINKD} add-genesis-account link15la35q37j2dcg427kfy4el2l0r227xwhuaapxd 9223372036854775807link,1stake
fi
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

