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
| modify_token_uri_perm_token    | to             | {ownerAddress}           |
| modify_token_uri_perm_token    | perm_resource  | {symbol}                 |
| modify_token_uri_perm_token    | perm_action    | modify_token_uri         |
| issue_token      | name           | {name}                   |
| issue_token      | symbol         | {symbol}                 |
| issue_token      | denom          | {symbol}                 |
| issue_token      | owner          | {ownerAddress}           |
| issue_token      | amount         | {amount}                 |
| issue_token      | mintable       | {mintable}               |
| issue_token      | decimals       | {decimals}               |
| issue_token      | token_type     | ft                       |
| mint_token       | amount         | {amount}{symbol}         |
| mint_token       | to             | {ownerAddress}           |

### MsgIssueCollection                                                   
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | issue_token_collection   |
| grant_perm       | to             | {ownerAddress}           |
| grant_perm       | perm_resource  | {symbol}                 |
| grant_perm       | perm_action    | mint                     |
| modify_token_uri_perm_token    | to             | {ownerAddress}           |
| modify_token_uri_perm_token    | perm_resource  | {symbol}                 |
| modify_token_uri_perm_token    | perm_action    | modify_token_uri         |
| issue_token      | name           | {name}                   |
| issue_token      | symbol         | {symbol}                 |
| issue_token      | denom          | {symbol}{tokenid}        |
| issue_token      | owner          | {ownerAddress}           |
| issue_token      | amount         | {amount}                 |
| issue_token      | mintable       | {mintable}               |
| issue_token      | decimals       | {decimals}               |
| issue_token      | token_uri      | {token_uri}              |
| issue_token      | token_type     | cft                      |
| mint_token       | amount         | {amount}{symbol}{tokenid}|
| mint_token       | to             | {ownerAddress}           |

### MsgIssueNFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | issue_nft                |
| modify_token_uri_perm_token    | to             | {ownerAddress}           |
| modify_token_uri_perm_token    | perm_resource  | {symbol}                 |
| modify_token_uri_perm_token    | perm_action    | modify_token_uri         |
| issue_token      | name           | {name}                   |
| issue_token      | symbol         | {symbol}                 |
| issue_token      | denom          | {symbol}                 |
| issue_token      | owner          | {ownerAddress}           |
| issue_token      | amount         | 1                        |
| issue_token      | token_uri      | {token_uri}              |
| issue_token      | token_type     | nft                      |
| mint_token       | amount         | {amount}{symbol}         |
| mint_token       | to             | {ownerAddress}           |

### MsgIssueNFTCollection                                                   
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | issue_nft_collection     |
| grant_perm       | to             | {ownerAddress}           |
| grant_perm       | perm_resource  | {symbol}                 |
| grant_perm       | perm_action    | mint                     |
| modify_token_uri_perm_token    | to             | {ownerAddress}           |
| modify_token_uri_perm_token    | perm_resource  | {symbol}                 |
| modify_token_uri_perm_token    | perm_action    | modify_token_uri         |
| issue_token      | name           | {name}                   |
| issue_token      | symbol         | {symbol}                 |
| issue_token      | denom          | {symbol}{tokenid}        |
| issue_token      | owner          | {ownerAddress}           |
| issue_token      | amount         | 1                        |
| issue_token      | token_uri      | {token_uri}              |
| issue_token      | token_type     | cnft                     |
| mint_token       | amount         | {amount}{symbol}{tokenid}|
| mint_token       | to             | {ownerAddress}           |

### MsgMint{Collection}
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | mint_token               |
| mint_token       | amount         | {amount}{denom}          |
| mint_token       | to             | {ownerAddress}           |

### MsgBurn{Collection}
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | burn_token               |
| burn_token       | amount         | {amount}{denom}          |
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

### MsgTransferFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | revoke_permission        |
| message          | module         | bank                     |
| transfer_ft      | from           | {fromAddress}            |
| transfer_ft      | to             | {toAddress}              |
| transfer_ft      | symbol         | {symbol}                 |
| transfer_ft      | amount         | {amount}                 |

### MsgTransferCFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | revoke_permission        |
| message          | module         | bank                     |
| transfer_cft     | from           | {fromAddress}            |
| transfer_cft     | to             | {toAddress}              |
| transfer_cft     | symbol         | {symbol}                 |
| transfer_cft     | token_id       | {token_id}               |
| transfer_cft     | amount         | {amount}                 |

### MsgTransferNFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | revoke_permission        |
| transfer_nft     | from           | {fromAddress}            |
| transfer_nft     | to             | {toAddress}              |
| transfer_nft     | symbol         | {symbol}                 |

### MsgTransferCNFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | revoke_permission        |
| transfer_cnft    | from           | {fromAddress}            |
| transfer_cnft    | to             | {toAddress}              |
| transfer_cnft    | symbol         | {symbol}                 |
| transfer_cnft    | token_id       | {token_id}               |

### MsgAttach
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | revoke_permission        |
| attach_token     | from           | {fromAddress}            |
| attach_token     | to_token       | {toTokenDenom}           |
| attach_token     | token          | {tokenDenom}             |

### MsgDetach
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | revoke_permission        |
| detach_token     | from           | {fromAddress}            |
| detach_token     | to             | {toAddress}              |
| detach_token     | token          | {tokenDenom}             |
