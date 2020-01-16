# Events
**Not fully documented yet** 
The token module emits the following events:


### MsgIssue          
| Type                           | Attribute Key  | Attribute Value          |
|--------------------------------|----------------|--------------------------|
| message                        | module         | token                    |
| message                        | sender         | {ownerAddress}           | 
| message                        | action         | issue_token              |
| grant_perm_token               | to             | {ownerAddress}           |
| grant_perm_token               | perm_resource  | {symbol}                 |
| grant_perm_token               | perm_action    | issue                    |
| grant_perm_token               | perm_action    | mint                     |
| modify_token_uri_perm_token    | to             | {ownerAddress}           |
| modify_token_uri_perm_token    | perm_resource  | {symbol}                 |
| modify_token_uri_perm_token    | perm_action    | modify_token_uri         |
| issue_token                    | name           | name                     |
| issue_token                    | symbol         | symbol                   |
| issue_token                    | owner          | {ownerAddress}           |
| issue_token                    | amount         | {amount}                 |
| issue_token                    | mintable       | {mintable}               |
| issue_token                    | decimals       | {decimals}               |
| issue_token                    | token_uri      | {token_uri}              |
| issue_token                    | token_type     | ft                       |
| mint_token                     | amount         | {amount}{symbol}         |
| mint_token                     | to             | {ownerAddress}           |
| occupy_symbol                  | symbol         | {symbol}                 |
| occupy_symbol                  | owner          | {ownerAddress}           |
                     
### MsgIssueCollection                                                   
| Type                           | Attribute Key  | Attribute Value          |
|--------------------------------|----------------|--------------------------|
| message                        | module         | token                    |
| message                        | sender         | {ownerAddress}           | 
| message                        | action         | issue_token_collection   |
| grant_perm_token               | to             | {ownerAddress}           |
| grant_perm_token               | perm_resource  | {symbol}                 |
| grant_perm_token               | perm_action    | issue                    |
| grant_perm_token               | perm_action    | mint                     |
| modify_token_uri_perm_token    | to             | {ownerAddress}           |
| modify_token_uri_perm_token    | perm_resource  | {symbol}                 |
| modify_token_uri_perm_token    | perm_action    | modify_token_uri         |
| issue_token                    | name           | name                     |
| issue_token                    | symbol         | {symbol}{tokenid}        |
| issue_token                    | owner          | {ownerAddress}           |
| issue_token                    | amount         | {amount}                 |
| issue_token                    | mintable       | {mintable}               |
| issue_token                    | decimals       | {decimals}               |
| issue_token                    | token_uri      | {token_uri}              |
| issue_token                    | token_type     | cft                      |
| mint_token                     | amount         | {amount}{symbol}{tokenid}|
| mint_token                     | to             | {ownerAddress}           |
| occupy_symbol                  | symbol         | {symbol}                 |
| occupy_symbol                  | owner          | {ownerAddress}           |
                                
### MsgIssueNFT
| Type                           | Attribute Key  | Attribute Value          |
|--------------------------------|----------------|--------------------------|
| message                        | module         | token                    |
| message                        | sender         | {ownerAddress}           | 
| message                        | action         | issue_nft                |
| grant_perm                     | to             | {ownerAddress}           |
| grant_perm                     | perm_resource  | {symbol}                 |
| grant_perm                     | perm_action    | issue                    |
| modify_token_uri_perm_token    | to             | {ownerAddress}           |
| modify_token_uri_perm_token    | perm_resource  | {symbol}                 |
| modify_token_uri_perm_token    | perm_action    | modify_token_uri         |
| issue_token                    | name           | name                     |
| issue_token                    | symbol         | symbol                   |
| issue_token                    | owner          | {ownerAddress}           |
| issue_token                    | amount         | 1                        |
| issue_token                    | mintable       | false                    |
| issue_token                    | decimals       | 0                        |
| issue_token                    | token_uri      | {token_uri}              |
| issue_token                    | token_type     | nft                      |
| mint_token                     | amount         | {amount}{symbol}         |
| mint_token                     | to             | {ownerAddress}           |
| occupy_symbol                  | symbol         | {symbol}                 |
| occupy_symbol                  | owner          | {ownerAddress}           |

### MsgIssueNFTCollection                                                   
| Type                           | Attribute Key  | Attribute Value          |
|--------------------------------|----------------|--------------------------|
| message                        | module         | token                    |
| message                        | sender         | {ownerAddress}           | 
| message                        | action         | issue_nft_collection     |
| grant_perm_token               | to             | {ownerAddress}           |         
| grant_perm_token               | perm_resource  | {symbol}                 |         
| grant_perm_token               | perm_action    | issue                    |         
| grant_perm_token               | perm_action    | mint                     |         
| modify_token_uri_perm_token    | to             | {ownerAddress}           |         
| modify_token_uri_perm_token    | perm_resource  | {symbol}                 |         
| modify_token_uri_perm_token    | perm_action    | modify_token_uri         |         
| issue_token                    | decimals       | 0                        |
| issue_token                    | token_uri      | {token_uri}              |
| issue_token                    | token_type     | cnft                     |
| mint_token                     | amount         | {amount}{symbol}{tokenid}|
| mint_token                     | to             | {ownerAddress}           |
| occupy_symbol                  | symbol         | {symbol}                 |
| occupy_symbol                  | owner          | {ownerAddress}           |

### MsgMint
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | mint_token               |
| mint_token       | amount         | {amount}{symbol}         |
| mint_token       | to             | {ownerAddress}           |

### MsgBurn
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | burn_token               |
| burn_token       | amount         | {amount}{symbol}         |
| burn_token       | from           | {ownerAddress}           |

### MsgGrantPermission
| Type                           | Attribute Key  | Attribute Value          |
|--------------------------------|----------------|--------------------------|
| message                        | module         | token                    |
| message                        | sender         | {fromAddress}            | 
| message                        | action         | grant_permission         |
| grant_perm_token               | from           | {fromAddress}            |
| grant_perm_token               | to             | {toAddress}              |
| grant_perm_token               | perm_resource  | {symbol}                 |
| grant_perm_token               | perm_action    | issue/burn               |
| modify_token_uri_perm_token    | to             | {ownerAddress}           |         
| modify_token_uri_perm_token    | perm_resource  | {symbol}                 |         
| modify_token_uri_perm_token    | perm_action    | modify_token_uri         |    

### MsgRevokePermission
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | revoke_permission        |
| revoke_perm      | from           | {fromAddress}            |
| revoke_perm      | perm_resource  | {symbol}                 |
| revoke_perm      | perm_action    | issue/burn               |

