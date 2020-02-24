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
| grant_perm       | perm_resource  | {symbol}                 |
| grant_perm       | perm_action    | mint                     |
| grant_perm       | to             | {ownerAddress}           |
| grant_perm       | perm_resource  | {symbol}                 |
| grant_perm       | perm_action    | modify                   |
| issue            | name           | {name}                   |
| issue            | symbol         | {symbol}                 |
| issue            | denom          | {symbol}                 |
| issue            | owner          | {ownerAddress}           |
| issue            | amount         | {amount}                 |
| issue            | mintable       | {mintable}               |
| issue            | decimals       | {decimals}               |

### MsgMint
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | mint                     |
| mint             | symbol         | {symbol}                 |
| mint             | amount         | {amount}                 |
| mint             | from           | {ownerAddress}           |
| mint             | to             | {toAddress}              |

### MsgBurn
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | burn                     |
| burn             | symbol         | {symbol}                 |
| burn             | amount         | {amount}                 |
| burn             | from           | {ownerAddress}           |

### MsgTransfer
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | transfer_ft              |
| transfer_ft      | from           | {fromAddress}            |
| transfer_ft      | to             | {toAddress}              |
| transfer_ft      | symbol         | {symbol}                 |
| transfer_ft      | amount         | {amount}                 |
