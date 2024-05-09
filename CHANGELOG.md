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
* (x/fswap) [\#1336](https://github.com/Finschia/finschia-sdk/pull/1336) Initialize fswap module
* (x/fswap) [\#1339](https://github.com/Finschia/finschia-sdk/pull/1339) Implement fswap module's genesis
* (x/fbridge) [\#1340](https://github.com/Finschia/finschia-sdk/pull/1340) Initialize fbridge module
* (x/fbridge) [\#1347](https://github.com/Finschia/finschia-sdk/pull/1347) Implement bridge transfer feature (sending side)
* (x/fbridge) [\#1351](https://github.com/Finschia/finschia-sdk/pull/1351) Map a sequence to block number for every bridge request (sending side)
* (x/fswap) [\#1345](https://github.com/Finschia/finschia-sdk/pull/1345) Implement fswap's basic functionality(MsgSwap, MsgSwapAll, Query, Proposal)
* (x/fbridge) [\#1350](https://github.com/Finschia/finschia-sdk/pull/1350) Add Role-based Access Control

### Improvements
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

### Bug Fixes
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
* (x/gov) [\#1304](https://github.com/Finschia/finschia-sdk/pull/1304) fetch a failed proposal tally from proposal.FinalTallyResult in the gprc query
* (x/staking) [\#1306](https://github.com/Finschia/finschia-sdk/pull/1306) add validation for potential slashing evasion during re-delegation
* (client/keys) [#1312](https://github.com/Finschia/finschia-sdk/pull/1312) ignore error when key not found in `keys delete`
* (store) [\#1310](https://github.com/Finschia/finschia-sdk/pull/1310) fix app-hash mismatch if upgrade migration commit is interrupted(backport cosmos/cosmos-sdk#13530)
* (types) [\#1313](https://github.com/Finschia/finschia-sdk/pull/1313) fix correctly coalesce coins even with repeated denominations(backport cosmos/cosmos-sdk#13265)
* (x/crypto) [\#1316](https://github.com/Finschia/finschia-sdk/pull/1316) error if incorrect ledger public key (backport cosmos/cosmos-sdk#14460, cosmos/cosmos-sdk#19691) 
* (x/auth) [#1319](https://github.com/Finschia/finschia-sdk/pull/1319) prevent signing from wrong key in multisig
* (x/mint, x/slashing) [\#1323](https://github.com/Finschia/finschia-sdk/pull/1323) add missing nil check for params validation
* (x/server) [\#1337](https://github.com/Finschia/finschia-sdk/pull/1337) fix panic when defining minimum gas config as `100stake;100uatom`. Use a `,` delimiter instead of `;`. Fixes the server config getter to use the correct delimiter (backport cosmos/cosmos-sdk#18537)
* (x/fbridge) [\#1361](https://github.com/Finschia/finschia-sdk/pull/1361) Fixes fbridge auth checking bug
* (x/fswap) [\#1365](https://github.com/Finschia/finschia-sdk/pull/1365) fix update swap keys for possibly overlapped keys(`(hello,world) should be different to (hel,loworld)`)

### Removed

### Breaking Changes
* (consensus) [\#1178](https://github.com/Finschia/finschia-sdk/pull/1178) change the consensus from Ostracon to Tendermint v0.34.24 

### State Machine Breaking

### Event Breaking Changes

### Build, CI
* (ci) [\#1078](https://github.com/Finschia/finschia-sdk/pull/1078) fix tag comments in github actions workflow docker.yml
* (repo) [\#1157](https://github.com/Finschia/finschia-sdk/pull/1157) setup CODEOWNERS and backport action
* (ci) [\#1160](https://github.com/Finschia/finschia-sdk/pull/1160) remove autopr ci
* (repo) [\#1356](https://github.com/Finschia/finschia-sdk/pull/1356) add v0.49.x for backport action

### Document Updates
* (docs) [\#1059](https://github.com/Finschia/finschia-sdk/pull/1059) create ERRORS.md for x/module
* (docs) [\#1083](https://github.com/Finschia/finschia-sdk/pull/1083) Add detailed explanation about default events
* (x/token,collection) [#1201](https://github.com/Finschia/finschia-sdk/pull/1201) Deprecate legacy features on x/token,collection
