# Changelog

## [v0.45.0-rc0](https://github.com/line/lbm-sdk/releases/tag/v0.45.0-rc0)

### Features

-   (x/wasm) [\#444](https://github.com/line/lbm-sdk/pull/444) Merge wasmd 0.19.0
    -   remove custom encoder from x/wasm/keeper.NewKeeper's arg. After the Token/collection module is added, it will be ported again.
-   (cosmos-sdk) [\#437](https://github.com/line/lbm-sdk/pull/437) dump up to cosmos-sdk v0.42.11
    -   [changelog of cosmos-sdk v0.42.11](https://github.com/cosmos/cosmos-sdk/blob/v0.42.11/CHANGELOG.md)
-   (feat) [\#434](https://github.com/line/lbm-sdk/pull/434) Revert signature mechanism using `sig_block_height`
-   (x/token) [\#416](https://github.com/line/lbm-sdk/pull/416) Migrate token module from line/link

### Bug Fixes

-   (x/wasm) [\#436](https://github.com/line/lbm-sdk/pull/436) remove `x/wasm/linkwasmd`

### Build, CI

-   (makefile, ci) [\#438](https://github.com/line/lbm-sdk/pull/438) fix `make proto-format` and `make proto-check-breaking` error

## [v0.44.0-rc0](https://github.com/line/lbm-sdk/releases/tag/v0.44.0-rc0)

### Features

-   (feat) [\#352] (https://github.com/line/lbm-sdk/pull/352) iavl, db & disk stats logging
-   (x/gov) [\#368](https://github.com/line/lbm-sdk/pull/368) Governance Split Votes, use `MsgWeightedVote` to send a split vote. Sending a regular `MsgVote` will convert the underlying vote option into a weighted vote with weight 1.
-   (x/upgrade) [\#377] (https://github.com/line/lbm-sdk/pull/377) To smoothen the update to the latest stable release, the SDK includes vesion map for managing migrations between SDK versions.
-   (x/wasm) [\#358] (https://github.com/line/lbm-sdk/pull/358) change wasm metrics method to using prometheus directly
-   (x/feegrant) [\#380] (https://github.com/line/lbm-sdk/pull/380) Feegrant module
-   (x/wasm) [\#395] (https://github.com/line/lbm-sdk/pull/395) Add the instantiate_permission in the CodeInfoResponse
-   (x/consortium) [\#406] (https://github.com/line/lbm-sdk/pull/406) Add CreateValidator access control feature
-   (x/bank) [\#400] (https://github.com/line/lbm-sdk/pull/400) add `bankplus` function to restrict to send coin to inactive smart contract.

### Improvements

-   (slashing) [\#347](https://github.com/line/lbm-sdk/pull/347) Introduce VoterSetCounter
-   (auth) [\#348](https://github.com/line/lbm-sdk/pull/348) Increase default valid_sig_block_period

### Bug Fixes

-   (x/feegrant) [\#383] (https://github.com/line/lbm-sdk/pull/383) Update allowance inside AllowedMsgAllowance
-   (tm-db) [\#388] (https://github.com/line/lbm-sdk/pull/388) Bump up tm-db fixing invalid memory reference
-   (swagger) [\#391] (https://github.com/line/lbm-sdk/pull/391) fix swagger's config path for wasm
-   (x/wasm) [\#393] (https://github.com/line/lbm-sdk/pull/393) fix bug where `StoreCodeAndInstantiateContract`, `UpdateContractStatus`, `UpdateContractStatusProposal` API does not work
-   (x/slashing) [\#407] (https://github.com/line/lbm-sdk/pull/407) Fix query signing infos command

### Breaking Changes

-   (x/consortium) [\#411] (https://github.com/line/lbm-sdk/pull/411) Validate validator addresses in update-validator-auths proposal

### Build, CI

-   (ci) [\#350](https://github.com/line/lbm-sdk/pull/350) Reduce sim test time
-   (ci) [\#351](https://github.com/line/lbm-sdk/pull/351) Remove diff condition from sim-normal

### Document Updates

-   (docs) [\#361](https://github.com/line/lbm-sdk/pull/361) Add sample command docs
-   (docs) [\#392](https://github.com/line/lbm-sdk/pull/392) Modify with latest version of swagger REST interface docs.

## [v0.43.1](https://github.com/line/lbm-sdk/releases/tag/v0.43.1)

### Bug Fixes

-   (distribution) [\#364](https://github.com/line/lbm-sdk/pull/364) Force genOrBroadcastFn even when max-msgs != 0
-   (bank) [\#366](https://github.com/line/lbm-sdk/pull/366) Check bech32 address format in bank query

## [v0.43.0](https://github.com/line/lbm-sdk/releases/tag/v0.43.0)

### Features

-   (global) [\#97](https://github.com/line/lbm-sdk/pull/97) Add codespace to query error
-   (config) [\#114](https://github.com/line/lbm-sdk/pull/114) Add idle-timeout to rest server and rpc server config
-   (x/wasm) [\#127](https://github.com/line/lbm-sdk/pull/127) Add wasm with Staragate migration completed.
-   (x/wasm) [\#151](https://github.com/line/lbm-sdk/pull/151) Add contract access control.
-   (x/wasm) [\#194](https://github.com/line/lbm-sdk/pull/194) Replace importing CosmWasm/wasmvm with line/wasmvm.
-   (x/auth) [\#176](https://github.com/line/lbm-sdk/pull/176) Add MsgEmpty to auth module
-   (metric) [\#184](https://github.com/line/lbm-sdk/pull/184) Add prometheus metrics for caches reverting telemetry metrics
-   (grpc) [\#291](https://github.com/line/lbm-sdk/pull/291) Add GRPC API the BlockByHash and BlockResultsByHeight, add prove parameter to GetTxsEvent

### Improvements

-   (bump-up) [\#93](https://github.com/line/lbm-sdk/pull/93) Adopt ostracon, line/tm-db and line/iavl
-   (bump-up) [\#107](https://github.com/line/lbm-sdk/pull/107) Bump up tm-db, iavl and ostracon
-   (script) [\#110](https://github.com/line/lbm-sdk/pull/110) Add script initializing simd
-   (bump-up) [\#118](https://github.com/line/lbm-sdk/pull/118) Bump up tm-db and remove Domain() call
-   (test) [\#128](https://github.com/line/lbm-sdk/pull/128) Allow creating new test network without init
-   (db) [\#136](https://github.com/line/lbm-sdk/pull/136) Fix DB_BACKEND configuration
-   (global) [\#140](https://github.com/line/lbm-sdk/pull/140) Modify default coin type, default address prefix
-   (perf) [\#141](https://github.com/line/lbm-sdk/pull/141) Concurrent `checkTx`
-   (perf) [\#142](https://github.com/line/lbm-sdk/pull/142) Implement `validateGasWanted()`
-   (perf) [\#143](https://github.com/line/lbm-sdk/pull/143) Signature verification cache
-   (global) [\#145](https://github.com/line/lbm-sdk/pull/145) Modify key type name
-   (perf) [\#155](https://github.com/line/lbm-sdk/pull/155) Concurrent recheckTx
-   (global) [\#158](https://github.com/line/lbm-sdk/pull/158) Remove tm-db dependency
-   (x/wasm) [\#162](https://github.com/line/lbm-sdk/pull/162) Add missed UpdateContractStatusProposal types
-   (perf) [\#164](https://github.com/line/lbm-sdk/pull/164) Sse fastcache
-   (build) [\#181](https://github.com/line/lbm-sdk/pull/181) Raise codecov-action version to 1.5.0
-   (build) [\#195](https://github.com/line/lbm-sdk/pull/195) Build properly when using libsecp256k1
-   (perf) [\#198](https://github.com/line/lbm-sdk/pull/198) Caching paramset
-   (global) [\#200](https://github.com/line/lbm-sdk/pull/200) Add a env prefix
-   (store) [\#202](https://github.com/line/lbm-sdk/pull/202) Param store doesn't use gas kv
-   (store) [\#203](https://github.com/line/lbm-sdk/pull/203) Remove transient store that is not used now
-   (perf) [\#204](https://github.com/line/lbm-sdk/pull/204) Apply rw mutex to cachekv
-   (perf) [\#208](https://github.com/line/lbm-sdk/pull/208) Use easyjson instead of amino when marshal abci logs
-   (perf) [\#209](https://github.com/line/lbm-sdk/pull/209) Apply async reactor ostracon
-   (proto) [\#212](https://github.com/line/lbm-sdk/pull/212) Reformat proto files and restore proto docs
-   (perf) [\#216](https://github.com/line/lbm-sdk/pull/216) Memoize bech32 encoding and decoding
-   (perf) [\#218](https://github.com/line/lbm-sdk/pull/218) Rootmulti store parallel commit
-   (perf) [\#219](https://github.com/line/lbm-sdk/pull/219) Fix bech32 cache to get bech32 from proper cache
-   (bump-up) [\#221](https://github.com/line/lbm-sdk/pull/221) Bump up iavl for parallel processing of batches
-   (perf) [\#224](https://github.com/line/lbm-sdk/pull/224) Updated log time to have milliseconds
-   (bump-up) [\#228](https://github.com/line/lbm-sdk/pull/228) Bump up ostracon to optimize checking the txs size
-   (global) [\#230](https://github.com/line/lbm-sdk/pull/230) Modify module name to lfb-sdk
-   (bump-up) [\#246](https://github.com/line/lbm-sdk/pull/246) Bump up ostracon to not flush wal when receive consensus msgs
-   (wasm) [\#250](https://github.com/line/lbm-sdk/pull/250) Migrate linkwasmd to the latest commit
-   (wasm) [\#253](https://github.com/line/lbm-sdk/pull/253) remove MaxGas const
-   (wasm) [\#254](https://github.com/line/lbm-sdk/pull/254) Specify wasm event types
-   (x) [\#255](https://github.com/line/lbm-sdk/pull/255) Remove legacy from modules
-   (perf) [\#320](https:/github.com/line/lbm-sdk/pull/320) internal objects optimization (BaseAccount, Balance & Supply)
-   (auth) [\#344](https://github.com/line/lbm-sdk/pull/344) move SigBlockHeight from TxBody into AuthInfo

### Bug Fixes

-   (test) [\#92](https://github.com/line/lbm-sdk/pull/92) Fix SendToModuleAccountTest
-   (store) [\#105](https://github.com/line/lbm-sdk/pull/105) Check `store == nil`
-   (test) [\#133](https://github.com/line/lbm-sdk/pull/133) Fix `Test_runImportCmd()`
-   (config) [\#138](https://github.com/line/lbm-sdk/pull/138) Fix getting coin type at running cmd
-   (race) [\#159](https://github.com/line/lbm-sdk/pull/159) Fix test-race failure
-   (test) [\#193](https://github.com/line/lbm-sdk/pull/193) Allow to add new validator in test network
-   (client) [\#286](https://github.com/line/lbm-sdk/pull/286) Fix invalid type casting for error
-   (test) [\#326](https://github.com/line/lbm-sdk/pull/326) Enable sim test and fix address related bug

### Breaking Changes

-   (global) [\#90](https://github.com/line/lbm-sdk/pull/90) Revise module path to `github.com/line/lfb-sdk`
-   (rpc) [\#97](https://github.com/line/lbm-sdk/pull/97) Send response with 404 status when quering non-exist account
-   (proto) [\#106](https://github.com/line/lbm-sdk/pull/106) Rename package of proto files
-   (api) [\#130](https://github.com/line/lbm-sdk/pull/130) Rename rest apis
-   (auth) [\#265](https://github.com/line/lbm-sdk/pull/265) Introduce sig block height for the new replay protection
-   (global) [\#298](https://github.com/line/lbm-sdk/pull/298) Treat addresses as strings
-   (ostracon) [\#317](https://github.com/line/lbm-sdk/pull/317) Integrate Ostracon including vrf election and voter concept
-   (global) [\#323](https://github.com/line/lfb-sdk/pull/323) Re-brand lfb-sdk to lbm-sdk
-   (proto) [\#338](https://github.com/line/lbm-sdk/pull/338) Upgrade proto buf from v1beta1 to v1

### Build, CI

-   (ci) [\#234](https://github.com/line/lbm-sdk/pull/234) Fix branch name in ci script
-   (docker) [\#264](https://github.com/line/lbm-sdk/pull/264) Remove docker publish
    <<<<<<< HEAD
    =======
-   (ci) [\#345](https://github.com/line/lbm-sdk/pull/345) Split long sim test into 3 parts
    > > > > > > > origin/main

### Document Updates

-   (docs) [\#205](https://github.com/line/lbm-sdk/pull/205) Renewal docs for open source
-   (docs) [\#207](https://github.com/line/lbm-sdk/pull/207) Fix license
-   (docs) [\#211](https://github.com/line/lbm-sdk/pull/211) Remove codeowners
-   (docs) [\#248](https://github.com/line/lbm-sdk/pull/248) Add PR procedure, apply main branch
-   (docs) [\#256](https://github.com/line/lbm-sdk/pull/256) Modify copyright and contributing
-   (docs) [\#259](https://github.com/line/lbm-sdk/pull/259) Modify copyright, verified from legal team
-   (docs) [\#260](https://github.com/line/lbm-sdk/pull/260) Remove gov, ibc and readme of wasm module
-   (docs) [\#262](https://github.com/line/lbm-sdk/pull/262) Fix link urls, remove invalid reference
-   (docs) [\#328](https://github.com/line/lbm-sdk/pull/328) Update quick start guide

## [cosmos-sdk v0.42.1] - 2021-03-15

Initial lbm-sdk is based on the cosmos-sdk v0.42.1

-   (cosmos-sdk) [v0.42.0](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.42.1).

Please refer [CHANGELOG_OF_COSMOS_SDK_v0.42.1](https://github.com/cosmos/cosmos-sdk/blob/v0.42.1/CHANGELOG.md)

<!-- Release links -->
