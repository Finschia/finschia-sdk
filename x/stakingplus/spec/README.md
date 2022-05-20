<!--
order: 0
title: Staking Plus Overview
parent:
  title: "staking plus"
-->

# `stakingplus`

## Abstract

This paper specifies the Staking Plus module of the LBM-SDK, which extends existing [Staking module](../../staking/spec/README.md) of the Cosmos-SDK.

The module enables LBM-SDK based blockchain to support an advanced Proof-of-Stake system. In this system, holders of the native staking token of the chain can become validators and can delegate tokens to validators, ultimately determining the effective validator set for the system.

This module is almost identical to the previous Staking module of the Cosmos-SDK, but introduces some breaking changes. For example, you must have x/foundation UpdateValidatorAuthsProposal passed before sending x/stakingplus MsgCreateValidator, or the message would fail.

## Contents

1. **[State](01_state.md)**
2. **[State Transitions](02_state_transitions.md)**
3. **[Messages](03_messages.md)**
    - [Msg/CreateValidator](03_messages.md#msgcreatevalidator)
4. **[Begin-Block](04_begin_block.md)**
5. **[End-Block ](05_end_block.md)**
6. **[Hooks](06_hooks.md)**
7. **[Events](07_events.md)**
8. **[Parameters](08_params.md)**
