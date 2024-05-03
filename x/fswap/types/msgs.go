package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/fswap/codec"
)

var _ sdk.Msg = &MsgSwap{}

// ValidateBasic Implements Msg.
func (m *MsgSwap) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid address (%s)", err)
	}

	if !m.FromCoinAmount.IsValid() {
		return sdkerrors.ErrInvalidCoins.Wrap(m.FromCoinAmount.String())
	}

	if !m.FromCoinAmount.IsPositive() {
		return sdkerrors.ErrInvalidCoins.Wrap(m.FromCoinAmount.String())
	}

	if len(m.GetToDenom()) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid denom (empty denom)")
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

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m *MsgSwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(m))
}

var _ sdk.Msg = &MsgSwapAll{}

// ValidateBasic Implements Msg.
func (m *MsgSwapAll) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid address (%s)", err)
	}

	if len(m.GetFromDenom()) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid denom (empty denom)")
	}

	if len(m.GetToDenom()) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid denom (empty denom)")
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

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m *MsgSwapAll) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(m))
}
