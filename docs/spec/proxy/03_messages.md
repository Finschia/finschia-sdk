# Messages

## Approve and DisApprove

`MsgProxyApproveCoins` and `MsgProxyDisapproveCoins` are to approve/disapprove coins 
to the proxy. The proxy may send approved coins on behalf of approving account.

`OnBehalfOf` account's signature is required.   

### MsgProxyApproveCoins
```golang
type MsgProxyApproveCoins struct {
    Proxy      sdk.AccAddress `json:"proxy"`
    OnBehalfOf sdk.AccAddress `json:"on_behalf_of"`
    Denom      string         `json:"denom"`
    Amount     sdk.Int        `json:"amount"`
}
```

### MsgProxyDisapproveCoins
```golang
type MsgProxyDisapproveCoins struct {
    Proxy      sdk.AccAddress `json:"proxy"`
    OnBehalfOf sdk.AccAddress `json:"on_behalf_of"`
    Denom      string         `json:"denom"`
    Amount     sdk.Int        `json:"amount"`
}
```

## Send Coins From

`MsgProxySendCoinsFrom` initiates the coin transfer on behalf of the original coin owner.

`Proxy` account's signature is required.  

### MsgProxySendCoinsFrom
```golang
type MsgProxySendCoinsFrom struct {
    Proxy      sdk.AccAddress   `json:"proxy"`
    OnBehalfOf sdk.AccAddress   `json:"on_behalf_of"`
    Denom      string           `json:"denom"`
    Amount     sdk.Int          `json:"amount"`
    ToAddress  sdk.AccAddress   `json:"to_address"`
}
```