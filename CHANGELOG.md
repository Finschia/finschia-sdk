# Changelog

## Unreleased

### Bug Fixes
* (x/wasm) [\#434](https://github.com/line/lbm-sdk/pull/434) remove `x/wasm/linkwasmd`
* (makefile) [\#438](https://github.com/line/lbm-sdk/pull/434) fix `make proto-format` and `make proto-check-breaking` error


## [v0.42.9](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.42.9) - 2021-08-04

### Bug Fixes

* [\#9835](https://github.com/cosmos/cosmos-sdk/pull/9835) Moved capability initialization logic to BeginBlocker to fix nondeterminsim issue mentioned in [\#9800](https://github.com/cosmos/cosmos-sdk/issues/9800). Applications must now include the capability module in their BeginBlocker order before any module that uses capabilities gets run.
* [\#9201](https://github.com/cosmos/cosmos-sdk/pull/9201) Fixed `<app> init --recover` flag.


### API Breaking Changes

* [\#9835](https://github.com/cosmos/cosmos-sdk/pull/9835) The `InitializeAndSeal` API has not changed, however it no longer initializes the in-memory state. `InitMemStore` has been introduced to serve this function, which will be called either in `InitChain` or `BeginBlock` (whichever is first after app start). Nodes may run this version on a network running 0.42.x, however, they must update their app.go files to include the capability module in their begin blockers.

### Client Breaking Changes

* [\#9781](https://github.com/cosmos/cosmos-sdk/pull/9781) Improve`withdraw-all-rewards` UX when broadcast mode `async` or `async` is used.

## [v0.42.8](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.42.8) - 2021-07-30

### Features

* [\#9750](https://github.com/cosmos/cosmos-sdk/pull/9750) Emit events for tx signature and sequence, so clients can now query txs by signature (`tx.signature='<base64_sig>'`) or by address and sequence combo (`tx.acc_seq='<addr>/<seq>'`).

### Improvements

* (cli) [\#9717](https://github.com/cosmos/cosmos-sdk/pull/9717) Added CLI flag `--output json/text` to `tx` cli commands.

### Bug Fixes

* [\#9766](https://github.com/cosmos/cosmos-sdk/pull/9766) Fix hardcoded ledger signing algorithm on `keys add` command.

## [v0.42.7](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.42.7) - 2021-07-09

### Improvements

* (baseapp) [\#9578](https://github.com/cosmos/cosmos-sdk/pull/9578) Return `Baseapp`'s `trace` value for logging error stack traces.
* (cli) [\#9593](https://github.com/cosmos/cosmos-sdk/pull/9593) Check if chain-id is blank before verifying signatures in multisign and error.

### Bug Fixes

* (x/ibc) [\#9640](https://github.com/cosmos/cosmos-sdk/pull/9640) Fix IBC Transfer Ack Success event as it was initially emitting opposite value.
* [\#9645](https://github.com/cosmos/cosmos-sdk/pull/9645) Use correct Prometheus format for metric labels.
* [\#9299](https://github.com/cosmos/cosmos-sdk/pull/9299) Fix `[appd] keys parse cosmos1...` freezing.
* (keyring) [\#9563](https://github.com/cosmos/cosmos-sdk/pull/9563) fix keyring kwallet backend when using with empty wallet.
* (x/capability) [\#9392](https://github.com/cosmos/cosmos-sdk/pull/9392) initialization fix, which fixes the consensus error when using statesync.


## [v0.42.6](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.42.6) - 2021-06-18

### Improvements

* [\#9428](https://github.com/cosmos/cosmos-sdk/pull/9428) Optimize bank InitGenesis. Added `k.initBalances`.
* [\#9429](https://github.com/cosmos/cosmos-sdk/pull/9429) Add `cosmos_sdk_version` to node_info
* [\#9541](https://github.com/cosmos/cosmos-sdk/pull/9541) Bump tendermint dependency to v0.34.11.

### Bug Fixes

* [\#9385](https://github.com/cosmos/cosmos-sdk/pull/9385) Fix IBC `query ibc client header` cli command. Support historical queries for query header/node-state commands.
* [\#9401](https://github.com/cosmos/cosmos-sdk/pull/9401) Fixes incorrect export of IBC identifier sequences. Previously, the next identifier sequence for clients/connections/channels was not set during genesis export. This resulted in the next identifiers being generated on the new chain to reuse old identifiers (the sequences began again from 0).
* [\#9408](https://github.com/cosmos/cosmos-sdk/pull/9408) Update simapp to use correct default broadcast mode.
* [\#9513](https://github.com/cosmos/cosmos-sdk/pull/9513) Fixes testnet CLI command. Testnet now updates the supply in genesis. Previously, when using add-genesis-account and testnet together, inconsistent genesis files would be produced, as only add-genesis-account was updating the supply.
* (x/gov) [\#8813](https://github.com/cosmos/cosmos-sdk/pull/8813) fix `GET /cosmos/gov/v1beta1/proposals/{proposal_id}/deposits` to include initial deposit

### Features

* [\#9383](https://github.com/cosmos/cosmos-sdk/pull/9383) New CLI command `query ibc-transfer escrow-address <port> <channel id>` to get the escrow address for a channel; can be used to then query balance of escrowed tokens
* (baseapp, types) [#\9390](https://github.com/cosmos/cosmos-sdk/pull/9390) Add current block header hash to `Context`
* (store) [\#9403](https://github.com/cosmos/cosmos-sdk/pull/9403) Add `RefundGas` function to `GasMeter` interface

## [v0.42.5](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.42.5) - 2021-05-18

### Bug Fixes

* [\#9514](https://github.com/cosmos/cosmos-sdk/issues/9514) Fix panic when retrieving the `BlockGasMeter` on `(Re)CheckTx` mode.
* [\#9235](https://github.com/cosmos/cosmos-sdk/pull/9235) CreateMembershipProof/CreateNonMembershipProof now returns an error
if input key is empty, or input data contains empty key.
* [\#9108](https://github.com/cosmos/cosmos-sdk/pull/9108) Fixed the bug with querying multisig account, which is not showing threshold and public_keys.
* [\#9345](https://github.com/cosmos/cosmos-sdk/pull/9345) Fix ARM support.
* [\#9040](https://github.com/cosmos/cosmos-sdk/pull/9040) Fix ENV variables binding to CLI flags for client config.

### Features

* [\#8953](https://github.com/cosmos/cosmos-sdk/pull/8953) Add the `config` CLI subcommand back to the SDK, which saves client-side configuration in a `client.toml` file.


## [v0.42.4](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.42.4) - 2021-04-08

### Client Breaking Changes

* [\#9026](https://github.com/cosmos/cosmos-sdk/pull/9026) By default, the `tx sign` and `tx sign-batch` CLI commands use SIGN_MODE_DIRECT to sign transactions for local pubkeys. For multisigs and ledger keys, the default LEGACY_AMINO_JSON is used.

### Bug Fixes

* (gRPC) [\#9015](https://github.com/cosmos/cosmos-sdk/pull/9015) Fix invalid status code when accessing gRPC endpoints.
* [\#9026](https://github.com/cosmos/cosmos-sdk/pull/9026) Fixed the bug that caused the `gentx` command to fail for Ledger keys.

### Improvements

* [\#9081](https://github.com/cosmos/cosmos-sdk/pull/9081) Upgrade Tendermint to v0.34.9 that includes a security issue fix for Tendermint light clients.

## [v0.42.3](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.42.3) - 2021-03-24

This release fixes a security vulnerability identified in x/bank.

## [v0.42.2](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.42.2) - 2021-03-19

### Improvements

* (grpc) [\#8815](https://github.com/cosmos/cosmos-sdk/pull/8815) Add orderBy parameter to `TxsByEvents` endpoint.
* (cli) [\#8826](https://github.com/cosmos/cosmos-sdk/pull/8826) Add trust to macOS Keychain for caller app by default.
* (store) [\#8811](https://github.com/cosmos/cosmos-sdk/pull/8811) store/cachekv: use typed types/kv.List instead of container/list.List

### Bug Fixes

* (crypto) [\#8841](https://github.com/cosmos/cosmos-sdk/pull/8841) Fix legacy multisig amino marshaling, allowing migrations to work between v0.39 and v0.40+.
* (cli) [\#8873](https://github.com/cosmos/cosmos-sdk/pull/8873) add --output-document to multisign-batch.

## [v0.42.1](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.42.1) - 2021-03-10

This release fixes security vulnerability identified in the simapp.


## [v0.42.0](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.42.0) - 2021-03-08

**IMPORTANT**: This release contains an important security fix for all non Cosmos Hub chains running Stargate version of the Cosmos SDK (>0.40). Non-hub chains should not be using any version of the SDK in the v0.40.x or v0.41.x release series. See [#8461](https://github.com/cosmos/cosmos-sdk/pull/8461) for more details.

### Improvements

* (x/ibc) [\#8624](https://github.com/cosmos/cosmos-sdk/pull/8624) Emit full header in IBC UpdateClient message.
* (x/crisis) [\#8621](https://github.com/cosmos/cosmos-sdk/issues/8621) crisis invariants names now print to loggers.

### Bug fixes

* (x/evidence) [\#8461](https://github.com/cosmos/cosmos-sdk/pull/8461) Fix bech32 prefix in evidence validator address conversion
* (x/gov) [\#8806](https://github.com/cosmos/cosmos-sdk/issues/8806) Fix q gov proposals command's mishandling of the --status parameter's values.

## [v0.41.4](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.41.3) - 2021-03-02

### Features

* [\#7787](https://github.com/cosmos/cosmos-sdk/pull/7787) Add multisign-batch command.

### Bug fixes

* [\#8730](https://github.com/cosmos/cosmos-sdk/pull/8730) Allow REST endpoint to query txs with multisig addresses.
* [\#8680](https://github.com/cosmos/cosmos-sdk/issues/8680) Fix missing timestamp in GetTxsEvent response [\#8732](https://github.com/cosmos/cosmos-sdk/pull/8732).
* [\#8681](https://github.com/cosmos/cosmos-sdk/issues/8681) Fix missing error message when calling GetTxsEvent [\#8732](https://github.com/cosmos/cosmos-sdk/pull/8732)
* (server) [\#8641](https://github.com/cosmos/cosmos-sdk/pull/8641) Fix Tendermint and application configuration reading from file
* (client/keys) [\#8639] (https://github.com/cosmos/cosmos-sdk/pull/8639) Fix keys migrate for mulitisig, offline, and ledger keys. The migrate command now takes a positional old_home_dir argument.

### Improvements

* (store/cachekv), (x/bank/types) [\#8719](https://github.com/cosmos/cosmos-sdk/pull/8719) algorithmically fix pathologically slow code
* [\#8701](https://github.com/cosmos/cosmos-sdk/pull/8701) Upgrade tendermint v0.34.8.
* [\#8714](https://github.com/cosmos/cosmos-sdk/pull/8714) Allow accounts to have a balance of 0 at genesis.

## [v0.41.3](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.41.3) - 2021-02-18

### Bug Fixes

* [\#8617](https://github.com/cosmos/cosmos-sdk/pull/8617) Fix build failures caused by a small API breakage introduced in tendermint v0.34.7.

## [v0.41.2](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.41.2) - 2021-02-18

### Improvements

* Bump tendermint dependency to v0.34.7.

## [v0.41.1](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.41.1) - 2021-02-17

### Bug Fixes

* (grpc) [\#8549](https://github.com/cosmos/cosmos-sdk/pull/8549) Make gRPC requests go through ABCI and disallow concurrency.
* (x/staking) [\#8546](https://github.com/cosmos/cosmos-sdk/pull/8546) Fix caching bug where concurrent calls to GetValidator could cause a node to crash
* (server) [\#8481](https://github.com/cosmos/cosmos-sdk/pull/8481) Don't create files when running `{appd} tendermint show-*` subcommands.
* (client/keys) [\#8436](https://github.com/cosmos/cosmos-sdk/pull/8436) Fix keybase->keyring keys migration.
* (crypto/hd) [\#8607](https://github.com/cosmos/cosmos-sdk/pull/8607) Make DerivePrivateKeyForPath error and not panic on trailing slashes.

### Improvements

* (x/ibc) [\#8458](https://github.com/cosmos/cosmos-sdk/pull/8458) Add `packet_connection` attribute to ibc events to enable relayer filtering
* [\#8396](https://github.com/cosmos/cosmos-sdk/pull/8396) Add support for ARM platform
* (x/bank) [\#8479](https://github.com/cosmos/cosmos-sdk/pull/8479) Aditional client denom metadata validation for `base` and `display` denoms.
* (codec/types) [\#8605](https://github.com/cosmos/cosmos-sdk/pull/8605) Avoid unnecessary allocations for NewAnyWithCustomTypeURL on error.

## [v0.41.0](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.41.0) - 2021-01-26

### State Machine Breaking

* (x/ibc) [\#8266](https://github.com/cosmos/cosmos-sdk/issues/8266) Add amino JSON for IBC messages in order to support Ledger text signing.
* (x/ibc) [\#8404](https://github.com/cosmos/cosmos-sdk/pull/8404) Reorder IBC `ChanOpenAck` and `ChanOpenConfirm` handler execution to perform core handler first, followed by application callbacks.

### Bug Fixes

* (simapp) [\#8418](https://github.com/cosmos/cosmos-sdk/pull/8418) Add balance coin to supply when adding a new genesis account
* (x/bank) [\#8417](https://github.com/cosmos/cosmos-sdk/pull/8417) Validate balances and coin denom metadata on genesis
* (x/staking) [\#8546](https://github.com/cosmos/cosmos-sdk/pull/8546) Fix caching bug where concurrent calls to GetValidator could cause a node to crash

## [v0.40.1](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.40.1) - 2021-01-19

### Improvements

* (x/bank) [\#8302](https://github.com/cosmos/cosmos-sdk/issues/8302) Add gRPC and CLI queries for client denomination metadata.
* (tendermint) Bump Tendermint version to [v0.34.3](https://github.com/tendermint/tendermint/releases/tag/v0.34.3).

### Bug Fixes

* [\#8085](https://github.com/cosmos/cosmos-sdk/pull/8058) fix zero time checks
* [\#8280](https://github.com/cosmos/cosmos-sdk/pull/8280) fix GET /upgrade/current query
* (x/auth) [\#8287](https://github.com/cosmos/cosmos-sdk/pull/8287) Fix `tx sign --signature-only` to return correct sequence value in signature.
* (build) [\8300](https://github.com/cosmos/cosmos-sdk/pull/8300), [\8301](https://github.com/cosmos/cosmos-sdk/pull/8301) Fix reproducible builds
* (types/errors) [\#8355][https://github.com/cosmos/cosmos-sdk/pull/8355] Fix errorWrap `Is` method.
* (x/ibc) [\#8341](https://github.com/cosmos/cosmos-sdk/pull/8341) Fix query latest consensus state.
* (proto) [\#8350][https://github.com/cosmos/cosmos-sdk/pull/8350], [\#8361](https://github.com/cosmos/cosmos-sdk/pull/8361) Update gogo proto deps with v1.3.2 security fixes
* (x/ibc) [\#8359](https://github.com/cosmos/cosmos-sdk/pull/8359) Add missing UnpackInterfaces functions to IBC Query Responses. Fixes 'cannot unpack Any' error for IBC types.
* (x/bank) [\#8317](https://github.com/cosmos/cosmos-sdk/pull/8317) Fix panic when querying for a not found client denomination metadata.


## [v0.40.0](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.40.0) - 2021-01-08

v0.40.0, known as the Stargate release of the Cosmos SDK, is one of the largest releases
of the Cosmos SDK since launch. Please read through this changelog and [release notes](./RELEASE_NOTES.md) to make
sure you are aware of any relevant breaking changes.

### Client Breaking Changes

* __Modules__

* __CLI__
  * (client/keys) [\#5889](https://github.com/cosmos/cosmos-sdk/pull/5889) remove `keys update` command.
  * (x/auth) [\#5844](https://github.com/cosmos/cosmos-sdk/pull/5844) `tx sign` command now returns an error when signing is attempted with offline/multisig keys.
  * (x/auth) [\#6108](https://github.com/cosmos/cosmos-sdk/pull/6108) `tx sign` command's `--validate-signatures` flag is migrated into a `tx validate-signatures` standalone command.
  * (x/auth) [#7788](https://github.com/cosmos/cosmos-sdk/pull/7788) Remove `tx auth` subcommands, all auth subcommands exist as `tx <subcommand>`
  * (x/genutil) [\#6651](https://github.com/cosmos/cosmos-sdk/pull/6651) The `gentx` command has been improved. No longer are `--from` and `--name` flags required. Instead, a single argument, `name`, is required which refers to the key pair in the Keyring. In addition, an optional
  `--moniker` flag can be provided to override the moniker found in `config.toml`.
  * (x/upgrade) [#7697](https://github.com/cosmos/cosmos-sdk/pull/7697) Rename flag name "--time" to "--upgrade-time", "--info" to "--upgrade-info", to keep it consistent with help message.
* __REST / Queriers__
  * (api) [\#6426](https://github.com/cosmos/cosmos-sdk/pull/6426) The ability to start an out-of-process API REST server has now been removed. Instead, the API server is now started in-process along with the application and Tendermint. Configuration options have been added to `app.toml` to enable/disable the API server along with additional HTTP server options.
  * (client) [\#7246](https://github.com/cosmos/cosmos-sdk/pull/7246) The rest server endpoint `/swagger-ui/` is replaced by `/swagger/`, and contains swagger documentation for gRPC Gateway routes in addition to legacy REST routes. Swagger API is exposed only if set in `app.toml`.
  * (x/auth) [\#5702](https://github.com/cosmos/cosmos-sdk/pull/5702) The `x/auth` querier route has changed from `"acc"` to `"auth"`.
  * (x/bank) [\#5572](https://github.com/cosmos/cosmos-sdk/pull/5572) The `/bank/balances/{address}` endpoint now returns all account balances or a single balance by denom when the `denom` query parameter is present.
  * (x/evidence) [\#5952](https://github.com/cosmos/cosmos-sdk/pull/5952) Remove CLI and REST handlers for querying `x/evidence` parameters.
  * (x/gov) [#6295](https://github.com/cosmos/cosmos-sdk/pull/6295) Fix typo in querying governance params.
* __General__
  * (baseapp) [\#6384](https://github.com/cosmos/cosmos-sdk/pull/6384) The `Result.Data` is now a Protocol Buffer encoded binary blob of type `TxData`. The `TxData` contains `Data` which contains a list of Protocol Buffer encoded message data and the corresponding message type.
  * (client) [\#5783](https://github.com/cosmos/cosmos-sdk/issues/5783) Unify all coins representations on JSON client requests for governance proposals.
  * (crypto) [\#7419](https://github.com/cosmos/cosmos-sdk/pull/7419) The SDK doesn't use Tendermint's `crypto.PubKey`
      interface anymore, and uses instead it's own `PubKey` interface, defined in `crypto/types`. Replace all instances of
      `crypto.PubKey` by `cryptotypes.Pubkey`.
  * (store/rootmulti) [\#6390](https://github.com/cosmos/cosmos-sdk/pull/6390) Proofs of empty stores are no longer supported.
  * (store/types) [\#5730](https://github.com/cosmos/cosmos-sdk/pull/5730) store.types.Cp() is removed in favour of types.CopyBytes().
  * (x/auth) [\#6054](https://github.com/cosmos/cosmos-sdk/pull/6054) Remove custom JSON marshaling for base accounts as multsigs cannot be bech32 decoded.
  * (x/auth/vesting) [\#6859](https://github.com/cosmos/cosmos-sdk/pull/6859) Custom JSON marshaling of vesting accounts was removed. Vesting accounts are now marshaled using their default proto or amino JSON representation.
  * (x/bank) [\#5785](https://github.com/cosmos/cosmos-sdk/issues/5785) In x/bank errors, JSON strings coerced to valid UTF-8 bytes at JSON marshalling time
  are now replaced by human-readable expressions. This change can potentially break compatibility with all those client side tools
  that parse log messages.
  * (x/evidence) [\#7538](https://github.com/cosmos/cosmos-sdk/pull/7538) The ABCI's `Result.Data` field for
    `MsgSubmitEvidence` responses does not contain the raw evidence's hash, but the protobuf encoded
    `MsgSubmitEvidenceResponse` struct.
  * (x/gov) [\#7533](https://github.com/cosmos/cosmos-sdk/pull/7533) The ABCI's `Result.Data` field for
    `MsgSubmitProposal` responses does not contain a raw binary encoding of the `proposalID`, but the protobuf encoded
    `MsgSubmitSubmitProposalResponse` struct.
  * (x/gov) [\#6859](https://github.com/cosmos/cosmos-sdk/pull/6859) `ProposalStatus` and `VoteOption` are now JSON serialized using its protobuf name, so expect names like `PROPOSAL_STATUS_DEPOSIT_PERIOD` as opposed to `DepositPeriod`.
  * (x/staking) [\#7499](https://github.com/cosmos/cosmos-sdk/pull/7499) `BondStatus` is now a protobuf `enum` instead
    of an `int32`, and JSON serialized using its protobuf name, so expect names like `BOND_STATUS_UNBONDING` as opposed
    to `Unbonding`.
  * (x/staking) [\#7556](https://github.com/cosmos/cosmos-sdk/pull/7556) The ABCI's `Result.Data` field for
    `MsgBeginRedelegate` and `MsgUndelegate` responses does not contain custom binary marshaled `completionTime`, but the
    protobuf encoded `MsgBeginRedelegateResponse` and `MsgUndelegateResponse` structs respectively

### API Breaking Changes

* __Baseapp / Client__
  * (AppModule) [\#7518](https://github.com/cosmos/cosmos-sdk/pull/7518) [\#7584](https://github.com/cosmos/cosmos-sdk/pull/7584) Rename `AppModule.RegisterQueryServices` to `AppModule.RegisterServices`, as this method now registers multiple services (the gRPC query service and the protobuf Msg service). A `Configurator` struct is used to hold the different services.
  * (baseapp) [\#5865](https://github.com/cosmos/cosmos-sdk/pull/5865) The `SimulationResponse` returned from tx simulation is now JSON encoded instead of Amino binary.
  * (client) [\#6290](https://github.com/cosmos/cosmos-sdk/pull/6290) `CLIContext` is renamed to `Context`. `Context` and all related methods have been moved from package context to client.
  * (client) [\#6525](https://github.com/cosmos/cosmos-sdk/pull/6525) Removed support for `indent` in JSON responses. Clients should consider piping to an external tool such as `jq`.
  * (client) [\#8107](https://github.com/cosmos/cosmos-sdk/pull/8107) Renamed `PrintOutput` and `PrintOutputLegacy`
      methods of the `context.Client` object to `PrintProto` and `PrintObjectLegacy`.
  * (client/flags) [\#6632](https://github.com/cosmos/cosmos-sdk/pull/6632) Remove NewCompletionCmd(), the function is now available in tendermint.
  * (client/input) [\#5904](https://github.com/cosmos/cosmos-sdk/pull/5904) Removal of unnecessary `GetCheckPassword`, `PrintPrefixed` functions.
  * (client/keys) [\#5889](https://github.com/cosmos/cosmos-sdk/pull/5889) Rename `NewKeyBaseFromDir()` -> `NewLegacyKeyBaseFromDir()`.
  * (client/keys) [\#5820](https://github.com/cosmos/cosmos-sdk/pull/5820/) Removed method CloseDB from Keybase interface.
  * (client/rpc) [\#6290](https://github.com/cosmos/cosmos-sdk/pull/6290) `client` package and subdirs reorganization.
  * (client/lcd) [\#6290](https://github.com/cosmos/cosmos-sdk/pull/6290) `CliCtx` of struct `RestServer` in package client/lcd has been renamed to `ClientCtx`.
  * (codec) [\#6330](https://github.com/cosmos/cosmos-sdk/pull/6330) `codec.RegisterCrypto` has been moved to the `crypto/codec` package and the global `codec.Cdc` Amino instance has been deprecated and moved to the `codec/legacy_global` package.
  * (codec) [\#8080](https://github.com/cosmos/cosmos-sdk/pull/8080) Updated the `codec.Marshaler` interface
    * Moved `MarshalAny` and `UnmarshalAny` helper functions to `codec.Marshaler` and renamed to `MarshalInterface` and
      `UnmarshalInterface` respectively. These functions must take interface as a parameter (not a concrete type nor `Any`
      object). Underneath they use `Any` wrapping for correct protobuf serialization.
  * (crypto) [\#6780](https://github.com/cosmos/cosmos-sdk/issues/6780) Move ledger code to its own package.
  * (crypto/types/multisig) [\#6373](https://github.com/cosmos/cosmos-sdk/pull/6373) `multisig.Multisignature` has been renamed  to `AminoMultisignature`
  * (codec) `*codec.LegacyAmino` is now a wrapper around Amino which provides backwards compatibility with protobuf `Any`. ALL legacy code should use `*codec.LegacyAmino` instead of `*amino.Codec` directly
  * (crypto) [\#5880](https://github.com/cosmos/cosmos-sdk/pull/5880) Merge `crypto/keys/mintkey` into `crypto`.
  * (crypto/hd) [\#5904](https://github.com/cosmos/cosmos-sdk/pull/5904) `crypto/keys/hd` moved to `crypto/hd`.
  * (crypto/keyring):
    * [\#5866](https://github.com/cosmos/cosmos-sdk/pull/5866) Rename `crypto/keys/` to `crypto/keyring/`.
    * [\#5904](https://github.com/cosmos/cosmos-sdk/pull/5904) `Keybase` -> `Keyring` interfaces migration. `LegacyKeybase` interface is added in order
  to guarantee limited backward compatibility with the old Keybase interface for the sole purpose of migrating keys across the new keyring backends. `NewLegacy`
  constructor is provided [\#5889](https://github.com/cosmos/cosmos-sdk/pull/5889) to allow for smooth migration of keys from the legacy LevelDB based implementation
  to new keyring backends. Plus, the package and the new keyring no longer depends on the sdk.Config singleton. Please consult the [package documentation](https://github.com/cosmos/cosmos-sdk/tree/master/crypto/keyring/doc.go) for more
  information on how to implement the new `Keyring` interface.
    * [\#5858](https://github.com/cosmos/cosmos-sdk/pull/5858) Make Keyring store keys by name and address's hexbytes representation.
  * (export) [\#5952](https://github.com/cosmos/cosmos-sdk/pull/5952) `AppExporter` now returns ABCI consensus parameters to be included in marshaled exported state. These parameters must be returned from the application via the `BaseApp`.
  * (simapp) Deprecating and renaming `MakeEncodingConfig` to `MakeTestEncodingConfig` (both in `simapp` and `simapp/params` packages).
  * (store) [\#5803](https://github.com/cosmos/cosmos-sdk/pull/5803) The `store.CommitMultiStore` interface now includes the new `snapshots.Snapshotter` interface as well.
  * (types) [\#5579](https://github.com/cosmos/cosmos-sdk/pull/5579) The `keepRecent` field has been removed from the `PruningOptions` type.
  The `PruningOptions` type now only includes fields `KeepEvery` and `SnapshotEvery`, where `KeepEvery`
  determines which committed heights are flushed to disk and `SnapshotEvery` determines which of these
  heights are kept after pruning. The `IsValid` method should be called whenever using these options. Methods
  `SnapshotVersion` and `FlushVersion` accept a version arugment and determine if the version should be
  flushed to disk or kept as a snapshot. Note, `KeepRecent` is automatically inferred from the options
  and provided directly the IAVL store.
  * (types) [\#5533](https://github.com/cosmos/cosmos-sdk/pull/5533) Refactored `AppModuleBasic` and `AppModuleGenesis`
  to now accept a `codec.JSONMarshaler` for modular serialization of genesis state.
  * (types/rest) [\#5779](https://github.com/cosmos/cosmos-sdk/pull/5779) Drop unused Parse{Int64OrReturnBadRequest,QueryParamBool}() functions.
* __Modules__
  * (modules) [\#7243](https://github.com/cosmos/cosmos-sdk/pull/7243) Rename `RegisterCodec` to `RegisterLegacyAminoCodec` and `codec.New()` is now renamed to `codec.NewLegacyAmino()`
  * (modules) [\#6564](https://github.com/cosmos/cosmos-sdk/pull/6564) Constant `DefaultParamspace` is removed from all modules, use ModuleName instead.
  * (modules) [\#5989](https://github.com/cosmos/cosmos-sdk/pull/5989) `AppModuleBasic.GetTxCmd` now takes a single `CLIContext` parameter.
  * (modules) [\#5664](https://github.com/cosmos/cosmos-sdk/pull/5664) Remove amino `Codec` from simulation `StoreDecoder`, which now returns a function closure in order to unmarshal the key-value pairs.
  * (modules) [\#5555](https://github.com/cosmos/cosmos-sdk/pull/5555) Move `x/auth/client/utils/` types and functions to `x/auth/client/`.
  * (modules) [\#5572](https://github.com/cosmos/cosmos-sdk/pull/5572) Move account balance logic and APIs from `x/auth` to `x/bank`.
  * (modules) [\#6326](https://github.com/cosmos/cosmos-sdk/pull/6326) `AppModuleBasic.GetQueryCmd` now takes a single `client.Context` parameter.
  * (modules) [\#6336](https://github.com/cosmos/cosmos-sdk/pull/6336) `AppModuleBasic.RegisterQueryService` method was added to support gRPC queries, and `QuerierRoute` and `NewQuerierHandler` were deprecated.
  * (modules) [\#6311](https://github.com/cosmos/cosmos-sdk/issues/6311) Remove `alias.go` usage
  * (modules) [\#6447](https://github.com/cosmos/cosmos-sdk/issues/6447) Rename `blacklistedAddrs` to `blockedAddrs`.
  * (modules) [\#6834](https://github.com/cosmos/cosmos-sdk/issues/6834) Add `RegisterInterfaces` method to `AppModuleBasic` to support registration of protobuf interface types.
  * (modules) [\#6734](https://github.com/cosmos/cosmos-sdk/issues/6834) Add `TxEncodingConfig` parameter to `AppModuleBasic.ValidateGenesis` command to support JSON tx decoding in `genutil`.
  * (modules) [#7764](https://github.com/cosmos/cosmos-sdk/pull/7764) Added module initialization options:
    * `server/types.AppExporter` requires extra argument: `AppOptions`.
    * `server.AddCommands` requires extra argument: `addStartFlags types.ModuleInitFlags`
    * `x/crisis.NewAppModule` has a new attribute: `skipGenesisInvariants`. [PR](https://github.com/cosmos/cosmos-sdk/pull/7764)
  * (types) [\#6327](https://github.com/cosmos/cosmos-sdk/pull/6327) `sdk.Msg` now inherits `proto.Message`, as a result all `sdk.Msg` types now use pointer semantics.
  * (types) [\#7032](https://github.com/cosmos/cosmos-sdk/pull/7032) All types ending with `ID` (e.g. `ProposalID`) now end with `Id` (e.g. `ProposalId`), to match default Protobuf generated format. Also see [\#7033](https://github.com/cosmos/cosmos-sdk/pull/7033) for more details.
  * (x/auth) [\#6029](https://github.com/cosmos/cosmos-sdk/pull/6029) Module accounts have been moved from `x/supply` to `x/auth`.
  * (x/auth) [\#6443](https://github.com/cosmos/cosmos-sdk/issues/6443) Move `FeeTx` and `TxWithMemo` interfaces from `x/auth/ante` to `types`.
  * (x/auth) [\#7006](https://github.com/cosmos/cosmos-sdk/pull/7006) All `AccountRetriever` methods now take `client.Context` as a parameter instead of as a struct member.
  * (x/auth) [\#6270](https://github.com/cosmos/cosmos-sdk/pull/6270) The passphrase argument has been removed from the signature of the following functions and methods: `BuildAndSign`, ` MakeSignature`, ` SignStdTx`, `TxBuilder.BuildAndSign`, `TxBuilder.Sign`, `TxBuilder.SignStdTx`
  * (x/auth) [\#6428](https://github.com/cosmos/cosmos-sdk/issues/6428):
    * `NewAnteHandler` and `NewSigVerificationDecorator` both now take a `SignModeHandler` parameter.
    * `SignatureVerificationGasConsumer` now has the signature: `func(meter sdk.GasMeter, sig signing.SignatureV2, params types.Params) error`.
    * The `SigVerifiableTx` interface now has a `GetSignaturesV2() ([]signing.SignatureV2, error)` method and no longer has the `GetSignBytes` method.
  * (x/auth/tx) [\#8106](https://github.com/cosmos/cosmos-sdk/pull/8106) change related to missing append functionality in
      client transaction signing
    + added `overwriteSig` argument to `x/auth/client.SignTx` and `client/tx.Sign` functions.
    + removed `x/auth/tx.go:wrapper.GetSignatures`. The `wrapper` provides `TxBuilder` functionality, and it's a private
      structure. That function was not used at all and it's not exposed through the `TxBuilder` interface.
  * (x/bank) [\#7327](https://github.com/cosmos/cosmos-sdk/pull/7327) AddCoins and SubtractCoins no longer return a resultingValue and will only return an error.
  * (x/capability) [#7918](https://github.com/cosmos/cosmos-sdk/pull/7918) Add x/capability safety checks:
    * All outward facing APIs will now check that capability is not nil and name is not empty before performing any state-machine changes
    * `SetIndex` has been renamed to `InitializeIndex`
  * (x/evidence) [\#7251](https://github.com/cosmos/cosmos-sdk/pull/7251) New evidence types and light client evidence handling. The module function names changed.
  * (x/evidence) [\#5952](https://github.com/cosmos/cosmos-sdk/pull/5952) Remove APIs for getting and setting `x/evidence` parameters. `BaseApp` now uses a `ParamStore` to manage Tendermint consensus parameters which is managed via the `x/params` `Substore` type.
  * (x/gov) [\#6147](https://github.com/cosmos/cosmos-sdk/pull/6147) The `Content` field on `Proposal` and `MsgSubmitProposal`
    is now `Any` in concordance with [ADR 019](docs/architecture/adr-019-protobuf-state-encoding.md) and `GetContent` should now
    be used to retrieve the actual proposal `Content`. Also the `NewMsgSubmitProposal` constructor now may return an `error`
  * (x/ibc) [\#6374](https://github.com/cosmos/cosmos-sdk/pull/6374) `VerifyMembership` and `VerifyNonMembership` now take a `specs []string` argument to specify the proof format used for verification. Most SDK chains can simply use `commitmenttypes.GetSDKSpecs()` for this argument.
  * (x/params) [\#5619](https://github.com/cosmos/cosmos-sdk/pull/5619) The `x/params` keeper now accepts a `codec.Marshaller` instead of
  a reference to an amino codec. Amino is still used for JSON serialization.
  * (x/staking) [\#6451](https://github.com/cosmos/cosmos-sdk/pull/6451) `DefaultParamspace` and `ParamKeyTable` in staking module are moved from keeper to types to enforce consistency.
  * (x/staking) [\#7419](https://github.com/cosmos/cosmos-sdk/pull/7419) The `TmConsPubKey` method on ValidatorI has been
      removed and replaced instead by `ConsPubKey` (which returns a SDK `cryptotypes.PubKey`) and `TmConsPublicKey` (which
      returns a Tendermint proto PublicKey).
  * (x/staking/types) [\#7447](https://github.com/cosmos/cosmos-sdk/issues/7447) Remove bech32 PubKey support:
    * `ValidatorI` interface update. `GetConsPubKey` renamed to `TmConsPubKey` (consensus public key must be a tendermint key). `TmConsPubKey`, `GetConsAddr` methods return error.
    * `Validator` update. Methods changed in `ValidatorI` (as described above) and `ToTmValidator` return error.
    * `Validator.ConsensusPubkey` type changed from `string` to `codectypes.Any`.
    * `MsgCreateValidator.Pubkey` type changed from `string` to `codectypes.Any`.
  * (x/supply) [\#6010](https://github.com/cosmos/cosmos-sdk/pull/6010) All `x/supply` types and APIs have been moved to `x/bank`.
  * [\#6409](https://github.com/cosmos/cosmos-sdk/pull/6409) Rename all IsEmpty methods to Empty across the codebase and enforce consistency.
  * [\#6231](https://github.com/cosmos/cosmos-sdk/pull/6231) Simplify `AppModule` interface, `Route` and `NewHandler` methods become only `Route`
  and returns a new `Route` type.
  * (x/slashing) [\#6212](https://github.com/cosmos/cosmos-sdk/pull/6212) Remove `Get*` prefixes from key construction functions
  * (server) [\#6079](https://github.com/cosmos/cosmos-sdk/pull/6079) Remove `UpgradeOldPrivValFile` (deprecated in Tendermint Core v0.28).
  * [\#5719](https://github.com/cosmos/cosmos-sdk/pull/5719) Bump Go requirement to 1.14+


### State Machine Breaking

* __General__
  * (client) [\#7268](https://github.com/cosmos/cosmos-sdk/pull/7268) / [\#7147](https://github.com/cosmos/cosmos-sdk/pull/7147) Introduce new protobuf based PubKeys, and migrate PubKey in BaseAccount to use this new protobuf based PubKey format

* __Modules__
  * (modules) [\#5572](https://github.com/cosmos/cosmos-sdk/pull/5572) Separate balance from accounts per ADR 004.
    * Account balances are now persisted and retrieved via the `x/bank` module.
    * Vesting account interface has been modified to account for changes.
    * Callers to `NewBaseVestingAccount` are responsible for verifying account balance in relation to
    the original vesting amount.
    * The `SendKeeper` and `ViewKeeper` interfaces in `x/bank` have been modified to account for changes.
  * (x/auth) [\#5533](https://github.com/cosmos/cosmos-sdk/pull/5533) Migrate the `x/auth` module to use Protocol Buffers for state
  serialization instead of Amino.
    * The `BaseAccount.PubKey` field is now represented as a Bech32 string instead of a `crypto.Pubkey`.
    * `NewBaseAccountWithAddress` now returns a reference to a `BaseAccount`.
    * The `x/auth` module now accepts a `Codec` interface which extends the `codec.Marshaler` interface by
    requiring a concrete codec to know how to serialize accounts.
    * The `AccountRetriever` type now accepts a `Codec` in its constructor in order to know how to
    serialize accounts.
  * (x/bank) [\#6518](https://github.com/cosmos/cosmos-sdk/pull/6518) Support for global and per-denomination send enabled flags.
    * Existing send_enabled global flag has been moved into a Params structure as `default_send_enabled`.
    * An array of: `{denom: string, enabled: bool}` is added to bank Params to support per-denomination override of global default value.
  * (x/distribution) [\#5610](https://github.com/cosmos/cosmos-sdk/pull/5610) Migrate the `x/distribution` module to use Protocol Buffers for state
  serialization instead of Amino. The exact codec used is `codec.HybridCodec` which utilizes Protobuf for binary encoding and Amino
  for JSON encoding.
    * `ValidatorHistoricalRewards.ReferenceCount` is now of types `uint32` instead of `uint16`.
    * `ValidatorSlashEvents` is now a struct with `slashevents`.
    * `ValidatorOutstandingRewards` is now a struct with `rewards`.
    * `ValidatorAccumulatedCommission` is now a struct with `commission`.
    * The `Keeper` constructor now takes a `codec.Marshaler` instead of a concrete Amino codec. This exact type
    provided is specified by `ModuleCdc`.
  * (x/evidence) [\#5634](https://github.com/cosmos/cosmos-sdk/pull/5634) Migrate the `x/evidence` module to use Protocol Buffers for state
  serialization instead of Amino.
    * The `internal` sub-package has been removed in order to expose the types proto file.
    * The module now accepts a `Codec` interface which extends the `codec.Marshaler` interface by
    requiring a concrete codec to know how to serialize `Evidence` types.
    * The `MsgSubmitEvidence` message has been removed in favor of `MsgSubmitEvidenceBase`. The application-level
    codec must now define the concrete `MsgSubmitEvidence` type which must implement the module's `MsgSubmitEvidence`
    interface.
  * (x/evidence) [\#5952](https://github.com/cosmos/cosmos-sdk/pull/5952) Remove parameters from `x/evidence` genesis and module state. The `x/evidence` module now solely uses Tendermint consensus parameters to determine of evidence is valid or not.
  * (x/gov) [\#5737](https://github.com/cosmos/cosmos-sdk/pull/5737) Migrate the `x/gov` module to use Protocol
  Buffers for state serialization instead of Amino.
    * `MsgSubmitProposal` will be removed in favor of the application-level proto-defined `MsgSubmitProposal` which
    implements the `MsgSubmitProposalI` interface. Applications should extend the `NewMsgSubmitProposalBase` type
    to define their own concrete `MsgSubmitProposal` types.
    * The module now accepts a `Codec` interface which extends the `codec.Marshaler` interface by
    requiring a concrete codec to know how to serialize `Proposal` types.
  * (x/mint) [\#5634](https://github.com/cosmos/cosmos-sdk/pull/5634) Migrate the `x/mint` module to use Protocol Buffers for state
  serialization instead of Amino.
    * The `internal` sub-package has been removed in order to expose the types proto file.
  * (x/slashing) [\#5627](https://github.com/cosmos/cosmos-sdk/pull/5627) Migrate the `x/slashing` module to use Protocol Buffers for state
  serialization instead of Amino. The exact codec used is `codec.HybridCodec` which utilizes Protobuf for binary encoding and Amino
  for JSON encoding.
    * The `Keeper` constructor now takes a `codec.Marshaler` instead of a concrete Amino codec. This exact type
    provided is specified by `ModuleCdc`.
  * (x/staking) [\#6844](https://github.com/cosmos/cosmos-sdk/pull/6844) Validators are now inserted into the unbonding queue based on their unbonding time and height. The relevant keeper APIs are modified to reflect these changes by now also requiring a height.
  * (x/staking) [\#6061](https://github.com/cosmos/cosmos-sdk/pull/6061) Allow a validator to immediately unjail when no signing info is present due to
  falling below their minimum self-delegation and never having been bonded. The validator may immediately unjail once they've met their minimum self-delegation.
  * (x/staking) [\#5600](https://github.com/cosmos/cosmos-sdk/pull/5600) Migrate the `x/staking` module to use Protocol Buffers for state
  serialization instead of Amino. The exact codec used is `codec.HybridCodec` which utilizes Protobuf for binary encoding and Amino
  for JSON encoding.
    * `BondStatus` is now of type `int32` instead of `byte`.
    * Types of `int16` in the `Params` type are now of type `int32`.
    * Every reference of `crypto.Pubkey` in context of a `Validator` is now of type string. `GetPubKeyFromBech32` must be used to get the `crypto.Pubkey`.
    * The `Keeper` constructor now takes a `codec.Marshaler` instead of a concrete Amino codec. This exact type
    provided is specified by `ModuleCdc`.
  * (x/staking) [\#7979](https://github.com/cosmos/cosmos-sdk/pull/7979) keeper pubkey storage serialization migration
      from bech32 to protobuf.
  * (x/supply) [\#6010](https://github.com/cosmos/cosmos-sdk/pull/6010) Removed the `x/supply` module by merging the existing types and APIs into the `x/bank` module.
  * (x/supply) [\#5533](https://github.com/cosmos/cosmos-sdk/pull/5533) Migrate the `x/supply` module to use Protocol Buffers for state
  serialization instead of Amino.
    * The `internal` sub-package has been removed in order to expose the types proto file.
    * The `x/supply` module now accepts a `Codec` interface which extends the `codec.Marshaler` interface by
    requiring a concrete codec to know how to serialize `SupplyI` types.
    * The `SupplyI` interface has been modified to no longer return `SupplyI` on methods. Instead the
    concrete type's receiver should modify the type.
  * (x/upgrade) [\#5659](https://github.com/cosmos/cosmos-sdk/pull/5659) Migrate the `x/upgrade` module to use Protocol
  Buffers for state serialization instead of Amino.
    * The `internal` sub-package has been removed in order to expose the types proto file.
    * The `x/upgrade` module now accepts a `codec.Marshaler` interface.

### Features

* __Baseapp / Client / REST__
  * (x/auth) [\#6213](https://github.com/cosmos/cosmos-sdk/issues/6213) Introduce new protobuf based path for transaction signing, see [ADR020](https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-020-protobuf-transaction-encoding.md) for more details
  * (x/auth) [\#6350](https://github.com/cosmos/cosmos-sdk/pull/6350) New sign-batch command to sign StdTx batch files.
  * (baseapp) [\#5803](https://github.com/cosmos/cosmos-sdk/pull/5803) Added support for taking state snapshots at regular height intervals, via options `snapshot-interval` and `snapshot-keep-recent`.
  * (baseapp) [\#7519](https://github.com/cosmos/cosmos-sdk/pull/7519) Add `ServiceMsgRouter` to BaseApp to handle routing of protobuf service `Msg`s. The two new types defined in ADR 031, `sdk.ServiceMsg` and `sdk.MsgRequest` are introduced with this router.
  * (client) [\#5921](https://github.com/cosmos/cosmos-sdk/issues/5921) Introduce new gRPC and gRPC Gateway based APIs for querying app & module data. See [ADR021](https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-021-protobuf-query-encoding.md) for more details
  * (cli) [\#7485](https://github.com/cosmos/cosmos-sdk/pull/7485) Introduce a new optional `--keyring-dir` flag that allows clients to specify a Keyring directory if it does not reside in the directory specified by `--home`.
  * (cli) [\#7221](https://github.com/cosmos/cosmos-sdk/pull/7221) Add the option of emitting amino encoded json from the CLI
  * (codec) [\#7519](https://github.com/cosmos/cosmos-sdk/pull/7519) `InterfaceRegistry` now inherits `jsonpb.AnyResolver`, and has a `RegisterCustomTypeURL` method to support ADR 031 packing of `Any`s. `AnyResolver` is now a required parameter to `RejectUnknownFields`.
  * (coin) [\#6755](https://github.com/cosmos/cosmos-sdk/pull/6755) Add custom regex validation for `Coin` denom by overwriting `CoinDenomRegex` when using `/types/coin.go`.
  * (config) [\#7265](https://github.com/cosmos/cosmos-sdk/pull/7265) Support Tendermint block pruning through a new `min-retain-blocks` configuration that can be set in either `app.toml` or via the CLI. This parameter is used in conjunction with other criteria to determine the height at which Tendermint should prune blocks.
  * (events) [\#7121](https://github.com/cosmos/cosmos-sdk/pull/7121) The application now derives what events are indexed by Tendermint via the `index-events` configuration in `app.toml`, which is a list of events taking the form `{eventType}.{attributeKey}`.
  * (tx) [\#6089](https://github.com/cosmos/cosmos-sdk/pull/6089) Transactions can now have a `TimeoutHeight` set which allows the transaction to be rejected if it's committed at a height greater than the timeout.
  * (rest) [\#6167](https://github.com/cosmos/cosmos-sdk/pull/6167) Support `max-body-bytes` CLI flag for the REST service.
  * (genesis) [\#7089](https://github.com/cosmos/cosmos-sdk/pull/7089) The `export` command now adds a `initial_height` field in the exported JSON. Baseapp's `CommitMultiStore` now also has a `SetInitialVersion` setter, so it can set the initial store version inside `InitChain` and start a new chain from a given height.
* __General__
  * (crypto/multisig) [\#6241](https://github.com/cosmos/cosmos-sdk/pull/6241) Add Multisig type directly to the repo. Previously this was in tendermint.
  * (codec/types) [\#8106](https://github.com/cosmos/cosmos-sdk/pull/8106) Adding `NewAnyWithCustomTypeURL` to correctly
     marshal Messages in TxBuilder.
  * (tests) [\#6489](https://github.com/cosmos/cosmos-sdk/pull/6489) Introduce package `testutil`, new in-process testing network framework for use in integration and unit tests.
  * (tx) Add new auth/tx gRPC & gRPC-Gateway endpoints for basic querying & broadcasting support
    * [\#7842](https://github.com/cosmos/cosmos-sdk/pull/7842) Add TxsByEvent gRPC endpoint
    * [\#7852](https://github.com/cosmos/cosmos-sdk/pull/7852) Add tx broadcast gRPC endpoint
  * (tx) [\#7688](https://github.com/cosmos/cosmos-sdk/pull/7688) Add a new Tx gRPC service with methods `Simulate` and `GetTx` (by hash).
  * (store) [\#5803](https://github.com/cosmos/cosmos-sdk/pull/5803) Added `rootmulti.Store` methods for taking and restoring snapshots, based on `iavl.Store` export/import.
  * (store) [\#6324](https://github.com/cosmos/cosmos-sdk/pull/6324) IAVL store query proofs now return CommitmentOp which wraps an ics23 CommitmentProof
  * (store) [\#6390](https://github.com/cosmos/cosmos-sdk/pull/6390) `RootMulti` store query proofs now return `CommitmentOp` which wraps `CommitmentProofs`
    * `store.Query` now only returns chained `ics23.CommitmentProof` wrapped in `merkle.Proof`
    * `ProofRuntime` only decodes and verifies `ics23.CommitmentProof`
* __Modules__
  * (modules) [\#5921](https://github.com/cosmos/cosmos-sdk/issues/5921) Introduction of Query gRPC service definitions along with REST annotations for gRPC Gateway for each module
  * (modules) [\#7540](https://github.com/cosmos/cosmos-sdk/issues/7540) Protobuf service definitions can now be used for
    packing `Msg`s in transactions as defined in [ADR 031](./docs/architecture/adr-031-msg-service.md). All modules now
    define a `Msg` protobuf service.
  * (x/auth/vesting) [\#7209](https://github.com/cosmos/cosmos-sdk/pull/7209) Create new `MsgCreateVestingAccount` message type along with CLI handler that allows for the creation of delayed and continuous vesting types.
  * (x/capability) [\#5828](https://github.com/cosmos/cosmos-sdk/pull/5828) Capability module integration as outlined in [ADR 3 - Dynamic Capability Store](https://github.com/cosmos/tree/master/docs/architecture/adr-003-dynamic-capability-store.md).
  * (x/crisis) `x/crisis` has a new function: `AddModuleInitFlags`, which will register optional crisis module flags for the start command.
  * (x/ibc) [\#5277](https://github.com/cosmos/cosmos-sdk/pull/5277) `x/ibc` changes from IBC alpha. For more details check the the [`x/ibc/core/spec`](https://github.com/cosmos/cosmos-sdk/tree/master/x/ibc/core/spec) directory, or the ICS specs below:
    * [ICS 002 - Client Semantics](https://github.com/cosmos/ics/tree/master/spec/ics-002-client-semantics) subpackage
    * [ICS 003 - Connection Semantics](https://github.com/cosmos/ics/blob/master/spec/ics-003-connection-semantics) subpackage
    * [ICS 004 - Channel and Packet Semantics](https://github.com/cosmos/ics/blob/master/spec/ics-004-channel-and-packet-semantics) subpackage
    * [ICS 005 - Port Allocation](https://github.com/cosmos/ics/blob/master/spec/ics-005-port-allocation) subpackage
    * [ICS 006 - Solo Machine Client](https://github.com/cosmos/ics/tree/master/spec/ics-006-solo-machine-client) subpackage
    * [ICS 007 - Tendermint Client](https://github.com/cosmos/ics/blob/master/spec/ics-007-tendermint-client) subpackage
    * [ICS 009 - Loopback Client](https://github.com/cosmos/ics/tree/master/spec/ics-009-loopback-client) subpackage
    * [ICS 020 - Fungible Token Transfer](https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer) subpackage
    * [ICS 023 - Vector Commitments](https://github.com/cosmos/ics/tree/master/spec/ics-023-vector-commitments) subpackage
    * [ICS 024 - Host State Machine Requirements](https://github.com/cosmos/ics/tree/master/spec/ics-024-host-requirements) subpackage
  * (x/ibc) [\#6374](https://github.com/cosmos/cosmos-sdk/pull/6374) ICS-23 Verify functions will now accept and verify ics23 CommitmentProofs exclusively
  * (x/params) [\#6005](https://github.com/cosmos/cosmos-sdk/pull/6005) Add new CLI command for querying raw x/params parameters by subspace and key.

### Bug Fixes

* __Baseapp / Client / REST__
  * (client) [\#5964](https://github.com/cosmos/cosmos-sdk/issues/5964) `--trust-node` is now false by default - for real. Users must ensure it is set to true if they don't want to enable the verifier.
  * (client) [\#6402](https://github.com/cosmos/cosmos-sdk/issues/6402) Fix `keys add` `--algo` flag which only worked for Tendermint's `secp256k1` default key signing algorithm.
  * (client) [\#7699](https://github.com/cosmos/cosmos-sdk/pull/7699) Fix panic in context when setting invalid nodeURI. `WithNodeURI` does not set the `Client` in the context.
  * (export) [\#6510](https://github.com/cosmos/cosmos-sdk/pull/6510/) Field TimeIotaMs now is included in genesis file while exporting.
  * (rest) [\#5906](https://github.com/cosmos/cosmos-sdk/pull/5906) Fix an issue that make some REST calls panic when sending invalid or incomplete requests.
  * (crypto) [\#7966](https://github.com/cosmos/cosmos-sdk/issues/7966) `Bip44Params` `String()` function now correctly
      returns the absolute HD path by adding the `m/` prefix.
  * (crypto/keyring) [\#5844](https://github.com/cosmos/cosmos-sdk/pull/5844) `Keyring.Sign()` methods no longer decode amino signatures when method receivers
  are offline/multisig keys.
  * (store) [\#7415](https://github.com/cosmos/cosmos-sdk/pull/7415) Allow new stores to be registered during on-chain upgrades.
* __Modules__
  * (modules) [\#5569](https://github.com/cosmos/cosmos-sdk/issues/5569) `InitGenesis`, for the relevant modules, now ensures module accounts exist.
  * (x/auth) [\#5892](https://github.com/cosmos/cosmos-sdk/pull/5892) Add `RegisterKeyTypeCodec` to register new
  types (eg. keys) to the `auth` module internal amino codec.
  * (x/bank) [\#6536](https://github.com/cosmos/cosmos-sdk/pull/6536) Fix bug in `WriteGeneratedTxResponse` function used by multiple
  REST endpoints. Now it writes a Tx in StdTx format.
  * (x/genutil) [\#5938](https://github.com/cosmos/cosmos-sdk/pull/5938) Fix `InitializeNodeValidatorFiles` error handling.
  * (x/gentx) [\#8183](https://github.com/cosmos/cosmos-sdk/pull/8183) change gentx cmd amount to arg from flag
  * (x/gov) [#7641](https://github.com/cosmos/cosmos-sdk/pull/7641) Fix tally calculation precision error.
  * (x/staking) [\#6529](https://github.com/cosmos/cosmos-sdk/pull/6529) Export validator addresses (previously was empty).
  * (x/staking) [\#5949](https://github.com/cosmos/cosmos-sdk/pull/5949) Skip staking `HistoricalInfoKey` in simulations as headers are not exported.
  * (x/staking) [\#6061](https://github.com/cosmos/cosmos-sdk/pull/6061) Allow a validator to immediately unjail when no signing info is present due to
falling below their minimum self-delegation and never having been bonded. The validator may immediately unjail once they've met their minimum self-delegation.
* __General__
  * (types) [\#7038](https://github.com/cosmos/cosmos-sdk/issues/7038) Fix infinite looping of `ApproxRoot` by including a hard-coded maximum iterations limit of 100.
  * (types) [\#7084](https://github.com/cosmos/cosmos-sdk/pull/7084) Fix panic when calling `BigInt()` on an uninitialized `Int`.
  * (simulation) [\#7129](https://github.com/cosmos/cosmos-sdk/issues/7129) Fix support for custom `Account` and key types on auth's simulation.


### Improvements
* __Baseapp / Client / REST__
  * (baseapp) [\#6186](https://github.com/cosmos/cosmos-sdk/issues/6186) Support emitting events during `AnteHandler` execution.
  * (baseapp) [\#6053](https://github.com/cosmos/cosmos-sdk/pull/6053) Customizable panic recovery handling added for `app.runTx()` method (as proposed in the [ADR 22](https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-022-custom-panic-handling.md)). Adds ability for developers to register custom panic handlers extending standard ones.
  * (client) [\#5810](https://github.com/cosmos/cosmos-sdk/pull/5810) Added a new `--offline` flag that allows commands to be executed without an
  internet connection. Previously, `--generate-only` served this purpose in addition to only allowing txs to be generated. Now, `--generate-only` solely
  allows txs to be generated without being broadcasted and disallows Keybase use and `--offline` allows the use of Keybase but does not allow any
  functionality that requires an online connection.
  * (cli) [#7764](https://github.com/cosmos/cosmos-sdk/pull/7764) Update x/banking and x/crisis InitChain to improve node startup time
  * (client) [\#5856](https://github.com/cosmos/cosmos-sdk/pull/5856) Added the possibility to set `--offline` flag with config command.
  * (client) [\#5895](https://github.com/cosmos/cosmos-sdk/issues/5895) show config options in the config command's help screen.
  * (client/keys) [\#8043](https://github.com/cosmos/cosmos-sdk/pull/8043) Add support for export of unarmored private key
  * (client/tx) [\#7801](https://github.com/cosmos/cosmos-sdk/pull/7801) Update sign-batch multisig to work online
  * (x/genutil) [\#8099](https://github.com/cosmos/cosmos-sdk/pull/8099) `init` now supports a `--recover` flag to recover
      the private validator key from a given mnemonic
* __Modules__
  * (x/auth) [\#5702](https://github.com/cosmos/cosmos-sdk/pull/5702) Add parameter querying support for `x/auth`.
  * (x/auth/ante) [\#6040](https://github.com/cosmos/cosmos-sdk/pull/6040) `AccountKeeper` interface used for `NewAnteHandler` and handler's decorators to add support of using custom `AccountKeeper` implementations.
  * (x/evidence) [\#5952](https://github.com/cosmos/cosmos-sdk/pull/5952) Tendermint Consensus parameters can now be changed via parameter change proposals through `x/gov`.
  * (x/evidence) [\#5961](https://github.com/cosmos/cosmos-sdk/issues/5961) Add `StoreDecoder` simulation for evidence module.
  * (x/ibc) [\#5948](https://github.com/cosmos/cosmos-sdk/issues/5948) Add `InitGenesis` and `ExportGenesis` functions for `ibc` module.
  * (x/ibc-transfer) [\#6871](https://github.com/cosmos/cosmos-sdk/pull/6871) Implement [ADR 001 - Coin Source Tracing](./docs/architecture/adr-001-coin-source-tracing.md).
  * (x/staking) [\#6059](https://github.com/cosmos/cosmos-sdk/pull/6059) Updated `HistoricalEntries` parameter default to 100.
  * (x/staking) [\#5584](https://github.com/cosmos/cosmos-sdk/pull/5584) Add util function `ToTmValidator` that converts a `staking.Validator` type to `*tmtypes.Validator`.
  * (x/staking) [\#6163](https://github.com/cosmos/cosmos-sdk/pull/6163) CLI and REST call to unbonding delegations and delegations now accept
  pagination.
  * (x/staking) [\#8178](https://github.com/cosmos/cosmos-sdk/pull/8178) Update default historical header number for stargate
* __General__
  * (crypto) [\#7987](https://github.com/cosmos/cosmos-sdk/pull/7987) Fix the inconsistency of CryptoCdc, only use
      `codec/legacy.Cdc`.
  * (logging) [\#8072](https://github.com/cosmos/cosmos-sdk/pull/8072) Refactor logging:
    * Use [zerolog](https://github.com/rs/zerolog) over Tendermint's go-kit logging wrapper.
    * Introduce Tendermint's `--log_format=plain|json` flag. Using format `json` allows for emitting structured JSON
    logs which can be consumed by an external logging facility (e.g. Loggly). Both formats log to STDERR.
    * The existing `--log_level` flag and it's default value now solely relates to the global logging
    level (e.g. `info`, `debug`, etc...) instead of `<module>:<level>`.
  * (rest) [#7649](https://github.com/cosmos/cosmos-sdk/pull/7649) Return an unsigned tx in legacy GET /tx endpoint when signature conversion fails
  * (simulation) [\#6002](https://github.com/cosmos/cosmos-sdk/pull/6002) Add randomized consensus params into simulation.
  * (store) [\#6481](https://github.com/cosmos/cosmos-sdk/pull/6481) Move `SimpleProofsFromMap` from Tendermint into the SDK.
  * (store) [\#6719](https://github.com/cosmos/cosmos-sdk/6754) Add validity checks to stores for nil and empty keys.
  * (SDK) Updated dependencies
    * Updated iavl dependency to v0.15.3
    * Update tendermint to v0.34.1
  * (types) [\#7027](https://github.com/cosmos/cosmos-sdk/pull/7027) `Coin(s)` and `DecCoin(s)` updates:
    * Bump denomination max length to 128
    * Allow uppercase letters and numbers in denominations to support [ADR 001](./docs/architecture/adr-001-coin-source-tracing.md)
    * Added `Validate` function that returns a descriptive error
  * (types) [\#5581](https://github.com/cosmos/cosmos-sdk/pull/5581) Add convenience functions {,Must}Bech32ifyAddressBytes.
  * (types/module) [\#5724](https://github.com/cosmos/cosmos-sdk/issues/5724) The `types/module` package does no longer depend on `x/simulation`.
  * (types) [\#5585](https://github.com/cosmos/cosmos-sdk/pull/5585) IBC additions:
    * `Coin` denomination max lenght has been increased to 32.
    * Added `CapabilityKey` alias for `StoreKey` to match IBC spec.
  * (types/rest) [\#5900](https://github.com/cosmos/cosmos-sdk/pull/5900) Add Check*Error function family to spare developers from replicating tons of boilerplate code.
  * (types) [\#6128](https://github.com/cosmos/cosmos-sdk/pull/6137) Add `String()` method to `GasMeter`.
  * (types) [\#6195](https://github.com/cosmos/cosmos-sdk/pull/6195) Add codespace to broadcast(sync/async) response.
  * (types) \#6897 Add KV type from tendermint to `types` directory.
  * (version) [\#7848](https://github.com/cosmos/cosmos-sdk/pull/7848) [\#7941](https://github.com/cosmos/cosmos-sdk/pull/7941)
    `version --long` output now shows the list of build dependencies and replaced build dependencies.

## [v0.39.1](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.39.1) - 2020-08-11

### Client Breaking

* (x/auth) [\#6861](https://github.com/cosmos/cosmos-sdk/pull/6861) Remove public key Bech32 encoding for all account types for JSON serialization, instead relying on direct Amino encoding. In addition, JSON serialization utilizes Amino instead of the Go stdlib, so integers are treated as strings.

### Improvements

* (client) [\#6853](https://github.com/cosmos/cosmos-sdk/pull/6853) Add --unsafe-cors flag.

## [v0.39.0](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.39.0) - 2020-07-20

### Improvements

* (deps) Bump IAVL version to [v0.14.0](https://github.com/cosmos/iavl/releases/tag/v0.14.0)
* (client) [\#5585](https://github.com/cosmos/cosmos-sdk/pull/5585) `CLIContext` additions:
  * Introduce `QueryABCI` that returns the full `abci.ResponseQuery` with inclusion Merkle proofs.
  * Added `prove` flag for Merkle proof verification.
* (x/staking) [\#6791)](https://github.com/cosmos/cosmos-sdk/pull/6791) Close {UBDQueue,RedelegationQueu}Iterator once used.

### API Breaking Changes

* (baseapp) [\#5837](https://github.com/cosmos/cosmos-sdk/issues/5837) Transaction simulation now returns a `SimulationResponse` which contains the `GasInfo` and `Result` from the execution.

### Client Breaking Changes

* (x/auth) [\#6745](https://github.com/cosmos/cosmos-sdk/issues/6745) Remove BaseAccount's custom JSON {,un}marshalling.

### Bug Fixes

* (store) [\#6475](https://github.com/cosmos/cosmos-sdk/pull/6475) Revert IAVL pruning functionality introduced in
[v0.13.0](https://github.com/cosmos/iavl/releases/tag/v0.13.0),
where the IAVL no longer keeps states in-memory in which it flushes periodically. IAVL now commits and
flushes every state to disk as it did pre-v0.13.0. The SDK's multi-store will track and ensure the proper
heights are pruned. The operator can set the pruning options via a `pruning` config via the CLI or
through `app.toml`. The `pruning` flag exposes `default|everything|nothing|custom` as options --
see docs for further details. If the operator chooses `custom`, they may provide granular pruning
options `pruning-keep-recent`, `pruning-keep-every`, and `pruning-interval`. The former two options
dictate how many recent versions are kept on disk and the offset of what versions are kept after that
respectively, and the latter defines the height interval in which versions are deleted in a batch.
**Note, there are some client-facing API breaking changes with regard to IAVL, stores, and pruning settings.**
* (x/distribution) [\#6210](https://github.com/cosmos/cosmos-sdk/pull/6210) Register `MsgFundCommunityPool` in distribution amino codec.
* (types) [\#5741](https://github.com/cosmos/cosmos-sdk/issues/5741) Prevent `ChainAnteDecorators()` from panicking when empty `AnteDecorator` slice is supplied.
* (baseapp) [\#6306](https://github.com/cosmos/cosmos-sdk/issues/6306) Prevent events emitted by the antehandler from being persisted between transactions.
* (client/keys) [\#5091](https://github.com/cosmos/cosmos-sdk/issues/5091) `keys parse` does not honor client app's configuration.
* (x/bank) [\#6674](https://github.com/cosmos/cosmos-sdk/pull/6674) Create account if recipient does not exist on handing `MsgMultiSend`.
* (x/auth) [\#6287](https://github.com/cosmos/cosmos-sdk/pull/6287) Fix nonce stuck when sending multiple transactions from an account in a same block.

## [v0.38.5] - 2020-07-02

### Improvements

* (tendermint) Bump Tendermint version to [v0.33.6](https://github.com/tendermint/tendermint/releases/tag/v0.33.6).

## [v0.38.4] - 2020-05-21

### Bug Fixes

* (x/auth) [\#5950](https://github.com/cosmos/cosmos-sdk/pull/5950) Fix `IncrementSequenceDecorator` to use is `IsReCheckTx` instead of `IsCheckTx` to allow account sequence incrementing.

## [v0.38.3] - 2020-04-09

### Improvements

* (tendermint) Bump Tendermint version to [v0.33.3](https://github.com/tendermint/tendermint/releases/tag/v0.33.3).

## [v0.38.2] - 2020-03-25

### Bug Fixes

* (baseapp) [\#5718](https://github.com/cosmos/cosmos-sdk/pull/5718) Remove call to `ctx.BlockGasMeter` during failed message validation which resulted in a panic when the tx execution mode was `CheckTx`.
* (x/genutil) [\#5775](https://github.com/cosmos/cosmos-sdk/pull/5775) Fix `ExportGenesis` in `x/genutil` to export default genesis state (`[]`) instead of `null`.
* (client) [\#5618](https://github.com/cosmos/cosmos-sdk/pull/5618) Fix crash on the client when the verifier is not set.
* (crypto/keys/mintkey) [\#5823](https://github.com/cosmos/cosmos-sdk/pull/5823) fix errors handling in `UnarmorPubKeyBytes` (underlying armoring function's return error was not being checked).
* (x/distribution) [\#5620](https://github.com/cosmos/cosmos-sdk/pull/5620) Fix nil pointer deref in distribution tax/reward validation helpers.

### Improvements

* (rest) [\#5648](https://github.com/cosmos/cosmos-sdk/pull/5648) Enhance /txs usability:
  * Add `tx.minheight` key to filter transaction with an inclusive minimum block height
  * Add `tx.maxheight` key to filter transaction with an inclusive maximum block height
* (crypto/keys) [\#5739](https://github.com/cosmos/cosmos-sdk/pull/5739) Print an error message if the password input failed.

## [v0.38.1] - 2020-02-11

### Improvements

* (modules) [\#5597](https://github.com/cosmos/cosmos-sdk/pull/5597) Add `amount` event attribute to the `complete_unbonding`
and `complete_redelegation` events that reflect the total balances of the completed unbondings and redelegations
respectively.

### Bug Fixes

* (types) [\#5579](https://github.com/cosmos/cosmos-sdk/pull/5579) The IAVL `Store#Commit` method has been refactored to
delete a flushed version if it is not a snapshot version. The root multi-store now keeps track of `commitInfo` instead
of `types.CommitID`. During `Commit` of the root multi-store, `lastCommitInfo` is updated from the saved state
and is only flushed to disk if it is a snapshot version. During `Query` of the root multi-store, if the request height
is the latest height, we'll use the store's `lastCommitInfo`. Otherwise, we fetch `commitInfo` from disk.
* (x/bank) [\#5531](https://github.com/cosmos/cosmos-sdk/issues/5531) Added missing amount event to MsgMultiSend, emitted for each output.
* (x/gov) [\#5622](https://github.com/cosmos/cosmos-sdk/pull/5622) Track any events emitted from a proposal's handler upon successful execution.

## [v0.38.0] - 2020-01-23

### State Machine Breaking

* (genesis) [\#5506](https://github.com/cosmos/cosmos-sdk/pull/5506) The `x/distribution` genesis state
  now includes `params` instead of individual parameters.
* (genesis) [\#5017](https://github.com/cosmos/cosmos-sdk/pull/5017) The `x/genaccounts` module has been
deprecated and all components removed except the `legacy/` package. This requires changes to the
genesis state. Namely, `accounts` now exist under `app_state.auth.accounts`. The corresponding migration
logic has been implemented for v0.38 target version. Applications can migrate via:
`$ {appd} migrate v0.38 genesis.json`.
* (modules) [\#5299](https://github.com/cosmos/cosmos-sdk/pull/5299) Handling of `ABCIEvidenceTypeDuplicateVote`
  during `BeginBlock` along with the corresponding parameters (`MaxEvidenceAge`) have moved from the
  `x/slashing` module to the `x/evidence` module.

### API Breaking Changes

* (modules) [\#5506](https://github.com/cosmos/cosmos-sdk/pull/5506) Remove individual setters of `x/distribution` parameters. Instead, follow the module spec in getting parameters, setting new value(s) and finally calling `SetParams`.
* (types) [\#5495](https://github.com/cosmos/cosmos-sdk/pull/5495) Remove redundant `(Must)Bech32ify*` and `(Must)Get*KeyBech32` functions in favor of `(Must)Bech32ifyPubKey` and `(Must)GetPubKeyFromBech32` respectively, both of which take a `Bech32PubKeyType` (string).
* (types) [\#5430](https://github.com/cosmos/cosmos-sdk/pull/5430) `DecCoins#Add` parameter changed from `DecCoins`
to `...DecCoin`, `Coins#Add` parameter changed from `Coins` to `...Coin`.
* (baseapp/types) [\#5421](https://github.com/cosmos/cosmos-sdk/pull/5421) The `Error` interface (`types/errors.go`)
has been removed in favor of the concrete type defined in `types/errors/` which implements the standard `error` interface.
  * As a result, the `Handler` and `Querier` implementations now return a standard `error`.
  Within `BaseApp`, `runTx` now returns a `(GasInfo, *Result, error)` tuple and `runMsgs` returns a
  `(*Result, error)` tuple. A reference to a `Result` is now used to indicate success whereas an error
  signals an invalid message or failed message execution. As a result, the fields `Code`, `Codespace`,
  `GasWanted`, and `GasUsed` have been removed the `Result` type. The latter two fields are now found
  in the `GasInfo` type which is always returned regardless of execution outcome.
  * Note to developers: Since all handlers and queriers must now return a standard `error`, the `types/errors/`
  package contains all the relevant and pre-registered errors that you typically work with. A typical
  error returned will look like `sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "...")`. You can retrieve
  relevant ABCI information from the error via `ABCIInfo`.
* (client) [\#5442](https://github.com/cosmos/cosmos-sdk/pull/5442) Remove client/alias.go as it's not necessary and
components can be imported directly from the packages.
* (store) [\#4748](https://github.com/cosmos/cosmos-sdk/pull/4748) The `CommitMultiStore` interface
now requires a `SetInterBlockCache` method. Applications that do not wish to support this can simply
have this method perform a no-op.
* (modules) [\#4665](https://github.com/cosmos/cosmos-sdk/issues/4665) Refactored `x/gov` module structure and dev-UX:
  * Prepare for module spec integration
  * Update gov keys to use big endian encoding instead of little endian
* (modules) [\#5017](https://github.com/cosmos/cosmos-sdk/pull/5017) The `x/genaccounts` module has been deprecated and all components removed except the `legacy/` package.
* [\#4486](https://github.com/cosmos/cosmos-sdk/issues/4486) Vesting account types decoupled from the `x/auth` module and now live under `x/auth/vesting`. Applications wishing to use vesting account types must be sure to register types via `RegisterCodec` under the new vesting package.
* [\#4486](https://github.com/cosmos/cosmos-sdk/issues/4486) The `NewBaseVestingAccount` constructor returns an error
if the provided arguments are invalid.
* (x/auth) [\#5006](https://github.com/cosmos/cosmos-sdk/pull/5006) Modular `AnteHandler` via composable decorators:
  * The `AnteHandler` interface now returns `(newCtx Context, err error)` instead of `(newCtx Context, result sdk.Result, abort bool)`
  * The `NewAnteHandler` function returns an `AnteHandler` function that returns the new `AnteHandler`
  interface and has been moved into the `auth/ante` directory.
  * `ValidateSigCount`, `ValidateMemo`, `ProcessPubKey`, `EnsureSufficientMempoolFee`, and `GetSignBytes`
  have all been removed as public functions.
  * Invalid Signatures may return `InvalidPubKey` instead of `Unauthorized` error, since the transaction
  will first hit `SetPubKeyDecorator` before the `SigVerificationDecorator` runs.
  * `StdTx#GetSignatures` will return an array of just signature byte slices `[][]byte` instead of
  returning an array of `StdSignature` structs. To replicate the old behavior, use the public field
  `StdTx.Signatures` to get back the array of StdSignatures `[]StdSignature`.
* (modules) [\#5299](https://github.com/cosmos/cosmos-sdk/pull/5299) `HandleDoubleSign` along with params `MaxEvidenceAge` and `DoubleSignJailEndTime` have moved from the `x/slashing` module to the `x/evidence` module.
* (keys) [\#4941](https://github.com/cosmos/cosmos-sdk/issues/4941) Keybase concrete types constructors such as `NewKeyBaseFromDir` and `NewInMemory` now accept optional parameters of type `KeybaseOption`. These
optional parameters are also added on the keys sub-commands functions, which are now public, and allows
these options to be set on the commands or ignored to default to previous behavior.
* [\#5547](https://github.com/cosmos/cosmos-sdk/pull/5547) `NewKeyBaseFromHomeFlag` constructor has been removed.
* [\#5439](https://github.com/cosmos/cosmos-sdk/pull/5439) Further modularization was done to the `keybase`
package to make it more suitable for use with different key formats and algorithms:
  * The `WithKeygenFunc` function added as a `KeybaseOption` which allows a custom bytes to key
    implementation to be defined when keys are created.
  * The `WithDeriveFunc` function added as a `KeybaseOption` allows custom logic for deriving a key
    from a mnemonic, bip39 password, and HD Path.
  * BIP44 is no longer build into `keybase.CreateAccount()`. It is however the default when using
    the `client/keys` add command.
  * `SupportedAlgos` and `SupportedAlgosLedger` functions return a slice of `SigningAlgo`s that are
    supported by the keybase and the ledger integration respectively.
* (simapp) [\#5419](https://github.com/cosmos/cosmos-sdk/pull/5419) The `helpers.GenTx()` now accepts a gas argument.
* (baseapp) [\#5455](https://github.com/cosmos/cosmos-sdk/issues/5455) A `sdk.Context` is now passed into the `router.Route()` function.

### Client Breaking Changes

* (rest) [\#5270](https://github.com/cosmos/cosmos-sdk/issues/5270) All account types now implement custom JSON serialization.
* (rest) [\#4783](https://github.com/cosmos/cosmos-sdk/issues/4783) The balance field in the DelegationResponse type is now sdk.Coin instead of sdk.Int
* (x/auth) [\#5006](https://github.com/cosmos/cosmos-sdk/pull/5006) The gas required to pass the `AnteHandler` has
increased significantly due to modular `AnteHandler` support. Increase GasLimit accordingly.
* (rest) [\#5336](https://github.com/cosmos/cosmos-sdk/issues/5336) `MsgEditValidator` uses `description` instead of `Description` as a JSON key.
* (keys) [\#5097](https://github.com/cosmos/cosmos-sdk/pull/5097) Due to the keybase -> keyring transition, keys need to be migrated. See `keys migrate` command for more info.
* (x/auth) [\#5424](https://github.com/cosmos/cosmos-sdk/issues/5424) Drop `decode-tx` command from x/auth/client/cli, duplicate of the `decode` command.

### Features

* (store) [\#5435](https://github.com/cosmos/cosmos-sdk/pull/5435) New iterator for paginated requests. Iterator limits DB reads to the range of the requested page.
* (x/evidence) [\#5240](https://github.com/cosmos/cosmos-sdk/pull/5240) Initial implementation of the `x/evidence` module.
* (cli) [\#5212](https://github.com/cosmos/cosmos-sdk/issues/5212) The `q gov proposals` command now supports pagination.
* (store) [\#4724](https://github.com/cosmos/cosmos-sdk/issues/4724) Multistore supports substore migrations upon load. New `rootmulti.Store.LoadLatestVersionAndUpgrade` method in
`Baseapp` supports `StoreLoader` to enable various upgrade strategies. It no
longer panics if the store to load contains substores that we didn't explicitly mount.
* [\#4972](https://github.com/cosmos/cosmos-sdk/issues/4972) A `TxResponse` with a corresponding code
and tx hash will be returned for specific Tendermint errors:
  * `CodeTxInMempoolCache`
  * `CodeMempoolIsFull`
  * `CodeTxTooLarge`
* [\#3872](https://github.com/cosmos/cosmos-sdk/issues/3872) Implement a RESTful endpoint and cli command to decode transactions.
* (keys) [\#4754](https://github.com/cosmos/cosmos-sdk/pull/4754) Introduce new Keybase implementation that can
leverage operating systems' built-in functionalities to securely store secrets. MacOS users may encounter
the following [issue](https://github.com/keybase/go-keychain/issues/47) with the `go-keychain` library. If
you encounter this issue, you must upgrade your xcode command line tools to version >= `10.2`. You can
upgrade via: `sudo rm -rf /Library/Developer/CommandLineTools; xcode-select --install`. Verify the
correct version via: `pkgutil --pkg-info=com.apple.pkg.CLTools_Executables`.
* [\#5355](https://github.com/cosmos/cosmos-sdk/pull/5355) Client commands accept a new `--keyring-backend` option through which users can specify which backend should be used
by the new key store:
  * `os`: use OS default credentials storage (default).
  * `file`: use encrypted file-based store.
  * `kwallet`: use [KDE Wallet](https://utils.kde.org/projects/kwalletmanager/) service.
  * `pass`: use the [pass](https://www.passwordstore.org/) command line password manager.
  * `test`: use password-less key store. *For testing purposes only. Use it at your own risk.*
* (keys) [\#5097](https://github.com/cosmos/cosmos-sdk/pull/5097) New `keys migrate` command to assist users migrate their keys
to the new keyring.
* (keys) [\#5366](https://github.com/cosmos/cosmos-sdk/pull/5366) `keys list` now accepts a `--list-names` option to list key names only, whilst the `keys delete`
command can delete multiple keys by passing their names as arguments. The aforementioned commands can then be piped together, e.g.
`appcli keys list -n | xargs appcli keys delete`
* (modules) [\#4233](https://github.com/cosmos/cosmos-sdk/pull/4233) Add upgrade module that coordinates software upgrades of live chains.
* [\#4486](https://github.com/cosmos/cosmos-sdk/issues/4486) Introduce new `PeriodicVestingAccount` vesting account type
that allows for arbitrary vesting periods.
* (baseapp) [\#5196](https://github.com/cosmos/cosmos-sdk/pull/5196) Baseapp has a new `runTxModeReCheck` to allow applications to skip expensive and unnecessary re-checking of transactions.
* (types) [\#5196](https://github.com/cosmos/cosmos-sdk/pull/5196) Context has new `IsRecheckTx() bool` and `WithIsReCheckTx(bool) Context` methods to to be used in the `AnteHandler`.
* (x/auth/ante) [\#5196](https://github.com/cosmos/cosmos-sdk/pull/5196) AnteDecorators have been updated to avoid unnecessary checks when `ctx.IsReCheckTx() == true`
* (x/auth) [\#5006](https://github.com/cosmos/cosmos-sdk/pull/5006) Modular `AnteHandler` via composable decorators:
  * The `AnteDecorator` interface has been introduced to allow users to implement modular `AnteHandler`
  functionality that can be composed together to create a single `AnteHandler` rather than implementing
  a custom `AnteHandler` completely from scratch, where each `AnteDecorator` allows for custom behavior in
  tightly defined and logically isolated manner. These custom `AnteDecorator` can then be chained together
  with default `AnteDecorator` or third-party `AnteDecorator` to create a modularized `AnteHandler`
  which will run each `AnteDecorator` in the order specified in `ChainAnteDecorators`. For details
  on the new architecture, refer to the [ADR](docs/architecture/adr-010-modular-antehandler.md).
  * `ChainAnteDecorators` function has been introduced to take in a list of `AnteDecorators` and chain
  them in sequence and return a single `AnteHandler`:
    * `SetUpContextDecorator`: Sets `GasMeter` in context and creates defer clause to recover from any
    `OutOfGas` panics in future AnteDecorators and return `OutOfGas` error to `BaseApp`. It MUST be the
    first `AnteDecorator` in the chain for any application that uses gas (or another one that sets the gas meter).
    * `ValidateBasicDecorator`: Calls tx.ValidateBasic and returns any non-nil error.
    * `ValidateMemoDecorator`: Validates tx memo with application parameters and returns any non-nil error.
    * `ConsumeGasTxSizeDecorator`: Consumes gas proportional to the tx size based on application parameters.
    * `MempoolFeeDecorator`: Checks if fee is above local mempool `minFee` parameter during `CheckTx`.
    * `DeductFeeDecorator`: Deducts the `FeeAmount` from first signer of the transaction.
    * `SetPubKeyDecorator`: Sets pubkey of account in any account that does not already have pubkey saved in state machine.
    * `SigGasConsumeDecorator`: Consume parameter-defined amount of gas for each signature.
    * `SigVerificationDecorator`: Verify each signature is valid, return if there is an error.
    * `ValidateSigCountDecorator`: Validate the number of signatures in tx based on app-parameters.
    * `IncrementSequenceDecorator`: Increments the account sequence for each signer to prevent replay attacks.
* (cli) [\#5223](https://github.com/cosmos/cosmos-sdk/issues/5223) Cosmos Ledger App v2.0.0 is now supported. The changes are backwards compatible and App v1.5.x is still supported.
* (x/staking) [\#5380](https://github.com/cosmos/cosmos-sdk/pull/5380) Introduced ability to store historical info entries in staking keeper, allows applications to introspect specified number of past headers and validator sets
  * Introduces new parameter `HistoricalEntries` which allows applications to determine how many recent historical info entries they want to persist in store. Default value is 0.
  * Introduces cli commands and rest routes to query historical information at a given height
* (modules) [\#5249](https://github.com/cosmos/cosmos-sdk/pull/5249) Funds are now allowed to be directly sent to the community pool (via the distribution module account).
* (keys) [\#4941](https://github.com/cosmos/cosmos-sdk/issues/4941) Introduce keybase option to allow overriding the default private key implementation of a key generated through the `keys add` cli command.
* (keys) [\#5439](https://github.com/cosmos/cosmos-sdk/pull/5439) Flags `--algo` and `--hd-path` are added to
  `keys add` command in order to make use of keybase modularized. By default, it uses (0, 0) bip44
  HD path and secp256k1 keys, so is non-breaking.
* (types) [\#5447](https://github.com/cosmos/cosmos-sdk/pull/5447) Added `ApproxRoot` function to sdk.Decimal type in order to get the nth root for a decimal number, where n is a positive integer.
  * An `ApproxSqrt` function was also added for convenience around the common case of n=2.

### Improvements

* (iavl) [\#5538](https://github.com/cosmos/cosmos-sdk/pull/5538) Remove manual IAVL pruning in favor of IAVL's internal pruning strategy.
* (server) [\#4215](https://github.com/cosmos/cosmos-sdk/issues/4215) The `--pruning` flag
has been moved to the configuration file, to allow easier node configuration.
* (cli) [\#5116](https://github.com/cosmos/cosmos-sdk/issues/5116) The `CLIContext` now supports multiple verifiers
when connecting to multiple chains. The connecting chain's `CLIContext` will have to have the correct
chain ID and node URI or client set. To use a `CLIContext` with a verifier for another chain:
  ```go
  // main or parent chain (chain as if you're running without IBC)
  mainCtx := context.NewCLIContext()

  // connecting IBC chain
  sideCtx := context.NewCLIContext().
    WithChainID(sideChainID).
    WithNodeURI(sideChainNodeURI) // or .WithClient(...)

  sideCtx = sideCtx.WithVerifier(
    context.CreateVerifier(sideCtx, context.DefaultVerifierCacheSize),
  )
  ```
* (modules) [\#5017](https://github.com/cosmos/cosmos-sdk/pull/5017) The `x/auth` package now supports
generalized genesis accounts through the `GenesisAccount` interface.
* (modules) [\#4762](https://github.com/cosmos/cosmos-sdk/issues/4762) Deprecate remove and add permissions in ModuleAccount.
* (modules) [\#4760](https://github.com/cosmos/cosmos-sdk/issues/4760) update `x/auth` to match module spec.
* (modules) [\#4814](https://github.com/cosmos/cosmos-sdk/issues/4814) Add security contact to Validator description.
* (modules) [\#4875](https://github.com/cosmos/cosmos-sdk/issues/4875) refactor integration tests to use SimApp and separate test package
* (sdk) [\#4566](https://github.com/cosmos/cosmos-sdk/issues/4566) Export simulation's parameters and app state to JSON in order to reproduce bugs and invariants.
* (sdk) [\#4640](https://github.com/cosmos/cosmos-sdk/issues/4640) improve import/export simulation errors by extending `DiffKVStores` to return an array of `KVPairs` that are then compared to check for inconsistencies.
* (sdk) [\#4717](https://github.com/cosmos/cosmos-sdk/issues/4717) refactor `x/slashing` to match the new module spec
* (sdk) [\#4758](https://github.com/cosmos/cosmos-sdk/issues/4758) update `x/genaccounts` to match module spec
* (simulation) [\#4824](https://github.com/cosmos/cosmos-sdk/issues/4824) `PrintAllInvariants` flag will print all failed invariants
* (simulation) [\#4490](https://github.com/cosmos/cosmos-sdk/issues/4490) add `InitialBlockHeight` flag to resume a simulation from a given block

  * Support exporting the simulation stats to a given JSON file
* (simulation) [\#4847](https://github.com/cosmos/cosmos-sdk/issues/4847), [\#4838](https://github.com/cosmos/cosmos-sdk/pull/4838) and [\#4869](https://github.com/cosmos/cosmos-sdk/pull/4869) `SimApp` and simulation refactors:
  * Implement `SimulationManager` for executing modules' simulation functionalities in a modularized way
  * Add `RegisterStoreDecoders` to the `SimulationManager` for decoding each module's types
  * Add `GenerateGenesisStates` to the `SimulationManager` to generate a randomized `GenState` for each module
  * Add `RandomizedParams` to the `SimulationManager` that registers each modules' parameters in order to
  simulate `ParamChangeProposal`s' `Content`s
  * Add `WeightedOperations` to the `SimulationManager` that define simulation operations (modules' `Msg`s) with their
  respective weights (i.e chance of being simulated).
  * Add `ProposalContents` to the `SimulationManager` to register each module's governance proposal `Content`s.
* (simulation) [\#4893](https://github.com/cosmos/cosmos-sdk/issues/4893) Change `SimApp` keepers to be public and add getter functions for keys and codec
* (simulation) [\#4906](https://github.com/cosmos/cosmos-sdk/issues/4906) Add simulation `Config` struct that wraps simulation flags
* (simulation) [\#4935](https://github.com/cosmos/cosmos-sdk/issues/4935) Update simulation to reflect a proper `ABCI` application without bypassing `BaseApp` semantics
* (simulation) [\#5378](https://github.com/cosmos/cosmos-sdk/pull/5378) Simulation tests refactor:
  * Add `App` interface for general SDK-based app's methods.
  * Refactor and cleanup simulation tests into util functions to simplify their implementation for other SDK apps.
* (store) [\#4792](https://github.com/cosmos/cosmos-sdk/issues/4792) panic on non-registered store
* (types) [\#4821](https://github.com/cosmos/cosmos-sdk/issues/4821) types/errors package added with support for stacktraces. It is meant as a more feature-rich replacement for sdk.Errors in the mid-term.
* (store) [\#1947](https://github.com/cosmos/cosmos-sdk/issues/1947) Implement inter-block (persistent)
caching through `CommitKVStoreCacheManager`. Any application wishing to utilize an inter-block cache
must set it in their app via a `BaseApp` option. The `BaseApp` docs have been drastically improved
to detail this new feature and how state transitions occur.
* (docs/spec) All module specs moved into their respective module dir in x/ (i.e. docs/spec/staking -->> x/staking/spec)
* (docs/) [\#5379](https://github.com/cosmos/cosmos-sdk/pull/5379) Major documentation refactor, including:
  * (docs/intro/) Add and improve introduction material for newcomers.
  * (docs/basics/) Add documentation about basic concepts of the cosmos sdk such as the anatomy of an SDK application, the transaction lifecycle or accounts.
  * (docs/core/) Add documentation about core conepts of the cosmos sdk such as `baseapp`, `server`, `store`s, `context` and more.
  * (docs/building-modules/) Add reference documentation on concepts relevant for module developers (`keeper`, `handler`, `messages`, `queries`,...).
  * (docs/interfaces/) Add documentation on building interfaces for the Cosmos SDK.
  * Redesigned user interface that features new dynamically generated sidebar, build-time code embedding from GitHub, new homepage as well as many other improvements.
* (types) [\#5428](https://github.com/cosmos/cosmos-sdk/pull/5428) Add `Mod` (modulo) method and `RelativePow` (exponentation) function for `Uint`.
* (modules) [\#5506](https://github.com/cosmos/cosmos-sdk/pull/5506) Remove redundancy in `x/distribution`s use of parameters. There
  now exists a single `Params` type with a getter and setter along with a getter for each individual parameter.

### Bug Fixes

* (client) [\#5303](https://github.com/cosmos/cosmos-sdk/issues/5303) Fix ignored error in tx generate only mode.
* (cli) [\#4763](https://github.com/cosmos/cosmos-sdk/issues/4763) Fix flag `--min-self-delegation` for staking `EditValidator`
* (keys) Fix ledger custom coin type support bug.
* (x/gov) [\#5107](https://github.com/cosmos/cosmos-sdk/pull/5107) Sum validator operator's all voting power when tally votes
* (rest) [\#5212](https://github.com/cosmos/cosmos-sdk/issues/5212) Fix pagination in the `/gov/proposals` handler.

## [v0.37.14] - 2020-08-12

### Improvements

* (tendermint) Bump Tendermint version to [v0.32.13](https://github.com/tendermint/tendermint/releases/tag/v0.32.13).


## [v0.37.13] - 2020-06-03

### Improvements

* (tendermint) Bump Tendermint version to [v0.32.12](https://github.com/tendermint/tendermint/releases/tag/v0.32.12).
* (cosmos-ledger-go) Bump Cosmos Ledger Wallet library version to [v0.11.1](https://github.com/cosmos/ledger-cosmos-go/releases/tag/v0.11.1).

## [v0.37.12] - 2020-05-05

### Improvements

* (tendermint) Bump Tendermint version to [v0.32.11](https://github.com/tendermint/tendermint/releases/tag/v0.32.11).

## [v0.37.11] - 2020-04-22

### Bug Fixes

* (x/staking) [\#6021](https://github.com/cosmos/cosmos-sdk/pull/6021) --trust-node's false default value prevents creation of the genesis transaction.

## [v0.37.10] - 2020-04-22

### Bug Fixes

* (client/context) [\#5964](https://github.com/cosmos/cosmos-sdk/issues/5964) Fix incorrect instantiation of tmlite verifier when --trust-node is off.

## [v0.37.9] - 2020-04-09

### Improvements

* (tendermint) Bump Tendermint version to [v0.32.10](https://github.com/tendermint/tendermint/releases/tag/v0.32.10).

## [v0.37.8] - 2020-03-11

### Bug Fixes

* (rest) [\#5508](https://github.com/cosmos/cosmos-sdk/pull/5508) Fix `x/distribution` endpoints to properly return height in the response.
* (x/genutil) [\#5499](https://github.com/cosmos/cosmos-sdk/pull/) Ensure `DefaultGenesis` returns valid and non-nil default genesis state.
* (x/genutil) [\#5775](https://github.com/cosmos/cosmos-sdk/pull/5775) Fix `ExportGenesis` in `x/genutil` to export default genesis state (`[]`) instead of `null`.
* (genesis) [\#5086](https://github.com/cosmos/cosmos-sdk/issues/5086) Ensure `gentxs` are always an empty array instead of `nil`.

### Improvements

* (rest) [\#5648](https://github.com/cosmos/cosmos-sdk/pull/5648) Enhance /txs usability:
  * Add `tx.minheight` key to filter transaction with an inclusive minimum block height
  * Add `tx.maxheight` key to filter transaction with an inclusive maximum block height

## [v0.37.7] - 2020-02-10

### Improvements

* (modules) [\#5597](https://github.com/cosmos/cosmos-sdk/pull/5597) Add `amount` event attribute to the `complete_unbonding`
and `complete_redelegation` events that reflect the total balances of the completed unbondings and redelegations
respectively.

### Bug Fixes

* (x/gov) [\#5622](https://github.com/cosmos/cosmos-sdk/pull/5622) Track any events emitted from a proposal's handler upon successful execution.
* (x/bank) [\#5531](https://github.com/cosmos/cosmos-sdk/issues/5531) Added missing amount event to MsgMultiSend, emitted for each output.

## [v0.37.6] - 2020-01-21

### Improvements

* (tendermint) Bump Tendermint version to [v0.32.9](https://github.com/tendermint/tendermint/releases/tag/v0.32.9)

## [v0.37.5] - 2020-01-07

### Features

* (types) [\#5360](https://github.com/cosmos/cosmos-sdk/pull/5360) Implement `SortableDecBytes` which
  allows the `Dec` type be sortable.

### Improvements

* (tendermint) Bump Tendermint version to [v0.32.8](https://github.com/tendermint/tendermint/releases/tag/v0.32.8)
* (cli) [\#5482](https://github.com/cosmos/cosmos-sdk/pull/5482) Remove old "tags" nomenclature from the `q txs` command in
  favor of the new events system. Functionality remains unchanged except that `=` is used instead of `:` to be
  consistent with the API's use of event queries.

### Bug Fixes

* (iavl) [\#5276](https://github.com/cosmos/cosmos-sdk/issues/5276) Fix potential race condition in `iavlIterator#Close`.
* (baseapp) [\#5350](https://github.com/cosmos/cosmos-sdk/issues/5350) Allow a node to restart successfully
  after a `halt-height` or `halt-time` has been triggered.
* (types) [\#5395](https://github.com/cosmos/cosmos-sdk/issues/5395) Fix `Uint#LTE`.
* (types) [\#5408](https://github.com/cosmos/cosmos-sdk/issues/5408) `NewDecCoins` constructor now sorts the coins.

## [v0.37.4] - 2019-11-04

### Improvements

* (tendermint) Bump Tendermint version to [v0.32.7](https://github.com/tendermint/tendermint/releases/tag/v0.32.7)
* (ledger) [\#4716](https://github.com/cosmos/cosmos-sdk/pull/4716) Fix ledger custom coin type support bug.

### Bug Fixes

* (baseapp) [\#5200](https://github.com/cosmos/cosmos-sdk/issues/5200) Remove duplicate events from previous messages.

## [v0.37.3] - 2019-10-10

### Bug Fixes

* (genesis) [\#5095](https://github.com/cosmos/cosmos-sdk/issues/5095) Fix genesis file migration from v0.34 to
v0.36/v0.37 not converting validator consensus pubkey to bech32 format.

### Improvements

* (tendermint) Bump Tendermint version to [v0.32.6](https://github.com/tendermint/tendermint/releases/tag/v0.32.6)

## [v0.37.1] - 2019-09-19

### Features

### Improvements

### Bug Fixes

### Breaking Changes

### Build, CI

### Document Updates
