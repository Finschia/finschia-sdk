package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ json.Marshaler = (*ProxyAllowance)(nil)
var _ json.Unmarshaler = (*ProxyAllowance)(nil)

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

func NewProxyAllowance(pxd ProxyDenom, amount sdk.Int) ProxyAllowance {
	return ProxyAllowance{pxd, amount}
}

// wrapper for nested structure as a work around to amino/json serializer mismatch
func (pa ProxyAllowance) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Proxy      sdk.AccAddress `json:"proxy"`
		OnBehalfOf sdk.AccAddress `json:"on_behalf_of"`
		Denom      string         `json:"denom"`
		Amount     sdk.Int        `json:"amount"`
	}{
		Proxy:      pa.Proxy,
		OnBehalfOf: pa.OnBehalfOf,
		Denom:      pa.Denom,
		Amount:     pa.Amount,
	})
}

func (pa *ProxyAllowance) UnmarshalJSON(data []byte) error {
	type msgAlias *ProxyAllowance
	return json.Unmarshal(data, msgAlias(pa))
}

func (pa ProxyAllowance) AddAllowance(pa2 ProxyAllowance) (ProxyAllowance, error) {
	if !pa.HasSameProxyDenom(pa2) {
		return ProxyAllowance{}, sdkerrors.Wrapf(ErrProxyDenomDoesNotMatch, "%v compared to %v", pa.ProxyDenom, pa2.ProxyDenom)
	}

	return NewProxyAllowance(pa.ProxyDenom, pa.Amount.Add(pa2.Amount)), nil
}

func (pa ProxyAllowance) SubAllowance(pa2 ProxyAllowance) (ProxyAllowance, error) {
	hasEnoughCoinAmount, err := pa.GTE(pa2)
	if err != nil {
		return ProxyAllowance{}, err
	}
	if !hasEnoughCoinAmount {
		return ProxyAllowance{}, sdkerrors.Wrapf(ErrProxyNotEnoughApprovedCoins, "Approved: %v, Requested: %v", pa.Amount, pa2.Amount)
	}

	return NewProxyAllowance(pa.ProxyDenom, pa.Amount.Sub(pa2.Amount)), nil
}

func (pa ProxyAllowance) GT(pa2 ProxyAllowance) (bool, error) {
	if !pa.HasSameProxyDenom(pa2) {
		return false, sdkerrors.Wrapf(ErrProxyDenomDoesNotMatch, "%v compared to %v", pa.ProxyDenom, pa2.ProxyDenom)
	}

	return pa.Amount.GT(pa2.Amount), nil
}

func (pa ProxyAllowance) GTE(pa2 ProxyAllowance) (bool, error) {
	if !pa.HasSameProxyDenom(pa2) {
		return false, sdkerrors.Wrapf(ErrProxyDenomDoesNotMatch, "%v compared to %v", pa.ProxyDenom, pa2.ProxyDenom)
	}

	return pa.Amount.GTE(pa2.Amount), nil
}

func (pa ProxyAllowance) LT(pa2 ProxyAllowance) (bool, error) {
	if !pa.HasSameProxyDenom(pa2) {
		return false, sdkerrors.Wrapf(ErrProxyDenomDoesNotMatch, "%v compared to %v", pa.ProxyDenom, pa2.ProxyDenom)
	}

	return pa.Amount.LT(pa2.Amount), nil
}

func (pa ProxyAllowance) LTE(pa2 ProxyAllowance) (bool, error) {
	if !pa.HasSameProxyDenom(pa2) {
		return false, sdkerrors.Wrapf(ErrProxyDenomDoesNotMatch, "%v compared to %v", pa.ProxyDenom, pa2.ProxyDenom)
	}

	return pa.Amount.LTE(pa2.Amount), nil
}

func (pa ProxyAllowance) HasSameProxyDenom(pa2 ProxyAllowance) bool {
	return pa.ProxyDenom.Equals(pa2.ProxyDenom)
}

func (pa ProxyAllowance) Equal(pa2 ProxyAllowance) bool {
	return pa.HasSameProxyDenom(pa2) && pa.Amount.Equal(pa2.Amount)
}

func (pa ProxyAllowance) String() string {
	return string(codec.MustMarshalJSONIndent(ModuleCdc, pa))
}
