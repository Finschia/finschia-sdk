#!/usr/bin/env bash

dredd
dredd_result=$?

echo "Terminating generate_tx_iteratively.sh"
pkill -f ./contract_test/testdata/generate_tx_iteratively.sh
sleep 1

./contract_test/testdata/stop_dredd_test.sh

exit ${dredd_result}