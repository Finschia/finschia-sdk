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
simd version    # you can see the version!
```

**Configure**
```
sh init_node.sh sim
```

**Run**
```
simd start      # run a node
```

**Visit with your browser**
* Node: http://localhost:26657/
* REST: http://localhost:1317/swagger-ui/


## Follow Guide


**Create new account**
```
simd keys add {new-account-name} --keyring-backend test
simd keys list --keyring-backend test                       # check if new account was added successfully
```

Let the new, user and validator account address be new-addr, user-addr and val-addr each.

**Send funds(Bank)**
```
simd query bank balances {new-addr}                 # balances: "0"
simd query bank balances {user-addr}                # balances: 100000000000stake, 100000000000ukrw
simd tx bank send {user-addr} {new-addr} 10000stake —keyring-backend test —chain-id sim
                                                    # send 10000stake to new-account from user-account
                                                            
simd query bank balances {new-addr}                 # balances: 10000stake
simd query bank balances {user-addr}                # balances: 99999990000stake, 100000000000ukrw
```

**Staking(deligate)**
```
simd query staking validators                       # operator_address is Bech32 of val-addr(let it be val-addr-Bech32)
simd tx staking delegate {val-addr-Bech32} 1000stake --from {new-addr} --keyring-backend test --chain-id sim
                                                    # deligate 1000stake to validator
simd query staking validators                       # check if deligation was successful

simd tx staking unbond {val-addr-Bech32} 1000stake --from {new-addr} --keyring-backend test --chain-id sim
                                                    # undeligate 1000stake from validator
simd query staking validators                       # check if undeligation was successful
```

Test different commands to get a broader understanding of lbm
