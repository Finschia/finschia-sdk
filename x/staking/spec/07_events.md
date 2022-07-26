<!--
order: 6
-->

# Events

The staking module emits the following events:

## EndBlocker

| Type                  | Attribute Key         | Attribute Value           |
|-----------------------|-----------------------|---------------------------|
| complete_unbonding    | amount                | {totalUnbondingAmount}    |
| complete_unbonding    | validator             | {validatorAddress}        |
| complete_unbonding    | delegator             | {delegatorAddress}        |
| complete_redelegation | amount                | {totalRedelegationAmount} |
| complete_redelegation | delegator             | {delegatorAddress}        |
| complete_redelegation | source_validator      | {srcValidatorAddress}     |
| complete_redelegation | destination_validator | {dstValidatorAddress}     |

## Msg's

### MsgCreateValidator

| Type             | Attribute Key | Attribute Value    |
|------------------|---------------|--------------------|
| create_validator | validator     | {validatorAddress} |
| create_validator | amount        | {delegationAmount} |
| message          | module        | staking            |
| message          | sender        | {senderAddress}    |

### MsgEditValidator

| Type           | Attribute Key       | Attribute Value     |
|----------------|---------------------|---------------------|
| edit_validator | commission_rate     | {commissionRate}    |
| edit_validator | min_self_delegation | {minSelfDelegation} |
| message        | module              | staking             |
| message        | sender              | {senderAddress}     |

### MsgDelegate

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| delegate | validator     | {validatorAddress} |
| delegate | amount        | {delegationAmount} |
| delegate | new_shares    | {newShares}        |
| message  | module        | staking            |
| message  | sender        | {senderAddress}    |

### MsgUndelegate

| Type    | Attribute Key       | Attribute Value    |
|---------|---------------------|--------------------|
| unbond  | validator           | {validatorAddress} |
| unbond  | amount              | {unbondAmount}     |
| unbond  | completion_time [0] | {completionTime}   |
| message | module              | staking            |
| message | sender              | {senderAddress}    |

- [0] Time is formatted in the RFC3339 standard

### MsgBeginRedelegate

| Type       | Attribute Key         | Attribute Value       |
|------------|-----------------------|-----------------------|
| redelegate | source_validator      | {srcValidatorAddress} |
| redelegate | destination_validator | {dstValidatorAddress} |
| redelegate | amount                | {unbondAmount}        |
| redelegate | completion_time [0]   | {completionTime}      |
| message    | module                | staking               |
| message    | sender                | {senderAddress}       |

- [0] Time is formatted in the RFC3339 standard
