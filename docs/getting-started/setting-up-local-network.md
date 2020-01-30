# Setting Up Local Network

## Binaries

- linkd: LINK full node daemon
- linkcli: CLI tools, LCD (light client daemon)


## Running a local network 

### Initializing config files and genesis file with a node name called moniker

```shell
> linkd init [MONIKER] --chain-id [CHAIN ID]

e.g.) 
> linkd init localnetnode --chain-id linklocal
```

then, the files are stored in default home dir of chain data. (default `~/.linkd/config`)

(It also can be changed by adding a flag `--home`)


### Adding genesis accounts

#### Creating new keys to add to genesis state

```shell
> linkcli keys add [NAME OF ACCOUNT]

e.g.) 
> linkcli keys add jack
> linkcli keys add alice
```
then, the keys are stored in default home dir (default `~/.linkcli/keys`) that `linkcli`  manages.

(You can change home dir by adding a flag `--home`)

#### Adding them as  genesis accounts  with initial coins

```shell
> linkd add-genesis-account [ADDRESS] [COMMA SEPARATED COIN FORM]

e.g.)
> linkd add-genesis-account $(linkcli keys show jack -a) 100link,100000000stake
> linkd add-genesis-account $(linkcli keys show alice -a) 100link,100000000stake
```

### Creating a genesis `create-validator` TX

For the network to run successfully, One or more validators should exist.
To be a validator, an account should be bonded by delegating coins.

The below command generates bonding(`create-validator`) tx as a genesis TX.
Which will be executed right after the network starts.

```shell
> linkd gentx --name [KEY NAME]

e.g.)
> linkd gentx --name jack --amount 100000000stake
```


### Collecting genesis TXs

The following step collects genesis TX to the genesis file.

```shell
> linkd collect-gentxs
```

And validates that the genesis file is the correct form.

```shell
> linkd validate-genesis
```

### Launching the network

```shell
> linkd start
```



## Running a node on an existing network

### Checking the network information and the genesis file

Basically, the full node supports querying the information through the HTTP API.

- `GET /genesis`

```shell
> curl http://[IP]:[PORT]/genesis
```

From the result of HTTP API, `result.genesis` object is the same as the genesis file

### Init config files and genesis file with another moniker and the chain id (network id) 

```shell
> linkd init [MONIKER] --chain-id [CHAIN ID]
``` 

### Overwrite the genesis file to `genesis.json` 

Overwrite the genesis file to `[LINKD HOME]/config/genesis.json`

```json
{
  "genesis_time": "2019-10-31T06:40:27.826184785Z",
  "chain_id": "k8s-chain-p2p-26656-rpc-26657-abci-26658-c2356069e",
  "consensus_params": {
   
  }
}
```


### Add peers to `config.toml`

From the genesis file, you can find the peer information([ID]@[PEER_IP]:[PORT] form) at the memo of the gentxs.

JSON path is `app_state.genutil.gentxs.item.value.memo`

```json
{
  "genesis_time": "2019-10-31T06:40:27.826184785Z",
  "chain_id": "k8s-chain-p2p-26656-rpc-26657-abci-26658-c2356069e",
  
  "app_state": {


      "genutil": {
        "gentxs": [
          {

            "memo": "0dc4cd6d7b719051b9cd6e64fc9a7a0f18ff55c0@[ip]:26656"
          }
        ]
      }
  }
}
```

Copy one or more peers to `config.toml`
```shell
> vi [LINKD HOME]/config/config.toml

persistent_peers = "id1@peer_ip1:26656,id2@peer_ip2:26656,..."

```

### Start the node
```shell
> linkd start
```



## Using the CLI tools

### Configure default flags for the CLI

```shell
# Configure your CLI to eliminate need for chain-id flag
> linkcli config chain-id linklocal
> linkcli config output json
> linkcli config indent true
> linkcli config trust-node true
```

### Sending coins

```shell
> linkcli tx send [FROM ADDRESS OR KEY NAME] [TO ADDRESS] [COIN FORM]

e.g.)
> linkcli tx send jack $(linkcli keys show alice -a) 100link
```

### Starting a LCD

```shell
> linkcli rest-server

# starts with a listen address
> linkcli rest-server --laddr tcp://0.0.0.0:1317

# starts with a node address
> linkcli rest-server --node tcp://[IP]:[PORT]
```

### More commands

```shell
> linkcli --help
```
