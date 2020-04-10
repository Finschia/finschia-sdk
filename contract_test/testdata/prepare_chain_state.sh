#!/usr/bin/env bash
source "./contract_test/testdata/common.sh"

pkill linkd || true
rm -rf /tmp/contract_test
mkdir /tmp/contract_test
./build/linkd init --home /tmp/contract_test/.linkd --chain-id lcd contract-test
sleep 2s

./build/linkcli --keyring-backend=test --home /tmp/contract_test/.linkcli keys add jack
./build/linkcli --keyring-backend=test --home /tmp/contract_test/.linkcli keys add alice
sleep 3s
./build/linkd --keyring-backend=test --home /tmp/contract_test/.linkd add-genesis-account $(./build/linkcli --keyring-backend=test --home /tmp/contract_test/.linkcli keys show jack -a) 100link,100000000stake
sleep 3s
./build/linkd --keyring-backend=test --home /tmp/contract_test/.linkd gentx --name jack --home-client /tmp/contract_test/.linkcli --amount 100000000stake
sleep 3s
./build/linkd --home /tmp/contract_test/.linkd  collect-gentxs
./build/linkd --home /tmp/contract_test/.linkd  validate-genesis
