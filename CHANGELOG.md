# Changelog

## [Unreleased]

### Features
* (global) [\#97](https://github.com/line/lfb-sdk/pull/97) Add codespace to query error
* (config) [\#114](https://github.com/line/lfb-sdk/pull/114) Add idle-timeout to rest server and rpc server config
* (x/wasm) [\#127](https://github.com/line/lfb-sdk/pull/127) Add wasm with Staragate migration completed.
* (x/wasm) [\#151](https://github.com/line/lfb-sdk/pull/151) Add contract access control.
* (x/wasm) [\#194](https://github.com/line/lfb-sdk/pull/194) Replace importing CosmWasm/wasmvm with line/wasmvm.
* (x/auth) [\#176](https://github.com/line/lfb-sdk/pull/176) Add MsgEmpty to auth module
* (metric) [\#184](https://github.com/line/lfb-sdk/pull/184) Add prometheus metrics for caches reverting telemetry metrics

### Improvements
* (bump-up) [\#93](https://github.com/line/lfb-sdk/pull/93) Adopt ostracon, line/tm-db and line/iavl
* (bump-up) [\#107](https://github.com/line/lfb-sdk/pull/107) Bump up tm-db, iavl and ostracon
* (script) [\#110](https://github.com/line/lfb-sdk/pull/110) Add script initializing simd
* (bump-up) [\#118](https://github.com/line/lfb-sdk/pull/118) Bump up tm-db and remove Domain() call
* (test) [\#128](https://github.com/line/lfb-sdk/pull/128) Allow creating new test network without init
* (db) [\#136](https://github.com/line/lfb-sdk/pull/136) Fix DB_BACKEND configuration
* (global) [\#140](https://github.com/line/lfb-sdk/pull/140) Modify default coin type, default address prefix
* (perf) [\#141](https://github.com/line/lfb-sdk/pull/141) Concurrent `checkTx`
* (perf) [\#142](https://github.com/line/lfb-sdk/pull/142) Implement `validateGasWanted()`
* (perf) [\#143](https://github.com/line/lfb-sdk/pull/143) Signature verification cache
* (global) [\#145](https://github.com/line/lfb-sdk/pull/145) Modify key type name
* (perf) [\#155](https://github.com/line/lfb-sdk/pull/155) Concurrent recheckTx
* (global) [\#158](https://github.com/line/lfb-sdk/pull/158) Remove tm-db dependency
* (x/wasm) [\#162](https://github.com/line/lfb-sdk/pull/162) Add missed UpdateContractStatusProposal types
* (perf) [\#164](https://github.com/line/lfb-sdk/pull/164) Sse fastcache
* (build) [\#181](https://github.com/line/lfb-sdk/pull/181) Raise codecov-action version to 1.5.0
* (build) [\#195](https://github.com/line/lfb-sdk/pull/195) Build properly when using libsecp256k1
* (perf) [\#198](https://github.com/line/lfb-sdk/pull/198) Caching paramset
* (global) [\#200](https://github.com/line/lfb-sdk/pull/200) Add a env prefix
* (store) [\#202](https://github.com/line/lfb-sdk/pull/202) param store doesn't use gas kv

### Bug Fixes
* (test) [\#92](https://github.com/line/lfb-sdk/pull/92) Fix SendToModuleAccountTest
* (store) [\#105](https://github.com/line/lfb-sdk/pull/105) Check `store == nil`
* (test) [\#133](https://github.com/line/lfb-sdk/pull/133) Fix `Test_runImportCmd()`
* (config) [\#138](https://github.com/line/lfb-sdk/pull/138) Fix getting coin type at running cmd 
* (race) [\#159](https://github.com/line/lfb-sdk/pull/159) Fix test-race failure
* (test) [\#193](https://github.com/line/lfb-sdk/pull/193) Allow to add new validator in test network

### Breaking Changes
* (global) [\#90](https://github.com/line/lfb-sdk/pull/90) Revise module path to `github.com/line/lfb-sdk`
* (rpc) [\#97](https://github.com/line/lfb-sdk/pull/97) Send response with 404 status when quering non-exist account
* (proto) [\#106](https://github.com/line/lfb-sdk/pull/106) Rename package of proto files
* (api) [\#130](https://github.com/line/lfb-sdk/pull/130) Rename rest apis

## [cosmos-sdk v0.42.1] - 2021-03-15
Initial lfb-sdk is based on the cosmos-sdk v0.42.1

* (cosmos-sdk) [v0.42.0](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.42.1).

Please refer [CHANGELOG_OF_COSMOS_SDK_v0.42.1](https://github.com/cosmos/cosmos-sdk/blob/v0.42.1/CHANGELOG.md)
<!-- Release links -->
