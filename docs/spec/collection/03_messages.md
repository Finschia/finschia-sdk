# Messages
## MsgIssue

**Issue token messages are to create a new token on Link Chain**
- The first issuer for the token symbol occupies the symbol and the issue permission is granted to the issuer
- An issuer who granted issue permission can issue collective tokens
- Mint permission is granted to the token issuer when the token is mintable
- The identifier for the collective token is defined by the concatenation of the symbol and the token id

### MsgCreate
```golang
type MsgCreate struct {
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	Owner    sdk.AccAddress `json:"owner"`
}
```

### MsgIssueFT
```golang
type MsgIssueFT struct {
	Symbol   string         `json:"symbol"`
	Name     string         `json:"name"`
	Owner    sdk.AccAddress `json:"owner"`
	TokenURI string         `json:"token_uri"`
	Amount   sdk.Int        `json:"amount"`
	Mintable bool           `json:"mintable"`
	Decimals sdk.Int        `json:"decimals"`
}
```


### MsgIssueNFT
```golang
type MsgIssueNFT struct {
	Symbol   string         `json:"symbol"`
	Owner    sdk.AccAddress `json:"owner"`
}
```


## Mint

**Mint message is to increase the total supply of the token**
- Signer(From) of this message must have permission 
- Minted token is added to the `To` account

### MsgMintFT

```golang
type MsgMintFT struct {
	From   sdk.AccAddress            `json:"from"`
	To     sdk.AccAddress            `json:"to"`
	Amount linktype.CoinWithTokenIDs `json:"amount"`
}

type CoinWithTokenIDs []CoinWithTokenID

type CoinWithTokenID struct {
	Symbol  string  `json:"symbol"`
	TokenID string  `json:"token_id"`
	Amount  sdk.Int `json:"amount"`
}
```

### MsgMintNFT
```golang
type MsgMintNFT struct {
	Symbol    string         `json:"symbol"`
	Name      string         `json:"name"`
	From      sdk.AccAddress `json:"from"`
	To        sdk.AccAddress `json:"to"`
	TokenURI  string         `json:"token_uri"`
	TokenType string         `json:"token_type"`
}
```

## Burn
**Burn message is to decrease the total supply of the token**
- Signer(From) of this message must have the amount of the tokens
- Token is subtracted from the `From` account 

### MsgBurnFT

```golang
type MsgBurnFT struct {
	From   sdk.AccAddress            `json:"from"`
	Amount linktype.CoinWithTokenIDs `json:"amount"`
}

type CoinWithTokenIDs []CoinWithTokenID

type CoinWithTokenID struct {
	Symbol  string  `json:"symbol"`
	TokenID string  `json:"token_id"`
	Amount  sdk.Int `json:"amount"`
}
```

### MsgBurnNFT
```golang
type MsgBurnNFT struct {
	From    sdk.AccAddress `json:"from"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
}
```

### MsgBurnFTFrom

```golang
type MsgBurnFTFrom struct {
	Proxy  sdk.AccAddress            `json:"proxy"`
	From   sdk.AccAddress            `json:"from"`
	Amount linktype.CoinWithTokenIDs `json:"amount"`
}

type CoinWithTokenIDs []CoinWithTokenID

