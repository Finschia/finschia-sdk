<!--
order: 2
-->

# State

## Params

As of now, the only parameter is on/off of the module. After turning off the module, the changes (might be breaking changes) applied on the other modules would not work, and the module cleans up its state.

- Params: `0x00 -> PropocolBuffer(Params)`

+++ https://github.com/line/lbm-sdk/blob/v0.44.0-rc0/proto/lbm/consortium/v1/consortium.proto#L9-L12

## ValidatorAuth

An operator must have been authorized before creating its validator node. `ValidatorAuth`s contain the authorization on the operators. One can authorize itself or other operator by the corresponding proposal, `lbm/consortium/v1/UpdateValidatorAuthsProposal`.

Note that if the chain starts with an empty list of it in the genesis, the module authorizes all the operators included in the list of validators, read from the staking module.

- ValidatorAuth: `0x01 -> ProtocolBuffer(ValidatorAuth)`

+++ https://github.com/line/lbm-sdk/blob/v0.44.0-rc0/proto/lbm/consortium/v1/consortium.proto#L14-L18
