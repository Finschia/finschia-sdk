# Messages
## MsgIssue

**Issue token messages are to create a new token on Link Chain**
- See [symbol rule](01_concept.md#Rule for defining symbols) for the details
- The first issuer for the token symbol occupies the symbol and the issue permission is granted to the issuer
- An issuer who granted issue permission can issue collective tokens
- Mint permission is granted to the token issuer when the token is mintable
- The identifier for the collective token is defined by the concatenation of the symbol and the token id

### MsgIssue
```golang
type MsgIssue struct {
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	Owner    sdk.AccAddress `json:"owner"`
	TokenURI string         `json:"token_uri"`
	Amount   sdk.Int        `json:"amount"`
	Mintable bool           `json:"mintable"`
	Decimals sdk.Int        `json:"decimals"`
}
```

### MsgIssueCollection
```golang
type MsgIssueCollection struct {
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	Owner    sdk.AccAddress `json:"owner"`
	TokenURI string         `json:"token_uri"`
	Amount   sdk.Int        `json:"amount"`
	Mintable bool           `json:"mintable"`
	Decimals sdk.Int        `json:"decimals"`
	TokenID string          `json:"token_id"`
}
```


### MsgIssueNFT
```golang
type MsgIssueNFT struct {
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	Owner    sdk.AccAddress `json:"owner"`
	TokenURI string         `json:"token_uri"`
}
```
### MsgIssueNFTCollection
```golang
type MsgIssueNFTCollection struct {
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	Owner    sdk.AccAddress `json:"owner"`
	TokenURI string         `json:"token_uri"`
	TokenID string          `json:"token_id"`
}
```


## MsgMint

```golang
type MsgMint struct {
	To     sdk.AccAddress `json:"to"`
	Amount sdk.Coins      `json:"amount"`
}
```

**Mint message is to increase the total supply of the token**
- Signer of this message must has permission 
- Minted token is added to the `To` account

## MsgBurn

```golang
type MsgBurn struct {
	From   sdk.AccAddress `json:"from"`
	Amount sdk.Coins      `json:"amount"`
}
```
**Burn message is to decrease the total supply of the token**
- Signer of this message must have the amount of the tokens
- Token is subtracted from the `From` account 

##MsgGrantPermission

```golang
type MsgGrantPermission struct {
	From       sdk.AccAddress `json:"from"`
	To         sdk.AccAddress `json:"to"`
	Permission Permission     `json:"permission"`
}
```

**Grant Permission is to give a permission to the `To` account**
- `From` account must has the permission

##MsgRevokePermission

```golang
type MsgRevokePermission struct {
	From       sdk.AccAddress `json:"from"`
	Permission Permission     `json:"permission"`
}
```

**Revoke Permission is to dump a permission from the `From` account**
- `From` account must has the permission
