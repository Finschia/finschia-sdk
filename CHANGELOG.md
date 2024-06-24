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

## [Unreleased](https://github.com/Finschia/finschia-sdk/compare/v0.49.0...HEAD)

### Features

### Improvements
* (x/fswap) [\#1415](https://github.com/Finschia/finschia-sdk/pull/1415) add more testcases for fswap module
* (style) [\#1427](https://github.com/Finschia/finschia-sdk/pull/1427) Lint all files based on latest setting

### Bug Fixes

### Removed

### Breaking Changes
* (server) [\#1423](https://github.com/Finschia/finschia-sdk/pull/1423) Remove grpc replace directive and refactor grpc-web/rosetta/grpc-gw

### State Machine Breaking

### Event Breaking Changes

### Build, CI
* (build, ci) [\#1410](https://github.com/Finschia/finschia-sdk/pull/1410) Bump Go from 1.20 to 1.22
* (build) [\#1413](https://github.com/Finschia/finschia-sdk/pull/1413) Update outdated dependencies for v0.49.x

### Document Updates
<<<<<<< HEAD
=======
* (docs) [\#1059](https://github.com/Finschia/finschia-sdk/pull/1059) create ERRORS.md for x/module
* (docs) [\#1083](https://github.com/Finschia/finschia-sdk/pull/1083) Add detailed explanation about default events
* (x/token,collection) [#1201](https://github.com/Finschia/finschia-sdk/pull/1201) Deprecate legacy features on x/token,collection
* (build) [\#1393](https://github.com/Finschia/finschia-sdk/pull/1393) add current directory as suffix for docker container
* (docs) [\#1419](https://github.com/Finschia/finschia-sdk/pull/1419) correct spec docs of fswap module
>>>>>>> 7bd6c8244 (docs: correct spec docs of fswap module (#1419))
