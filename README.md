The LBM SDK is a framework for building blockchain applications in Golang.
It is being used to build [`LBM`](https://github.com/line/lbm), the first implementation of the LINE Blockchain Mainnet.
This is forked from [`cosmos-sdk`](https://github.com/cosmos/cosmos-sdk) at 2021-03-15.

**WARNING**: Breaking changes may occur because this repository is still in the pre-release development phase.

**Note**: Requires [Go 1.15+](https://golang.org/dl/)

## What is the LBM SDK?

The [LBM SDK](https://github.com/line/lbm-sdk) is an open-source framework for building multi-asset public Proof-of-Stake (PoS) <df value="blockchain">blockchains</df>, as well as permissioned Proof-Of-Authority (PoA) blockchains. Blockchains built with the Cosmos SDK are generally referred to as **application-specific blockchains**. 

The purpose of `LBM SDK` is to succeed to [the objectives of `Cosmos sdk`](https://github.com/cosmos/cosmos-sdk/blob/master/docs/intro/overview.md) while helping develop blockchains that requires faster transaction processing to be applied to reality.

## Why the LBM SDK?

Cosmos-sdk, which created the philosophy of application-specific blockchain, established its status as a framework for various application blockchain development. `LBM SDK` inherited this `cosmos-sdk` philosophy, addressing slow transaction processing problem that was difficult for cosmos-sdk to apply in real financial system. Real financial systems require thousands of processing performance per second, with LBM SDK adding many performance improvements to meet that demand.
The following work was carried out to improve performance.

- Concurrent checkTx, deliverTx
- Use [fastcache](https://github.com/victoriametrics/fastcache) for inter block cache and nodedb cache of iavl
- Lock granularity enhancement

In addition, the following functions were added:

- Virtual machine using `cosmwasm` that makes smart contracts possible to be executed 
- Use [Ostracon](https://github.com/line/ostracon) as consensus engine instead of `Tendermint`


To learn about Cosmos SDK, please refer [Cosmos SDK Docs](https://github.com/cosmos/cosmos-sdk/blob/master/docs).

## Quick Start


**Build**
```
make build
make install

# you can see the version!
simd version
```

**Configure**
```
sh init_node.sh sim {N(number of nodes), default=1}
```

**Run**
```
# run a node
simd start --home ~/.simapp/simapp0

# If N is larger than 1, run all node.
# simapp0 has other nodes as persistant_peer. 
simd start --home ~/.simapp/simapp0
simd start --home ~/.simapp/simapp1
...
```

**Visit with your browser**
* Node: http://localhost:26657/
* REST: http://localhost:1317/swagger-ui/


## Follow Guide


**Create new account**
```
simd keys add user0 --keyring-backend test --home ~/.simapp/simapp0

# check if new account was added successfully
simd keys list --keyring-backend test --home ~/.simapp/simapp0               
```

Let the user0 and validator0 **account address** be each 
* **user0: link1lu5hgjp2gyvgdpf674aklzrpdeuyhjr4fsuqrj**
* **validator0: link146asaycmtydq45kxc8evntqfgepagygelel00h**

If you run multi node, home option's value can be ~/.simapp/simapp1, ~/.simapp/simapp2, ...
You can get same result whatever --home option you use

**Send funds(Bank)**
```
# user0 balances: "0"
simd query bank balances link1lu5hgjp2gyvgdpf674aklzrpdeuyhjr4fsuqrj --home ~/.simapp/simapp0

# validator0 balances: 90000000000stake, 100000000000ukrw
simd query bank balances link146asaycmtydq45kxc8evntqfgepagygelel00h --home ~/.simapp/simapp0

# send 10000stake from validator0 to user0
simd tx bank send link146asaycmtydq45kxc8evntqfgepagygelel00h link1lu5hgjp2gyvgdpf674aklzrpdeuyhjr4fsuqrj 10000000000stake --keyring-backend test --chain-id sim --home ~/.simapp/simapp0

# user0 balances: 10000000000stake
simd query bank balances link1lu5hgjp2gyvgdpf674aklzrpdeuyhjr4fsuqrj --home ~/.simapp/simapp0

# validator0 balances: 80000000000stake, 100000000000ukrw
simd query bank balances link146asaycmtydq45kxc8evntqfgepagygelel00h --home ~/.simapp/simapp0
```

**Staking(deligate)**
```
# Bech32 Val is operator address of validator0
simd debug addr link146asaycmtydq45kxc8evntqfgepagygelel00h --home ~/.simapp/simapp0
```
Let the validator0 operator address be linkvaloper146asaycmtydq45kxc8evntqfgepagygeddajpy

```
# deligate 10000000000stake to validator0
simd tx staking delegate linkvaloper146asaycmtydq45kxc8evntqfgepagygeddajpy 10000000000stake 
--from link1lu5hgjp2gyvgdpf674aklzrpdeuyhjr4fsuqrj --keyring-backend test --chain-id sim --home ~/.simapp/simapp0

# check if deligation was successful
simd query staking validators --chain-id sim --home ~/.simapp/simapp0

# undeligate 10000000000stake from validator
simd tx staking unbond linkvaloper146asaycmtydq45kxc8evntqfgepagygeddajpy 10000000000stake --from link1lu5hgjp2gyvgdpf674aklzrpdeuyhjr4fsuqrj --keyring-backend test --chain-id sim --home ~/.simapp/simapp0

# check if undeligation was successful
simd query staking validators --chain-id sim --home ~/.simapp/simapp0
```

Test different commands to get a broader understanding of lbm
