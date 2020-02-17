#!/usr/bin/env bash
set -e
PASSWORD=1234567890

# proxy module not in use as of 2019/2/10
#REPLACE_PROXY_ADDR="link1pxjrsfqam3nuf75v2mh4atshaze2nlmjnqc99m"
#REPLACE_ON_BEHALF_OF_ADDR="link1n5vsmtppfs4g4ue6zvuj7eqt5ajteu3qrvt0h2"

# safetybox module not in use as of 2019/2/14
#REPLACE_OPERATOR_ADDR="linkoperatormpp92x9hyzz9wrgf94r6j9h5f06pxxv"
REPLACE_ALLOCATOR_ADDR="linkallocatorpp92x9hyzz9wrgf94r6j9h5f06pxxv"
#REPLACE_ISSUER_ADDR="linkissuetormpp92x9hyzz9wrgf94r6j9h5f06pxxv"
#REPLACE_RETURNER_ADDR="linkreturnormpp92x9hyzz9wrgf94r6j9h5f06pxxv"

REPLACE_JACK_ADDR="link16xyempempp92x9hyzz9wrgf94r6j9h5f06pxxv"
REPLACE_SOMEBODY_ADDR="linksomebodympp92x9hyzz9wrgf94r6j9h5f06pxxv"
REPLACE_TX_HASH="BCBE20E8D46758B96AE5883B792858296AC06E51435490FBDCAE25A72B3CC76B"
REPLACE_TOKEN_SYMBOL="conyxxv"
REPLACE_COLLECTION_SYMBOL="con2xxv"
REPLACE_NFT_SYMBOL="con3xxv"
REPLACE_MSG_EXAMPLES="\"MsgExamples\""
CHAIN_ID="lcd"
HOME="/tmp/contract_test/.linkcli"
SWAGGER='/tmp/contract_test/swagger.yaml'
