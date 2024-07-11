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
* (x/fbridge) [\#1430](https://github.com/Finschia/finschia-sdk/pull/1430) Add CLI, gRPC, MsgServer tests

### Bug Fixes
<<<<<<< HEAD
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
* (x/fswap, x/fbridge) [\#1378](https://github.com/Finschia/finschia-sdk/pull/1378) Fix bug where amino is not supported in fswap and fbridge
* (x/fswap) [\#1379](https://github.com/Finschia/finschia-sdk/pull/1379) add missing router registration
* (x/fswap) [\#1385](https://github.com/Finschia/finschia-sdk/pull/1385) add accidentally deleted event emissions(EventSetSwap, EventAddDenomMetadata)
* (x/fswap) [\#1392](https://github.com/Finschia/finschia-sdk/pull/1392) fix dummy denom coin data for test in fswap
* (style) [\#1414](https://github.com/Finschia/finschia-sdk/pull/1414) improve code quality with new linters
* (ci) [\#1431](https://github.com/Finschia/finschia-sdk/pull/1431) Fix goreleaser config and replace deprecated docker image with official one
>>>>>>> 8e266d837 (bug: Fix goreleaser config (#1431))

### Removed

### Breaking Changes

### State Machine Breaking

### Event Breaking Changes

### Build, CI

### Document Updates
* (docs) [\#1429](https://github.com/Finschia/finschia-sdk/pull/1429) correct spec docs of fswap module
