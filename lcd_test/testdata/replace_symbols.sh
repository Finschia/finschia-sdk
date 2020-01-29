#!/usr/bin/env bash
source "./lcd_test/testdata/common.sh"

JACK_ADDR="$(./build/linkcli --home /tmp/contract_tests/.linkcli keys show jack -a)"
REPLACE_OPTION=$1

replace_address () {
  ACTUAL_ADDR="$(./build/linkcli --home /tmp/contract_tests/.linkcli keys show "$1" -a)"
  sed -i.bak -e "s/$2/${ACTUAL_ADDR}/g" ${SWAGGER}
  echo "Replaced dummy with actual ADDR of $1 : ${ACTUAL_ADDR}"
}

replace_token_symbol () {
  sed -i.bak -e "s/$2/$1/g" ${SWAGGER}
  echo "Replaced dummy with actual TOKEN_SYMBOL of $3 : $1"
}

replace_address jack ${REPLACE_JACK_ADDR}
replace_address operator ${REPLACE_OPERATOR_ADDR}
replace_address allocator ${REPLACE_ALLOCATOR_ADDR}
replace_address issuer ${REPLACE_ISSUER_ADDR}
replace_address returner ${REPLACE_RETURNER_ADDR}
replace_address proxy ${REPLACE_PROXY_ADDR}
replace_address on_behalf_of ${REPLACE_ON_BEHALF_OF_ADDR}

replace_token_symbol "alcd"${JACK_ADDR:40} ${REPLACE_TOKEN_SYMBOL} FT
replace_token_symbol "blcd"${JACK_ADDR:40} ${REPLACE_COLLECTION_SYMBOL} FT_COLLECTION
replace_token_symbol "clcd"${JACK_ADDR:40} ${REPLACE_NFT_SYMBOL} NFT
replace_token_symbol "dlcd"${JACK_ADDR:40} ${REPLACE_NFT_COLLECTION_SYMBOL} NFT_COLLECTION

if [ "${REPLACE_OPTION}" == "--replace_tx_hash" ]
then
 ALL_MSG_TX_HASH=$(cat '/tmp/contract_tests/all_msg_tx_hash.txt')
 sed -i.bak -e "s/${REPLACE_TX_HASH}/${ALL_MSG_TX_HASH}/g" ${SWAGGER}
 echo "Replaced dummy with actual ALL_MSG_TX_HASH ${ALL_MSG_TX_HASH}"
fi
