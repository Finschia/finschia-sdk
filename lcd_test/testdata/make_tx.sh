#!/usr/bin/env bash
source "./lcd_test/testdata/common.sh"

SEND_TX_HASH=$(echo ${PASSWORD} | ./build/linkcli tx send jack $(./build/linkcli --home /tmp/contract_tests/.linkcli keys show alice -a) 100link --home /tmp/contract_tests/.linkcli --chain-id lcd --yes  | awk '/txhash.*/{print $2}')
echo ${SEND_TX_HASH}
