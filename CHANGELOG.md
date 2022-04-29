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
* (x/wasm) [\#470](https://github.com/line/lbm-sdk/pull/470) remove contract activation control by actor
* (x/wasm) [\#513](https://github.com/line/lbm-sdk/pull/513) fix message representation for signing

### Improvements

* (refactor) [\#493](https://github.com/line/lbm-sdk/pull/493) restructure x/consortium

### Bug Fixes
* (x/wasm) [\#453](https://github.com/line/lbm-sdk/pull/453) modify wasm grpc query api path
* (client) [\#476](https://github.com/line/lbm-sdk/pull/476) change the default value of the client output format in the config
* (server/grpc) [\#516](https://github.com/line/lbm-sdk/pull/516) restore build norace flag

### Breaking Changes

### Build, CI

* (ci) [\#457](https://github.com/line/lbm-sdk/pull/457), [\#471](https://github.com/line/lbm-sdk/pull/471) add swagger check
* (ci) [\#469](https://github.com/line/lbm-sdk/pull/469) publish docker image on tag push

### Document Updates

* (docs) [\#483](https://github.com/line/lbm-sdk/pull/483) update documents on x/stakingplus
* (docs) [\#490](https://github.com/line/lbm-sdk/pull/490) update documents on x/consortium
