# Governance

This document gives an overview of how the various governance
proposals interact with the CosmWasm contract lifecycle. It is
a high-level, technical introduction meant to provide context before
looking into the code, or constructing proposals.

## Proposal Types
We have added 5 new wasm specific proposal types that cover the contract's live cycle and authorization:

* `StoreCodeProposal` - upload a wasm binary
* `InstantiateContractProposal` - instantiate a wasm contract
* `MigrateContractProposal` - migrate a wasm contract to a new code version
* `UpdateAdminProposal` - set a new admin for a contract
* `ClearAdminProposal` - clear admin for a contract to prevent further migrations

For details, see the proposal type [implementation](internal/types/proposal.go)

A wasm message but no proposal type:
* `ExecuteContract` - execute a command on a wasm contract

And you can use `Parameter Change Proposal` to change wasm parameters.
These parameters are as following.

* `UploadAccess` - who can upload wasm codes
* `DefaultInstantiatePermission` - who can instantiate contracts from a code in default
* `MaxWasmCodeSize` - max size of wasm code to be uploaded

### Unit tests
[Proposal type validations](internal/types/proposal_test.go)

## Proposal Handler
The [wasm proposal_handler](internal/keeper/proposal_handler.go) implements the `gov.Handler` function
and executes the wasm proposal types after a successful tally.

The proposal handler uses a [`GovAuthorizationPolicy`](internal/keeper/authz_policy.go#L29) to bypass the existing contract's authorization policy.

### Tests
* [Integration: Submit and execute proposal](internal/keeper/proposal_integration_test.go)

## Gov Integration
The wasm proposal handler can be added to the gov router in the [abci app](linkwasmd/app/app.go#L240)
to receive proposal execution calls.
```go
govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.wasmKeeper, wasm.EnableAllProposals))
```

## Wasm Authorization Settings

Settings via sdk `params` module:
- `code_upload_access` - who can upload a wasm binary: `Nobody`, `Everybody`, `OnlyAddress`
- `instantiate_default_permission` - platform default, who can instantiate a wasm binary when the code owner has not set it

See [params.go](internal/types/params.go)

### Init Params Via Genesis

```json
    "wasm": {
      "params": {
        "code_upload_access": {
          "permission": "Everybody"
        },
        "instantiate_default_permission": "Everybody"
      }
    },
```

The values can be updated via gov proposal implemented in the `params` module.

### Enable gov proposals
Gov proposals authorization policy needs to be specified with `enabledProposalTypes` which is an argument of NewWasmProposalHandler in [proposal_handler.go](internal/keeper/proposal_handler.go)

### Tests
* [params validation unit tests](internal/types/params_test.go)
* [genesis validation tests](internal/types/genesis_test.go)
* [policy integration tests](internal/keeper/keeper_test.go)

## CLI

```shell script
  wasmcli tx gov submit-proposal [command]

Available Commands:
  wasm-store           Submit a wasm binary proposal
  instantiate-contract Submit an instantiate wasm contract proposal
  migrate-contract     Submit a migrate wasm contract to a new code version proposal
  set-contract-admin   Submit a new admin for a contract proposal
  clear-contract-admin Submit a clear admin for a contract to prevent further migrations proposal
...
```
## Rest
New [`ProposalHandlers`](client/proposal_handler.go)

* Integration
```shell script
gov.NewAppModuleBasic(append(wasmclient.ProposalHandlers, paramsclient.ProposalHandler,)...),
```
In [abci app](linkwasmd/app/app.go)

### Tests
* [Rest Unit tests](client/proposal_handler_test.go)
* [CLI tests](linkwasmd/cli_test/cli_test.go)
