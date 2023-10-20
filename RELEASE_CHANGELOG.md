# Changelog

## [v0.47.0](https://github.com/Finschia/finschia-sdk/releases/tag/v0.47.0) - 2023-04-24

### Features
* (baseapp) [\#840](https://github.com/Finschia/finschia-sdk/pull/840) allow querying the state based on `CheckState`.
* (x/foundation) [\#848](https://github.com/Finschia/finschia-sdk/pull/848) remove `gov mint` for x/foundation proposal
* (x/wasm) [\#850](https://github.com/Finschia/finschia-sdk/pull/850) remove `x/wasm` module in lbm-sdk
* (log) [\#883](https://github.com/Finschia/finschia-sdk/pull/883) add zerolog based rolling log system
* (Ostracon) [\#887](https://github.com/Finschia/finschia-sdk/pull/887) apply the changes of vrf location in Ostracon
* (x/upgrade) [\#889](https://github.com/Finschia/finschia-sdk/pull/889) remove time based upgrade
* (all) [\#970](https://github.com/Finschia/finschia-sdk/pull/970) change import path to `github.com/Finschia/finschia-sdk` and update license

### Improvements
* (cosmovisor) [\#792](https://github.com/Finschia/finschia-sdk/pull/792) Use upstream's cosmovisor
* (server) [\#821](https://github.com/Finschia/finschia-sdk/pull/821) Get validator pubkey considering KMS
* (client) [\#890](https://github.com/Finschia/finschia-sdk/pull/890) Map Ostracon:ErrTxInMap to lbm-sdk:ErrTxInMempoolCache
* (x/collection) [\#894](https://github.com/Finschia/finschia-sdk/pull/894) Change the default params of x/collection
* (ante) [\#895](https://github.com/Finschia/finschia-sdk/pull/895) Remove max gas validation
* (x/collection,token) [\#900](https://github.com/Finschia/finschia-sdk/pull/900) Add uri for MsgModify and deprecate the old ones
* (x/foundation) [\#912](https://github.com/Finschia/finschia-sdk/pull/912) Introduce censorship into x/foundation
* (x/foundation) [\#933](https://github.com/Finschia/finschia-sdk/pull/933) Clean up x/foundation apis
* (x/collection) [\#938](https://github.com/Finschia/finschia-sdk/pull/938) Add progress log into x/collection import-genesis
* (x/foundation) [\#952](https://github.com/Finschia/finschia-sdk/pull/952) Address generation of the empty coins in x/foundation
* (x/collection,token,foundation) [\#963](https://github.com/Finschia/finschia-sdk/pull/963) Check event determinism on original modules
* (x/collection) [\#965](https://github.com/Finschia/finschia-sdk/pull/965) Provide specific error messages on x/collection queries
* (x/collection,token) [\#980](https://github.com/Finschia/finschia-sdk/pull/980) refactor x/token,collection query errors

### Bug Fixes
* (client) [\#817](https://github.com/Finschia/finschia-sdk/pull/817) remove support for composite (BLS) type
* (x/foundation) [\#834](https://github.com/Finschia/finschia-sdk/pull/834) Apply foundation audit
* (x/collection,token) [\#849](https://github.com/Finschia/finschia-sdk/pull/849) Introduce codespace into x/collection,token
* (x/token,collection) [\#863](https://github.com/Finschia/finschia-sdk/pull/863) Update x/collection,token proto
* (x/collection,token) [\#866](https://github.com/Finschia/finschia-sdk/pull/866) Do not create account on x/token,collection
* (x/collection,token) [\#881](https://github.com/Finschia/finschia-sdk/pull/881) Remove some x/token,collection queries on listable collections
* (swagger) [\#898](https://github.com/Finschia/finschia-sdk/pull/898) fix a bug not added `lbm.tx.v1beta1.Service/GetBlockWithTxs` in swagger
* (x/collection) [\#911](https://github.com/Finschia/finschia-sdk/pull/911) Add missing command(TxCmdModify) for CLI
* (x/foundation) [\#922](https://github.com/Finschia/finschia-sdk/pull/922) Propagate events in x/foundation through sdk.Results
* (x/foundation) [\#946](https://github.com/Finschia/finschia-sdk/pull/946) Fix broken x/foundation invariant on treasury
* (x/foundation) [\#947](https://github.com/Finschia/finschia-sdk/pull/947) Unpack proposals in x/foundation import-genesis
* (x/collection) [\#953](https://github.com/Finschia/finschia-sdk/pull/953) Allow zero amount of coin in x/collection Query/Balance
* (x/collection) [\#954](https://github.com/Finschia/finschia-sdk/pull/954) Remove duplicated events in x/collection Msg/Modify
* (x/collection) [\#955](https://github.com/Finschia/finschia-sdk/pull/955) Return nil where the parent not exists in x/collection Query/Parent
* (x/collection) [\#959](https://github.com/Finschia/finschia-sdk/pull/959) Revert #955 and add Query/HasParent into x/collection
* (x/collection) [\#960](https://github.com/Finschia/finschia-sdk/pull/960) Fix default next class ids of x/collection
* (x/collection) [\#961](https://github.com/Finschia/finschia-sdk/pull/961) Do not loop enum in x/collection
* (x/collection,token) [\#957](https://github.com/Finschia/finschia-sdk/pull/957) Refactor queries of x/collection and x/token
* (x/auth) [\#982](https://github.com/Finschia/finschia-sdk/pull/982) Fix not to emit error when no txs in block while querying `GetBlockWithTxs`
* (x/foundation) [\#984](https://github.com/Finschia/finschia-sdk/pull/984) Revert #952

### Removed
* [\#853](https://github.com/Finschia/finschia-sdk/pull/853) remove useless stub BeginBlock, EndBlock methods from modules below
  * ibc, authz, collection, feegrant, ibc, token, wasm
* (x/ibc) [\#858](https://github.com/Finschia/finschia-sdk/pull/858) detach ibc module(repo: [finschia/ibc-go](https://github.com/Finschia/ibc-go))
* (x/collection,token) [\#966](https://github.com/Finschia/finschia-sdk/pull/966) Remove legacy events on x/collection and x/token

### Breaking Changes
* (rest) [\#807](https://github.com/Finschia/finschia-sdk/pull/807) remove legacy REST API
* (codec) [\#833](https://github.com/Finschia/finschia-sdk/pull/833) Fix foundation amino codec
* (ostracon) [\#869](https://github.com/Finschia/finschia-sdk/pull/869) apply changes to replace Ostracon proto message with Tendermint
* (x/bank) [\#876](https://github.com/Finschia/finschia-sdk/pull/876) Add `MultiSend` deactivation
* (x/auth) [\#891](https://github.com/Finschia/finschia-sdk/pull/891) deprecate `cosmos.tx.v1beta1.Service/GetBlockWithTxs` and add `lbm.tx.v1beta1.Service/GetBlockWithTxs` for lbm
* (abci) [\#892](https://github.com/Finschia/finschia-sdk/pull/892) remove the incompatible field `index=14` in `TxResponse`
* (proto) [\#923](https://github.com/Finschia/finschia-sdk/pull/923) deprecate broadcast mode `block`
* (x/collection,token) [\#956](https://github.com/Finschia/finschia-sdk/pull/956) Replace query errors on the original modules into gRPC ones

### Build, CI
* (ci) [\#829](https://github.com/Finschia/finschia-sdk/pull/829) automate release process
* (build) [\#872](https://github.com/Finschia/finschia-sdk/pull/872) Retract v1.0.0
* (ci, build) [\#901](https://github.com/Finschia/finschia-sdk/pull/901) Update release pipeline to match non-wasm env
* (ci) [\#983](https://github.com/Finschia/finschia-sdk/pull/983) update docker action to fit new repository

### Document Updates
* (x/foundation) [\#934](https://github.com/Finschia/finschia-sdk/pull/934) Update permlinks in x/foundation documents
* (x/collection,token) [\#944](https://github.com/Finschia/finschia-sdk/pull/944) Update comments in the x/token,collection events proto


## [v0.46.0](https://github.com/Finschia/finschia-sdk/releases/tag/v0.46.0) - 2022-11-28

This version based on cosmos-sdk v0.45.10, wasmd v0.27.0 and Ostracon v1.0.7.

### Features
* (global) [\#783](https://github.com/Finschia/finschia-sdk/pull/783) bump up github.com/cosmos/cosmos-sdk to v0.45.10
* (build) [\#793](https://github.com/Finschia/finschia-sdk/pull/793) enable to use libsodium version ostracon
* (x/wasm) [\#570](https://github.com/Finschia/finschia-sdk/pull/570) Merge wasmd 0.27.0
* (x/wasm) [\#470](https://github.com/Finschia/finschia-sdk/pull/470) remove contract activation control by actor
* (x/wasm) [\#513](https://github.com/Finschia/finschia-sdk/pull/513) fix message representation for signing
* (x/foundation) [\#518](https://github.com/Finschia/finschia-sdk/pull/518) add foundation treasury feature to x/foundation
* (x/foundation) [\#528](https://github.com/Finschia/finschia-sdk/pull/528) add a feature of whitelist for /lbm.foundation.v1.MsgWithdrawFromTreasury
* (proto) [\#584](https://github.com/Finschia/finschia-sdk/pull/564) remove `prove` field in the `GetTxsEventRequest` of `tx` proto
* (x/collection) [\#571](https://github.com/Finschia/finschia-sdk/pull/571) add x/collection proto
* (x/collection) [\#574](https://github.com/Finschia/finschia-sdk/pull/574) implement x/collection
* (store) [\#605](https://github.com/Finschia/finschia-sdk/pull/605) replace line/iavl and line/tm-db with cosmos/iavl and tendermint/tm-db.
* (server/grpc) [\#607](https://github.com/Finschia/finschia-sdk/pull/607) revert gRPC block height header.
* (global) [\#611](https://github.com/Finschia/finschia-sdk/pull/611) bump github.com/cosmos/cosmos-sdk from v0.45.1 to v0.45.6
* (simapp) [\#620](https://github.com/Finschia/finschia-sdk/pull/620) chore: add iterator feature for simapp
* (x/collection) [\#622](https://github.com/Finschia/finschia-sdk/pull/622) add Query/TokenClassTypeName
* (x/bank) [\#629](https://github.com/Finschia/finschia-sdk/pull/629) remove unsafe balance changing methods from bank keeper such as `SetBalance` and `SetSupply`.
* (x/wasm) [\#649](https://github.com/Finschia/finschia-sdk/pull/649) fix: wasm module's FIXME in the snapshotter.go file
* (x/ibc) [\#651](https://github.com/Finschia/finschia-sdk/pull/651) feat: update x/ibc to support github.com/cosmos/ibc-go@v3.0.0
* (config) [\#665](https://github.com/Finschia/finschia-sdk/pull/665) remove bech32-cache-size
* (x/foundation) [\#709](https://github.com/Finschia/finschia-sdk/pull/709) add `gov mint` for x/foundation proposal
* (iavl) [\#738](https://github.com/Finschia/finschia-sdk/pull/738) bump github.com/cosmos/iavl from v0.17.3 to v0.19.3
* (baseapp) [\#756](https://github.com/Finschia/finschia-sdk/pull/756) Change to create chCheckTx with the value set in app config
* (x/foundation) [\#758](https://github.com/Finschia/finschia-sdk/pull/758) add invariants to x/foundation

### Improvements
* (ostracon) [\#811](https://github.com/Finschia/finschia-sdk/pull/811) bump up Ostracon from `2231684789c6e754744e3f10f8b41c96c4575ad1` to v1.0.7
* (x/auth) [\#776](https://github.com/Finschia/finschia-sdk/pull/776) remove unused MsgEmpty
* (refactor) [\#493](https://github.com/Finschia/finschia-sdk/pull/493) restructure x/consortium
* (server/grpc) [\#526](https://github.com/Finschia/finschia-sdk/pull/526) add index field into TxResponse
* (cli) [\#535](https://github.com/Finschia/finschia-sdk/pull/536) updated ostracon to v1.0.5; `unsafe-reset-all` command has been moved to the `ostracon` sub-command.
* (x/foundation) [\#597](https://github.com/Finschia/finschia-sdk/pull/597) tidy up x/foundation
* (x/collection) [\#604](https://github.com/Finschia/finschia-sdk/pull/604) add EventOwnerChanged and EventRootChanged
* (x/collection) [\#608](https://github.com/Finschia/finschia-sdk/pull/608) remove new APIs on x/collection
* (x/token) [\#609](https://github.com/Finschia/finschia-sdk/pull/609) remove new APIs on x/token
* (x/collection) [\#621](https://github.com/Finschia/finschia-sdk/pull/621) add additional information into EventXXXChanged
* (x/token) [\#636](https://github.com/Finschia/finschia-sdk/pull/636) add creator into x/token EventIssue
* (x/token) [\#637](https://github.com/Finschia/finschia-sdk/pull/637) rename x/token events
* (x/collection) [\#639](https://github.com/Finschia/finschia-sdk/pull/639) rename x/collection events
* (x/wasm) [\#661](https://github.com/Finschia/finschia-sdk/pull/661) x/wasm refactoring - detaching the custom wasm proto part of lbm-sdk. (apply changes of [\#625](https://github.com/Finschia/finschia-sdk/pull/625) and [\#655](https://github.com/Finschia/finschia-sdk/pull/655))
* (refactor) [\#685](https://github.com/Finschia/finschia-sdk/pull/685) remove x/foundation UpdateValidatorAuthsProposal
* (x/foundation) [\#686](https://github.com/Finschia/finschia-sdk/pull/686) remove `Minthreshold` and `MinPercentage` from x/foundation config
* (x/foundation) [\#693](https://github.com/Finschia/finschia-sdk/pull/693) add pool to the state of x/foundation
* (x/auth, client) [\#699](https://github.com/Finschia/finschia-sdk/pull/699) Improvement on input validation of `req.Hash`
* (x/wasm,distribution) [\#696](https://github.com/Finschia/finschia-sdk/pull/696) x/wasm,distribution - add checking a file size before reading it
* (x/foundation) [\#698](https://github.com/Finschia/finschia-sdk/pull/698) update x/group relevant logic in x/foundation
* (x/auth,bank,foundation,wasm) [\#691](https://github.com/Finschia/finschia-sdk/pull/691) change AccAddressFromBech32 to MustAccAddressFromBech32
* (x/wasm) [\#690](https://github.com/Finschia/finschia-sdk/pull/690) fix to prevent accepting file name
* (cli) [\#708](https://github.com/Finschia/finschia-sdk/pull/708) In CLI, allow 1 SIGN_MODE_DIRECT signer in transactions with multiple signers.
* (x/modules) [\#722](https://github.com/Finschia/finschia-sdk/pull/722) Check error for `RegisterQueryHandlerClient` in all modules `RegisterGRPCGatewayRoutes`
* (x/bank) [\#716](https://github.com/Finschia/finschia-sdk/pull/716) remove useless DenomMetadata key function
* (x/foundation) [\#704](https://github.com/Finschia/finschia-sdk/pull/704) update x/foundation params
* (x/wasm)  [\#695](https://github.com/Finschia/finschia-sdk/pull/695) fix to prevent external filesystem dependency of simulation
* (x/foundation) [\#729](https://github.com/Finschia/finschia-sdk/pull/729) add UpdateParams to x/foundation
* (amino) [\#736](https://github.com/Finschia/finschia-sdk/pull/736) apply the missing amino codec registratoin of cosmos-sdk
* (x/foundation) [\#744](https://github.com/Finschia/finschia-sdk/pull/744) revisit foundation operator
* (store,x/wasm) [\#742](https://github.com/Finschia/finschia-sdk/pull/742) fix to add error message in GetByteCode()
* (amino) [\#745](https://github.com/Finschia/finschia-sdk/pull/745) apply the missing amino codec of `x/token`, `x/collection`, `x/wasm` and `x/foundation`
* (x/foundation) [\#757](https://github.com/Finschia/finschia-sdk/pull/757) remove redundant granter from x/foundation events

### Bug Fixes
* (x/foundation) [\#772](https://github.com/Finschia/finschia-sdk/pull/772) export x/foundation pool
* (baseapp) [\#781](https://github.com/Finschia/finschia-sdk/pull/781) implement method `SetOption()` in baseapp
* (global) [\#782](https://github.com/Finschia/finschia-sdk/pull/782) add unhandled return error handling
* (x/collection,x/token) [\#798](https://github.com/Finschia/finschia-sdk/pull/798) Fix x/collection ModifyContract
* (ci) [\#803](https://github.com/Finschia/finschia-sdk/pull/803) fix test flow to install libsodium
* (x/collection,token) [\#784](https://github.com/Finschia/finschia-sdk/pull/784) Make field key matching in x/token & x/collection case-sensitive
* (x/wasm) [\#453](https://github.com/Finschia/finschia-sdk/pull/453) modify wasm grpc query api path
* (client) [\#476](https://github.com/Finschia/finschia-sdk/pull/476) change the default value of the client output format in the config
* (server/grpc) [\#516](https://github.com/Finschia/finschia-sdk/pull/516) restore build norace flag
* (genesis) [\#517](https://github.com/Finschia/finschia-sdk/pull/517) fix genesis auth account format(cosmos-sdk style -> lbm-sdk style)
* (x/token) [\#539](https://github.com/Finschia/finschia-sdk/pull/539) fix the compatibility issues with daphne
* (x/foundation) [\#545](https://github.com/Finschia/finschia-sdk/pull/545) fix genesis and support abstain
* (x/auth) [\#563](https://github.com/Finschia/finschia-sdk/pull/563) fix unmarshal bug of `BaseAccountJSON`
* (client) [\#565](https://github.com/Finschia/finschia-sdk/pull/565) fix the data race problem in `TestQueryABCIHeight`
* (x/token) [\#589](https://github.com/Finschia/finschia-sdk/pull/589) fix naming collision in x/token enums
* (x/token) [\#599](https://github.com/Finschia/finschia-sdk/pull/599) fix the order of events
* (x/wasm) [\#640](https://github.com/Finschia/finschia-sdk/pull/640) remove legacy codes of wasm
* (amino) [\#635](https://github.com/Finschia/finschia-sdk/pull/635) change some minor things that haven't been fixed in #549
* (store) [\#666](https://github.com/Finschia/finschia-sdk/pull/666) change default `iavl-cache-size` and description
* (x/auth) [\#673](https://github.com/Finschia/finschia-sdk/pull/673) fix max gas validation
* (simapp) [\#679](https://github.com/Finschia/finschia-sdk/pull/679) fix the bug not setting `iavl-cache-size` value of `app.toml`
* (x/foundation) [\#687](https://github.com/Finschia/finschia-sdk/pull/687) fix bugs on aborting x/foundation proposals
* (global) [\#694](https://github.com/Finschia/finschia-sdk/pull/694) replace deprecated functions since go 1.16 or 1.17
* (x/bankplus) [\#705](https://github.com/Finschia/finschia-sdk/pull/705) add missing blockedAddr checking in bankplus
* (x/foundation) [\#712](https://github.com/Finschia/finschia-sdk/pull/712) fix x/foundation EndBlocker
* (x/feegrant) [\#720](https://github.com/Finschia/finschia-sdk/pull/720) remove potential runtime panic in x/feegrant
* (baseapp) [\#724](https://github.com/Finschia/finschia-sdk/pull/724) add checking pubkey type from validator params
* (x/staking) [\#726](https://github.com/Finschia/finschia-sdk/pull/726) check allowedList size in StakeAuthorization.Accept()
* (x/staking) [\#728](https://github.com/Finschia/finschia-sdk/pull/728) fix typo in unbondingToUnbonded() panic
* (crypto) [\#731](https://github.com/Finschia/finschia-sdk/pull/731) remove VRFProve function
* (x/foundation) [\#732](https://github.com/Finschia/finschia-sdk/pull/732) add verification on accounts into x/foundation Grants cli
* (x/foundation) [\#730](https://github.com/Finschia/finschia-sdk/pull/730) prune stale x/foundation proposals at voting period end
* (cli) [\#734](https://github.com/Finschia/finschia-sdk/pull/734) add restrictions on the number of args in the CLIs
* (client) [\#737](https://github.com/Finschia/finschia-sdk/pull/737) check multisig key list to prevent unexpected key deletion
* (simapp) [\#752](https://github.com/Finschia/finschia-sdk/pull/752) add x/distribution's module account into blockedAddr
* (x/auth) [\#754](https://github.com/Finschia/finschia-sdk/pull/754) Fix wrong sequences in `sign-batch`
* (x/foundation) [\#761](https://github.com/Finschia/finschia-sdk/pull/761) restore build norace flag
* (server) [\#763](https://github.com/Finschia/finschia-sdk/pull/763) start telemetry independently from the API server

### Breaking Changes
* (cli) [\#773](https://github.com/Finschia/finschia-sdk/pull/773) guide users to use generate-only in messages for x/foundation authority
* (x/foundation) [\#790](https://github.com/Finschia/finschia-sdk/pull/790) fix case of gov_mint_left_count in x/foundation
* (proto) [\#564](https://github.com/Finschia/finschia-sdk/pull/564) change gRPC path to original cosmos path
* (global) [\#603](https://github.com/Finschia/finschia-sdk/pull/603) apply types/address.go from cosmos-sdk@v0.45.1
* (amino) [\#600](https://github.com/Finschia/finschia-sdk/pull/600) change amino codec path from `lbm-sdk/` to `cosmos-sdk/`
* (ostracon) [\#610](https://github.com/Finschia/finschia-sdk/pull/610) apply change of prefix of key name in ostracon
* (ostracon) [\#614](https://github.com/Finschia/finschia-sdk/pull/614) apply Ostracon's changes that replace `StakingPower` with `VotingPower` and `StakingPower` with `VotingPower`
* (proto) [\#617](https://github.com/Finschia/finschia-sdk/pull/617) change wasm gRPC path to original `cosmwasm` path.
* (proto) [\#627](https://github.com/Finschia/finschia-sdk/pull/627) revert changes in x/slashing proto

### Build, CI
* (ci) [\#779](https://github.com/Finschia/finschia-sdk/pull/779) change github action trigger rules for `release/*` and `rc*/*` branches
* (ci) [\#457](https://github.com/Finschia/finschia-sdk/pull/457), [\#471](https://github.com/Finschia/finschia-sdk/pull/471) add swagger check
* (ci) [\#469](https://github.com/Finschia/finschia-sdk/pull/469) publish docker image on tag push
* (ci) [\#580](https://github.com/Finschia/finschia-sdk/pull/580) fix the problem that the registered docker image couldn't run on M1.
* (simapp) [\#591](https://github.com/Finschia/finschia-sdk/pull/591) chore: add x/wasm module to simapp
* (ci) [\#618](https://github.com/Finschia/finschia-sdk/pull/618) remove stale action
* (ci) [\#619](https://github.com/Finschia/finschia-sdk/pull/619) change the Dockerfile to use the downloaded static library

### Document Updates
* (docs) [\#766](https://github.com/Finschia/finschia-sdk/pull/766) fix submit-proposal command on x/foundation
* (docs) [\#802](https://github.com/Finschia/finschia-sdk/pull/802) update x/foundation documentation
* (docs) [\#483](https://github.com/Finschia/finschia-sdk/pull/483) update documents on x/stakingplus
* (docs) [\#490](https://github.com/Finschia/finschia-sdk/pull/490) update documents on x/consortium
* (docs) [\#602](https://github.com/Finschia/finschia-sdk/pull/602) update outdated events in specs
* (docs) [\#721](https://github.com/Finschia/finschia-sdk/pull/721) update x/foundation specification
* (docs) [\#748](https://github.com/Finschia/finschia-sdk/pull/748) add `GovMint` to x/foundation specification


## [v0.45.0](https://github.com/Finschia/finschia-sdk/releases/tag/v0.45.0)

### Features
* (x/wasm) [\#444](https://github.com/Finschia/finschia-sdk/pull/444) Merge wasmd 0.19.0
    * remove custom encoder from x/wasm/keeper.NewKeeper's arg. After the Token/collection module is added, it will be ported again.
* (cosmos-sdk) [\#437](https://github.com/Finschia/finschia-sdk/pull/437) dump up to cosmos-sdk v0.42.11
    * [changelog of cosmos-sdk v0.42.11](https://github.com/cosmos/cosmos-sdk/blob/v0.42.11/CHANGELOG.md)
* (feat) [\#434](https://github.com/Finschia/finschia-sdk/pull/434) Revert signature mechanism using `sig_block_height`
* (x/token) [\#416](https://github.com/Finschia/finschia-sdk/pull/416) Migrate token module from line/link

### Bug Fixes
* (x/wasm) [\#453](https://github.com/Finschia/finschia-sdk/pull/453) modify wasm grpc query api path
* (x/wasm) [\#436](https://github.com/Finschia/finschia-sdk/pull/436) remove `x/wasm/linkwasmd`

### Build, CI
* (makefile, ci) [\#438](https://github.com/Finschia/finschia-sdk/pull/438) fix `make proto-format` and `make proto-check-breaking` error


## [v0.44.0-rc0](https://github.com/Finschia/finschia-sdk/releases/tag/v0.44.0-rc0)

### Features
* (feat) [\#352] (https://github.com/Finschia/finschia-sdk/pull/352) iavl, db & disk stats logging
* (x/gov) [\#368](https://github.com/Finschia/finschia-sdk/pull/368) Governance Split Votes, use `MsgWeightedVote` to send a split vote. Sending a regular `MsgVote` will convert the underlying vote option into a weighted vote with weight 1.
* (x/upgrade) [\#377] (https://github.com/Finschia/finschia-sdk/pull/377) To smoothen the update to the latest stable release, the SDK includes vesion map for managing migrations between SDK versions.
* (x/wasm) [\#358] (https://github.com/Finschia/finschia-sdk/pull/358) change wasm metrics method to using prometheus directly
* (x/feegrant) [\#380] (https://github.com/Finschia/finschia-sdk/pull/380) Feegrant module
* (x/wasm) [\#395] (https://github.com/Finschia/finschia-sdk/pull/395) Add the instantiate_permission in the CodeInfoResponse
* (x/consortium) [\#406] (https://github.com/Finschia/finschia-sdk/pull/406) Add CreateValidator access control feature
* (x/bank) [\#400] (https://github.com/Finschia/finschia-sdk/pull/400) add `bankplus` function to restrict to send coin to inactive smart contract.

### Improvements
* (slashing) [\#347](https://github.com/Finschia/finschia-sdk/pull/347) Introduce VoterSetCounter
* (auth) [\#348](https://github.com/Finschia/finschia-sdk/pull/348) Increase default valid_sig_block_period

### Bug Fixes
* (x/feegrant) [\#383] (https://github.com/Finschia/finschia-sdk/pull/383) Update allowance inside AllowedMsgAllowance
* (tm-db) [\#388] (https://github.com/Finschia/finschia-sdk/pull/388) Bump up tm-db fixing invalid memory reference
* (swagger) [\#391] (https://github.com/Finschia/finschia-sdk/pull/391) fix swagger's config path for wasm
* (x/wasm) [\#393] (https://github.com/Finschia/finschia-sdk/pull/393) fix bug where `StoreCodeAndInstantiateContract`, `UpdateContractStatus`, `UpdateContractStatusProposal` API does not work
* (x/slashing) [\#407] (https://github.com/Finschia/finschia-sdk/pull/407) Fix query signing infos command

### Breaking Changes
* (x/consortium) [\#411] (https://github.com/Finschia/finschia-sdk/pull/411) Validate validator addresses in update-validator-auths proposal

### Build, CI
* (ci) [\#350](https://github.com/Finschia/finschia-sdk/pull/350) Reduce sim test time
* (ci) [\#351](https://github.com/Finschia/finschia-sdk/pull/351) Remove diff condition from sim-normal

### Document Updates
* (docs) [\#361](https://github.com/Finschia/finschia-sdk/pull/361) Add sample command docs
* (docs) [\#392](https://github.com/Finschia/finschia-sdk/pull/392) Modify with latest version of swagger REST interface docs.


## [v0.43.1](https://github.com/Finschia/finschia-sdk/releases/tag/v0.43.1)

### Bug Fixes
* (distribution) [\#364](https://github.com/Finschia/finschia-sdk/pull/364) Force genOrBroadcastFn even when max-msgs != 0
* (bank) [\#366](https://github.com/Finschia/finschia-sdk/pull/366) Check bech32 address format in bank query

## [v0.43.0](https://github.com/Finschia/finschia-sdk/releases/tag/v0.43.0)

### Features
* (global) [\#97](https://github.com/Finschia/finschia-sdk/pull/97) Add codespace to query error
* (config) [\#114](https://github.com/Finschia/finschia-sdk/pull/114) Add idle-timeout to rest server and rpc server config
* (x/wasm) [\#127](https://github.com/Finschia/finschia-sdk/pull/127) Add wasm with Staragate migration completed.
* (x/wasm) [\#151](https://github.com/Finschia/finschia-sdk/pull/151) Add contract access control.
* (x/wasm) [\#194](https://github.com/Finschia/finschia-sdk/pull/194) Replace importing CosmWasm/wasmvm with line/wasmvm.
* (x/auth) [\#176](https://github.com/Finschia/finschia-sdk/pull/176) Add MsgEmpty to auth module
* (metric) [\#184](https://github.com/Finschia/finschia-sdk/pull/184) Add prometheus metrics for caches reverting telemetry metrics
* (grpc) [\#291](https://github.com/Finschia/finschia-sdk/pull/291) Add GRPC API the BlockByHash and BlockResultsByHeight, add prove parameter to GetTxsEvent

### Improvements
* (bump-up) [\#93](https://github.com/Finschia/finschia-sdk/pull/93) Adopt ostracon, line/tm-db and line/iavl
* (bump-up) [\#107](https://github.com/Finschia/finschia-sdk/pull/107) Bump up tm-db, iavl and ostracon
* (script) [\#110](https://github.com/Finschia/finschia-sdk/pull/110) Add script initializing simd
* (bump-up) [\#118](https://github.com/Finschia/finschia-sdk/pull/118) Bump up tm-db and remove Domain() call
* (test) [\#128](https://github.com/Finschia/finschia-sdk/pull/128) Allow creating new test network without init
* (db) [\#136](https://github.com/Finschia/finschia-sdk/pull/136) Fix DB_BACKEND configuration
* (global) [\#140](https://github.com/Finschia/finschia-sdk/pull/140) Modify default coin type, default address prefix
* (perf) [\#141](https://github.com/Finschia/finschia-sdk/pull/141) Concurrent `checkTx`
* (perf) [\#142](https://github.com/Finschia/finschia-sdk/pull/142) Implement `validateGasWanted()`
* (perf) [\#143](https://github.com/Finschia/finschia-sdk/pull/143) Signature verification cache
* (global) [\#145](https://github.com/Finschia/finschia-sdk/pull/145) Modify key type name
* (perf) [\#155](https://github.com/Finschia/finschia-sdk/pull/155) Concurrent recheckTx
* (global) [\#158](https://github.com/Finschia/finschia-sdk/pull/158) Remove tm-db dependency
* (x/wasm) [\#162](https://github.com/Finschia/finschia-sdk/pull/162) Add missed UpdateContractStatusProposal types
* (perf) [\#164](https://github.com/Finschia/finschia-sdk/pull/164) Sse fastcache
* (build) [\#181](https://github.com/Finschia/finschia-sdk/pull/181) Raise codecov-action version to 1.5.0
* (build) [\#195](https://github.com/Finschia/finschia-sdk/pull/195) Build properly when using libsecp256k1
* (perf) [\#198](https://github.com/Finschia/finschia-sdk/pull/198) Caching paramset
* (global) [\#200](https://github.com/Finschia/finschia-sdk/pull/200) Add a env prefix
* (store) [\#202](https://github.com/Finschia/finschia-sdk/pull/202) Param store doesn't use gas kv
* (store) [\#203](https://github.com/Finschia/finschia-sdk/pull/203) Remove transient store that is not used now
* (perf) [\#204](https://github.com/Finschia/finschia-sdk/pull/204) Apply rw mutex to cachekv
* (perf) [\#208](https://github.com/Finschia/finschia-sdk/pull/208) Use easyjson instead of amino when marshal abci logs 
* (perf) [\#209](https://github.com/Finschia/finschia-sdk/pull/209) Apply async reactor ostracon
* (proto) [\#212](https://github.com/Finschia/finschia-sdk/pull/212) Reformat proto files and restore proto docs
* (perf) [\#216](https://github.com/Finschia/finschia-sdk/pull/216) Memoize bech32 encoding and decoding
* (perf) [\#218](https://github.com/Finschia/finschia-sdk/pull/218) Rootmulti store parallel commit
* (perf) [\#219](https://github.com/Finschia/finschia-sdk/pull/219) Fix bech32 cache to get bech32 from proper cache
* (bump-up) [\#221](https://github.com/Finschia/finschia-sdk/pull/221) Bump up iavl for parallel processing of batches
* (perf) [\#224](https://github.com/Finschia/finschia-sdk/pull/224) Updated log time to have milliseconds
* (bump-up) [\#228](https://github.com/Finschia/finschia-sdk/pull/228) Bump up ostracon to optimize checking the txs size
* (global) [\#230](https://github.com/Finschia/finschia-sdk/pull/230) Modify module name to lfb-sdk
* (bump-up) [\#246](https://github.com/Finschia/finschia-sdk/pull/246) Bump up ostracon to not flush wal when receive consensus msgs
* (wasm) [\#250](https://github.com/Finschia/finschia-sdk/pull/250) Migrate linkwasmd to the latest commit
* (wasm) [\#253](https://github.com/Finschia/finschia-sdk/pull/253) remove MaxGas const
* (wasm) [\#254](https://github.com/Finschia/finschia-sdk/pull/254) Specify wasm event types
* (x) [\#255](https://github.com/Finschia/finschia-sdk/pull/255) Remove legacy from modules
* (perf) [\#320](https:/github.com/Finschia/finschia-sdk/pull/320) internal objects optimization (BaseAccount, Balance & Supply)
* (auth) [\#344](https://github.com/Finschia/finschia-sdk/pull/344) move SigBlockHeight from TxBody into AuthInfo

### Bug Fixes
* (test) [\#92](https://github.com/Finschia/finschia-sdk/pull/92) Fix SendToModuleAccountTest
* (store) [\#105](https://github.com/Finschia/finschia-sdk/pull/105) Check `store == nil`
* (test) [\#133](https://github.com/Finschia/finschia-sdk/pull/133) Fix `Test_runImportCmd()`
* (config) [\#138](https://github.com/Finschia/finschia-sdk/pull/138) Fix getting coin type at running cmd 
* (race) [\#159](https://github.com/Finschia/finschia-sdk/pull/159) Fix test-race failure
* (test) [\#193](https://github.com/Finschia/finschia-sdk/pull/193) Allow to add new validator in test network
* (client) [\#286](https://github.com/Finschia/finschia-sdk/pull/286) Fix invalid type casting for error
* (test) [\#326](https://github.com/Finschia/finschia-sdk/pull/326) Enable sim test and fix address related bug
 
### Breaking Changes
* (global) [\#90](https://github.com/Finschia/finschia-sdk/pull/90) Revise module path to `github.com/line/lfb-sdk`
* (rpc) [\#97](https://github.com/Finschia/finschia-sdk/pull/97) Send response with 404 status when quering non-exist account
* (proto) [\#106](https://github.com/Finschia/finschia-sdk/pull/106) Rename package of proto files
* (api) [\#130](https://github.com/Finschia/finschia-sdk/pull/130) Rename rest apis
* (auth) [\#265](https://github.com/Finschia/finschia-sdk/pull/265) Introduce sig block height for the new replay protection
* (global) [\#298](https://github.com/Finschia/finschia-sdk/pull/298) Treat addresses as strings
* (ostracon) [\#317](https://github.com/Finschia/finschia-sdk/pull/317) Integrate Ostracon including vrf election and voter concept
* (global) [\#323](https://github.com/Finschia/finschia-sdk/pull/323) Re-brand lfb-sdk to lbm-sdk
* (proto) [\#338](https://github.com/Finschia/finschia-sdk/pull/338) Upgrade proto buf from v1beta1 to v1

### Build, CI
* (ci) [\#234](https://github.com/Finschia/finschia-sdk/pull/234) Fix branch name in ci script
* (docker) [\#264](https://github.com/Finschia/finschia-sdk/pull/264) Remove docker publish
* (ci) [\#345](https://github.com/Finschia/finschia-sdk/pull/345) Split long sim test into 3 parts
 
### Document Updates
* (docs) [\#205](https://github.com/Finschia/finschia-sdk/pull/205) Renewal docs for open source
* (docs) [\#207](https://github.com/Finschia/finschia-sdk/pull/207) Fix license
* (docs) [\#211](https://github.com/Finschia/finschia-sdk/pull/211) Remove codeowners
* (docs) [\#248](https://github.com/Finschia/finschia-sdk/pull/248) Add PR procedure, apply main branch
* (docs) [\#256](https://github.com/Finschia/finschia-sdk/pull/256) Modify copyright and contributing
* (docs) [\#259](https://github.com/Finschia/finschia-sdk/pull/259) Modify copyright, verified from legal team
* (docs) [\#260](https://github.com/Finschia/finschia-sdk/pull/260) Remove gov, ibc and readme of wasm module
* (docs) [\#262](https://github.com/Finschia/finschia-sdk/pull/262) Fix link urls, remove invalid reference
* (docs) [\#328](https://github.com/Finschia/finschia-sdk/pull/328) Update quick start guide

## [cosmos-sdk v0.42.1] - 2021-03-15
Initial lbm-sdk is based on the cosmos-sdk v0.42.1

* (cosmos-sdk) [v0.42.0](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.42.1).

Please refer [CHANGELOG_OF_COSMOS_SDK_v0.42.1](https://github.com/cosmos/cosmos-sdk/blob/v0.42.1/CHANGELOG.md)
<!-- Release links -->
