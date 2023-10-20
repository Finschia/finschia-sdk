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

## [Unreleased](https://github.com/Finschia/finschia-sdk/compare/v0.46.0...HEAD)

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
* (server, client) [\#1152](https://github.com/Finschia/finschia-sdk/pull/1152) remove grpc replace directive

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
* (x/auth) [\#982](https://github.com/Finschia/finschia-sdk/pull/957) Fix not to emit error when no txs in block while querying `GetBlockWithTxs`
* (x/foundation) [\#984](https://github.com/Finschia/finschia-sdk/pull/984) Revert #952

### Removed
* [\#853](https://github.com/Finschia/finschia-sdk/pull/853) remove useless stub BeginBlock, EndBlock methods from modules below
  * ibc, authz, collection, feegrant, ibc, token, wasm
* (x/ibc) [\#858](https://github.com/Finschia/finschia-sdk/pull/858) detach ibc module(repo: [line/ibc-go](https://github.com/Finschia/ibc-go))
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
