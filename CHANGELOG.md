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

## [Unreleased](https://github.com/line/lbm-sdk/compare/v0.46.0-rc9...HEAD)

### Features
* (global) [\#783](https://github.com/line/lbm-sdk/pull/783) bump up github.com/cosmos/cosmos-sdk to v0.45.10
* (build) [\#793](https://github.com/line/lbm-sdk/pull/793) enable to use libsodium version ostracon

### Improvements
* (x/auth) [\#776](https://github.com/line/lbm-sdk/pull/776) remove unused MsgEmpty

### Bug Fixes
* (x/foundation) [\#772](https://github.com/line/lbm-sdk/pull/772) export x/foundation pool
* (baseapp) [\#781](https://github.com/line/lbm-sdk/pull/781) implement method `SetOption()` in baseapp
* (global) [\#782](https://github.com/line/lbm-sdk/pull/782) add unhandled return error handling
* (x/collection,x/token) [\#798](https://github.com/line/lbm-sdk/pull/798) Fix x/collection ModifyContract

### Breaking Changes
* (cli) [\#773](https://github.com/line/lbm-sdk/pull/773) guide users to use generate-only in messages for x/foundation authority
* (x/foundation) [\#790](https://github.com/line/lbm-sdk/pull/790) fix case of gov_mint_left_count in x/foundation

### Build, CI
* (ci) [\#779](https://github.com/line/lbm-sdk/pull/779) change github action trigger rules for `release/*` and `rc*/*` branches

### Document Updates
* (docs) [\#766](https://github.com/line/lbm-sdk/pull/766) fix submit-proposal command on x/foundation
