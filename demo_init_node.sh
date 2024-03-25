#!/usr/bin/env bash
# This script is intended for demo purposes, so please do not use it in production.

set -ex

mode="mainnet"
FNSAD=${FNSAD:-simd}

# initialize
rm -rf ~/.simapp

# Initialize configuration files and genesis file
# moniker is the name of your node
${FNSAD} init solo --chain-id=sim

# Please do not use the TEST_MNEMONIC for production purpose
TEST_MNEMONIC="mind flame tobacco sense move hammer drift crime ring globe art gaze cinnamon helmet cruise special produce notable negative wait path scrap recall have"

${FNSAD} keys add jack --keyring-backend=test --recover --account=0 <<< ${TEST_MNEMONIC}
${FNSAD} keys add alice --keyring-backend=test --recover --account=1 <<< ${TEST_MNEMONIC}
${FNSAD} keys add bob --keyring-backend=test --recover --account=2 <<< ${TEST_MNEMONIC}
${FNSAD} keys add rinah --keyring-backend=test --recover --account=3 <<< ${TEST_MNEMONIC}
${FNSAD} keys add sam --keyring-backend=test --recover --account=4 <<< ${TEST_MNEMONIC}
${FNSAD} keys add evelyn --keyring-backend=test --recover --account=5 <<< ${TEST_MNEMONIC}

${FNSAD} add-genesis-account $(${FNSAD} keys show jack -a --keyring-backend=test) 1000link,1000000000000stake
${FNSAD} add-genesis-account $(${FNSAD} keys show alice -a --keyring-backend=test) 1000link,1000000000000stake
${FNSAD} add-genesis-account $(${FNSAD} keys show bob -a --keyring-backend=test) 1000link,1000000000000stake,1000000000cony
${FNSAD} add-genesis-account $(${FNSAD} keys show rinah -a --keyring-backend=test) 1000link,1000000000000stake
${FNSAD} add-genesis-account $(${FNSAD} keys show sam -a --keyring-backend=test) 1000link,1000000000000stake
${FNSAD} add-genesis-account $(${FNSAD} keys show evelyn -a --keyring-backend=test) 10000000000000link,1000000000000stake

${FNSAD} gentx jack 100000000stake --keyring-backend=test --chain-id=sim

${FNSAD} collect-gentxs

${FNSAD} validate-genesis

# enable unsafe cors and rest server
sed -i -e s/"enable = false"/"enable = true"/g ${HOME}/.simapp/config/app.toml
sed -i -e s/"enabled-unsafe-cors = false"/"enabled-unsafe-cors = true"/g ${HOME}/.simapp/config/app.toml
sed -i -e s/"enable-unsafe-cors = false"/"enable-unsafe-cors = true"/g ${HOME}/.simapp/config/app.toml
sed -i -e s/"cors_allowed_origins = \[\]"/"cors_allowed_origins = \[\"*\"\]"/g ${HOME}/.simapp/config/config.toml

docker run -i -p 26656:26656 -p 26657:26657 -p 1317:1317 -v ${HOME}/.simapp:/root/.simapp simapp simd start --log_level *:debug --rpc.laddr=tcp://0.0.0.0:26657 --p2p.laddr=tcp://0.0.0.0:26656
