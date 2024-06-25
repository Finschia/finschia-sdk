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

## [Unreleased](https://github.com/Finschia/finschia-sdk/compare/v0.49.1...HEAD)

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
* (x/fswap) [\#1363](https://github.com/Finschia/finschia-sdk/pull/1363) introduce new event for MakeSwapProposal
* (x/fbridge) [\#1366](https://github.com/Finschia/finschia-sdk/pull/1366) Set target denom as module parameters
* (x/fbridge) [\#1369](https://github.com/Finschia/finschia-sdk/pull/1369) Add the event of `SetBridgeStatus`
* (x/fswap) [\#1372](https://github.com/Finschia/finschia-sdk/pull/1372) support message based proposals
* (x/fswap) [\#1387](https://github.com/Finschia/finschia-sdk/pull/1387) add new Swap query to get a single swap
* (x/fswap) [\#1382](https://github.com/Finschia/finschia-sdk/pull/1382) add validation & unit tests in fswap module 
* (x/fbridge) [\#1395](https://github.com/Finschia/finschia-sdk/pull/1395) Return error instead of panic for behaviors triggered by client
* (x/fswap) [\#1396](https://github.com/Finschia/finschia-sdk/pull/1396) refactor to use snake_case in proto
* (x/fswap) [\#1391](https://github.com/Finschia/finschia-sdk/pull/1391) add cli_test for fswap module
* (x/fbridge) [\#1405](https://github.com/Finschia/finschia-sdk/pull/1405) Add CLI, gRPC, MsgServer tests
* (x/fswap) [\#1415](https://github.com/Finschia/finschia-sdk/pull/1415) add more testcases for fswap module
>>>>>>> 97219d5ee (test: add CLI, gRPC, MsgServer tests (#1405))

### Bug Fixes

### Removed

### Breaking Changes

### State Machine Breaking

### Event Breaking Changes

### Build, CI

### Document Updates
* (docs) [\#1429](https://github.com/Finschia/finschia-sdk/pull/1429) correct spec docs of fswap module
