<!--
order: 4
-->

# Events

The feegrant module emits the following events:

<<<<<<< HEAD
# Msg Server

### MsgGrantAllowance
=======
# Handlers

### MsgGrantFeeAllowance
>>>>>>> 0af248b95 (Add specs for feegrant (#8496))

| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |
| message  | action        | set_feegrant       |
| message  | granter       | {granterAddress}   |
| message  | grantee       | {granteeAddress}   |

<<<<<<< HEAD
### MsgRevokeAllowance
=======
### MsgRevokeFeeAllowance
>>>>>>> 0af248b95 (Add specs for feegrant (#8496))

| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |
| message  | action        | revoke_feegrant    |
| message  | granter       | {granterAddress}   |
| message  | grantee       | {granteeAddress}   |

### Exec fee allowance

| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |
| message  | action        | use_feegrant       |
| message  | granter       | {granterAddress}   |
<<<<<<< HEAD
| message  | grantee       | {granteeAddress}   |
=======
| message  | grantee       | {granteeAddress}   |
>>>>>>> 0af248b95 (Add specs for feegrant (#8496))
