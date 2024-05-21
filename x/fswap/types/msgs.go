package types

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/foundation/codec"
)

var (
	_ sdk.Msg = &MsgSwap{}
	_ sdk.Msg = &MsgSwapAll{}
	_ sdk.Msg = &MsgSetSwap{}
)

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

	if err := sdk.ValidateDenom(m.GetToDenom()); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrap(err.Error())
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
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// Type implements the LegacyMsg.Type method.
func (m *MsgSwap) Type() string {
	return sdk.MsgTypeURL(m)
}

// Route implements the LegacyMsg.Route method.
func (m *MsgSwap) Route() string {
	return RouterKey
}

func (m *MsgSwap) GetSignerBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic Implements Msg.
func (m *MsgSwapAll) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("Invalid address (%s)", err)
	}

	if err := sdk.ValidateDenom(m.FromDenom); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrap(err.Error())
	}

	if err := sdk.ValidateDenom(m.ToDenom); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrap(err.Error())
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

// Type implements the LegacyMsg.Type method.
func (m *MsgSwapAll) Type() string {
	return sdk.MsgTypeURL(m)
}

// Route implements the LegacyMsg.Route method.
func (m *MsgSwapAll) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m *MsgSwapAll) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgSetSwap) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", m.Authority)
	}

	if err := m.Swap.ValidateBasic(); err != nil {
		return err
	}
	if err := m.ToDenomMetadata.Validate(); err != nil {
		return err
	}
	if m.Swap.ToDenom != m.ToDenomMetadata.Base {
		return sdkerrors.ErrInvalidRequest.Wrapf("denomination does not match %s != %s", m.Swap.ToDenom, m.ToDenomMetadata.Base)
	}

	return nil
}

func (m *MsgSetSwap) GetSigners() []sdk.AccAddress {
	signer := sdk.MustAccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m *MsgSetSwap) Type() string {
	return sdk.MsgTypeURL(m)
}

// Route implements the LegacyMsg.Route method.
func (m *MsgSetSwap) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m *MsgSetSwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *QuerySwapRequest) Validate() error {
	if m == nil {
		return status.Error(codes.InvalidArgument, "empty QuerySwapRequest is not allowed")
	}
	if m.GetFromDenom() == m.GetToDenom() {
		return status.Error(codes.InvalidArgument, "fromDenom and toDenom cannot be the same")
	}

	if m.GetFromDenom() == "" {
		return status.Error(codes.InvalidArgument, "from denom is empty")
	}
	if m.GetToDenom() == "" {
		return status.Error(codes.InvalidArgument, "to denom is empty")
	}
	return nil
}
