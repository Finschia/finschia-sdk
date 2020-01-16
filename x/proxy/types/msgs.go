package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ json.Marshaler = (*MsgProxyAllowance)(nil)
var _ json.Unmarshaler = (*MsgProxyAllowance)(nil)
var _ json.Marshaler = (*MsgProxySendCoinsFrom)(nil)
var _ json.Unmarshaler = (*MsgProxySendCoinsFrom)(nil)
var _ json.Marshaler = (*MsgProxyApproveCoins)(nil)
var _ json.Unmarshaler = (*MsgProxyApproveCoins)(nil)
var _ json.Marshaler = (*MsgProxyDisapproveCoins)(nil)
var _ json.Unmarshaler = (*MsgProxyDisapproveCoins)(nil)

type MsgProxyDenom struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	OnBehalfOf sdk.AccAddress `json:"on_behalf_of"`
	Denom      string         `json:"denom"`
}

func NewMsgProxyDenom(proxy, onBehalfOf sdk.AccAddress, denom string) MsgProxyDenom {
	return MsgProxyDenom{proxy, onBehalfOf, denom}
}

type MsgProxyAllowance struct {
	MsgProxyDenom
	Amount sdk.Int `json:"amount"`
}

func NewMsgProxyAllowance(proxy, onBehalfOf sdk.AccAddress, denom string, amount sdk.Int) MsgProxyAllowance {
	return MsgProxyAllowance{NewMsgProxyDenom(proxy, onBehalfOf, denom), amount}
}

// wrapper for nested structure as a work around to amino/json serializer mismatch
func (msgProxyAllowance MsgProxyAllowance) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Proxy      sdk.AccAddress `json:"proxy"`
		OnBehalfOf sdk.AccAddress `json:"on_behalf_of"`
		Denom      string         `json:"denom"`
		Amount     sdk.Int        `json:"amount"`
	}{
		Proxy:      msgProxyAllowance.Proxy,
		OnBehalfOf: msgProxyAllowance.OnBehalfOf,
		Denom:      msgProxyAllowance.Denom,
		Amount:     msgProxyAllowance.Amount,
	})
}

func (msgProxyAllowance *MsgProxyAllowance) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgProxyAllowance
	return json.Unmarshal(data, msgAlias(msgProxyAllowance))
}

type MsgProxyApproveCoins struct {
	MsgProxyAllowance
}

func NewMsgProxyApproveCoins(proxy, onBehalfOf sdk.AccAddress, denom string, amount sdk.Int) MsgProxyApproveCoins {
	return MsgProxyApproveCoins{NewMsgProxyAllowance(proxy, onBehalfOf, denom, amount)}
}

// wrapper for nested structure as a work around to amino/json serializer mismatch
func (msgPxApproveCoins MsgProxyApproveCoins) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Proxy      sdk.AccAddress `json:"proxy"`
		OnBehalfOf sdk.AccAddress `json:"on_behalf_of"`
		Denom      string         `json:"denom"`
		Amount     sdk.Int        `json:"amount"`
	}{
		Proxy:      msgPxApproveCoins.Proxy,
		OnBehalfOf: msgPxApproveCoins.OnBehalfOf,
		Denom:      msgPxApproveCoins.Denom,
		Amount:     msgPxApproveCoins.Amount,
	})
}

func (msgPxApproveCoins *MsgProxyApproveCoins) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgProxyApproveCoins
	return json.Unmarshal(data, msgAlias(msgPxApproveCoins))
}

func (msgPxApproveCoins MsgProxyApproveCoins) Route() string { return RouterKey }

func (msgPxApproveCoins MsgProxyApproveCoins) Type() string { return MsgTypeProxyApproveCoins }

func (msgPxApproveCoins MsgProxyApproveCoins) ValidateBasic() sdk.Error {
	return validateMsgProxyAllowance(msgPxApproveCoins.MsgProxyAllowance)
}

func (msgPxApproveCoins MsgProxyApproveCoins) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgPxApproveCoins))
}

func (msgPxApproveCoins MsgProxyApproveCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgPxApproveCoins.OnBehalfOf}
}

type MsgProxyDisapproveCoins struct {
	MsgProxyAllowance
}

// wrapper for nested structure as a work around to amino/json serializer mismatch
func (msgPxDisapproveCoins MsgProxyDisapproveCoins) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Proxy      sdk.AccAddress `json:"proxy"`
		OnBehalfOf sdk.AccAddress `json:"on_behalf_of"`
		Denom      string         `json:"denom"`
		Amount     sdk.Int        `json:"amount"`
	}{
		Proxy:      msgPxDisapproveCoins.Proxy,
		OnBehalfOf: msgPxDisapproveCoins.OnBehalfOf,
		Denom:      msgPxDisapproveCoins.Denom,
		Amount:     msgPxDisapproveCoins.Amount,
	})
}

