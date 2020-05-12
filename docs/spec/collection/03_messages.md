# Messages
## MsgIssue

**Issue token messages are to create a new token on Link Chain**
- The new contract id is generated while issuer issues and the issue permission is granted to the issuer
- An issuer who granted issue permission can issue collective tokens
- Mint permission is granted to the token issuer when the token is mintable
- The identifier for the collective token is defined by the concatenation of the contract_id and the token id

### MsgCreate
```golang
type MsgCreateCollection struct {
	Owner      sdk.AccAddress `json:"owner"`
	Name       string         `json:"name"`
	Meta       string         `json:"meta"`
	BaseImgURI string         `json:"base_img_uri"`
}
```

### MsgIssueFT
```golang
type MsgIssueFT struct {
	Owner      sdk.AccAddress `json:"owner"`
	ContractID string         `json:"contract_id"`
	To         sdk.AccAddress `json:"to"`
	Name       string         `json:"name"`
	Meta       string         `json:"meta"`
	Amount     sdk.Int        `json:"amount"`
	Mintable   bool           `json:"mintable"`
	Decimals   sdk.Int        `json:"decimals"`
}
```


### MsgIssueNFT
```golang
type MsgIssueNFT struct {
	Owner      sdk.AccAddress `json:"owner"`
	ContractID string         `json:"contract_id"`
	Name       string         `json:"name"`
	Meta       string         `json:"meta"`
}
```


## Mint

**Mint message is to increase the total supply of the token**
- Signer(From) of this message must have permission 
- Minted token is added to the `To` account

### MsgMintFT

```golang
type MsgMintFT struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	To         sdk.AccAddress `json:"to"`
	Amount     Coins          `json:"amount"`
}

type Coin struct {
	Denom  string   `json:"token_id"`
	Amount sdk.Int  `json:"amount"`
}

type Coins []Coin
```


### MsgMintNFT
```golang
type MintNFTParam struct {
	Name      string `json:"name"`
	Meta      string `json:"meta"`
	TokenType string `json:"token_type"`
}

type MsgMintNFT struct {
	From          sdk.AccAddress `json:"from"`
	ContractID    string         `json:"contract_id"`
	To            sdk.AccAddress `json:"to"`
	MintNFTParams []MintNFTParam `json:"params"`
}
```

## Burn
**Burn message is to decrease the total supply of the token**
- Signer(From) of this message must have the amount of the tokens
- Token is subtracted from the `From` account 

### MsgBurnFT

```golang
type MsgBurnFT struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	Amount     Coins          `json:"amount"`
}
```

### MsgBurnNFT
```golang
type MsgBurnNFT struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	TokenIDs   []string       `json:"token_ids"`
}
```

### MsgBurnFTFrom

```golang
type MsgBurnFTFrom struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
	From       sdk.AccAddress `json:"from"`
	Amount     Coins          `json:"amount"`
}
```

### MsgBurnNFTFrom
```golang
type MsgBurnNFTFrom struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
	From       sdk.AccAddress `json:"from"`
	TokenIDs   []string       `json:"token_ids"`
}
```

## MsgGrantPermission

```golang
type MsgGrantPermission struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	To         sdk.AccAddress `json:"to"`
	Permission Permission     `json:"permission"`
}
```

**Grant Permission is to give a permission to the `To` account**
- `From` account must has the permission

## MsgRevokePermission

```golang
type MsgRevokePermission struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	Permission Permission     `json:"permission"`
}
```

**Revoke Permission is to dump a permission from the `From` account**
- `From` account must has the permission


## MsgTransferFT
```golang
type MsgTransferFT struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	To         sdk.AccAddress `json:"to"`
	Amount     Coins          `json:"amount"`
}
```

**TransferFT message is to transfer a collective non-reserved fungible token**
- Signer of this message must have the amount of the tokens
- Token is subtracted from the `From` account
- Token is added to the `To` account


## MsgTransferNFT

