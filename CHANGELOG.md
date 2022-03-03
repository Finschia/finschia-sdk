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

## [Unreleased](https://github.com/line/lbm-sdk/compare/v0.44.0-rc0...HEAD)

### Features
* (x/wasm) [\#444](https://github.com/line/lbm-sdk/pull/444) Merge wasmd 0.19.0
  * remove custom encoder from x/wasm/keeper.NewKeeper's arg. After the Token/collection module is added, it will be ported again.
* (cosmos-sdk) [\#437](https://github.com/line/lbm-sdk/pull/437) dump up to cosmos-sdk v0.42.11
  * [changelog of cosmos-sdk v0.42.11](https://github.com/cosmos/cosmos-sdk/blob/v0.42.11/CHANGELOG.md)
* (feat) [\#434](https://github.com/line/lbm-sdk/pull/434) Revert signature mechanism using `sig_block_height`
* (x/token) [\#416](https://github.com/line/lbm-sdk/pull/416) Migrate token module from line/link

### Improvements

### Bug Fixes
* (x/wasm) [\#436](https://github.com/line/lbm-sdk/pull/436) remove `x/wasm/linkwasmd`

### Breaking Changes

### Build, CI
* (makefile, ci) [\#438](https://github.com/line/lbm-sdk/pull/438) fix `make proto-format` and `make proto-check-breaking` error

### Document Updates
