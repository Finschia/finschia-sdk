# LINK Network Version2

[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://www.repostatus.org/badges/latest/wip.svg)](https://www.repostatus.org/#wip)

This repository hosts `LINK`, alternative implementation of the LINK Network.

**Node**: Requires [Go 1.12+](https://golang.org/dl/)

**Warnings**: Initial development is in progress, but there has not yet been a stable.

## Quick Start

### Prerequisite
```
make get-tools                  # install tools
```
### Build & Install LINK
```
make install                    # build and install binaries
```

### Test
```
make check-unit                 # unit test
make check-race                 # run unit test with -race option
make check-build                # integration test (/cli_test)
```

### Solo Node
```
./initialize.sh
```

### Local Test Net
```
make build-linux                # Cross-compile the binaries for linux/amd64
make build-docker-testnet      # Build docker image for testnet
make build-conf-testnet        # Build configurations for testnet
make start-testnet             # Boot up testnet network with 4 validator nodes
make stop-testnet              # Stop the testnet
```


### Current Status
The most of development is in progress for testing tendermint/cosmos-sdk.


## API documentation

### Run the API documentation server
```shell script
make install && linkcli rest-server --trust-node=true
```
then visit [http://localhost:1317/swagger-ui/#/](http://localhost:1317/swagger-ui/#/)

### Update the API documentation
Update [swagger.yaml](./client/lcd/swagger-ui/swagger.yaml) then:
```shell script
make update-swagger-docs && make install
```