```golang
type MsgTransferNFT struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	To         sdk.AccAddress `json:"to"`
	TokenIDs   []string       `json:"token_ids"`
}
```

**TransferNFT message is to transfer a collective non-fungible token**
- Signer of this message must have the token
- Token is subtracted from the `From` account
- Token is added to the `To` account


## MsgTransferFTFrom

```golang
type MsgTransferFTFrom struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
	From       sdk.AccAddress `json:"from"`
	To         sdk.AccAddress `json:"to"`
	Amount     Coins          `json:"amount"`
}
```

**TransferFTFrom message is for `Proxy` to transfer a collective non-reserved fungible token owned by `From`**
- Signer(`Proxy`) of this message must have been approved for the collection
- Token is subtracted from the `From` account
- Token is added to the `To` account


## MsgTransferNFTFrom

```golang
type MsgTransferNFTFrom struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
	From       sdk.AccAddress `json:"from"`
	To         sdk.AccAddress `json:"to"`
	TokenIDs   []string       `json:"token_ids"`
}
```

**TransferNFT message is for `Proxy` to transfer a collective non-fungible token owned by `From`**
- Signer(`Proxy`) of this message must have been approved for the collection
- Token is subtracted from the `From` account
- Token is added to the `To` account


## MsgAttach

```golang
type MsgAttach struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	ToTokenID  string         `json:"to_token_id"`
	TokenID    string         `json:"token_id"`
}
```

**Attach message is to attach a non-fungible token to another non-fungible token**
- Signer(`From`) of this message must have the token
- The token having `TokenID` is attached to the token having `ToTokenID`
- If the owner of the `ToToken` is different with `From`, the owner of the Token is changed to the owner of `ToToken`
- Cannot attach a child token of some other to any token


## MsgDetach

```golang
type MsgDetach struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	TokenID    string         `json:"token_id"`
}
```

**Detach message is to detach a non-fungible token from another parent token**
- Signer of this message must have the token
- Cannot detach a non-child token from any token


## MsgAttachFrom

```golang
type MsgAttachFrom struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
	From       sdk.AccAddress `json:"from"`
	ToTokenID  string         `json:"to_token_id"`
	TokenID    string         `json:"token_id"`
}
```

**Attach message is for a proxy to attach a non-fungible token to another non-fungible token**
- Signer(Proxy) of this message must have been approved by From having the token
- The token having TokenID is attached to the token having ToTokenID
- If the owner of the ToToken is different with From, the owner of the Token is changed to the owner of ToToken
- Cannot attach a child token of some other to any token


## MsgDetachFrom

```golang
type MsgDetachFrom struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
	From       sdk.AccAddress `json:"from"`
	TokenID    string         `json:"token_id"`
}
```

**Detach message is for a proxy to detach a non-fungible token from another parent token**
- Signer(`Proxy`) of this message must have been approved by From having the token
- Cannot detach a non-child token from any token


## MsgApprove

```golang
type MsgApprove struct {
	Approver   sdk.AccAddress `json:"approver"`
	ContractID string         `json:"contract_id"`
	Proxy      sdk.AccAddress `json:"proxy"`
}
```

**Approve message is to approve a proxy to transfer, attach/detach tokens of a collection**
- `Approver` is the signer


## MsgDisapprove

```golang
type MsgDisapprove struct {
	Approver   sdk.AccAddress `json:"approver"`
	ContractID string         `json:"contract_id"`
	Proxy      sdk.AccAddress `json:"proxy"`
}
```

**Disapprove message is to withdraw proxy's approval for a collection**
- `Approver` is the signer

## MsgModify

```golang
type MsgModify struct {
	Owner      sdk.AccAddress   `json:"owner"`
	ContractID string           `json:"contract_id"`
	TokenType  string           `json:"token_type"`
	TokenIndex string           `json:"token_index"`
	Changes    linktype.Changes `json:"changes"`
}
```

**Modify message is to modify fields of collection, token type, CFT or CNFT**
- `Owner` is the signer
