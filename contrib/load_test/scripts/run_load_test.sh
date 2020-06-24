#!/usr/bin/env bash
set -e

docker run --rm --name master --network link_localnet -v $(pwd)/contrib/load_test/:/link-load-tester/contrib/load_test/ link-load-tester link-load-tester prepare || docker stop $(docker ps -f name=slave -q); false
docker run --rm --name master --network link_localnet -v $(pwd)/contrib/load_test/:/link-load-tester/contrib/load_test/ link-load-tester link-load-tester start || docker stop $(docker ps -f name=slave -q); false
docker stop $(docker ps -f name=slave -q)
