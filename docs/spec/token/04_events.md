# Events

The token module emits the following events:


### MsgPublishToken

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| transfer | recipient     | {recipientAddress} |
| transfer | amount        | {amount}           |
| message  | module        | token              |
| message  | sender        | {senderAddress}    | *fixme* https://github.com/line/link/issues/118
| message  | action        | publich_token      |
| publish_token | name     | name               |
| publish_token | symbol   | token symbol       |
| publish_token | owner    | {ownerAddress}     |
| publish_token | amount   | {amount}           |
| publish_token | mintable | {boolean}          |

### MsgMint
| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| transfer | recipient     | {recipientAddress} |
| transfer | amount        | {amount}           |
| message  | module        | token              |
| message  | sender        | {senderAddress}    | 
| message  | action        | mint_token         |
| mint_token | to          | {recipientAddress} |
| mint_token | amount      | {amount}           |

### MsgBurn
| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| transfer | recipient     | {recipientAddress} |
| transfer | amount        | {amount}           |
| message  | module        | token              |
| message  | sender        | {senderAddress}    | 
| message  | action        | burn_token         |
| burn_token | from        | {senderAddress}    |
| burn_token | amount      | {amount}           |

### MsgGrantPermission
| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| message  | module        | token              |
| message  | action        | grant_perm         |
| grant_perm | from        | {senderAddress}    |
| grant_perm | to          | {recipientAddress} |
| grant_perm | resource    | {token}            |
| grant_perm | action      | {mint/burn}        |

### MsgRevokePermission
| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| message  | module        | token              |
| message  | action        | revoke_perm        |
| revoke_perm | from       | {senderAddress}    |
| revoke_perm | resource   | {token}            |
| revoke_perm | action     | {mint/burn}        |
