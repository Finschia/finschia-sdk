#!/usr/bin/env bash
source "./contract_test/testdata/common.sh"
ADDR="$(./build/linkcli --keyring-backend=test --home /tmp/contract_test/.linkcli keys show jack -a)"
memo="."
while true
do
  ./build/linkcli --keyring-backend=test tx send --home ${HOME} ${ADDR} ${ADDR} 1link -b async --chain-id ${CHAIN_ID} --yes --memo ${memo} > /dev/null
	memo=${memo}.
	sleep 0.1
done
