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

## [Unreleased](https://github.com/line/lbm-sdk/compare/v0.45.0-rc0...HEAD)

### Features
* (x/wasm) [\#570](https://github.com/line/lbm-sdk/pull/570) Merge wasmd 0.27.0
* (x/wasm) [\#470](https://github.com/line/lbm-sdk/pull/470) remove contract activation control by actor
* (x/wasm) [\#513](https://github.com/line/lbm-sdk/pull/513) fix message representation for signing
* (x/foundation) [\#518](https://github.com/line/lbm-sdk/pull/518) add foundation treasury feature to x/foundation
* (x/foundation) [\#528](https://github.com/line/lbm-sdk/pull/528) add a feature of whitelist for /lbm.foundation.v1.MsgWithdrawFromTreasury
* (proto) [\#584](https://github.com/line/lbm-sdk/pull/564) remove `prove` field in the `GetTxsEventRequest` of `tx` proto

### Improvements

* (refactor) [\#493](https://github.com/line/lbm-sdk/pull/493) restructure x/consortium
* (server/grpc) [\#526](https://github.com/line/lbm-sdk/pull/526) add index field into TxResponse
* (cli) [\#535](https://github.com/line/lbm-sdk/pull/536) updated ostracon to v1.0.5; `unsafe-reset-all` command has been moved to the `ostracon` sub-command.

### Bug Fixes
* (x/wasm) [\#453](https://github.com/line/lbm-sdk/pull/453) modify wasm grpc query api path
* (client) [\#476](https://github.com/line/lbm-sdk/pull/476) change the default value of the client output format in the config
* (server/grpc) [\#516](https://github.com/line/lbm-sdk/pull/516) restore build norace flag
* (genesis) [\#517](https://github.com/line/lbm-sdk/pull/517) fix genesis auth account format(cosmos-sdk style -> lbm-sdk style)
* (x/token) [\#539](https://github.com/line/lbm-sdk/pull/539) fix the compatibility issues with daphne
* (x/foundation) [\#545](https://github.com/line/lbm-sdk/pull/545) fix genesis and support abstain
* (x/auth) [\#563](https://github.com/line/lbm-sdk/pull/563) fix unmarshal bug of `BaseAccountJSON`
* (client) [\#565](https://github.com/line/lbm-sdk/pull/565) fix the data race problem in `TestQueryABCIHeight`
* (client) [\#817](https://github.com/line/lbm-sdk/pull/817) remove support for composite (BLS) type

### Breaking Changes
* (proto) [\#564](https://github.com/line/lbm-sdk/pull/564) change gRPC path to original cosmos path

### Build, CI

* (ci) [\#457](https://github.com/line/lbm-sdk/pull/457), [\#471](https://github.com/line/lbm-sdk/pull/471) add swagger check
* (ci) [\#469](https://github.com/line/lbm-sdk/pull/469) publish docker image on tag push
* (ci) [\#580](https://github.com/line/lbm-sdk/pull/580) fix the problem that the registered docker image couldn't  run on M1.

### Document Updates

* (docs) [\#483](https://github.com/line/lbm-sdk/pull/483) update documents on x/stakingplus
* (docs) [\#490](https://github.com/line/lbm-sdk/pull/490) update documents on x/consortium
