# CHANGELOG

## Unreleased Version
### Added
- Fork CosmWasm/wasmd/x/wasm (v0.10.0) to add wasm module (#1)
- Implement the token module encoder (#5)
- Implement a querier to get approvers of a token (#7)
- Add linkwasmd and cli tests (#20)
- Implement the collection module encoder (#22)
- Add cli tests for token (#42)
- Add performance reporting feature in CI (#43)
- Add cli tests for collection (#50)
- Add tests for managing max contract size (#51)

### Changed
- Replace links to original CosmWasm code with ours (#244)
- Change used marshal/unmarshal from json.(Un)Marshal -> codec.(Un)Marshal (#37)
- Change the max size of contract managed with the parameter module (#44)
- Update linkwasmd to follow CosmWasm's wasmd v0.11.1 (#65)
- Rewrite Governonce.md (#73, #75)
- Rename QueryXxxParam to XxxParam (#81)
- Change Query total interface (#81)

### Fixed
- Fix linkwasmd's wasmKeeper for #5 (#32)
- Solve a TODO in wasm's cli test (#53)
- Fix CI error on develop branch (#71)
- Fix init params first in InitGenesis (cherry-pick CosmWasm/wasmd@ae169ce) (#76)
