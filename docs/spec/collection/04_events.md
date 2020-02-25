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
| message          | action         | issue_ft                 |
| grant_perm       | to             | {ownerAddress}           |
| grant_perm       | perm_resource  | {symbol}{token_id}       |
| grant_perm       | perm_action    | mint                     |
| grant_perm       | to             | {ownerAddress}           |
| grant_perm       | perm_resource  | {symbol}{token_id}       |
| grant_perm       | perm_action    | modify                   |
| issue_ft         | name           | {name}                   |
| issue_ft         | symbol         | {symbol}                 |
| issue_ft         | token_id       | {tokenid}                |
| issue_ft         | owner          | {ownerAddress}           |
| issue_ft         | amount         | {amount}                 |
| issue_ft         | mintable       | {mintable}               |
| issue_ft         | decimals       | {decimals}               |
| issue_ft         | token_uri      | {token_uri}              |

### MsgMintFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {ownerAddress}           | 
| message          | action         | mint_ft                  |
| mint_ft          | amount         | {amount}{symbol}{tokenid}|
| mint_ft          | from           | {fromAddress}            |
| mint_ft          | to             | {toAddress}              |

### MsgBurnFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | burn_ft                  |
| burn_ft          | from           | {fromAddress}            |
| burn_ft          | amount         | {amount}{symbol}{tokenid}|

### MsgBurnFTFrom
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {proxyAddress}           | 
| message          | action         | burn_ft                  |
| burn_ft _from    | proxy          | {proxyAddress}           |
| burn_ft _from    | from           | {fromAddress}            |
| burn_ft _from    | amount         | {amount}{symbol}{tokenid}|

### MsgIssueNFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | issue_nft                |
| grant_perm       | to             | {toAddress}              |
| grant_perm       | perm_resource  | {symbol}                 |
| grant_perm       | perm_action    | mint                     |
| issue_nft        | symbol         | {symbol}                 |
| issue_nft        | token_type     | {tokentype}              |

### MsgMintNFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | mint_nft                 |
| grant_perm       | to             | {toAddress}              |
| grant_perm       | perm_resource  | {symbol}{token_id}       |
| grant_perm       | perm_action    | modify                   |
| mint_nft         | name           | {name}                   |
| mint_nft         | symbol         | {symbol}                 |
| mint_nft         | token_id       | {tokenid}                |
| mint_nft         | from           | {fromAddress}            |
| mint_nft         | to             | {toAddress}              |
| mint_nft         | token_uri      | {token_uri}              |

### MsgBurnNFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | burn_nft                 |
| burn_nft         | from           | {fromAddress}            |
| burn_nft         | symbol         | {symbol}                 |
| burn_nft         | token_id       | {token_id}               |

### MsgBurnNFTFrom
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {proxyAddress}           | 
| message          | action         | burn_nft _from           |
| burn_nft _from   | proxy          | {proxyAddress}           |
| burn_nft _from   | from           | {fromAddress}            |
| burn_nft _from   | symbol         | {symbol}                 |
| burn_nft _from   | token_id       | {token_id}               |

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
| message          | action         | transfer_ft              |
| transfer_ft      | from           | {fromAddress}            |
| transfer_ft      | to             | {toAddress}              |
| transfer_ft      | symbol         | {symbol}                 |
| transfer_ft      | token_id       | {token_id}               |
| transfer_ft      | amount         | {amount}                 |

### MsgTransferFTFrom
| Type              | Attribute Key  | Attribute Value         |
|-------------------|----------------|-------------------------|
| message           | module         | collection              |
| message           | sender         | {proxyAddress}          | 
| message           | action         | transfer_ft _from       |
| transfer_ft _from | proxy          | {proxyAddress}          |
| transfer_ft _from | from           | {fromAddress}           |
| transfer_ft _from | to             | {toAddress}             |
| transfer_ft _from | symbol         | {symbol}                |
| transfer_ft _from | token_id       | {token_id}              |
| transfer_ft _from | amount         | {amount}                |

### MsgTransferNFT
| Type                | Attribute Key  | Attribute Value       |
|---------------------|----------------|-----------------------|
| message             | module         | collection            |
| message             | sender         | {fromAddress}         | 
| message             | action         | transfer_nft          |
| transfer_nft        | from           | {fromAddress}         |
| transfer_nft        | to             | {toAddress}           |
| transfer_nft        | symbol         | {symbol}              |
| transfer_nft        | token_id       | {token_id}            |
| transfer_nft _child | child_token_id | {child_token_id}      |

### MsgTransferNFTFrom
| Type                | Attribute Key  | Attribute Value       |
|---------------------|----------------|-----------------------|
| message             | module         | collection            |
| message             | sender         | {proxyAddress}        | 
| message             | action         | transfer_nft _from    |
| transfer_nft _from  | proxy          | {proxyAddress}        |
| transfer_nft _from  | from           | {fromAddress}         |
| transfer_nft _from  | to             | {toAddress}           |
| transfer_nft _from  | symbol         | {symbol}              |
| transfer_nft _from  | token_id       | {token_id}            |
| transfer_nft _child | child_token_id | {child_token_id}      |

### MsgAttach
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | attach                   |
| attach_nft       | from           | {fromAddress}            |
| attach_nft       | symbol         | {symbol}                 |
| attach_nft       | to_token_id    | {to_token_id}            |
| attach_nft       | token_id       | {token_id}               |

### MsgDetach
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | detach                   |
| detach_nft       | from           | {fromAddress}            |
| detach_nft       | symbol         | {symbol}                 |
| detach_nft       | from_token_id  | {from_token_id}          |
| detach_nft       | token_id       | {token_id}               |

### MsgAttachFrom
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {proxyAddress}           | 
| message          | action         | attach_from              |
| attach_nft _from | proxy          | {proxyAddress}           |
| attach_nft _from | from           | {fromAddress}            |
| attach_nft _from | symbol         | {symbol}                 |
| attach_nft _from | to_token_id    | {to_token_id}            |
| attach_nft _from | token_id       | {token_id}               |

### MsgDetachFrom
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {proxyAddress}           | 
| message          | action         | detach_from              |
| detach_nft _from | proxy          | {proxyAddress}           |
| detach_nft _from | from           | {fromAddress}            |
| detach_nft _from | symbol         | {symbol}                 |
| detach_nft _from | from_token_id  | {from_token_id}          |
| detach_nft _from | token_id       | {token_id}               |

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
