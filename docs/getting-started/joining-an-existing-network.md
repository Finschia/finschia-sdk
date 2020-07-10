# Joining an existing network
This document describes how to join an existing LINK network. The goal is to run a new node connecting to the network.

---

## Check the network information and the genesis file

Basically, the full node supports querying the information through the HTTP API.

- `GET /genesis`

```shell
$ curl http://[IP]:[PORT]/genesis
```

From the result of HTTP API, `result.genesis` object is the same as the genesis file

## Initialize config files and genesis file with another moniker and the chain id (network id) 

```shell
$ linkd init [MONIKER] --chain-id [CHAIN_ID]
``` 

## Overwrite the genesis file to `genesis.json` 

Overwrite the genesis file to `/path/to/link-home/config/genesis.json`. The default path to link-home is `$HOME/.linkd`.

```json
{
  "genesis_time": "2019-10-31T06:40:27.826184785Z",
  "chain_id": "k8s-chain-p2p-26656-rpc-26657-abci-26658-c2356069e",
  "consensus_params": {
    // skip
  }
}
```


## Add peers to `config.toml`

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
$ vi [LINKD HOME]/config/config.toml

persistent_peers = "id1@peer_ip1:26656,id2@peer_ip2:26656,..."

```

## Start the node
```shell
$ linkd start
```

