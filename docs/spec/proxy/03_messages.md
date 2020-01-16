# Messages

## Base Msg

`MsgProxy` is base msg for other msgs in the proxy module. 
However, `MsgProxy` isn't being used as it is.   

### MsgProxy
```golang
type MsgProxy struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	OnBehalfOf sdk.AccAddress `json:"on_behalf_of"`
	Denom      string         `json:"denom"`
	Amount     sdk.Int        `json:"amount"`
}
```

## Approve and DisApprove

`MsgProxyApproveCoins` and `MsgProxyDisapproveCoins` are to approve/disapprove coins 
to the proxy. The proxy may send approved coins on behalf of approving account.

`OnBehalfOf` account's signature is required.   

### MsgProxyApproveCoins
```golang
type MsgProxyApproveCoins struct {
	MsgProxy
}
```

### MsgProxyDisapproveCoins
```golang
type MsgProxyDisapproveCoins struct {
	MsgProxy
}
```

## Send Coins From

`MsgProxySendCoinsFrom` initiates the coin transfer on behalf of the original coin owner.

`Proxy` account's signature is required.  

### MsgProxySendCoinsFrom
```golang
type MsgProxySendCoinsFrom struct {
	MsgProxy
	ToAddress sdk.AccAddress `json:"to_address"`
}
```