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
| mint             | amount         | {amount}{symbol}         |
| mint             | from           | {ownerAddress}           |
| mint             | to             | {toAddress}              |

### MsgBurn
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | burn                     |
| burn             | amount         | {amount}{symbol}         |
| burn             | from           | {ownerAddress}           |

### MsgCreateCollection
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | create_collection        |
| grant_perm       | to             | {ownerAddress}           |
| grant_perm       | perm_resource  | {symbol}                 |
| grant_perm       | perm_action    | issue                    |
| create_collection| name           | {name}                   |
| create_collection| symbol         | {symbol}                 |
| create_collection| owner          | {ownerAddress}           |

### MsgIssueCollection                                                   
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | issue_ft_collection      |
| grant_perm       | to             | {ownerAddress}           |
| grant_perm       | perm_resource  | {symbol}{token_id}       |
| grant_perm       | perm_action    | mint                     |
| grant_perm       | to             | {ownerAddress}           |
| grant_perm       | perm_resource  | {symbol}{token_id}       |
| grant_perm       | perm_action    | modify                   |
| issue_cft        | name           | {name}                   |
| issue_cft        | symbol         | {symbol}                 |
| issue_cft        | token_id       | {tokenid}                |
| issue_cft        | owner          | {ownerAddress}           |
| issue_cft        | amount         | {amount}                 |
| issue_cft        | mintable       | {mintable}               |
| issue_cft        | decimals       | {decimals}               |
| issue_cft        | token_uri      | {token_uri}              |

### MsgMintCollection
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | mint_token               |
| mint_cft         | amount         | {amount}{symbol}{tokenid}|
| mint_cft         | from           | {fromAddress}            |
| mint_cft         | to             | {toAddress}              |

### MsgBurnCollection
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {ownerAddress}           | 
| message          | action         | burn_token               |
| burn_cft         | amount         | {amount}{symbol}{tokenid}|
| burn_cft         | from           | {fromAddress}            |

### MsgIssueNFTCollection                                                   
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | issue_nft_collection     |
| grant_perm       | to             | {toAddress}              |
| grant_perm       | perm_resource  | {symbol}                 |
| grant_perm       | perm_action    | mint                     |
| issue_cnft       | symbol         | {symbol}                 |
| issue_cnft       | token_type     | {tokentype}              |

### MsgMintNFTCollection
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | mint_token               |
| grant_perm       | to             | {toAddress}              |
| grant_perm       | perm_resource  | {symbol}{token_id}       |
| grant_perm       | perm_action    | modify                   |
| mint_cnft        | name           | {name}                   |
| mint_cnft        | symbol         | {symbol}                 |
| mint_cnft        | token_id       | {tokenid}                |
| mint_cnft        | from           | {fromAddress}            |
| mint_cnft        | to             | {toAddress}              |
| mint_cnft        | token_uri      | {token_uri}              |


### MsgGrantPermission
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | grant_permission         |
| grant_perm       | from           | {fromAddress}            |
| grant_perm       | to             | {toAddress}              |
| grant_perm       | perm_resource  | {resource}               |
| grant_perm       | perm_action    | issue/mint/modify        |

### MsgRevokePermission
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | revoke_permission        |
| revoke_perm      | from           | {fromAddress}            |
| revoke_perm      | perm_resource  | {resource}               |
| revoke_perm      | perm_action    | issue/mint/modify        |

### MsgTransferFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | transfer-ft              |
| transfer_ft      | from           | {fromAddress}            |
| transfer_ft      | to             | {toAddress}              |
| transfer_ft      | symbol         | {symbol}                 |
| transfer_ft      | amount         | {amount}                 |

### MsgTransferCFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | transfer-cft             |
| transfer_cft     | from           | {fromAddress}            |
| transfer_cft     | to             | {toAddress}              |
| transfer_cft     | symbol         | {symbol}                 |
| transfer_cft     | token_id       | {token_id}               |
| transfer_cft     | amount         | {amount}                 |

### MsgTransferCNFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | transfer-cnft            |
| transfer_cnft    | from           | {fromAddress}            |
| transfer_cnft    | to             | {toAddress}              |
| transfer_cnft    | symbol         | {symbol}                 |
| transfer_cnft    | token_id       | {token_id}               |

### MsgAttach
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | attach                   |
| attach_cnft      | from           | {fromAddress}            |
| attach_cnft      | symbol         | {symbol}                 |
| attach_cnft      | to_token_id    | {to_token_id}            |
| attach_cnft      | token_id       | {token_id}               |

### MsgDetach
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | token                    |
| message          | sender         | {fromAddress}            | 
| message          | action         | detach                   |
| detach_cnft      | from           | {fromAddress}            |
| detach_cnft      | to             | {toAddress}              |
| detach_cnft      | symbol         | {symbol}                 |
| detach_cnft      | token_id       | {token_id}               |

