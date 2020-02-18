#!/usr/bin/env bash
source "./contract_test/testdata/common.sh"

echo "Install node"
./contrib/get_node.sh
echo "Validate swagger docs"
swagger-cli validate ./client/lcd/swagger-ui/swagger.yaml
echo "Run contract test hook server"
./build/contract_test_hook -port 61322 &