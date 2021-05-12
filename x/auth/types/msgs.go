package types

import (
	sdk "github.com/line/lbm-sdk/v2/types"
	sdkerrors "github.com/line/lbm-sdk/v2/types/errors"
)

// auth message types
const (
	TypeMsgEmpty = "empty"
)

var _ sdk.Msg = &MsgEmpty{}

// NewMsgEmpty creates a new MsgEmpty object
//nolint:interfacer
func NewMsgEmpty(fromAddr sdk.AccAddress) *MsgEmpty {
	return &MsgEmpty{FromAddress: fromAddr.String()}
}

func NewServiceMsgEmpty(fromAddr sdk.AccAddress) sdk.ServiceMsg {
	return sdk.ServiceMsg{
		MethodName: "/lbm.auth.v1beta1.Msg/Empty",
		Request:    NewMsgEmpty(fromAddr),
	}
}

func (msg MsgEmpty) Route() string { return ModuleName }

func (msg MsgEmpty) Type() string { return TypeMsgEmpty }

func (msg MsgEmpty) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgEmpty) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgEmpty) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
