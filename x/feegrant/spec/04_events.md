<!--
order: 4
-->

# Events

The feegrant module emits the following events:

# Msg Server

### MsgGrantAllowance

| Type         | Attribute Key | Attribute Value  |
|--------------|---------------|------------------|
| set_feegrant | granter       | {granterAddress} |
| set_feegrant | grantee       | {granteeAddress} |

### MsgRevokeAllowance

| Type            | Attribute Key | Attribute Value  |
|-----------------|---------------|------------------|
| revoke_feegrant | granter       | {granterAddress} |
| revoke_feegrant | grantee       | {granteeAddress} |

## Exec fee allowance

| Type         | Attribute Key | Attribute Value  |
|--------------|---------------|------------------|
| use_feegrant | granter       | {granterAddress} |
| use_feegrant | grantee       | {granteeAddress} |
