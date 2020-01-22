#!/usr/bin/env bash

dredd
dredd_result=$?

echo "Terminating generate_tx_iteratively.sh"
pkill -f ./lcd_test/testdata/generate_tx_iteratively.sh
sleep 1

echo "Terminating linkcli"
pkill linkcli

echo "Terminating linkd"
pkill -9 linkd

exit ${dredd_result}