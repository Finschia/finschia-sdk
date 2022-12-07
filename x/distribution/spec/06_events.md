<!--
order: 6
-->

# Events

The distribution module emits the following events:

## BeginBlocker

| Type            | Attribute Key | Attribute Value    |
|-----------------|---------------|--------------------|
| proposer_reward | amount        | {proposerReward}   |
| proposer_reward | validator     | {validatorAddress} |
| commission      | amount        | {commissionAmount} |
| commission      | validator     | {validatorAddress} |
| rewards         | amount        | {rewardAmount}     |
| rewards         | validator     | {validatorAddress} |

## Handlers

### MsgSetWithdrawAddress

| Type                 | Attribute Key    | Attribute Value   |
|----------------------|------------------|-------------------|
| set_withdraw_address | withdraw_address | {withdrawAddress} |
| message              | module           | distribution      |
| message              | sender           | {senderAddress}   |

### MsgWithdrawDelegatorReward

| Type             | Attribute Key | Attribute Value    |
|------------------|---------------|--------------------|
| withdraw_rewards | amount        | {rewardAmount}     |
| withdraw_rewards | validator     | {validatorAddress} |
| message          | module        | distribution       |
| message          | sender        | {senderAddress}    |

### MsgWithdrawValidatorCommission

| Type                | Attribute Key | Attribute Value    |
|---------------------|---------------|--------------------|
| withdraw_commission | amount        | {commissionAmount} |
| message             | module        | distribution       |
| message             | sender        | {senderAddress}    |
