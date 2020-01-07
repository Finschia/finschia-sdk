#!/usr/bin/env bash
source "./lcd_test/testdata/common.sh"

DYNAMIC_JACK_ADDR="$(./build/linkcli --home /tmp/contract_tests/.linkcli keys show jack -a)"
VALIDATOR="linkvaloper1qe6nfrpqsgj4g0sgg3v0y29sa7r275yzdy6tc5"
AMOUNT="1link"
CHAIN_ID="lcd"
PROPOSALID="2"
HOME="/tmp/contract_tests/.linkcli"
SWAGGER='/tmp/contract_tests/swagger.yaml'

# sleeping a whole second between each step is a conservative precaution
cp client/lcd/swagger-ui/swagger_OSS_2_0.yaml /tmp/contract_tests/swagger.yaml
## start setup transactions
sed -i.bak -e "s/${REPLACE_ADDR}/${DYNAMIC_JACK_ADDR}/g" ${SWAGGER}
sed -i.bak -e "s/${REPLACE_ADDR}/${DYNAMIC_JACK_ADDR}/g" ${SWAGGER}
sleep 3s
SEND_TX_HASH=$(echo ${PASSWORD} | ./build/linkcli tx send --home ${HOME} ${DYNAMIC_JACK_ADDR} ${DYNAMIC_JACK_ADDR} ${AMOUNT} --chain-id ${CHAIN_ID} --yes -b block | awk '/txhash.*/{print $2}')
sleep 3s
PUBLISH_TX_HASH=$(echo ${PASSWORD} | ./build/linkcli --home ${HOME} tx token issue jack cony lcdtoken --decimals=0 --mintable=true --total-supply=10000 --chain-id ${CHAIN_ID} --yes -b block | awk '/txhash.*/{print $2}')
sed -i.bak -e "s/${REPLACE_PUBLISH_TX_HASH}/${PUBLISH_TX_HASH}/g" ${SWAGGER}
echo "Replaced dummy with actual PUBLISH_TX_HASH ${PUBLISH_TX_HASH}"
sleep 3s
echo ${PASSWORD} | ./build/linkcli tx staking unbond --home ${HOME} --from ${DYNAMIC_JACK_ADDR} ${VALIDATOR} 100stake --yes --chain-id ${CHAIN_ID}
