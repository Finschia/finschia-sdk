package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
)

var _ sdk.Msg = (*MsgApprove)(nil)

type MsgApprove struct {
	Approver sdk.AccAddress `json:"approver"`
	Proxy    sdk.AccAddress `json:"proxy"`
	Symbol   string         `json:"symbol"`
}

func NewMsgApprove(approver sdk.AccAddress, proxy sdk.AccAddress, symbol string) MsgApprove {
	return MsgApprove{
		Approver: approver,
		Proxy:    proxy,
		Symbol:   symbol,
	}
}

func (msg MsgApprove) ValidateBasic() sdk.Error {
	if msg.Approver.Empty() {
		return sdk.ErrInvalidAddress("Approver cannot be empty")
	}
	if msg.Proxy.Empty() {
		return sdk.ErrInvalidAddress("Proxy cannot be empty")
	}
	if msg.Approver.Equals(msg.Proxy) {
		return ErrApproverProxySame(DefaultCodespace, msg.Approver.String())
	}
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	return nil
}

func (MsgApprove) Route() string { return RouterKey }
func (MsgApprove) Type() string  { return "approve_collection" }
func (msg MsgApprove) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Approver}
}
func (msg MsgApprove) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

var _ sdk.Msg = (*MsgDisapprove)(nil)

type MsgDisapprove struct {
	Approver sdk.AccAddress `json:"approver"`
	Proxy    sdk.AccAddress `json:"proxy"`
	Symbol   string         `json:"symbol"`
}

func NewMsgDisapprove(approver sdk.AccAddress, proxy sdk.AccAddress, symbol string) MsgDisapprove {
	return MsgDisapprove{
		Approver: approver,
		Proxy:    proxy,
		Symbol:   symbol,
	}
}

func (msg MsgDisapprove) ValidateBasic() sdk.Error {
	if msg.Approver.Empty() {
		return sdk.ErrInvalidAddress("Approver cannot be empty")
	}
	if msg.Proxy.Empty() {
		return sdk.ErrInvalidAddress("Proxy cannot be empty")
	}
	if msg.Approver.Equals(msg.Proxy) {
		return ErrApproverProxySame(DefaultCodespace, msg.Approver.String())
	}
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	return nil
}

func (MsgDisapprove) Route() string { return RouterKey }
func (MsgDisapprove) Type() string  { return "disapprove_collection" }
func (msg MsgDisapprove) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Approver}
}
func (msg MsgDisapprove) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
