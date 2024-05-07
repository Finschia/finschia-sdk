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
* (x/fswap) [\#1336](https://github.com/Finschia/finschia-sdk/pull/1336) Initialize fswap module
* (x/fswap) [\#1339](https://github.com/Finschia/finschia-sdk/pull/1339) Implement fswap module's genesis
* (x/fbridge) [\#1340](https://github.com/Finschia/finschia-sdk/pull/1340) Initialize fbridge module
* (x/fbridge) [\#1347](https://github.com/Finschia/finschia-sdk/pull/1347) Implement bridge transfer feature (sending side)
* (x/fbridge) [\#1351](https://github.com/Finschia/finschia-sdk/pull/1351) Map a sequence to block number for every bridge request (sending side)
* (x/fswap) [\#1345](https://github.com/Finschia/finschia-sdk/pull/1345) Implement fswap's basic functionality(MsgSwap, MsgSwapAll, Query, Proposal)
* (x/fbridge) [\#1360](https://github.com/Finschia/finschia-sdk/pull/1360) Add Role-based Access Control

### Improvements
* (types) [\#1314](https://github.com/Finschia/finschia-sdk/pull/1314) replace IsEqual with Equal

### Bug Fixes
* (x/auth) [#1281](https://github.com/Finschia/finschia-sdk/pull/1281) `ModuleAccount.Validate` now reports a nil `.BaseAccount` instead of panicking. (backport #1274)
* (x/foundation) [\#1283](https://github.com/Finschia/finschia-sdk/pull/1283) add init logic of foundation module accounts to InitGenesis in order to eliminate potential panic (backport #1277)
* (x/collection) [\#1282](https://github.com/Finschia/finschia-sdk/pull/1282) eliminates potential risk for Insufficient Sanity Check of tokenID in Genesis (backport #1276)
* (x/collection) [\#1290](https://github.com/Finschia/finschia-sdk/pull/1290) export x/collection params into genesis (backport #1268)
* (x/foundation) [\#1295](https://github.com/Finschia/finschia-sdk/pull/1295) add missing error handling for migration
* (sec) [\#1302](https://github.com/Finschia/finschia-sdk/pull/1302) remove map iteration non-determinism with keys + sorting
* (client) [\#1303](https://github.com/Finschia/finschia-sdk/pull/1303) fix possible overflow in BuildUnsignedTx
* (types) [\#1299](https://github.com/Finschia/finschia-sdk/pull/1299) add missing nil checks
* (x/staking) [\#1301](https://github.com/Finschia/finschia-sdk/pull/1301) Use bytes instead of string comparison in delete validator queue (backport cosmos/cosmos-sdk#12303)
* (client/keys) [#1312](https://github.com/Finschia/finschia-sdk/pull/1312) ignore error when key not found in `keys delete`
* (store) [\#1310](https://github.com/Finschia/finschia-sdk/pull/1310) fix app-hash mismatch if upgrade migration commit is interrupted(backport cosmos/cosmos-sdk#13530)
* (types) [\#1313](https://github.com/Finschia/finschia-sdk/pull/1313) fix correctly coalesce coins even with repeated denominations(backport cosmos/cosmos-sdk#13265)
* (x/crypto) [\#1316](https://github.com/Finschia/finschia-sdk/pull/1316) error if incorrect ledger public key (backport cosmos/cosmos-sdk#14460, cosmos/cosmos-sdk#19691)
* (x/auth) [#1319](https://github.com/Finschia/finschia-sdk/pull/1319) prevent signing from wrong key in multisig
* (x/mint, x/slashing) [\#1323](https://github.com/Finschia/finschia-sdk/pull/1323) add missing nil check for params validation
* (x/server) [\#1337](https://github.com/Finschia/finschia-sdk/pull/1337) fix panic when defining minimum gas config as `100stake;100uatom`. Use a `,` delimiter instead of `;`. Fixes the server config getter to use the correct delimiter (backport cosmos/cosmos-sdk#18537)

### Removed

### Breaking Changes

### State Machine Breaking

### Event Breaking Changes

### Build, CI
* (build) [#1298](https://github.com/Finschia/finschia-sdk/pull/1298) Set Finschia/ostracon version

### Document Updates
* (x/token,collection) [#1201](https://github.com/Finschia/finschia-sdk/pull/1201) Deprecate legacy features on x/token,collection
