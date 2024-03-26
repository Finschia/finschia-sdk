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

### Bug Fixes
<<<<<<< HEAD
* (x/auth) [#1281](https://github.com/Finschia/finschia-sdk/pull/1281) `ModuleAccount.Validate` now reports a nil `.BaseAccount` instead of panicking. (backport #1274)
* (x/foundation) [\#1283](https://github.com/Finschia/finschia-sdk/pull/1283) add init logic of foundation module accounts to InitGenesis in order to eliminate potential panic (backport #1277)
* (x/collection) [\#1282](https://github.com/Finschia/finschia-sdk/pull/1282) eliminates potential risk for Insufficient Sanity Check of tokenID in Genesis (backport #1276)
* (x/collection) [\#1290](https://github.com/Finschia/finschia-sdk/pull/1290) export x/collection params into genesis (backport #1268)
* (x/foundation) [\#1295](https://github.com/Finschia/finschia-sdk/pull/1295) add missing error handling for migration
=======
* chore(deps) [\#1141](https://github.com/Finschia/finschia-sdk/pull/1141) Bump github.com/cosmos/ledger-cosmos-go from 0.12.2 to 0.13.2 to fix ledger signing issue
* (x/auth, x/slashing) [\#1179](https://github.com/Finschia/finschia-sdk/pull/1179) modify missing changes of converting to tendermint
* (x/auth) [#1274](https://github.com/Finschia/finschia-sdk/pull/1274) `ModuleAccount.Validate` now reports a nil `.BaseAccount` instead of panicking.
* (x/collection) [\#1276](https://github.com/Finschia/finschia-sdk/pull/1276) eliminates potential risk for Insufficient Sanity Check of tokenID in Genesis 
* (x/foundation) [\#1277](https://github.com/Finschia/finschia-sdk/pull/1277) add init logic of foundation module accounts to InitGenesis in order to eliminate potential panic
* (x/collection, x/token) [\#1288](https://github.com/Finschia/finschia-sdk/pull/1288) use accAddress to compare in validatebasic function in collection & token modules
* (x/collection) [\#1268](https://github.com/Finschia/finschia-sdk/pull/1268) export x/collection params into genesis
* (x/collection) [\#1294](https://github.com/Finschia/finschia-sdk/pull/1294) reject NFT coins on FT APIs
* (sec) [\#1302](https://github.com/Finschia/finschia-sdk/pull/1302) remove map iteration non-determinism with keys + sorting
* (client) [\#1303](https://github.com/Finschia/finschia-sdk/pull/1303) fix possible overflow in BuildUnsignedTx 
* (types) [\#1299](https://github.com/Finschia/finschia-sdk/pull/1299) add missing nil checks
* (x/staking) [\#1301](https://github.com/Finschia/finschia-sdk/pull/1301) Use bytes instead of string comparison in delete validator queue (backport cosmos/cosmos-sdk#12303)
>>>>>>> 1038a394c (fix: Use bytes instead of string comparison in delete validator queue (backport cosmos/cosmos-sdk#12303) (#1301))

### Removed

### Breaking Changes

### State Machine Breaking

### Event Breaking Changes

### Build, CI
* (build) [#1298](https://github.com/Finschia/finschia-sdk/pull/1298) Set Finschia/ostracon version

### Document Updates
* (x/token,collection) [#1201](https://github.com/Finschia/finschia-sdk/pull/1201) Deprecate legacy features on x/token,collection
