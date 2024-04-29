package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

var _ sdk.Msg = &MsgSwapRequest{}

// NewMsgSwapRequest - construct a msg to swap amounts of old coin to new coin
//
//nolint:interfacer
func NewMsgSwapRequest(fromAddr, toAddr sdk.AccAddress, amount sdk.Coin) *MsgSwapRequest {
	return &MsgSwapRequest{FromAddress: fromAddr.String(), Amount: amount}
}

// ValidateBasic Implements Msg.
func (m *MsgSwapRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err)
	}

	if !m.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	return nil
}

// GetSigners Implements Msg.
func (m *MsgSwapRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

var _ sdk.Msg = &MsgSwapAllRequest{}

// NewMsgSwapRequest - construct a msg to swap all old coin to new coin
//
//nolint:interfacer
func NewMsgSwapAllRequest(fromAddr, toAddr sdk.AccAddress) *MsgSwapAllRequest {
	return &MsgSwapAllRequest{FromAddress: fromAddr.String()}
}

// ValidateBasic Implements Msg.
func (m *MsgSwapAllRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err)
	}

	return nil
}

// GetSigners Implements Msg.
func (m *MsgSwapAllRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
