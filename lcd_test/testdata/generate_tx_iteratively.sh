#!/usr/bin/env bash
source "./lcd_test/testdata/common.sh"
ADDR="$(./build/linkcli --home /tmp/contract_tests/.linkcli keys show jack -a)"
memo="."
while true
do
  echo ${PASSWORD} | ./build/linkcli tx send --home ${HOME} ${ADDR} ${ADDR} 1link -b async --chain-id ${CHAIN_ID} --yes --memo ${memo} > /dev/null
	memo=${memo}.
	sleep 0.1
done