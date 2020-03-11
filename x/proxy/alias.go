package proxy

import (
	"github.com/line/link/x/proxy/keeper"
	"github.com/line/link/x/proxy/types"
)

const (
	ModuleName = types.ModuleName
	StoreKey   = types.StoreKey
	RouterKey  = types.RouterKey
)

type (
	Keeper    = keeper.Keeper
	Allowance = types.ProxyAllowance
)

var (
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
	NewKeeper     = keeper.NewKeeper

	EventProxyApproveCoins    = types.EventProxyApproveCoins
	EventProxyDisapproveCoins = types.EventProxyDisapproveCoins
	EventProxySendCoinsFrom   = types.EventProxySendCoinsFrom

	ErrProxyNotExist               = types.ErrProxyNotExist
	ErrProxyNotEnoughApprovedCoins = types.ErrProxyNotEnoughApprovedCoins

	AttributeKeyProxyAddress           = types.AttributeKeyProxyAddress
	AttributeKeyProxyOnBehalfOfAddress = types.AttributeKeyProxyOnBehalfOfAddress
	AttributeKeyProxyToAddress         = types.AttributeKeyProxyToAddress
	AttributeKeyProxyDenom             = types.AttributeKeyProxyDenom
	AttributeKeyProxyAmount            = types.AttributeKeyProxyAmount

	AttributeValueCategory = types.AttributeValueCategory
)
