#!/usr/bin/env bash
source "./lcd_test/testdata/common.sh"


JACK_ADDR="$(./build/linkcli --home /tmp/contract_tests/.linkcli keys show jack -a)"

set_test_address () {
  # create address
  echo ${PASSWORD} | echo ${PASSWORD} | ./build/linkcli --home /tmp/contract_tests/.linkcli keys add "$1"
  ACTUAL_ADDR="$(./build/linkcli --home /tmp/contract_tests/.linkcli keys show "$1" -a)"

  # register the account
  SEND_TX_HASH=$(echo ${PASSWORD} | ./build/linkcli tx send --home ${HOME} ${JACK_ADDR} ${ACTUAL_ADDR} 10link --chain-id ${CHAIN_ID} --yes -b block | awk '/txhash.*/{print $2}')
  echo "Send token: ${SEND_TX_HASH}"
}

# sleeping a whole second between each step is a conservative precaution
sleep 3s

# prepare test files
ALL_MSG_TX='/tmp/contract_tests/all_msg_tx.json'
SIGNED_TX='/tmp/contract_tests/signed_tx.json'
TMP_TX_RESULT='/tmp/contract_tests/tmp_result.txt'
cp client/lcd/swagger-ui/swagger.yaml ${SWAGGER}
cp ./lcd_test/testdata/all_msg_tx.json ${ALL_MSG_TX}

set_test_address operator ${REPLACE_OPERATOR_ADDR}
set_test_address allocator ${REPLACE_ALLOCATOR_ADDR}
set_test_address issuer ${REPLACE_ISSUER_ADDR}
set_test_address returner ${REPLACE_RETURNER_ADDR}

./lcd_test/testdata/replace_symbols.sh

# copy MsgExamples from swagger.yaml to all_msg_tx.json
sed -i.bak -e "s%${REPLACE_MSG_EXAMPLES}%$(yq read -j ${SWAGGER} components.examples.MsgExamples.value)%g" ${ALL_MSG_TX}

# sign transaction that has all messages
echo ${PASSWORD} | ./build/linkcli tx sign --home ${HOME} ${ALL_MSG_TX} --from jack --chain-id ${CHAIN_ID} --output-document ${SIGNED_TX}
echo ${PASSWORD} | ./build/linkcli tx sign --home ${HOME} ${SIGNED_TX} --append --from operator --chain-id ${CHAIN_ID}  --output-document ${SIGNED_TX}
echo ${PASSWORD} | ./build/linkcli tx sign --home ${HOME} ${SIGNED_TX} --append --from allocator --chain-id ${CHAIN_ID}  --output-document ${SIGNED_TX}
echo ${PASSWORD} | ./build/linkcli tx sign --home ${HOME} ${SIGNED_TX} --append --from issuer --chain-id ${CHAIN_ID}  --output-document ${SIGNED_TX}
echo ${PASSWORD} | ./build/linkcli tx sign --home ${HOME} ${SIGNED_TX} --append --from returner --chain-id ${CHAIN_ID}  --output-document ${SIGNED_TX}
echo "Request: $(cat ${SIGNED_TX})"

# broadcast transaction that has all messages
./build/linkcli tx broadcast --home ${HOME} ${SIGNED_TX} --chain-id ${CHAIN_ID} --yes -b block > ${TMP_TX_RESULT}
if [ "$(cat ${TMP_TX_RESULT} | awk '/code:/{print $2}')" -ne "0" ]
then
  echo "ERROR: $(cat ${TMP_TX_RESULT})"
  exit 1
fi
ALL_MSG_TX_HASH=$(cat ${TMP_TX_RESULT}  | awk '/txhash/{print $2}')
echo "All messages have been processed: ${ALL_MSG_TX_HASH}"
echo ${ALL_MSG_TX_HASH} > '/tmp/contract_tests/all_msg_tx_hash.txt'

# change tx hash to check message format in dredd test (/txs/{hash})
sed -i.bak -e "s/${REPLACE_TX_HASH}/${ALL_MSG_TX_HASH}/g" ${SWAGGER}
echo "Replaced dummy with actual ALL_MSG_TX_HASH ${ALL_MSG_TX_HASH}"

sleep 3s