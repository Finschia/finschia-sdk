package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgCreateAccount - create account transaction of the account module
type MsgCreateAccount struct {
	FromAddress   sdk.AccAddress `json:"from_address" yaml:"from_address"`
	TargetAddress sdk.AccAddress `json:"target_address" yaml:"target_address"`
}

var _ sdk.Msg = MsgCreateAccount{}

// NewMsgCreateAccount - construct create account msg.
func NewMsgCreateAccount(fromAddr, targetAddr sdk.AccAddress) MsgCreateAccount {
	return MsgCreateAccount{FromAddress: fromAddr, TargetAddress: targetAddr}
}

// Route Implements Msg.
func (msg MsgCreateAccount) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCreateAccount) Type() string { return MsgTypeCreateAccount }

// ValidateBasic Implements Msg.
func (msg MsgCreateAccount) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("missing signer address")
	}
	if msg.TargetAddress.Empty() {
		return sdk.ErrInvalidAddress("missing target address to be created")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgCreateAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}
