#!/usr/bin/env bash
source "./contract_test/testdata/common.sh"


JACK_ADDR="$(./build/linkcli --keyring-backend=test --home /tmp/contract_test/.linkcli keys show jack -a)"

create_only_address () {
  # create address
  ./build/linkcli --keyring-backend=test --home /tmp/contract_test/.linkcli keys add "$1"
  ACTUAL_ADDR="$(./build/linkcli --keyring-backend=test --home /tmp/contract_test/.linkcli keys show "$1" -a)"
}

set_test_address () {
  # create address
  ./build/linkcli --keyring-backend=test --home /tmp/contract_test/.linkcli keys add "$1"
  ACTUAL_ADDR="$(./build/linkcli --keyring-backend=test --home /tmp/contract_test/.linkcli keys show "$1" -a)"

  # register the account
  SEND_TX_HASH=$(./build/linkcli --keyring-backend=test tx send --home ${HOME} ${JACK_ADDR} ${ACTUAL_ADDR} 10link --chain-id ${CHAIN_ID} --yes -b block | awk '/txhash.*/{print $2}')
  echo "Send token: ${SEND_TX_HASH}"
}

# sleeping a whole second between each step is a conservative precaution
sleep 3s

# prepare test files
ALL_MSG_TX='/tmp/contract_test/all_msg_tx.json'
SIGNED_TX='/tmp/contract_test/signed_tx.json'
TMP_TX_RESULT='/tmp/contract_test/tmp_result.txt'
cp client/lcd/static_resources/swagger-ui/swagger.yaml ${SWAGGER}
cp ./contract_test/testdata/all_msg_tx.json ${ALL_MSG_TX}

create_only_address somebody ${REPLACE_SOMEBODY_ADDR}

set_test_address allocator ${REPLACE_ALLOCATOR_ADDR}
set_test_address issuer ${REPLACE_ISSUER_ADDR}
set_test_address returner ${REPLACE_RETURNER_ADDR}

# proxy module not in use as of 2019/2/10
#set_test_address proxy ${REPLACE_PROXY_ADDR}
#set_test_address on_behalf_of ${REPLACE_ON_BEHALF_OF_ADDR}

./contract_test/testdata/replace_symbols.sh

# copy MsgExamples from swagger.yaml to all_msg_tx.json
sed -i.bak -e "s%${REPLACE_MSG_EXAMPLES}%$(yq read -j ${SWAGGER} components.examples.MsgExamples.value)%g" ${ALL_MSG_TX}

# sign transaction that has all messages
./build/linkcli --keyring-backend=test tx sign --home ${HOME} ${ALL_MSG_TX} --from jack --chain-id ${CHAIN_ID} --output-document ${SIGNED_TX}
./build/linkcli --keyring-backend=test tx sign --home ${HOME} ${SIGNED_TX} --append --from allocator --chain-id ${CHAIN_ID}  --output-document ${SIGNED_TX}
# proxy module not in use as of 2019/2/10
#echo ${PASSWORD} | ./build/linkcli tx sign --home ${HOME} ${SIGNED_TX} --append --from on_behalf_of --chain-id ${CHAIN_ID}  --output-document ${SIGNED_TX}
#echo ${PASSWORD} | ./build/linkcli tx sign --home ${HOME} ${SIGNED_TX} --append --from proxy --chain-id ${CHAIN_ID}  --output-document ${SIGNED_TX}

echo "Request: $(cat ${SIGNED_TX})"

# broadcast transaction that has all messages
  ./build/linkcli --keyring-backend=test tx broadcast --home ${HOME} ${SIGNED_TX} --chain-id ${CHAIN_ID} --yes -b block > ${TMP_TX_RESULT}
if [ "$(cat ${TMP_TX_RESULT} | awk '/code:/{print $2}')" -ne "0" ]
then
  echo "ERROR: $(cat ${TMP_TX_RESULT})"
  exit 1
fi
ALL_MSG_TX_HASH=$(cat ${TMP_TX_RESULT}  | awk '/txhash/{print $2}')
echo "All messages have been processed: ${ALL_MSG_TX_HASH}"
echo ${ALL_MSG_TX_HASH} > '/tmp/contract_test/all_msg_tx_hash.txt'

# change tx hash to check message format in dredd test (/txs/{hash})
sed -i.bak -e "s/${REPLACE_TX_HASH}/${ALL_MSG_TX_HASH}/g" ${SWAGGER}
echo "Replaced dummy with actual ALL_MSG_TX_HASH ${ALL_MSG_TX_HASH}"

sleep 3s
