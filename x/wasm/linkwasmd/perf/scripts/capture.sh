#!/usr/bin/env sh
num_txs=$1

# capture with running binary
# log format:
# t:7 => tx count of the block
# s:2020-11-16 11:11:00.000 => time before execution
# e:2020-11-16 11:11:02.000 => time after execution
perf record --call-graph dwarf -F 500 -D 10000 -o linkwasmd.dwarf.perf \
linkwasmd start --log_level="state:info,consensus:info,*:error" \
| awk '{if (/Finalizing commit/) {print "t:" substr($13,3,4) "\ns:" substr($1,3,10) " " substr($1,14,12)} else if(/Committed state/) {print "e:" substr($1,3,10) " " substr($1,14,12)}}'  \
> linkwasmd.log &

echo "start linkwasmd"

sleep 7s # wait for booting

echo "fire txs"
./perf/scripts/fire.sh $num_txs

sleep 6s

echo "terminate"
killall -SIGTERM linkwasmd

sleep 5s
mkdir reports
perf script -i linkwasmd.dwarf.perf | perf/tools/FlameGraph/stackcollapse-perf.pl | perf/tools/FlameGraph/flamegraph.pl > reports/linkwasmd.dwarf.perf.svg

./perf/scripts/report.sh linkwasmd.log reports/report.txt
