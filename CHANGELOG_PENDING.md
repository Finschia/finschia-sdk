# Unreleased Changes

## lbm-sdk v2

Write the changes made after branching from cosmos-sdk v0.42.1.

### BREAKING CHANGES

- CLI/RPC/Config
* (rpc) [\#97](https://github.com/line/lbm-sdk/pull/97) Send response with 404 status when quering non-exist account

- Apps

- P2P Protocol

- Go API
* (global) [\#90](https://github.com/line/lbm-sdk/pull/90) Rename module path to `github.com/line/lbm-sdk/v2`

- Blockchain Protocol

### FEATURES
* (global) [\#97](https://github.com/line/lbm-sdk/pull/97) Add codespace to response
* (global) [\#97](https://github.com/line/lbm-sdk/pull/97) Add codespace to query error
* (config) [\#114](https://github.com/line/lbm-sdk/pull/114) Add idle-timeout to rest server and rpc server config
* (x/wasm) [\#127](https://github.com/line/lbm-sdk/pull/127) Add wasm with Staragate migration completed.
* (x/wasm) [\#151](https://github.com/line/lbm-sdk/pull/151) Add contract access control.
* (x/wasm) [\#194](https://github.com/line/lbm-sdk/pull/194) Replace importing CosmWasm/wasmvm with line/wasmvm.

### IMPROVEMENTS

### BUG FIXES

