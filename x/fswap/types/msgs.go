package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

var _ sdk.Msg = &MsgSwapRequest{}

// ValidateBasic Implements Msg.
func (m *MsgSwapRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err)
	}

	if !m.FromCoinAmount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.FromCoinAmount.String())
	}

	if !m.FromCoinAmount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.FromCoinAmount.String())
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
