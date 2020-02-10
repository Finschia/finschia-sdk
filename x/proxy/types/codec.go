package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(ProxyAllowance{}, "proxy/ProxyAllowance", nil)
	cdc.RegisterConcrete(MsgProxyDenom{}, "proxy/MsgProxyDenom", nil)
	cdc.RegisterConcrete(MsgProxyAllowance{}, "proxy/MsgProxyAllowance", nil)
	cdc.RegisterConcrete(MsgProxyApproveCoins{}, "proxy/MsgProxyApproveCoins", nil)
	cdc.RegisterConcrete(MsgProxyDisapproveCoins{}, "proxy/MsgProxyDisapproveCoins", nil)
	cdc.RegisterConcrete(MsgProxySendCoinsFrom{}, "proxy/MsgProxySendCoinsFrom", nil)
}
