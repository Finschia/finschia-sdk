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
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	Owner    sdk.AccAddress `json:"owner"`
	TokenURI string         `json:"token_uri"`
	Amount   sdk.Int        `json:"amount"`
	Mintable bool           `json:"mintable"`
	Decimals sdk.Int        `json:"decimals"`
}
```

### MsgCreateCollection
```golang
type MsgCreateCollection struct {
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	Owner    sdk.AccAddress `json:"owner"`
}
```

### MsgIssueCFT
```golang
type MsgIssueCFT struct {
	Symbol   string         `json:"symbol"`
	Name     string         `json:"name"`
	Owner    sdk.AccAddress `json:"owner"`
	TokenURI string         `json:"token_uri"`
	Amount   sdk.Int        `json:"amount"`
	Mintable bool           `json:"mintable"`
	Decimals sdk.Int        `json:"decimals"`
}
```


### MsgIssueCNFT
```golang
type MsgIssueCNFT struct {
	Symbol   string         `json:"symbol"`
	Owner    sdk.AccAddress `json:"owner"`
}
```


## Mint

**Mint message is to increase the total supply of the token**
- Signer(From) of this message must have permission 
- Minted token is added to the `To` account

### MsgMint

```golang
type MsgMint struct {
	From   sdk.AccAddress `json:"from"`
	To     sdk.AccAddress `json:"to"`
	Amount sdk.Coins      `json:"amount"`
}
```
### MsgMintCFT

```golang
type MsgMintCFT struct {
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

### MsgMintCNFT
```golang
type MsgMintCNFT struct {
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

### MsgBurn

```golang
type MsgBurn struct {
	From   sdk.AccAddress `json:"from"`
	Amount sdk.Coins      `json:"amount"`
}
```
### MsgBurnCFT

```golang
type MsgBurnCFT struct {
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

### MsgBurnCNFT
```golang
type MsgBurnCNFT struct {
	Symbol  string         `json:"symbol"`
	From    sdk.AccAddress `json:"from"`
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
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address" yaml:"to_address"`
	TokenSymbol string         `json:"token_symbol"`
	Amount      sdk.Int        `json:"amount" yaml:"amount"`
}
```

**Transfer message is to transfer a non-reserved fungible token**
- Signer of this message must have the amount of the tokens
- Token is subtracted from the `FromAddress` account
- Token is added to the `ToAddress` account


## MsgTransferCFT

```golang
type MsgTransferCFT struct {
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address" yaml:"to_address"`
	TokenSymbol string         `json:"token_symbol"`
	TokenID     string         `json:"token_id"`
	Amount      sdk.Int        `json:"amount" yaml:"amount"`
}
```

**Transfer message is to transfer a collective non-reserved fungible token**
- Signer of this message must have the amount of the tokens
- Token is subtracted from the `FromAddress` account
- Token is added to the `ToAddress` account


## MsgTransferCNFT

```golang
type MsgTransferCNFT struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	TokenSymbol string         `json:"token_symbol"`
	TokenID     string         `json:"token_id"`
}
```

**TransferCNFT message is to transfer a collective non-fungible token**
- Signer of this message must have the token
- Token is subtracted from the `FromAddress` account
- Token is added to the `ToAddress` account


## MsgAttach

```golang
type MsgAttach struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	Symbol      string         `json:"symbol"`
	ToTokenID   string         `json:"to_token_id"`
	TokenID     string         `json:"token_id"`
}
```

**Attach message is to attach a non-fungible token to another non-fungible token**
- Signer of this message must have the token
- The token having TokenID is attached to the token having ToTokenID
- If the owner of the ToToken is different with From, the owner of the Token is changed to the owner of ToToken
- Cannot attach a child token of some other to any token


## MsgDetach

```golang
type MsgDetach struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	To          sdk.AccAddress `json:"to"`
	Symbol      string         `json:"symbol"`
	TokenID     string         `json:"token_id"`
}
```

**Detach message is to detach a non-fungible token from another parent token**
- Signer of this message must have the token
- The token of TokenID will be owned by To
- Cannot detach a non-child token from any token
