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

## [Unreleased](https://github.com/line/lbm-sdk/compare/v0.46.0...HEAD)

### Features
* (baseapp) [\#840](https://github.com/line/lbm-sdk/pull/840) allow querying the state based on `CheckState`.
* (x/foundation) [\#848](https://github.com/line/lbm-sdk/pull/848) remove `gov mint` for x/foundation proposal
* (x/wasm) [\#850](https://github.com/line/lbm-sdk/pull/850) remove `x/wasm` module in lbm-sdk

### Improvements
* (cosmovisor) [\#792](https://github.com/line/lbm-sdk/pull/792) Use upstream's cosmovisor
* (server) [\#821](https://github.com/line/lbm-sdk/pull/821) Get validator pubkey considering KMS

### Bug Fixes
* (client) [\#817](https://github.com/line/lbm-sdk/pull/817) remove support for composite (BLS) type
* (x/foundation) [\#834](https://github.com/line/lbm-sdk/pull/834) Apply foundation audit
* (x/collection,token) [\#849](https://github.com/line/lbm-sdk/pull/849) Introduce codespace into x/collection,token
* (x/token,collection) [\#863](https://github.com/line/lbm-sdk/pull/863) Update x/collection,token proto
* (x/collection,token) [\#866](https://github.com/line/lbm-sdk/pull/866) Do not create account on x/token,collection

### Removed
* [\#853](https://github.com/line/lbm-sdk/pull/853) remove useless stub BeginBlock, EndBlock methods from modules below
  * ibc, authz, collection, feegrant, ibc, token, wasm
* (x/ibc) [\#858](https://github.com/line/lbm-sdk/pull/858) detach ibc module(repo: [line/ibc-go](https://github.com/line/ibc-go))

### Breaking Changes
* (rest) [\#807](https://github.com/line/lbm-sdk/pull/807) remove legacy REST API
* (codec) [\#833](https://github.com/line/lbm-sdk/pull/833) Fix foundation amino codec

### Build, CI
* (ci) [\#829](https://github.com/line/lbm-sdk/pull/829) automate release process

### Document Updates
