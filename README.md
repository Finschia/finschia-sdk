The Finschia SDK is a framework for building blockchain applications in Golang.
It is being used to build [`Finschia`](https://github.com/Finschia/finschia), the first implementation of the Finschia Mainnet.
This is forked from [`cosmos-sdk`](https://github.com/cosmos/cosmos-sdk) at 2021-03-15.

**WARNING**: Breaking changes may occur because this repository is still in the pre-release development phase.

**Note**: Requires [Go 1.20+](https://golang.org/dl/)

## What is the Finschia SDK?

The [Finschia SDK](https://github.com/Finschia/finschia-sdk) is an open-source framework for building multi-asset public Proof-of-Stake (PoS) <df value="blockchain">blockchains</df>, as well as permissioned Proof-Of-Authority (PoA) blockchains. Blockchains built with the Cosmos SDK are generally referred to as **application-specific blockchains**. 

The purpose of `Finschia SDK` is to succeed to [the objectives of `Cosmos sdk`](https://github.com/cosmos/cosmos-sdk/blob/master/docs/intro/overview.md) while helping develop blockchains that require faster transaction processing to be applied to reality.

## Why the Finschia SDK?

Cosmos-sdk, which created the philosophy of application-specific blockchain, established its status as a framework for various application blockchain development. `Finschia SDK` inherited this `cosmos-sdk` philosophy, addressing slow transaction processing problem that was difficult for cosmos-sdk to apply in real financial system. Real financial systems require thousands of processing performance per second, with Finschia SDK adding many performance improvements to meet that demand.
The following work was carried out to improve performance.

- Concurrent checkTx, deliverTx
- Use [fastcache](https://github.com/victoriametrics/fastcache) for inter block cache and nodedb cache of iavl
- Lock granularity enhancement

In addition, the following functions were added:

- Virtual machine using `cosmwasm` that makes smart contracts possible to be executed 
- Use [Ostracon](https://github.com/Finschia/ostracon) as consensus engine instead of `Tendermint`


To learn about Cosmos SDK, please refer [Cosmos SDK Docs](https://github.com/cosmos/cosmos-sdk/blob/master/docs).

## Quick Start

**Build**
```
make build
make install

# you can see the version!
simd version
```

&nbsp;

**Configure**
```
zsh init_node.sh sim {N(number of nodes), default=1}
```

&nbsp;

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
* REST: http://localhost:1317/swagger/

&nbsp;

## Follow Guide
You can refer to the sample tx commands [here](docs/sample-tx.md). 
Test different commands to get a broader understanding of Finschia

