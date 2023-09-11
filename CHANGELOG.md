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

## [Unreleased](https://github.com/Finschia/finschia-sdk/compare/v0.47.0...HEAD)

### Features
* (x/auth) [\#1011](https://github.com/Finschia/finschia-sdk/pull/1011) add the api for querying next account number
* (server/grpc) [\#1017](https://github.com/Finschia/finschia-sdk/pull/1017) support custom r/w gRPC options (backport cosmos/cosmos-sdk#11889)

### Improvements
* (third_party/proto) [\#1037](https://github.com/Finschia/finschia-sdk/pull/1037) change the proof.proto path to third_party/proto/confio
* (ostracon) [\#1057](https://github.com/Finschia/finschia-sdk/pull/1057) Bump up Ostracon from to v1.1.1
* (docs) [\#1120](https://github.com/Finschia/finschia-sdk/pull/1120) Update links in x/foundation README.md

### Bug Fixes
* (ledger) [\#1040](https://github.com/Finschia/finschia-sdk/pull/1040) fix a bug(unable to connect nano S plus ledger on ubuntu)
* (x/foundation) [\#1053](https://github.com/Finschia/finschia-sdk/pull/1053) Make x/foundation MsgExec propagate events

### Removed

### Breaking Changes
* (x/foundation) [\#999](https://github.com/Finschia/finschia-sdk/pull/999) migrate x/foundation FoundationTax into x/params

### Build, CI
* (build,ci) [\#1043](https://github.com/Finschia/finschia-sdk/pull/1043) Update golang version to 1.20
* (ci) [\#1078](https://github.com/Finschia/finschia-sdk/pull/1078) fix tag comments in github actions workflow docker.yml

### Document Updates
* (readme) [\#997](https://github.com/finschia/finschia-sdk/pull/997) fix swagger url
* (docs) [\#1059](https://github.com/Finschia/finschia-sdk/pull/1059) create ERRORS.md for x/module
