package types

import (
	"github.com/gogo/protobuf/proto"

	"github.com/line/lbm-sdk/codec/legacy"
	"github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

var (
	_, _ sdk.Msg = &MsgGrantFeeAllowance{}, &MsgRevokeFeeAllowance{}

	_ types.UnpackInterfacesMessage = &MsgGrantFeeAllowance{}
)

// NewMsgGrantFeeAllowance creates a new MsgGrantFeeAllowance.
//nolint:interfacer
func NewMsgGrantFeeAllowance(feeAllowance FeeAllowanceI, granter, grantee sdk.AccAddress) (*MsgGrantFeeAllowance, error) {
	msg, ok := feeAllowance.(proto.Message)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPackAny, "cannot proto marshal %T", msg)
	}
	any, err := types.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return &MsgGrantFeeAllowance{
		Granter:   granter.String(),
		Grantee:   grantee.String(),
		Allowance: any,
	}, nil
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgGrantFeeAllowance) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(msg.Granter); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid granter address: %s", err)
	}
	if err := sdk.ValidateAccAddress(msg.Grantee); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid grantee address: %s", err)
	}
	if msg.Grantee == msg.Granter {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "cannot self-grant fee authorization")
	}
	allowance, err := msg.GetFeeAllowanceI()
	if err != nil {
		return err
	}

	return allowance.ValidateBasic()
}

// GetSigners gets the granter account associated with an allowance
func (msg MsgGrantFeeAllowance) GetSigners() []sdk.AccAddress {
	granter := sdk.AccAddress(msg.Granter)
	return []sdk.AccAddress{granter}
}

// Type implements the LegacyMsg.Type method.
func (msg MsgGrantFeeAllowance) Type() string {
	return sdk.MsgTypeURL(&msg)
}

// Route implements the LegacyMsg.Route method.
func (msg MsgGrantFeeAllowance) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (msg MsgGrantFeeAllowance) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&msg))
}

// GetFeeAllowanceI returns unpacked FeeAllowance
func (msg MsgGrantFeeAllowance) GetFeeAllowanceI() (FeeAllowanceI, error) {
	allowance, ok := msg.Allowance.GetCachedValue().(FeeAllowanceI)
	if !ok {
		return nil, sdkerrors.Wrap(ErrNoAllowance, "failed to get allowance")
	}

	return allowance, nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgGrantFeeAllowance) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var allowance FeeAllowanceI
	return unpacker.UnpackAny(msg.Allowance, &allowance)
}

// NewMsgRevokeFeeAllowance returns a message to revoke a fee allowance for a given
// granter and grantee
//nolint:interfacer
func NewMsgRevokeFeeAllowance(granter sdk.AccAddress, grantee sdk.AccAddress) MsgRevokeFeeAllowance {
	return MsgRevokeFeeAllowance{Granter: granter.String(), Grantee: grantee.String()}
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgRevokeFeeAllowance) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(msg.Granter); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid granter address: %s", err)
	}
	if err := sdk.ValidateAccAddress(msg.Grantee); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid grantee address: %s", err)
	}
	if msg.Grantee == msg.Granter {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "addresses must be different")
	}

	return nil
}

// GetSigners gets the granter address associated with an Allowance
// to revoke.
func (msg MsgRevokeFeeAllowance) GetSigners() []sdk.AccAddress {
	granter := sdk.AccAddress(msg.Granter)
	return []sdk.AccAddress{granter}
}

// Type implements the LegacyMsg.Type method.
func (msg MsgRevokeFeeAllowance) Type() string {
	return sdk.MsgTypeURL(&msg)
}

// Route implements the LegacyMsg.Route method.
func (msg MsgRevokeFeeAllowance) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (msg MsgRevokeFeeAllowance) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&msg))
}
