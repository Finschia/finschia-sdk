#!/usr/bin/env sh
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

PASSWORD="1234567890"

# initialize
rm -rf ~/.linkd ~/.linkcli

# Configure your CLI to eliminate need for chain-id flag
${LINKCLI} config chain-id link
${LINKCLI} config output json
${LINKCLI} config indent true
${LINKCLI} config trust-node true

# Initialize configuration files and genesis file
# moniker is the name of your node
${LINKD} init solo --chain-id link

# configure for testnet
if [[ $mode == "testnet" ]]
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

echo ${PASSWORD} | echo ${PASSWORD} | ${LINKCLI} keys add jack
echo ${PASSWORD} | echo ${PASSWORD} | ${LINKCLI} keys add alice
echo ${PASSWORD} | echo ${PASSWORD} | ${LINKCLI} keys add bob
echo ${PASSWORD} | echo ${PASSWORD} | ${LINKCLI} keys add rinah
echo ${PASSWORD} | echo ${PASSWORD} | ${LINKCLI} keys add sam
echo ${PASSWORD} | echo ${PASSWORD} | ${LINKCLI} keys add evelyn

# Add both accounts, with coins to the genesis file
${LINKD} add-genesis-account $(${LINKCLI} keys show jack -a) 1000link,100000000stake
${LINKD} add-genesis-account $(${LINKCLI} keys show alice -a) 1000link,100000000stake
${LINKD} add-genesis-account $(${LINKCLI} keys show bob -a) 1000link,100000000stake
${LINKD} add-genesis-account $(${LINKCLI} keys show rinah -a) 1000link,100000000stake
${LINKD} add-genesis-account $(${LINKCLI} keys show sam -a) 1000link,100000000stake
${LINKD} add-genesis-account $(${LINKCLI} keys show evelyn -a) 1000link,100000000stake

echo ${PASSWORD} | ${LINKD} gentx --name jack

${LINKD} collect-gentxs

${LINKD} validate-genesis

# ${LINKD} start --log_level *:debug --rpc.laddr=tcp://0.0.0.0:26657 --p2p.laddr=tcp://0.0.0.0:26656
