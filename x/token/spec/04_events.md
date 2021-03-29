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

### MsgBurnFrom
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {proxyAddress}           | 
| message          | action         | burn                     |
| burn_from        | contract_id    | {contractID}             |
| burn_from        | proxy          | {proxyAddress}           |
| burn_from        | from           | {fromAddress}            |
| burn_from        | amount         | {amount}                 |

### MsgBurnFrom
| Type             | Attribute Key  | Attribute Value              |
|------------------|----------------|------------------------------|
| message          | module         | token                        |
| message          | sender         | {proxyAddress}               | 
| message          | action         | burn_ft                      |
| burn_ft_from     | contract_id    | {contractID}                 |
| burn_ft_from     | proxy          | {proxyAddress}               |
| burn_ft_from     | from           | {fromAddress}                |
| burn_ft_from     | amount         | {amount}{contractID}{tokenID}|


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
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | grant_permission         |
| grant_perm       | from           | {fromAddress}            |
| grant_perm       | to             | {toAddress}              |
| grant_perm       | contract_id    | {resource}               |
| grant_perm       | perm           | issue/mint/burn/modify   |

### MsgRevokePermission
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | revoke_permission        |
| revoke_perm      | from           | {fromAddress}            |
| revoke_perm      | contract_id    | {resource}               |
| revoke_perm      | perm           | issue/mint/burn/modify   |

### MsgApprove
| Type          | Attribute Key  | Attribute Value   |
|---------------|----------------|-------------------|
| message       | module         | token             |
| message       | sender         | {approverAddress} | 
| message       | action         | approve_token     |
| approve_token | contract_id    | {contractID}      |
| approve_token | approver       | {approverAddress} |
| approve_token | proxy          | {proxyAddress}    |

