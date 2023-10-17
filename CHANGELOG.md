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
"Event Breaking" for breaking events.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.
Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## [Unreleased](https://github.com/Finschia/finschia-sdk/compare/v0.47.0...HEAD)

### Features
* (x/auth) [\#1011](https://github.com/Finschia/finschia-sdk/pull/1011) add the api for querying next account number
* (server/grpc) [\#1017](https://github.com/Finschia/finschia-sdk/pull/1017) support custom r/w gRPC options (backport cosmos/cosmos-sdk#11889)
* (x/auth) [\#1085](https://github.com/Finschia/finschia-sdk/pull/1085) rollback GetBlockWithTxs of cosmos service.proto for compatibility with cosmos-sdk APIs
* (proto) [\#1087](https://github.com/Finschia/finschia-sdk/pull/1087) add tendermint query apis for compatibility with cosmos-sdk

### Improvements
* (third_party/proto) [\#1037](https://github.com/Finschia/finschia-sdk/pull/1037) change the proof.proto path to third_party/proto/confio
* (ostracon) [\#1057](https://github.com/Finschia/finschia-sdk/pull/1057) Bump up Ostracon from to v1.1.1
* (x/foundation) [\#1072](https://github.com/Finschia/finschia-sdk/pull/1072) Address generation of the empty coins in x/foundation (backport #952)
* (cli) [\#1086](https://github.com/Finschia/finschia-sdk/pull/1086) Fix for redundant key generation. With running kms, generating priv-key is unnecessary.
* (ostracon) [\#1089](https://github.com/Finschia/finschia-sdk/pull/1089) Bump up ostracon from v1.1.1 to v1.1.1-449aa3148b12
* (ostracon) [\#1099](https://github.com/Finschia/finschia-sdk/pull/1099) Remove libsodium vrf library.
* (refactor) [\#1114](https://github.com/Finschia/finschia-sdk/pull/1114) Check statistics and balance on x/collection mint and burn operations
* (x/token) [\#1128](https://github.com/Finschia/finschia-sdk/pull/1128) add more unittest for MsgIssue of x/token
* (x/token) [\#1129](https://github.com/Finschia/finschia-sdk/pull/1129) add more unittest for `MsgGrantPermission` and `MsgRevokePermission` of x/token
* (x/token) [\#1130](https://github.com/Finschia/finschia-sdk/pull/1130) Add more unittest for MsgMint, MsgBurn, MsgOperatorBurn, MsgModify of x/token
* (x/collection) [\#1131](https://github.com/Finschia/finschia-sdk/pull/1131) add additional unittest of x/collection(`MsgIssueFT`, `MsgMintFT`, `MsgBurnFT`)
* (x/collection) [\#1133](https://github.com/Finschia/finschia-sdk/pull/1133) Refactor unittest for `SendFT`, `OperatorSendFT`, `AuthorizeOperator`, and `RevokeOperator` to check more states
* (x/token) [\#1137](https://github.com/Finschia/finschia-sdk/pull/1137) Add test for event compatibility
* (ostracon) [\#1142](https://github.com/Finschia/finschia-sdk/pull/1142) Bump up ostracon from v1.1.2-0.20230822110903-449aa3148b12 to v1.1.2

### Bug Fixes
* (ledger) [\#1040](https://github.com/Finschia/finschia-sdk/pull/1040) Fix a bug(unable to connect nano S plus ledger on ubuntu)
* (x/foundation) [\#1053](https://github.com/Finschia/finschia-sdk/pull/1053) Make x/foundation MsgExec propagate events
* (baseapp) [\#1091](https://github.com/finschia/finschia-sdk/pull/1091) Add `events.GetAttributes` and `event.GetAttribute` methods to simplify the retrieval of an attribute from event(s) (backport #1075)
* (x/foundation) [\#1108](https://github.com/Finschia/finschia-sdk/pull/1108) Rollback MsgUpdateParams parts from #999

### Removed

### Breaking Changes

### State Machine Breaking
* (x/foundation) [\#999](https://github.com/Finschia/finschia-sdk/pull/999) migrate x/foundation FoundationTax into x/params
* (x/collection) [\#1102](https://github.com/finschia/finschia-sdk/pull/1102) Reject modifying NFT class with token index filled in MsgModify
* (x/collection) [\#1105](https://github.com/Finschia/finschia-sdk/pull/1105) Add minted coins to balance in x/collection MsgMintFT
* (x/collection) [\#1106](https://github.com/Finschia/finschia-sdk/pull/1106) Support x/collection migration on old chains

### Event Breaking Changes
* (refactor) [\#1090](https://github.com/Finschia/finschia-sdk/pull/1090) Automate EventTypeMessage inclusion in every message execution (backport #1063)
* (x/bank) [#1093](https://github.com/Finschia/finschia-sdk/pull/1093) Remove message events including `sender` attribute whose information is already present in the relevant events (backport #1066)
* (baseapp) [\#1092](https://github.com/finschia/finschia-sdk/pull/1092) Do not add `module` attribute in case of ibc messages (backport #1079)

### Build, CI
* (build,ci) [\#1043](https://github.com/Finschia/finschia-sdk/pull/1043) Update golang version to 1.20

### Document Updates
* (readme) [\#997](https://github.com/finschia/finschia-sdk/pull/997) fix swagger url
* (docs) [\#1094](https://github.com/Finschia/finschia-sdk/pull/1094) Document default events (backport #1081)
