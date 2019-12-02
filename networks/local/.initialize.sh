#!/usr/bin/env sh

rm -rf ./build

LINKCLI=${LINKCLI:-linkcli}
LINKD=${LINKD:-linkd}
LINKCLIHOME=./build/linkcli
LINKDHOME=linkd/config
GENTXHOME=./build/gentxs
HOMEPREFIX_VALIDATOR=./build/node_validator
HOMEPREFIX_SENTRY=./build/node_sentry
HOMEPREFIX_SEED=./build/node_seed

PASSWORD="1234567890"

NUMOFVALIDATORS=4
SET=$(seq 0 `expr ${NUMOFVALIDATORS} - 1`)


# Create Accounts
for i in $SET
do
    echo ${PASSWORD} | echo ${PASSWORD} | ${LINKCLI} --home ${LINKCLIHOME} keys add account${i}
done

# Initialize Validator, Sentry, Seed nodes
for i in ${SET}
do
    ${LINKD} --home ${HOMEPREFIX_VALIDATOR}${i}/linkd init node${i} --chain-id testnet1000
    ${LINKD} --home ${HOMEPREFIX_SENTRY}${i}/linkd init node${i} --chain-id testnet1000
done
${LINKD} --home ${HOMEPREFIX_SEED}/linkd init node_seed --chain-id testnet1000

# Add Genesis Account for all validators
for i in ${SET}
do
    for j in ${SET}
    do
        ${LINKD} --home ${HOMEPREFIX_VALIDATOR}${i}/linkd add-genesis-account $(${LINKCLI} --home ${LINKCLIHOME} keys show -a account${j}) 100000000stake
    done
done

# Create Validator And Collect gentxs
for i in ${SET}
do
    mkdir -p ${GENTXHOME}
    echo ${PASSWORD} | echo ${PASSWORD} | ${LINKD} --home ${HOMEPREFIX_VALIDATOR}${i}/linkd --home-client ${LINKCLIHOME} gentx --name account${i} --output-document ${GENTXHOME}/node${i}.json --ip 192.168.10.`expr ${i} + 2`  --node-id $(${LINKD} --home ${HOMEPREFIX_SENTRY}${i}/linkd tendermint show-node-id)
done

${LINKD} --home ${HOMEPREFIX_VALIDATOR}0/linkd collect-gentxs --gentx-dir ${GENTXHOME}


# Copy genesis.json to all nodes
cp ${HOMEPREFIX_VALIDATOR}0/LINKDHOME/genesis.json ./build/genesis.json
for i in ${SET}
do
    cp ./build/genesis.json ${HOMEPREFIX_VALIDATOR}${i}/LINKDHOME/genesis.json
    cp ./build/genesis.json ${HOMEPREFIX_SENTRY}${i}/LINKDHOME/genesis.json
done
cp ./build/genesis.json ${HOMEPREFIX_SEED}/LINKDHOME/genesis.json


# Setup configuration according to node types

# SEED
seed_id=$(${LINKD} --home ${HOMEPREFIX_SEED}/linkd tendermint show-node-id)
seed_ip=192.168.10.200:26656

sed -i -e "s/external_address.*/external_address = \"tcp\:\/\/${seed_ip}\"/g" ${HOMEPREFIX_SEED}/LINKDHOME/config.toml
sed -i -e "s/addr_book_strict.*/addr_book_strict = \"false\"/g" ${HOMEPREFIX_SEED}/LINKDHOME/config.toml
sed -i -e "s/pex.*/pex = \"true\"/g" ${HOMEPREFIX_SEED}/LINKDHOME/config.toml
sed -i -e "s/seed_mode.*/seed_mode = \"true\"/g" ${HOMEPREFIX_SEED}/LINKDHOME/config.toml
sed -i -e "s/persistent_peers.*/persistent_peers = \"\"/g" ${HOMEPREFIX_SEED}/LINKDHOME/config.toml
for i in ${SET}
do
    val_id=$(${LINKD} --home ${HOMEPREFIX_VALIDATOR}${i}/linkd tendermint show-node-id)
    val_ip=192.168.10.`expr ${i} + 2 + 100`:26656
    sentry_id=$(${LINKD} --home ${HOMEPREFIX_SENTRY}${i}/linkd tendermint show-node-id)
    sentry_ip=192.168.10.`expr ${i} + 2`:26656

    # VALIDATOR
    sed -i -e "s/external_address.*/external_address = \"tcp\:\/\/${val_ip}\"/g" ${HOMEPREFIX_VALIDATOR}${i}/LINKDHOME/config.toml
    sed -i -e "s/addr_book_strict.*/addr_book_strict = \"false\"/g" ${HOMEPREFIX_VALIDATOR}${i}/LINKDHOME/config.toml
    sed -i -e "s/pex.*/pex = \"false\"/g" ${HOMEPREFIX_VALIDATOR}${i}/LINKDHOME/config.toml
    sed -i -e "s/seed_mode.*/seed_mode = \"false\"/g" ${HOMEPREFIX_VALIDATOR}${i}/LINKDHOME/config.toml
    sed -i -e "s/persistent_peers.*/persistent_peers = \"${sentry_id}@${sentry_ip}\"/g" ${HOMEPREFIX_VALIDATOR}${i}/LINKDHOME/config.toml

    # SENTRY
    sed -i -e "s/external_address.*/external_address = \"tcp\:\/\/${sentry_ip}\"/g" ${HOMEPREFIX_SENTRY}${i}/LINKDHOME/config.toml
    sed -i -e "s/addr_book_strict.*/addr_book_strict = \"false\"/g" ${HOMEPREFIX_SENTRY}${i}/LINKDHOME/config.toml
    sed -i -e "s/pex.*/pex = \"true\"/g" ${HOMEPREFIX_SENTRY}${i}/LINKDHOME/config.toml
    sed -i -e "s/seed_mode.*/seed_mode = \"false\"/g" ${HOMEPREFIX_SENTRY}${i}/LINKDHOME/config.toml
    sed -i -e "s/persistent_peers.*/persistent_peers = \"\"/g" ${HOMEPREFIX_SENTRY}${i}/LINKDHOME/config.toml
    sed -i -e "s/private_peer_ids.*/private_peer_ids = \"${val_id}@tcp\:\/\/${val_ip}\"/g" ${HOMEPREFIX_SENTRY}${i}/LINKDHOME/config.toml
    sed -i -e "s/seeds.*/seeds = \"${seed_id}@${seed_ip}\"/g" ${HOMEPREFIX_SENTRY}${i}/LINKDHOME/config.toml
done

