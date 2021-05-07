<!--
order: 4
-->

# Events

The feegrant module emits the following events:

# Msg Server

### MsgGrantAllowance

| Type    | Attribute Key | Attribute Value  |
| ------- | ------------- | ---------------- |
| message | action        | set_feegrant     |
| message | granter       | {granterAddress} |
| message | grantee       | {granteeAddress} |

### MsgRevokeAllowance
| ------- | ------------- | ---------------- |
| message | granter       | {granterAddress} |
| message | grantee       | {granteeAddress} |

## Exec fee allowance

| Type    | Attribute Key | Attribute Value  |
| ------- | ------------- | ---------------- |
| message | action        | use_feegrant     |
| message | granter       | {granterAddress} |
| message | grantee       | {granteeAddress} |
