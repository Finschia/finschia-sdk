#!/usr/bin/env sh
rm -rf ~/.linkwasmd ~/.linkwasmcli
 
linkwasmcli config chain-id perf
linkwasmcli config output json
linkwasmcli config indent true
linkwasmcli config trust-node true
linkwasmcli config keyring-backend test

linkwasmcli keys add jack

linkwasmd init solo --chain-id perf
linkwasmd add-genesis-account $(linkwasmcli keys show jack -a) 100000000stake
linkwasmd --keyring-backend=test gentx --name jack
linkwasmd collect-gentxs
linkwasmd validate-genesis
