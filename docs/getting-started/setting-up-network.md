# Setting Up a Network

## Node Types 
A node can be classified as one or more types among "node, validator, sentry and seed".

### Node
```toml
[p2p]
external_address = "tcp://[SELFIP]:26656"
seeds = "[ID]@[IP]:26656"
addr_book_strict = false
seed_mode = false
pex = true
persistent_peers = "[ID]@[IP]:26656" # list of sentry nodes
```


### Validator
A validator's address should not be exposed to others except the sentry node.
Hence, disable pex and no need to connect to others but the sentry
```toml
[p2p]
external_address = "tcp://[SELFIP]:26656"
seeds = "" # keep empty
addr_book_strict = false
seed_mode = false
pex = false
persistent_peers = "[ID]@[IP]:26656" # list of sentry nodes
```

### Sentry
A sentry node connect to others like a normal node but keep the validators as private
```toml
[p2p]
external_address = "tcp://[SELFIP]:26656"
seeds = "[ID]@[IP]:26656"
addr_book_strict = false
seed_mode = false
pex = true
private_peer_ids = "[ID]@[IP]:26656" # list of validator nodes
```

### Seed
Every node can be a seed node if the `seed_mode` turned on.
It returns a list of known active peers. And its `pex` reactor operated in `crawler` mode and continuously explores others

```toml
[p2p]
external_address = "tcp://[SELF IP]:26656"
seeds = "" # Can be set to other seed nodes
addr_book_strict = false
seed_mode = true
pex = true
```

## Prepare Validators

### Private Validator Key
```bash
linkd init [moniker] --chain-id [chain-id]
```

Then you can get the `priv_validator_key.json` file. 
This file actually contains a private key, and should thus be kept absolutely secret.

### Generate Operating Account 

```bash
linkcli keys add [name]
```
Generate accounts by using the cli wallet.
You need accounts more than the number of validators.

To make it easy to manage the private keys, you can use HD wallet.

```bash
linkcli keys mnemonic
```
Get the mnemonic for the master seed.

```bash
linkcli keys add [account name] --index [index] --recover
```
You can generate and recover child accounts by using the mnemonic.
just remember the mnemonic and the indexes.


## Default on-chain data
Set up the default value for the `consensus parameters` and `application state`
It is a genesis state every node start on. 
```json
{
  "consensus_params": {
    "block": {
      "max_bytes": "22020096",
      "max_gas": "-1",
      "time_iota_ms": "1000"
    },
    "evidence": {
      "max_age": "100000"
    },
    "validator": {
      "pub_key_types": [
        "ed25519"
      ]
    }
  },
  "app_state": {
    "bank": {
      "send_enabled": true
    },
    "staking": {
      "params": {
        "unbonding_time": "1814400000000000",
        "max_validators": 100,
        "max_entries": 7,
        "bond_denom": "stake"
      },
      "last_total_power": "0",
      "last_validator_powers": null,
      "validators": null,
      "delegations": null,
      "unbonding_delegations": null,
      "redelegations": null,
      "exported": false
    },
    "auth": {
      "params": {
        "max_memo_characters": "256",
        "tx_sig_limit": "7",
        "tx_size_cost_per_byte": "10",
        "sig_verify_cost_ed25519": "590",
        "sig_verify_cost_secp256k1": "1000"
      }
    },
    "supply": {
      "supply": []
    }
  }
}
```

## Gathering Create-Validator Transactions
Generally, this step is done by individual organizations operating a validator node.
Hence, if you are planning to set up multiple validator nodes by your-self, repeat bellow for the number of validators you have 

### Initialize Node
```bash
rm -rf ${HOME}/.linkd ${HOME}/.linkcli
linkd init [moniker] --chain-id [chain id]
# ex) linkd init node0 --chain-id testnet1000
```

### Operating Account
```bash
linkcli keys add [name]

# or by using mnemonic
linkcli keys add [name] --index [index] --recover

```
An operator account needs to be added to cli wallet to make signature for the transaction

```bash
# Create Account to genesis.json
linkd add-genesis-account [address] [amount]
# ex) linkd add-genesis-account link1f853nx2ecr92x0w005ts2yckemd0zuj0cfnlmu 100000000stake
# ex) linkd add-genesis-account $(linkcli keys show account0 -a) 100000000stake
```
Add the account to `genesis.json` with initial amount of coin

### Generate create-validator transaction
```bash
# Create genesis transaction for creating validator
linkd gentx --name account0
```
Now, you can generate the transaction and sign on it with the private key of the account.
The generated transaction is stored as a json file with node-id prefix.

```bash
ls ${HOME}/.linkd/config/gentx/gentx-6f03a785e1b200346e386f98bf87604fe9ab2dc3.json
```

```json
{"type":"cosmos-sdk/StdTx","value":{"msg":[{"type":"cosmos-sdk/MsgCreateValidator","value":{"description":{"moniker":"solo","identity":"","website":"","details":""},"commission":{"rate":"0.100000000000000000","max_rate":"0.200000000000000000","max_change_rate":"0.010000000000000000"},"min_self_delegation":"1","delegator_address":"link1ukygv7l4dew9gvl5xmwpeq7uyuuy8hvwgglae0","validator_address":"linkvaloper1ukygv7l4dew9gvl5xmwpeq7uyuuy8hvw6uaqhu","pubkey":"linkvalconspub1zcjduepqqy5hzwnx3uunph8qet75mrse3mfwp9jxzfwk7v9vt37me4wu44fqhpwpxf","value":{"denom":"stake","amount":"100000000"}}}],"fee":{"amount":[],"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"A0yvrn1FxBfr8Z2WgxqhpXSIV8cQZsjrFE27URPpViff"},"signature":"MWTmZBUPsJ6V87gRWXFx9PxicdrE3BloU0Tulcl6DKZKSovxRLmP0iIJdNp0XGv6XL+tDJtW+etbBM5t/ylulw=="}],"memo":"6f03a785e1b200346e386f98bf87604fe9ab2dc3@192.168.21.181:26656"}}
```
The `memo` field is used to configure `persistent_peers` for the `config.toml`. But it is not mandatory. 


## Collect Create-Validator Transactions And Distribute genesis.json
There should be a collector who have a role for collecting the transactions and distribute `genesis.json` even if the network is public.
The collector should know all of the genesis accounts and their initial coins.
Bellow is how to generate the `genesis.json`

### Add genesis accounts

```bash
# Create Account to genesis.json
linkd add-genesis-account [address] [amount]
```
All of the account should be registered in cli wallet in advance.

### Collect gentxs and store to genesis.json
```bash
linkd collect-gentxs --gentx-dir [path]
linkd validate-genesis
```
The `[path]` should contains all of create-validator transaction files formed with json

Distribute the `genesis.json` which contains `genesis accounts` and `create validator transactions`

## Operate Nodes
Configure `config.toml` according to type of the nodes.

### Non-Validator
```bash
rm -rf ${HOME}/.linkd ${HOME}/.linkcli
linkd init [moniker] --chain-id [chain id]
```
initialize node and its configuration.
```bash
cp [genesis.json] ${HOME}/.linkd/config/genesis.json
```
Replace genesis.json file by the distributed one.

```bash
linkd start --p2p.seed_mode=[seed node address]
```
Start a node and let it trying to connect others.

### Validator

Same as above. But replace `priv_validator_key.json` file.

```bash
linkd start
```
