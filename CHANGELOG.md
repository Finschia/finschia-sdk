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

## [Unreleased](https://github.com/Finschia/finschia-sdk/compare/v0.47.2...HEAD)

### Features

### Improvements

### Bug Fixes
* (x/foundation) [\#1061](https://github.com/Finschia/finschia-sdk/pull/1061) Make x/foundation MsgExec propagate events (backport #1053)
* (baseapp) [\#1075](https://github.com/cosmos/cosmos-sdk/pull/1075) Add `events.GetAttributes` and `event.GetAttribute` methods to simplify the retrieval of an attribute from event(s) (backport cosmos/cosmos-sdk#14356)
* (baseapp) [\#1079](https://github.com/cosmos/cosmos-sdk/pull/1079) Do not add `module` attribute in case of ibc messages

### Removed

### Breaking Changes
* (refactor) [\#1063](https://github.com/Finschia/finschia-sdk/pull/1063) Automate EventTypeMessage inclusion in every message execution (backport cosmos/cosmos-sdk#13532)
* (x/bank) [#1066](https://github.com/Finschia/finschia-sdk/pull/1066) Remove message events including `sender` attribute whose information is already present in the relevant events (backport cosmos/cosmos#17273)

### Build, CI

### Document Updates
* (docs) [\#1081](https://github.com/Finschia/finschia-sdk/pull/1081) Document default events
