<!--
order: 3
-->

# Messages

In this section we describe only the changes introduced by LBM-SDK. Refer to the [original document](../../staking/spec/03_messages.md) for more information.

## Msg/CreateValidator

A validator is created using the `Msg/CreateValidator` service message.

+++ https://github.com/line/lbm-sdk/blob/main/proto/cosmos/staking/v1beta1/tx.proto#L16-L17

+++ https://github.com/line/lbm-sdk/blob/main/proto/cosmos/staking/v1beta1/tx.proto#L35-L51

This service message is expected to fail if:

- one of the conditions described in the staking module of the Cosmos-SDK is met.
- the operator address is not registered on x/foundation through UpdateValidatorAuthsProposal. TODO: add a ref to x/foundation spec file.

The other [statements](../../staking/spec/03_messages.md#msgcreatevalidator) on this message in the exising document are still valid.