func (msgPxDisapproveCoins *MsgProxyDisapproveCoins) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgProxyDisapproveCoins
	return json.Unmarshal(data, msgAlias(msgPxDisapproveCoins))
}

func NewMsgProxyDisapproveCoins(proxy, onBehalfOf sdk.AccAddress, denom string, amount sdk.Int) MsgProxyDisapproveCoins {
	return MsgProxyDisapproveCoins{NewMsgProxyAllowance(proxy, onBehalfOf, denom, amount)}
}

func (msgPxDisapproveCoins MsgProxyDisapproveCoins) Route() string { return RouterKey }

func (msgPxDisapproveCoins MsgProxyDisapproveCoins) Type() string { return MsgTypeProxyDisapproveCoins }

func (msgPxDisapproveCoins MsgProxyDisapproveCoins) ValidateBasic() sdk.Error {
	return validateMsgProxyAllowance(msgPxDisapproveCoins.MsgProxyAllowance)
}

func (msgPxDisapproveCoins MsgProxyDisapproveCoins) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgPxDisapproveCoins))
}

func (msgPxDisapproveCoins MsgProxyDisapproveCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgPxDisapproveCoins.OnBehalfOf}
}

type MsgProxySendCoinsFrom struct {
	MsgProxyAllowance
	ToAddress sdk.AccAddress `json:"to_address"`
}

// wrapper for nested structure as a work around to amino/json serializer mismatch
func (msgPxSendCoinsFrom MsgProxySendCoinsFrom) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Proxy      sdk.AccAddress `json:"proxy"`
		OnBehalfOf sdk.AccAddress `json:"on_behalf_of"`
		Denom      string         `json:"denom"`
		Amount     sdk.Int        `json:"amount"`
		ToAddress  sdk.AccAddress `json:"to_address"`
	}{
		Proxy:      msgPxSendCoinsFrom.Proxy,
		OnBehalfOf: msgPxSendCoinsFrom.OnBehalfOf,
		Denom:      msgPxSendCoinsFrom.Denom,
		Amount:     msgPxSendCoinsFrom.Amount,
		ToAddress:  msgPxSendCoinsFrom.ToAddress,
	})
}

func (msgPxSendCoinsFrom *MsgProxySendCoinsFrom) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgProxySendCoinsFrom
	return json.Unmarshal(data, msgAlias(msgPxSendCoinsFrom))
}

func NewMsgProxySendCoinsFrom(proxy, onBehalfOf, to sdk.AccAddress, denom string, amount sdk.Int) MsgProxySendCoinsFrom {
	return MsgProxySendCoinsFrom{NewMsgProxyAllowance(proxy, onBehalfOf, denom, amount), to}
}

func (msgPxSendCoinsFrom MsgProxySendCoinsFrom) Route() string { return RouterKey }

func (msgPxSendCoinsFrom MsgProxySendCoinsFrom) Type() string { return MsgTypeProxySendCoinsFrom }

func (msgPxSendCoinsFrom MsgProxySendCoinsFrom) ValidateBasic() sdk.Error {
	if msgPxSendCoinsFrom.ToAddress.Empty() {
		return ErrProxyToAddressRequired(DefaultCodespace)
	}

	return validateMsgProxyAllowance(msgPxSendCoinsFrom.MsgProxyAllowance)
}

func (msgPxSendCoinsFrom MsgProxySendCoinsFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgPxSendCoinsFrom))
}

// SendCoinsFrom is signed by the proxy
func (msgPxSendCoinsFrom MsgProxySendCoinsFrom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msgPxSendCoinsFrom.Proxy}
}

func validateMsgProxyAllowance(msg MsgProxyAllowance) sdk.Error {
	if msg.Proxy.Empty() {
		return ErrProxyAddressRequired(DefaultCodespace)
	}
	if msg.OnBehalfOf.Empty() {
		return ErrProxyOnBehalfOfAddressRequired(DefaultCodespace)
	}
	if len(msg.Denom) == 0 {
		return ErrProxyDenomRequired(DefaultCodespace)
	}
	if msg.Amount.LT(sdk.NewInt(1)) {
		return ErrProxyAmountMustBePositiveInteger(DefaultCodespace)
	}
	return nil
}
