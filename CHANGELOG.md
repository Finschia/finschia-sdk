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

## [Unreleased](https://github.com/Finschia/finschia-sdk/compare/v0.48.0...HEAD)

### Features
* (consensus) [\#1178](https://github.com/Finschia/finschia-sdk/pull/1178) change the consensus from Ostracon to Tendermint v0.34.24

### Improvements
* (docs) [\#1120](https://github.com/Finschia/finschia-sdk/pull/1120) Update links in x/foundation README.md
* (feat) [\#1121](https://github.com/Finschia/finschia-sdk/pull/1121) Add update-censorship cmd to x/foundation cli
* (server) [#1153](https://github.com/Finschia/finschia-sdk/pull/1153) remove grpc replace directive
* (crypto) [\#1163](https://github.com/Finschia/finschia-sdk/pull/1163) Update some secp256k1 logics with latest `dcrec`
* (x/crisis) [#1167](https://github.com/Finschia/finschia-sdk/pull/1167) Use `CacheContext()` in `AssertInvariants()`
* (chore) [\#1168](https://github.com/Finschia/finschia-sdk/pull/1168) Replace `ExactArgs(0)` with `NoArgs()` in `x/upgrade` module
* (server) [\#1175](https://github.com/Finschia/finschia-sdk/pull/1175) Use go embed for swagger
* (x/collection) [\#1287](https://github.com/Finschia/finschia-sdk/pull/1287) add nft id validation to MsgSendNFT

### Bug Fixes
* chore(deps) [\#1141](https://github.com/Finschia/finschia-sdk/pull/1141) Bump github.com/cosmos/ledger-cosmos-go from 0.12.2 to 0.13.2 to fix ledger signing issue
* (x/auth, x/slashing) [\#1179](https://github.com/Finschia/finschia-sdk/pull/1179) modify missing changes of converting to tendermint
* (x/auth) [#1274](https://github.com/Finschia/finschia-sdk/pull/1274) `ModuleAccount.Validate` now reports a nil `.BaseAccount` instead of panicking.
* (x/collection) [\#1276](https://github.com/Finschia/finschia-sdk/pull/1276) eliminates potential risk for Insufficient Sanity Check of tokenID in Genesis 
* (x/foundation) [\#1277](https://github.com/Finschia/finschia-sdk/pull/1277) add init logic of foundation module accounts to InitGenesis in order to eliminate potential panic

### Removed

### Breaking Changes
* (consensus) [\#1178](https://github.com/Finschia/finschia-sdk/pull/1178) change the consensus from Ostracon to Tendermint v0.34.24 

### State Machine Breaking

### Event Breaking Changes

### Build, CI
* (ci) [\#1078](https://github.com/Finschia/finschia-sdk/pull/1078) fix tag comments in github actions workflow docker.yml
* (repo) [\#1157](https://github.com/Finschia/finschia-sdk/pull/1157) setup CODEOWNERS and backport action
* (ci) [\#1160](https://github.com/Finschia/finschia-sdk/pull/1160) remove autopr ci

### Document Updates
* (docs) [\#1059](https://github.com/Finschia/finschia-sdk/pull/1059) create ERRORS.md for x/module
* (docs) [\#1083](https://github.com/Finschia/finschia-sdk/pull/1083) Add detailed explanation about default events
* (x/token,collection) [#1201](https://github.com/Finschia/finschia-sdk/pull/1201) Deprecate legacy features on x/token,collection
