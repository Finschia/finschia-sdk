<!--
order: 2
-->

# State

## Params

As of now, the only parameter is on/off of the module. After turning off the module, the changes (might be breaking changes) applied on the other modules would not work, and the module cleans up its state.

- Params: `0x00 -> PropocolBuffer(Params)`

+++ https://github.com/line/lbm-sdk/blob/v0.44.0-rc0/proto/lbm/foundation/v1/foundation.proto#L9-L12
```go
// Params defines the parameters for the foundation module.
message Params {
  bool enabled = 1 [(gogoproto.moretags) = "yaml:\"enabled\""];
}
```

## ValidatorAuth

An operator must have been authorized before creating its validator node. `ValidatorAuth`s contain the authorization on the operators. One can authorize itself or other operator by the corresponding proposal, `lbm/foundation/v1/UpdateValidatorAuthsProposal`.

Note that if the chain starts with an empty list of it in the genesis, the module authorizes all the operators included in the list of validators, read from the staking module.

- ValidatorAuth: `0x01 -> ProtocolBuffer(ValidatorAuth)`

+++ https://github.com/line/lbm-sdk/blob/v0.44.0-rc0/proto/lbm/foundation/v1/foundation.proto#L14-L18
```go
// ValidatorAuth defines authorization info of a validator.
message ValidatorAuth {
  string operator_address = 1 [(gogoproto.moretags) = "yaml:\"operator_address\""];
  bool   creation_allowed = 2 [(gogoproto.moretags) = "yaml:\"creation_allowed\""];
}
```