type CoinWithTokenID struct {
	Symbol  string  `json:"symbol"`
	TokenID string  `json:"token_id"`
	Amount  sdk.Int `json:"amount"`
}
```

### MsgBurnNFTFrom
```golang
type MsgBurnNFTFrom struct {
	Proxy  sdk.AccAddress  `json:"proxy"`
	From    sdk.AccAddress `json:"from"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
}
```

## MsgGrantPermission

```golang
type MsgGrantPermission struct {
	From       sdk.AccAddress `json:"from"`
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
	Permission Permission     `json:"permission"`
}
```

**Revoke Permission is to dump a permission from the `From` account**
- `From` account must has the permission


## MsgTransferFT
```golang
type MsgTransferFT struct {
	From        sdk.AccAddress `json:"from" yaml:"from"`
	To          sdk.AccAddress `json:"to" yaml:"to"`
	TokenSymbol string         `json:"token_symbol"`
	TokenID     string         `json:"token_id"`
	Amount      sdk.Int        `json:"amount" yaml:"amount"`
}
```

**TransferFT message is to transfer a collective non-reserved fungible token**
- Signer of this message must have the amount of the tokens
- Token is subtracted from the `From` account
- Token is added to the `To` account


## MsgTransferNFT

```golang
type MsgTransferNFT struct {
	From        sdk.AccAddress `json:"from"`
	To          sdk.AccAddress `json:"to"`
	TokenSymbol string         `json:"token_symbol"`
	TokenID     string         `json:"token_id"`
}
```

**TransferNFT message is to transfer a collective non-fungible token**
- Signer of this message must have the token
- Token is subtracted from the `From` account
- Token is added to the `To` account


## MsgTransferFTFrom

```golang
type MsgTransferFTFrom struct {
	Proxy   sdk.AccAddress `json:"proxy"`
	From    sdk.AccAddress `json:"from"`
	To      sdk.AccAddress `json:"to"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
	Amount  sdk.Int        `json:"amount"`
}
```

**TransferFTFrom message is for `Proxy` to transfer a collective non-reserved fungible token owned by `From`**
- Signer(`Proxy`) of this message must have been approved for the collection
- Token is subtracted from the `From` account
- Token is added to the `To` account


## MsgTransferNFTFrom

```golang
type MsgTransferNFTFrom struct {
	Proxy   sdk.AccAddress `json:"proxy"`
	From    sdk.AccAddress `json:"from"`
	To      sdk.AccAddress `json:"to"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
}
```

**TransferNFT message is for `Proxy` to transfer a collective non-fungible token owned by `From`**
- Signer(`Proxy`) of this message must have been approved for the collection
- Token is subtracted from the `From` account
- Token is added to the `To` account


## MsgAttach

```golang
type MsgAttach struct {
	From      sdk.AccAddress `json:"from"`
	Symbol    string         `json:"symbol"`
	ToTokenID string         `json:"to_token_id"`
	TokenID   string         `json:"token_id"`
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
	From    sdk.AccAddress `json:"from"`
	To      sdk.AccAddress `json:"to"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
}
```

**Detach message is to detach a non-fungible token from another parent token**
- Signer of this message must have the token
- The token of TokenID will be owned by To
- Cannot detach a non-child token from any token


## MsgAttachFrom

```golang
type MsgAttachFrom struct {
	Proxy     sdk.AccAddress `json:"proxy"`
	From      sdk.AccAddress `json:"from"`
	Symbol    string         `json:"symbol"`
	ToTokenID string         `json:"to_token_id"`
	TokenID   string         `json:"token_id"`
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
	Proxy    sdk.AccAddress `json:"proxy"`
	From     sdk.AccAddress `json:"from"`
	To       sdk.AccAddress `json:"to"`
	Symbol   string         `json:"symbol"`
	TokenID  string         `json:"token_id"`
}
```

**Detach message is for a proxy to detach a non-fungible token from another parent token**
- Signer(`Proxy`) of this message must have been approved by From having the token
- The token of TokenID will be owned by To
- Cannot detach a non-child token from any token


## MsgApprove

```golang
type MsgApprove struct {
	Approver sdk.AccAddress `json:"approver"`
	Proxy    sdk.AccAddress `json:"proxy"`
	Symbol   string         `json:"symbol"`
}
```

**Approve message is to approve a proxy to transfer, attach/detach tokens of a collection**
- `Approver` is the signer


## MsgDisapprove

```golang
type MsgDisapprove struct {
	Approver sdk.AccAddress `json:"approver"`
	Proxy    sdk.AccAddress `json:"proxy"`
	Symbol   string         `json:"symbol"`
}
```

**Disapprove message is to withdraw proxy's approval for a collection**
- `Approver` is the signer
