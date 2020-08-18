#!/usr/bin/env bash
set -e

docker run --rm --name master --network link_localnet -v $(pwd)/contrib/load_test/:/link-load-tester/contrib/load_test/ link-load-tester link-load-tester prepare || $(docker stop $(docker ps -f name=slave -q); false)
docker run --rm --name master --network link_localnet -v $(pwd)/contrib/load_test/:/link-load-tester/contrib/load_test/ link-load-tester link-load-tester start || $(docker stop $(docker ps -f name=slave -q); false)

echo "Wait until query is possible"
sleep 10
while [ "$(curl -s -o /dev/null -w ''%{http_code}'' http://localhost:1317/blocks/latest)" != "200" ];
do sleep 1;
done;
docker run --rm --name master --network link_localnet -v $(pwd)/contrib/load_test/:/link-load-tester/contrib/load_test/ link-load-tester link-load-tester report || $(docker stop $(docker ps -f name=slave -q); false)

docker stop $(docker ps -f name=slave -q)
rm $(pwd)/contrib/load_test/result_data.json && true
rm $(pwd)/contrib/load_test/Latency.png && true
rm $(pwd)/contrib/load_test/TPS.png && true
rm $(pwd)/contrib/load_test/test_params.json && true
