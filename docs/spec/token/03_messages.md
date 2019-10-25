# Messages

## MsgPublishToken


```golang
type MsgPublishToken struct {
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	Amount   int64          `json:"amount"`
	Owner    sdk.AccAddress `json:"owner"`
	Mintable bool           `json:"mintable"`
}
```

**Publishing token is to create a new token on Link Chain**
- Mint/Burn for the token is granted to the token owner
- Token is issued and added to the token owner


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
- Signer of this message must has permission 
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
