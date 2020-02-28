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
| grant_perm       | perm_action    | mint                     |
| grant_perm       | perm_action    | burn                     |
| grant_perm       | perm_action    | modify                   |
| create_collection| contract_id    | {contractID}             |
| create_collection| name           | {name}                   |
| create_collection| owner          | {ownerAddress}           |

### MsgIssueFT                                                   
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {ownerAddress}           | 
| message          | action         | issue_ft                 |
| issue_ft         | contract_id    | {contractID}             |
| issue_ft         | name           | {name}                   |
| issue_ft         | token_id       | {tokenid}                |
| issue_ft         | owner          | {ownerAddress}           |
| issue_ft         | amount         | {amount}                 |
| issue_ft         | mintable       | {mintable}               |
| issue_ft         | decimals       | {decimals}               |

### MsgMintFT
| Type             | Attribute Key  | Attribute Value              |
|------------------|----------------|------------------------------|
| message          | module         | collection                   |
| message          | sender         | {ownerAddress}               | 
| message          | action         | mint_ft                      |
| mint_ft          | contract_id    | {contractID}                 |
| mint_ft          | amount         | {amount}{contractID}{tokenid}|
| mint_ft          | from           | {fromAddress}                |
| mint_ft          | to             | {toAddress}                  |

### MsgBurnFT
| Type             | Attribute Key  | Attribute Value              |
|------------------|----------------|------------------------------|
| message          | module         | collection                   |
| message          | sender         | {fromAddress}                | 
| message          | action         | burn_ft                      |
| burn_ft          | contract_id    | {contractID}                 |
| burn_ft          | from           | {fromAddress}                |
| burn_ft          | amount         | {amount}{contractID}{tokenid}|

### MsgBurnFTFrom
| Type             | Attribute Key  | Attribute Value              |
|------------------|----------------|------------------------------|
| message          | module         | collection                   |
| message          | sender         | {proxyAddress}               | 
| message          | action         | burn_ft                      |
| burn_ft_from     | contract_id    | {contractID}                 |
| burn_ft_from     | proxy          | {proxyAddress}               |
| burn_ft_from     | from           | {fromAddress}                |
| burn_ft_from     | amount         | {amount}{contractID}{tokenid}|

### MsgIssueNFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | issue_nft                |
| issue_nft        | contract_id    | {contractID}             |
| issue_nft        | token_type     | {tokentype}              |

### MsgMintNFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | mint_nft                 |
| mint_nft         | contract_id    | {contractID}             |
| mint_nft         | name           | {name}                   |
| mint_nft         | token_id       | {tokenid}                |
| mint_nft         | from           | {fromAddress}            |
| mint_nft         | to             | {toAddress}              |

### MsgBurnNFT
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | burn_nft                 |
| burn_nft         | from           | {fromAddress}            |
| burn_nft         | contract_id    | {contractID}             |
| burn_nft         | token_id       | {token_id}               |

### MsgBurnNFTFrom
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {proxyAddress}           | 
| message          | action         | burn_nft _from           |
| burn_nft_from    | contract_id    | {contractID}             |
| burn_nft_from    | proxy          | {proxyAddress}           |
| burn_nft_from    | from           | {fromAddress}            |
| burn_nft_from    | token_id       | {token_id}               |

### MsgGrantPermission
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | grant_permission         |
| grant_perm       | from           | {fromAddress}            |
| grant_perm       | to             | {toAddress}              |
| grant_perm       | perm_resource  | {resource}               |
| grant_perm       | perm_action    | issue/mint/burn/modify   |

### MsgRevokePermission
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | revoke_permission        |
| revoke_perm      | from           | {fromAddress}            |
| revoke_perm      | perm_resource  | {resource}               |
| revoke_perm      | perm_action    | issue/mint/burn/modify   |

### MsgTransferFT
| Type             | Attribute Key  | Attribute Value               |
|------------------|----------------|-------------------------------|
| message          | module         | collection                    |
| message          | sender         | {fromAddress}                 | 
| message          | action         | transfer_ft                   |
| transfer_ft      | contract_id    | {contractID}                  |
| transfer_ft      | from           | {fromAddress}                 |
| transfer_ft      | to             | {toAddress}                   |
| transfer_ft      | amount         | {amount}{contractID}{tokenID} |

