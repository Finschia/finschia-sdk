package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ProxyDenom struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	OnBehalfOf sdk.AccAddress `json:"on_behalf_of"`
	Denom      string         `json:"denom"`
}

func NewProxyDenom(proxy, onBehalfOf sdk.AccAddress, denom string) ProxyDenom {
	return ProxyDenom{proxy, onBehalfOf, denom}
}

func (pd ProxyDenom) Equals(pd2 ProxyDenom) bool {
	return pd.Proxy.Equals(pd2.Proxy) &&
		pd.OnBehalfOf.Equals(pd2.OnBehalfOf) &&
		pd.Denom == pd2.Denom
}

func (pd ProxyDenom) String() string {
	return string(codec.MustMarshalJSONIndent(ModuleCdc, pd))
}

type ProxyAllowance struct {
	ProxyDenom
	Amount sdk.Int `json:"amount"`
}

func NewProxyAllowance(proxy, onBehalfOf sdk.AccAddress, denom string, amount sdk.Int) ProxyAllowance {
	return ProxyAllowance{NewProxyDenom(proxy, onBehalfOf, denom), amount}
}

func (pa ProxyAllowance) AddAllowance(pa2 ProxyAllowance) (ProxyAllowance, sdk.Error) {
	if !pa.HasSameProxyDenom(pa2) {
		return ProxyAllowance{}, ErrProxyDenomDoesNotMatch(DefaultCodespace, pa, pa2)
	}

	return NewProxyAllowance(pa.Proxy, pa.OnBehalfOf, pa.Denom, pa.Amount.Add(pa2.Amount)), nil
}

func (pa ProxyAllowance) SubAllowance(pa2 ProxyAllowance) (ProxyAllowance, sdk.Error) {
	hasEnoughCoinAmount, err := pa.GTE(pa2)
	if err != nil {
		return ProxyAllowance{}, err
	}
	if !hasEnoughCoinAmount {
		return ProxyAllowance{}, ErrProxyNotEnoughApprovedCoins(DefaultCodespace, pa.Amount, pa2.Amount)
	}

	return NewProxyAllowance(pa.Proxy, pa.OnBehalfOf, pa.Denom, pa.Amount.Sub(pa2.Amount)), nil
}

func (pa ProxyAllowance) GT(pa2 ProxyAllowance) (bool, sdk.Error) {
	if !pa.HasSameProxyDenom(pa2) {
		return false, ErrProxyDenomDoesNotMatch(DefaultCodespace, pa, pa2)
	}

	return pa.Amount.GT(pa2.Amount), nil
}

func (pa ProxyAllowance) GTE(pa2 ProxyAllowance) (bool, sdk.Error) {
	if !pa.HasSameProxyDenom(pa2) {
		return false, ErrProxyDenomDoesNotMatch(DefaultCodespace, pa, pa2)
	}

	return pa.Amount.GTE(pa2.Amount), nil
}

func (pa ProxyAllowance) LT(pa2 ProxyAllowance) (bool, sdk.Error) {
	if !pa.HasSameProxyDenom(pa2) {
		return false, ErrProxyDenomDoesNotMatch(DefaultCodespace, pa, pa2)
	}

	return pa.Amount.LT(pa2.Amount), nil
}

func (pa ProxyAllowance) LTE(pa2 ProxyAllowance) (bool, sdk.Error) {
	if !pa.HasSameProxyDenom(pa2) {
		return false, ErrProxyDenomDoesNotMatch(DefaultCodespace, pa, pa2)
	}

	return pa.Amount.LTE(pa2.Amount), nil
}

// all properties are the same except the amount
func (pa ProxyAllowance) HasSameProxyDenom(pa2 ProxyAllowance) bool {
	return pa.ProxyDenom.Equals(pa2.ProxyDenom)
}

func (pa ProxyAllowance) Equal(pa2 ProxyAllowance) bool {
	return pa.HasSameProxyDenom(pa2) && pa.Amount.Equal(pa2.Amount)
}

func (pa ProxyAllowance) String() string {
	return string(codec.MustMarshalJSONIndent(ModuleCdc, pa))
}
