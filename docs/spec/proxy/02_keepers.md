# Keepers

## Common Types

### Proxy

```golang
type Proxy struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	OnBehalfOf sdk.AccAddress `json:"onBehalfOf"`
	Action     string         `json:"action"`
	Denom      string         `json:"denom"`
	Amount     sdk.Int        `json:"amount"`
}
```


## Keeper

```golang
type ProxyKeeper interface {
    ApproveCoins(ctx sdk.Context, msg types.MsgProxyApproveCoins) sdk.Error
    DisapproveCoins(ctx sdk.Context, msg types.MsgProxyDisapproveCoins) sdk.Error
    SendCoinsFrom(ctx sdk.Context, msg types.MsgProxySendCoinsFrom) sdk.Error
}
```
