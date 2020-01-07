#!/usr/bin/env bash
source "./lcd_test/testdata/common.sh"

pkill linkd || true
rm -rf /tmp/contract_tests
mkdir /tmp/contract_tests
./build/linkd init --home /tmp/contract_tests/.linkd --chain-id lcd contract-tests
sleep 2s
echo ${PASSWORD} | echo ${PASSWORD} | ./build/linkcli --home /tmp/contract_tests/.linkcli keys add jack
echo ${PASSWORD} | echo ${PASSWORD} | ./build/linkcli --home /tmp/contract_tests/.linkcli keys add alice
./build/linkd --home /tmp/contract_tests/.linkd add-genesis-account $(./build/linkcli --home /tmp/contract_tests/.linkcli keys show jack -a) 100link,100000000stake
sleep 3s
echo ${PASSWORD} | ./build/linkd --home /tmp/contract_tests/.linkd gentx --name jack --home-client /tmp/contract_tests/.linkcli --amount 100000000stake
sleep 3s
./build/linkd --home /tmp/contract_tests/.linkd  collect-gentxs
./build/linkd --home /tmp/contract_tests/.linkd  validate-genesis