#!/usr/bin/env sh

num_txs=$1

start_time="$(date -u +%s)"
for i in $(seq 0 $num_txs)
do
    linkwasmcli tx broadcast tx$i.json
done
end_time="$(date -u +%s)"
elapsed="$(($end_time-$start_time))"
echo "$elapsed seconds elapsed for firing txs"
