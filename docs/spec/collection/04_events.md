# Events
**Not fully documented yet** 
The token module emits the following events:


### MsgCreate
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {ownerAddress}           | 
| message          | action         | create_collection        |
| grant_perm       | to             | {ownerAddress}           |
| grant_perm       | perm_resource  | {symbol}                 |
| grant_perm       | perm_action    | issue                    |
| create_collection| name           | {name}                   |
| create_collection| symbol         | {symbol}                 |
| create_collection| owner          | {ownerAddress}           |

### MsgIssueFT                                                   
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {ownerAddress}           | 
| message          | action         | issue_cft                |
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

### MsgMintFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {ownerAddress}           | 
| message          | action         | mint_cft                 |
| mint_cft         | amount         | {amount}{symbol}{tokenid}|
| mint_cft         | from           | {fromAddress}            |
| mint_cft         | to             | {toAddress}              |

### MsgBurnFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | burn_cft                 |
| burn_cft         | from           | {fromAddress}            |
| burn_cft         | amount         | {amount}{symbol}{tokenid}|

### MsgBurnFTFrom
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {proxyAddress}           | 
| message          | action         | burn_cft                 |
| burn_cft_from    | proxy          | {proxyAddress}           |
| burn_cft_from    | from           | {fromAddress}            |
| burn_cft_from    | amount         | {amount}{symbol}{tokenid}|

### MsgIssueNFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | issue_cnft               |
| grant_perm       | to             | {toAddress}              |
| grant_perm       | perm_resource  | {symbol}                 |
| grant_perm       | perm_action    | mint                     |
| issue_cnft       | symbol         | {symbol}                 |
| issue_cnft       | token_type     | {tokentype}              |

### MsgMintNFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | mint_cnft                |
| grant_perm       | to             | {toAddress}              |
| grant_perm       | perm_resource  | {symbol}{token_id}       |
| grant_perm       | perm_action    | modify                   |
| mint_cnft        | name           | {name}                   |
| mint_cnft        | symbol         | {symbol}                 |
| mint_cnft        | token_id       | {tokenid}                |
| mint_cnft        | from           | {fromAddress}            |
| mint_cnft        | to             | {toAddress}              |
| mint_cnft        | token_uri      | {token_uri}              |

### MsgBurnNFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | burn_cnft                |
| burn_cnft        | from           | {fromAddress}            |
| burn_cnft        | symbol         | {symbol}                 |
| burn_cnft        | token_id       | {token_id}               |

### MsgBurnNFTFrom
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {proxyAddress}           | 
| message          | action         | burn_cnft_from           |
| burn_cnft_from   | proxy          | {proxyAddress}           |
| burn_cnft_from   | from           | {fromAddress}            |
| burn_cnft_from   | symbol         | {symbol}                 |
| burn_cnft_from   | token_id       | {token_id}               |

### MsgGrantPermission
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | grant_permission         |
| grant_perm       | from           | {fromAddress}            |
| grant_perm       | to             | {toAddress}              |
| grant_perm       | perm_resource  | {resource}               |
| grant_perm       | perm_action    | issue/mint/modify        |

### MsgRevokePermission
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | revoke_permission        |
| revoke_perm      | from           | {fromAddress}            |
| revoke_perm      | perm_resource  | {resource}               |
| revoke_perm      | perm_action    | issue/mint/modify        |

### MsgTransferFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | transfer_cft             |
| transfer_cft     | from           | {fromAddress}            |
| transfer_cft     | to             | {toAddress}              |
| transfer_cft     | symbol         | {symbol}                 |
| transfer_cft     | token_id       | {token_id}               |
| transfer_cft     | amount         | {amount}                 |

### MsgTransferFTFrom
| Type              | Attribute Key  | Attribute Value         |
|-------------------|----------------|-------------------------|
| message           | module         | collection              |
| message           | sender         | {proxyAddress}          | 
| message           | action         | transfer_cft_from       |
| transfer_cft_from | proxy          | {proxyAddress}          |
| transfer_cft_from | from           | {fromAddress}           |
| transfer_cft_from | to             | {toAddress}             |
| transfer_cft_from | symbol         | {symbol}                |
| transfer_cft_from | token_id       | {token_id}              |
| transfer_cft_from | amount         | {amount}                |

### MsgTransferNFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | transfer_cnft            |
| transfer_cnft    | from           | {fromAddress}            |
| transfer_cnft    | to             | {toAddress}              |
| transfer_cnft    | symbol         | {symbol}                 |
| transfer_cnft    | token_id       | {token_id}               |

### MsgTransferNFTFrom
| Type               | Attribute Key  | Attribute Value        |
|--------------------|----------------|------------------------|
| message            | module         | collection             |
| message            | sender         | {proxyAddress}         | 
| message            | action         | transfer_cnft_from     |
| transfer_cnft_from | proxy          | {proxyAddress}         |
| transfer_cnft_from | from           | {fromAddress}          |
| transfer_cnft_from | to             | {toAddress}            |
| transfer_cnft_from | symbol         | {symbol}               |
| transfer_cnft_from | token_id       | {token_id}             |

### MsgAttach
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | attach                   |
| attach_cnft      | from           | {fromAddress}            |
| attach_cnft      | symbol         | {symbol}                 |
| attach_cnft      | to_token_id    | {to_token_id}            |
| attach_cnft      | token_id       | {token_id}               |

### MsgDetach
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | detach                   |
| detach_cnft      | from           | {fromAddress}            |
| detach_cnft      | to             | {toAddress}              |
| detach_cnft      | symbol         | {symbol}                 |
| detach_cnft      | token_id       | {token_id}               |

### MsgAttachFrom
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {proxyAddress}           | 
| message          | action         | attach_from              |
| attach_cnft_from | proxy          | {proxyAddress}           |
| attach_cnft_from | from           | {fromAddress}            |
| attach_cnft_from | symbol         | {symbol}                 |
| attach_cnft_from | to_token_id    | {to_token_id}            |
| attach_cnft_from | token_id       | {token_id}               |

### MsgDetachFrom
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {proxyAddress}           | 
| message          | action         | detach_from              |
| detach_cnft_from | proxy          | {proxyAddress}           |
| detach_cnft_from | from           | {fromAddress}            |
| detach_cnft_from | to             | {toAddress}              |
| detach_cnft_from | symbol         | {symbol}                 |
| detach_cnft_from | token_id       | {token_id}               |

### MsgApprove
| Type               | Attribute Key  | Attribute Value        |
|--------------------|----------------|------------------------|
| message            | module         | collection             |
| message            | sender         | {approverAddress}      | 
| message            | action         | approve_collection     |
| approve_collection | approver       | {approverAddress}      |
| approve_collection | proxy          | {proxyAddress}         |
| approve_collection | symbol         | {symbol}               |

### MsgDisapprove
| Type                  | Attribute Key  | Attribute Value       |
|-----------------------|----------------|-----------------------|
| message               | module         | collection            |
| message               | sender         | {approverAddress}     | 
| message               | action         | disapprove_collection |
| disapprove_collection | approver       | {approverAddress}     |
| disapprove_collection | proxy          | {proxyAddress}        |
| disapprove_collection | symbol         | {symbol}              |
