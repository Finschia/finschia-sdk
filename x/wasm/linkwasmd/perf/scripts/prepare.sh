#!/usr/bin/env sh
wasm_path=$1
num_txs=$2

linkwasmd start > /dev/null &

# creates account as much as txs that will be fired
echo "create accounts"
for i in $(seq 0 $num_txs)
do
    linkwasmcli keys add $i
done

code_id=$(linkwasmcli tx wasm store $wasm_path --from jack --gas 200000000 -y  --broadcast-mode=block \
| jq -cr '.logs[0].events[0].attributes | map(select(.key=="code_id"))[0].value')

echo "code_id: $code_id"

msg=$(jq -nc \
--arg address $(linkwasmcli keys show jack -a) \
--arg amount 10000000000 \
'{decimals:6,initial_balances:[{address:$address,amount:$amount}],name:"TKN",symbol:"TKN"}')

contract_address=$(linkwasmcli tx wasm instantiate $code_id $msg --label TKN --from jack --gas 200000000 --broadcast-mode=block -y \
| jq -cr '.logs[0].events[0].attributes | map(select(.key=="contract_address"))[0].value')

echo "contract_address: $contract_address"

account_number=$(linkwasmcli query account $(linkwasmcli keys show jack -a) | jq -cr '.value.account_number')
sequence=$(linkwasmcli query account $(linkwasmcli keys show jack -a) | jq -cr '.value.sequence')

start_time="$(date -u +%s)"
for i in $(seq 0 $num_txs)
do
    msg=$(jq -nc \
    --arg amount 1 \
    --arg recipient $(linkwasmcli keys show $i -a) \
    '{transfer:{amount:$amount,recipient:$recipient}}')

    linkwasmcli tx wasm execute $contract_address $msg --from $(linkwasmcli keys show jack -a) --gas 200000000 --generate-only > tx$i.unsigned.json
    linkwasmcli tx sign tx$i.unsigned.json --from jack --offline -a $account_number -s $(expr $sequence + $i) > tx$i.json &
done
end_time="$(date -u +%s)"
elapsed="$(($end_time-$start_time))"
echo "$elapsed seconds elapsed for generating txs"

killall -SIGTERM linkwasmd
