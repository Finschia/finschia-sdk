#!/usr/bin/env bash
source "./contract_test/testdata/common.sh"

pkill linkd || true
rm -rf /tmp/contract_test
mkdir /tmp/contract_test
./build/linkd init --home /tmp/contract_test/.linkd --chain-id lcd contract-test
sleep 2s
echo ${PASSWORD} | echo ${PASSWORD} | ./build/linkcli --home /tmp/contract_test/.linkcli keys add jack
echo ${PASSWORD} | echo ${PASSWORD} | ./build/linkcli --home /tmp/contract_test/.linkcli keys add alice
./build/linkd --home /tmp/contract_test/.linkd add-genesis-account $(./build/linkcli --home /tmp/contract_test/.linkcli keys show jack -a) 100link,100000000stake
sleep 3s
echo ${PASSWORD} | ./build/linkd --home /tmp/contract_test/.linkd gentx --name jack --home-client /tmp/contract_test/.linkcli --amount 100000000stake
sleep 3s
./build/linkd --home /tmp/contract_test/.linkd  collect-gentxs
./build/linkd --home /tmp/contract_test/.linkd  validate-genesis