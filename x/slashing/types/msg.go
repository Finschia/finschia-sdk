package types

import (
	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/types/errors"
)

// slashing message types
const (
	TypeMsgUnjail = "unjail"
)

// verify interface at compile time
var _ sdk.Msg = &MsgUnjail{}

// NewMsgUnjail creates a new MsgUnjail instance
//nolint:interfacer
func NewMsgUnjail(validatorAddr sdk.ValAddress) *MsgUnjail {
	return &MsgUnjail{
		ValidatorAddr: validatorAddr.String(),
	}
}

func (msg MsgUnjail) Route() string { return RouterKey }
func (msg MsgUnjail) Type() string  { return TypeMsgUnjail }
func (msg MsgUnjail) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.ValAddress(msg.ValidatorAddr).ToAccAddress()}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgUnjail) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgUnjail) ValidateBasic() error {
	if msg.ValidatorAddr == "" {
		return ErrBadValidatorAddr
	}
	if err := sdk.ValidateValAddress(msg.ValidatorAddr); err != nil {
		return errors.Wrapf(errors.ErrInvalidAddress, "Invalid validator address (%s)", err)
	}

	return nil
}
