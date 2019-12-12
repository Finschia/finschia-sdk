#!/usr/bin/make -f

# Run a 4-node testnet locally
testnet-start: testnet-stop
	@if ! [ -f build/node0/linkd/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/root:Z line/link linkd testnet --v 4 -o . --starting-ip-address 192.168.10.2 ; fi
	docker-compose -f $(CURDIR)/docker-compose.yml up -d

# Stop testnet
testnet-stop:
	docker-compose down

testnet-test:
	${CURDIR}/localnet-blocks-test.sh 40 5 10 localhost

clean:
	rm -rf $(CURDIR)/build

.PHONY: start-testnet stop-testnet clean
