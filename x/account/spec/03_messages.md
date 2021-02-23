# Messages

## MsgEmpty

```golang
type MsgEmpty struct {
	From string         `json:"name"`
}
```

**Do nothing**
- This message does nothing
- Signer(`From`) of this message is required
- You can pay the fee for any other message with this message

### MsgCreateAccount

```golang
type MsgCreateAccount struct {
	From   sdk.AccAddress `json:"from"`
	Target sdk.AccAddress `json:"target"`
}
```

**Create an account**
- Signer(`FromAddress`) of this message must already been registered before 
- `TargetAddress` must never been registered before

# Syntax
| Message/Attributes | Tag | Type |
| ---- | ---- | ---- |
| Message | account/MsgCreateAccount | github.com/line/link/x/account/internal/types.MsgCreateAccount |  
 | Attributes | from | []uint8 |  
 | Attributes | target | []uint8 |  
| Message | account/MsgEmpty | github.com/line/link/x/account/internal/types.MsgEmpty |  
 | Attributes | from | []uint8 |  
