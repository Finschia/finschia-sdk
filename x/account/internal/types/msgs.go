package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgEmpty struct {
	From sdk.AccAddress `json:"from"`
}

var _ sdk.Msg = MsgEmpty{}

// NewMsgCreateAccount - construct create account msg.
func NewMsgEmpty(from sdk.AccAddress) MsgEmpty {
	return MsgEmpty{From: from}
}

// Route Implements Msg.
func (msg MsgEmpty) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgEmpty) Type() string { return MsgTypeEmpty }

// ValidateBasic Implements Msg.
func (msg MsgEmpty) ValidateBasic() error {
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing signer address")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgEmpty) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgEmpty) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// MsgCreateAccount - create account transaction of the account module
type MsgCreateAccount struct {
	From   sdk.AccAddress `json:"from"`
	Target sdk.AccAddress `json:"target"`
}

var _ sdk.Msg = MsgCreateAccount{}

// NewMsgCreateAccount - construct create account msg.
func NewMsgCreateAccount(fromAddr, targetAddr sdk.AccAddress) MsgCreateAccount {
	return MsgCreateAccount{From: fromAddr, Target: targetAddr}
}

// Route Implements Msg.
func (msg MsgCreateAccount) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCreateAccount) Type() string { return MsgTypeCreateAccount }

// ValidateBasic Implements Msg.
func (msg MsgCreateAccount) ValidateBasic() error {
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing signer address")
	}
	if msg.Target.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing target address to be created")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgCreateAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}
