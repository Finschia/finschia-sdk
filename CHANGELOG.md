# Changelog

## [Unreleased]

### Features
* (feat) [\#352] (https://github.com/line/lbm-sdk/pull/352) iavl, db & disk stats logging
* (x/gov) [\#368](https://github.com/line/lbm-sdk/pull/368) Governance Split Votes, use `MsgWeightedVote` to send a split vote. Sending a regular `MsgVote` will convert the underlying vote option into a weighted vote with weight 1.
* (x/upgrade) [\#377] (https://github.com/line/lbm-sdk/pull/377) To smoothen the update to the latest stable release, the SDK includes vesion map for managing migrations between SDK versions.
* (x/wasm) [\#358] (https://github.com/line/lbm-sdk/pull/358) change wasm metrics method to using prometheus directly
* (x/feegrant) [\#380] (https://github.com/line/lbm-sdk/pull/380) Feegrant module
* (x/wasm) [\#395] (https://github.com/line/lbm-sdk/pull/395) Add the instantiate_permission in the CodeInfoResponse
* (x/consortium) [\#406] (https://github.com/line/lbm-sdk/pull/406) Add CreateValidator access control feature

### Improvements
* (slashing) [\#347](https://github.com/line/lbm-sdk/pull/347) Introduce VoterSetCounter
* (auth) [\#348](https://github.com/line/lbm-sdk/pull/348) Increase default valid_sig_block_period

### Bug Fixes
* (x/feegrant) [\#383] (https://github.com/line/lbm-sdk/pull/383) Update allowance inside AllowedMsgAllowance
* (tm-db) [\#388] (https://github.com/line/lbm-sdk/pull/388) Bump up tm-db fixing invalid memory reference
* (swagger) [\#391] (https://github.com/line/lbm-sdk/pull/391) fix swagger's config path for wasm
* (x/wasm) [\#393] (https://github.com/line/lbm-sdk/pull/393) fix bug where `StoreCodeAndInstantiateContract`, `UpdateContractStatus`, `UpdateContractStatusProposal` API does not work

### Breaking Changes

### Build, CI
* (ci) [\#350](https://github.com/line/lbm-sdk/pull/350) Reduce sim test time
* (ci) [\#351](https://github.com/line/lbm-sdk/pull/351) Remove diff condition from sim-normal

### Document Updates
* (docs) [\#361](https://github.com/line/lbm-sdk/pull/361) Add sample command docs
* (docs) [\#392](https://github.com/line/lbm-sdk/pull/392) Modify with latest version of swagger REST interface docs.
