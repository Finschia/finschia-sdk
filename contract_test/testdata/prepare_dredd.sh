#!/usr/bin/env bash
source "./contract_test/testdata/common.sh"

echo "Install dredd"
./contrib/get_node.sh
./build/contract_test_hook -port 61322 &