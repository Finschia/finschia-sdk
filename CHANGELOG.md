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

## [Unreleased](https://github.com/Finschia/finschia-sdk/compare/v0.48.1...HEAD)

### Features

### Improvements
<<<<<<< HEAD
=======
* (docs) [\#1120](https://github.com/Finschia/finschia-sdk/pull/1120) Update links in x/foundation README.md
* (feat) [\#1121](https://github.com/Finschia/finschia-sdk/pull/1121) Add update-censorship cmd to x/foundation cli
* (server) [#1153](https://github.com/Finschia/finschia-sdk/pull/1153) remove grpc replace directive
* (crypto) [\#1163](https://github.com/Finschia/finschia-sdk/pull/1163) Update some secp256k1 logics with latest `dcrec`
* (x/crisis) [#1167](https://github.com/Finschia/finschia-sdk/pull/1167) Use `CacheContext()` in `AssertInvariants()`
* (chore) [\#1168](https://github.com/Finschia/finschia-sdk/pull/1168) Replace `ExactArgs(0)` with `NoArgs()` in `x/upgrade` module
* (server) [\#1175](https://github.com/Finschia/finschia-sdk/pull/1175) Use go embed for swagger
* (x/collection) [\#1287](https://github.com/Finschia/finschia-sdk/pull/1287) add nft id validation to MsgSendNFT
* (types) [\#1314](https://github.com/Finschia/finschia-sdk/pull/1314) replace IsEqual with Equal
>>>>>>> 143a76304 (fix: replace IsEqual with Equal (#1314))

### Bug Fixes
* (x/auth) [#1281](https://github.com/Finschia/finschia-sdk/pull/1281) `ModuleAccount.Validate` now reports a nil `.BaseAccount` instead of panicking. (backport #1274)
* (x/foundation) [\#1283](https://github.com/Finschia/finschia-sdk/pull/1283) add init logic of foundation module accounts to InitGenesis in order to eliminate potential panic (backport #1277)
* (x/collection) [\#1282](https://github.com/Finschia/finschia-sdk/pull/1282) eliminates potential risk for Insufficient Sanity Check of tokenID in Genesis (backport #1276)
* (x/collection) [\#1290](https://github.com/Finschia/finschia-sdk/pull/1290) export x/collection params into genesis (backport #1268)
* (x/foundation) [\#1295](https://github.com/Finschia/finschia-sdk/pull/1295) add missing error handling for migration

### Removed

### Breaking Changes

### State Machine Breaking

### Event Breaking Changes

### Build, CI
* (build) [#1298](https://github.com/Finschia/finschia-sdk/pull/1298) Set Finschia/ostracon version

### Document Updates
* (x/token,collection) [#1201](https://github.com/Finschia/finschia-sdk/pull/1201) Deprecate legacy features on x/token,collection
