package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

var _ sdk.Msg = &MsgSwap{}

// ValidateBasic Implements Msg.
func (m *MsgSwap) ValidateBasic() error {
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

	if len(m.GetToDenom()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.FromCoinAmount.String())
	}

	return nil
}

// GetSigners Implements Msg.
func (m *MsgSwap) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

var _ sdk.Msg = &MsgSwapAll{}

// ValidateBasic Implements Msg.
func (m *MsgSwapAll) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err)
	}

	if len(m.GetFromDenom()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Invalid Denom")
	}

	if len(m.GetToDenom()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Invalid Denom")
	}

	return nil
}

// GetSigners Implements Msg.
func (m *MsgSwapAll) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
