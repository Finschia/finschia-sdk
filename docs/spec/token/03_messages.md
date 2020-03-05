# Messages
## MsgIssue

**Issue token messages are to create a new token on Link Chain**
- See [symbol rule](01_concept.md#rule-for-defining-symbols) for the details
- The first issuer for the token symbol occupies the symbol and the issue permission is granted to the issuer
- An issuer who granted issue permission can issue collective tokens
- Mint permission is granted to the token issuer when the token is mintable
- The identifier for the collective token is defined by the concatenation of the symbol and the token id

### MsgIssue
```golang
type MsgIssue struct {
	Owner    sdk.AccAddress `json:"owner"`
	To       sdk.AccAddress `json:"to"`
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	ImageURI string         `json:"image_uri"`
	Amount   sdk.Int        `json:"amount"`
	Mintable bool           `json:"mintable"`
	Decimals sdk.Int        `json:"decimals"`
}
```
## Mint

**Mint message is to increase the total supply of the token**
- Signer(From) of this message must have permission 
- Minted token is added to the `To` account

### MsgMint

```golang
type MsgMint struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`     
	To         sdk.AccAddress `json:"to"`
	Amount     sdk.Int        `json:"amount"`
}
```

## Burn
**Burn message is to decrease the total supply of the token**
- Signer(From) of this message must have the amount of the tokens
- Token is subtracted from the `From` account 

### MsgBurn

```golang
type MsgBurn struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`   
	Amount     sdk.Int        `json:"amount"`
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


## MsgTransfer

```golang
type MsgTransfer struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	To         sdk.AccAddress `json:"to"`
	Amount     sdk.Int        `json:"amount"`
}
```

**Transfer message is to transfer a non-reserved fungible token**
- Signer of this message must have the amount of the tokens
- Token is subtracted from the `From` account
- Token is added to the `To` account

## MsgModify

```golang
type MsgModify struct {
	Owner      sdk.AccAddress   `json:"owner"`
	ContractID string           `json:"contract_id"`
	Changes    linktype.Changes `json:"changes"`
}
```

**Modify message is to modify fields of token**
- `Owner` is the signer
