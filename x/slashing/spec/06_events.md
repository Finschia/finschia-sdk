<!--
order: 6
-->

# Tags

The slashing module emits the following events/tags:

## MsgServer

### MsgUnjail

| Type    | Attribute Key | Attribute Value    |
|---------|---------------|--------------------|
| message | module        | slashing           |
| message | sender        | {validatorAddress} |

## Keeper

## BeginBlocker: HandleValidatorSignature

| Type     | Attribute Key | Attribute Value             |
|----------|---------------|-----------------------------|
| liveness | address       | {validatorConsensusAddress} |
| liveness | missed_blocks | {missedBlocksCounter}       |
| liveness | height        | {blockHeight}               |
| slash    | address       | {validatorConsensusAddress} |
| slash    | power         | {validatorPower}            |
| slash    | reason        | {slashReason}               |
| slash    | jailed [0]    | {validatorConsensusAddress} |

- [0] Only included if the validator is jailed.

### Slash

| Type  | Attribute Key | Attribute Value             |
|-------|---------------|-----------------------------|
| slash | address       | {validatorConsensusAddress} |
| slash | power         | {validatorPower}            |
| slash | reason        | {slashReason}               |

### Jail

| Type  | Attribute Key | Attribute Value    |
|-------|---------------|--------------------|
| slash | jailed        | {validatorAddress} |
