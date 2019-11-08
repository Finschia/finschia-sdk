# LINK Network Version2

[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://www.repostatus.org/badges/latest/wip.svg)](https://www.repostatus.org/#wip)

[![codecov](https://codecov.io/gh/line/link/branch/develop/graph/badge.svg?token=JFFuUevpzJ)](https://codecov.io/gh/line/link)

This repository hosts `LINK`, alternative implementation of the LINK Network.

**Node**: Requires [Go 1.12+](https://golang.org/dl/)

**Warnings**: Initial development is in progress, but there has not yet been a stable.

# Quick Start
**Build Docker Image**
```
make build-docker                # build docker image
```

**Configure**
```
./.initialize.sh docker          # prepare keys, validators, initial state, etc.
```

**Run**
```
docker-compose up                # Run a Node and Rest
```

**visit with your browser**
* Node: http://localhost:26657/
* REST: http://localhost:1317/swagger-ui/

# Step by Step

## Build and Test
**Prerequisite**
```
make get-tools                   # install tools
```
**Build & Install LINK**
```
make install                     # build and install binaries
```

**Test**
```
make test-unit                   # unit test
make test-unit-race              # run unit test with -race option
make test-integration            # integration test (/cli_test#cli_test)
make test-integration-multi-node # integration test (/cli_test#cli_multi_node_test)
```

## Configure

Configuring to run a node should be done following steps.

- init chain
- add genesis account
- make the node be a validator
-  genesis TXs

For more informations [here](docs/getting-started/setting-up-local-network.md#running-a-local-network).

<br />

The configuring steps are also in a script file(`.initialize.sh`).

```shell
./.initialize.sh                 # prepare keys, validators, initial state, etc.
```
**WARNING**: Do not use it for production. Use it only for local testing 

<br />

If the configuring steps are done, the following files should be in the home dir.

- Check the home of linkd**

  ```
  ls ${HOME}/.linkd/config
  ```
- You must have these files
  ```
  app.toml	config.toml	genesis.json	gentx	node_key.json	priv_validator_key.json
  ```

## Run the node

**Start the Node**
```
linkd start                 # Start a validator
```
Check Node: http://localhost:26657/

**Start Rest Server**
```
linkcli rest-server         # Start a rest server connecting to the validator
```
Check Rest Server: http://localhost:1317/swagger-ui/

**Query/SendTx with cli**
```
linkcli status                                              # check the status of node
linkcli tx send jack $(linkcli keys show alice -a) 1link -y # password: 1234567890
linkcli query account $(linkcli keys show jack -a)          # Get account
```

# Local Test Network

## local test network with 4 validators

**Build Docker Image**
```
make build-docker
```
**Start the testnet**
```
make testnet-start          
```
**Test the liveness**
```
make testnet-test
```
**Stop the testnet**
```
make testnet-stop
```

# Current Status
The most of development is in progress for testing tendermint/cosmos-sdk.
