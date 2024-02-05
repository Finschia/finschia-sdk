---
sidebar_position: 1
---

# `x/stakingplus`

## Abstract

This paper specifies the Staking Plus module of the Finschia-sdk, which extends existing [Staking module](https://github.com/cosmos/cosmos-sdk/blob/v0.50.2/x/staking/README.md) of the Cosmos-SDK.

The module enables Finschia-sdk based blockchain to support an advanced Proof-of-Stake system. In this system, holders of the native staking token of the chain can become validators and can delegate tokens to validators, ultimately determining the effective validator set for the system.

This module is almost identical to the previous Staking module of the Cosmos-SDK, but introduces some breaking changes. For example, you must have x/foundation UpdateValidatorAuthsProposal passed before sending x/stakingplus MsgCreateValidator, or the message would fail.

In this document, we describe only the changes introduced by Finschia-SDK. Refer to the [original document](https://github.com/cosmos/cosmos-sdk/blob/v0.50.2/x/staking/README.md) for more information.

# Messages

## Msg/CreateValidator

A validator is created using the `Msg/CreateValidator` service message.

+++ https://github.com/cosmos/cosmos-sdk/blob/v0.50.2/proto/cosmos/staking/v1beta1/tx.proto#L20-L21

+++ https://github.com/cosmos/cosmos-sdk/blob/v0.50.2/proto/cosmos/staking/v1beta1/tx.proto#L50-L73

This service message is expected to fail if:

- one of the conditions described in the staking module of the Cosmos-SDK is met.
- the operator address is not registered on x/foundation through [MsgGrant](https://github.com/Finschia/finschia-sdk/tree/main/x/foundation#msggrant) with `CreateValidatorAuthorization`. 

The other [statements](https://github.com/cosmos/cosmos-sdk/blob/v0.50.2/x/staking/README.md#msgcreatevalidator) on this message in the exising document are still valid.
