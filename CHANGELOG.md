<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking Protobuf, gRPC and REST routes used by end-users.
"CLI Breaking" for breaking CLI commands.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.
Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## [Unreleased]
* (x/wasm) [\#835](https://github.com/line/lbm-sdk/pull/835) add calculate instantiate cost function into get contract env api
* (global) [\#764](https://github.com/line/lbm-sdk/pull/764) update dependencies of packages in dynamic link branch
* (x/wasm) [\#669](https://github.com/line/lbm-sdk/pull/669) update wasmvm and update contracts for tests
* (x/wasm) [\#656](https://github.com/line/lbm-sdk/pull/656) add unit test for dynamic link when callee fails
* (global) [\#658](https://github.com/line/lbm-sdk/pull/658) Add benchmarking job to ci
* (x/wasm) [\#659](https://github.com/line/lbm-sdk/pull/659) Adjust gas cost for GetContractEnv of CosmWasmAPI
* (x/wasm) [\#633](https://github.com/line/lbm-sdk/pull/633) Add benchmarks for cosmwasm APIs
* (x/wasm) [\#509](https://github.com/line/lbm-sdk/pull/590) Add integration tests for dynamic call


### Features

### Improvements
* (cosmovisor) [\#792](https://github.com/line/lbm-sdk/pull/792) Use upstream's cosmovisor
* (server) [\#821](https://github.com/line/lbm-sdk/pull/821) Get validator pubkey considering KMS

### Bug Fixes
* (client) [\#817](https://github.com/line/lbm-sdk/pull/817) remove support for composite (BLS) type

* (x/foundation) [#834](https://github.com/line/lbm-sdk/pull/834) Apply foundation audit

### Breaking Changes
* (rest) [\#807](https://github.com/line/lbm-sdk/pull/807) remove legacy REST API

### Build, CI
* (ci) [\#829](https://github.com/line/lbm-sdk/pull/829) automate release process

### Document Updates