### MsgTransferFTFrom
| Type              | Attribute Key  | Attribute Value               |
|-------------------|----------------|-------------------------------|
| message           | module         | collection                    |
| message           | sender         | {proxyAddress}                | 
| message           | action         | transfer_ft _from             |
| transfer_ft_from  | contract_id    | {contractID}                  |
| transfer_ft_from  | proxy          | {proxyAddress}                |
| transfer_ft_from  | from           | {fromAddress}                 |
| transfer_ft_from  | to             | {toAddress}                   |
| transfer_ft_from  | amount         | {amount}{contractID}{tokenID} |

### MsgTransferNFT
| Type                | Attribute Key  | Attribute Value       |
|---------------------|----------------|-----------------------|
| message             | module         | collection            |
| message             | sender         | {fromAddress}         | 
| message             | action         | transfer_nft          |
| transfer_nft        | contract_id    | {contractID}          |
| transfer_nft        | from           | {fromAddress}         |
| transfer_nft        | to             | {toAddress}           |
| transfer_nft        | token_id       | {token_id}            |

### MsgTransferNFTFrom
| Type                | Attribute Key  | Attribute Value       |
|---------------------|----------------|-----------------------|
| message             | module         | collection            |
| message             | sender         | {proxyAddress}        | 
| message             | action         | transfer_nft _from    |
| transfer_nft _from  | contract_id    | {contractID}          |
| transfer_nft _from  | proxy          | {proxyAddress}        |
| transfer_nft _from  | from           | {fromAddress}         |
| transfer_nft _from  | to             | {toAddress}           |
| transfer_nft _from  | token_id       | {token_id}            |

### MsgAttach
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | attach                   |
| attach_nft       | contract_id    | {contractID}             |
| attach_nft       | from           | {fromAddress}            |
| attach_nft       | to_token_id    | {to_token_id}            |
| attach_nft       | token_id       | {token_id}               |

### MsgDetach
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {fromAddress}            | 
| message          | action         | detach                   |
| detach_nft       | contract_id    | {contractID}             |
| detach_nft       | from           | {fromAddress}            |
| detach_nft       | from_token_id  | {from_token_id}          |
| detach_nft       | token_id       | {token_id}               |

### MsgAttachFrom
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {proxyAddress}           | 
| message          | action         | attach_from              |
| attach_nft_from  | contract_id    | {contractID}             |
| attach_nft_from  | proxy          | {proxyAddress}           |
| attach_nft_from  | from           | {fromAddress}            |
| attach_nft_from  | to_token_id    | {to_token_id}            |
| attach_nft_from  | token_id       | {token_id}               |

### MsgDetachFrom
| Type             | Attribute Key  | Attribute Value          |
|------------------|----------------|--------------------------|
| message          | module         | collection               |
| message          | sender         | {proxyAddress}           | 
| message          | action         | detach_from              |
| detach_nft_from  | contract_id    | {contractID}             |
| detach_nft_from  | proxy          | {proxyAddress}           |
| detach_nft_from  | from           | {fromAddress}            |
| detach_nft_from  | from_token_id  | {from_token_id}          |
| detach_nft_from  | token_id       | {token_id}               |

### MsgApprove
| Type               | Attribute Key  | Attribute Value        |
|--------------------|----------------|------------------------|
| message            | module         | collection             |
| message            | sender         | {approverAddress}      | 
| message            | action         | approve_collection     |
| approve_collection | contract_id    | {contractID}           |
| approve_collection | approver       | {approverAddress}      |
| approve_collection | proxy          | {proxyAddress}         |

### MsgDisapprove
| Type                  | Attribute Key  | Attribute Value       |
|-----------------------|----------------|-----------------------|
| message               | module         | collection            |
| message               | sender         | {approverAddress}     | 
| message               | action         | disapprove_collection |
| disapprove_collection | contract_id    | {contractID}          |
| disapprove_collection | approver       | {approverAddress}     |
| disapprove_collection | proxy          | {proxyAddress}        |

### MsgModify
| Type                  | Attribute Key  | Attribute Value       |
|-----------------------|----------------|-----------------------|
| message               | module         | collection            |
| message               | sender         | {ownerAddress}        | 
| message               | action         | modify_collection     |
| message               | action         | modify_token_type     |
| message               | action         | modify_token          |
| modify_collection     | contract_id    | {contract_id}         |
| modify_collection     | {modifiedField}| {modifiedValue}       |
| modify_token_type     | contract_id    | {contract_id}         |
| modify_token_type     | token_type     | {token_type}          |
| modify_token_type     | {modifiedField}| {modifiedValue}       |
| modify_token          | contract_id    | {contract_id}         |
| modify_token          | token_id       | {token_id}            |
| modify_token          | {modifiedField}| {modifiedValue}       |
