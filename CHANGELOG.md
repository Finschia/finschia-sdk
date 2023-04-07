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

## [Unreleased](https://github.com/line/lbm-sdk/compare/v0.47.0-alpha1...HEAD)

### Features

### Improvements
* (x/collection,token) [\#900](https://github.com/line/lbm-sdk/pull/900) Add uri for MsgModify and deprecate the old ones
* (x/foundation) [\#912](https://github.com/line/lbm-sdk/pull/912) Introduce censorship into x/foundation
* (x/foundation) [\#933](https://github.com/line/lbm-sdk/pull/933) Clean up x/foundation apis
* (x/collection) [\#938](https://github.com/line/lbm-sdk/pull/938) Add progress log into x/collection import-genesis

### Bug Fixes
* (swagger) [\#898](https://github.com/line/lbm-sdk/pull/898) fix a bug not added `lbm.tx.v1beta1.Service/GetBlockWithTxs` in swagger
* (x/collection) [\#911](https://github.com/line/lbm-sdk/pull/911) Add missing command(TxCmdModify) for CLI
* (x/foundation) [\#922](https://github.com/line/lbm-sdk/pull/922) Propagate events in x/foundation through sdk.Results
* (x/foundation) [\#946](https://github.com/line/lbm-sdk/pull/946) Fix broken x/foundation invariant on treasury
* (x/foundation) [\#947](https://github.com/line/lbm-sdk/pull/947) Unpack proposals in x/foundation import-genesis
* (x/collection) [\#953](https://github.com/line/lbm-sdk/pull/953) Allow zero amount of coin in x/collection Query/Balance
* (x/collection) [\#954](https://github.com/line/lbm-sdk/pull/954) Remove duplicated events in x/collection Msg/Modify

### Removed

### Breaking Changes
* (proto) [\#923](https://github.com/line/lbm-sdk/pull/923) deprecate broadcast mode `block`

### Build, CI
* (ci, build) [\#901](https://github.com/line/lbm-sdk/pull/901) Update release pipeline to match non-wasm env

### Document Updates
* (x/foundation) [\#934](https://github.com/line/lbm-sdk/pull/934) Update permlinks in x/foundation documents
* (x/collection,token) [\#944](https://github.com/line/lbm-sdk/pull/944) Update comments in the x/token,collection events proto
