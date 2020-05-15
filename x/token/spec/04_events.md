# Events
**Not fully documented yet** 
The token module emits the following events:


### MsgIssue
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | issue_token              |
| grant_perm       | to             | {ownerAddress}           |
| grant_perm       | contract_id    | {contractID}             |
| grant_perm       | perm           | mint                     |
| grant_perm       | to             | {ownerAddress}           |
| grant_perm       | contract_id    | {contractID}             |
| grant_perm       | perm           | modify                   |
| issue            | contract_id    | {contractID}             |
| issue            | name           | {name}                   |
| issue            | symbol         | {symbol}                 |
| issue            | owner          | {ownerAddress}           |
| issue            | to             | {toAddress}              |
| issue            | amount         | {amount}                 |
| issue            | mintable       | {mintable}               |
| issue            | decimals       | {decimals}               |

### MsgMint
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | mint                     |
| mint             | contract_id    | {contractID}             |
| mint             | amount         | {amount}                 |
| mint             | from           | {ownerAddress}           |
| mint             | to             | {toAddress}              |

### MsgBurn
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | burn                     |
| burn             | contract_id    | {contractID}             |
| burn             | amount         | {amount}                 |
| burn             | from           | {ownerAddress}           |

### MsgTransfer
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | transfer_ft              |
| transfer_ft      | contract_id    | {contractID}             |
| transfer_ft      | from           | {fromAddress}            |
| transfer_ft      | to             | {toAddress}              |
| transfer_ft      | amount         | {amount}                 |

### MsgModify
| Type                  | Attribute Key  | Attribute Value       |
|-----------------------|----------------|-----------------------|
| message               | module         | token                 |
| message               | sender         | {ownerAddress}        | 
| message               | action         | modify_token          |
| modify_token          | contract_id    | {contract_id}         |
| modify_token          | {modifiedField}| {modifiedValue}       |

### MsgGrantPermission
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | grant_permission         |
| grant_perm       | from           | {fromAddress}            |
| grant_perm       | to             | {toAddress}              |
| grant_perm       | contract_id    | {resource}               |
| grant_perm       | perm           | issue/mint/burn/modify   |

### MsgRevokePermission
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | revoke_permission        |
| revoke_perm      | from           | {fromAddress}            |
| revoke_perm      | contract_id    | {resource}               |
| revoke_perm      | perm           | issue/mint/burn/modify   |