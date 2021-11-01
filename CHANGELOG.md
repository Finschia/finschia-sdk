# Changelog

## [Unreleased]

### Features
* (feat) [\#352] (https://github.com/line/lbm-sdk/pull/352) iavl, db & disk stats logging
* (x/gov) [\#368](https://github.com/line/lbm-sdk/pull/368) Governance Split Votes, use `MsgWeightedVote` to send a split vote. Sending a regular `MsgVote` will convert the underlying vote option into a weighted vote with weight 1.

### Improvements
* (slashing) [\#347](https://github.com/line/lbm-sdk/pull/347) Introduce VoterSetCounter 
* (auth) [\#348](https://github.com/line/lbm-sdk/pull/348) Increase default valid_sig_block_period

### Bug Fixes

### Breaking Changes

### Build, CI
* (ci) [\#350](https://github.com/line/lbm-sdk/pull/350) Reduce sim test time
* (ci) [\#351](https://github.com/line/lbm-sdk/pull/351) Remove diff condition from sim-normal
 
### Document Updates
* (docs) [\#361](https://github.com/line/lbm-sdk/pull/361) Add sample command docs
