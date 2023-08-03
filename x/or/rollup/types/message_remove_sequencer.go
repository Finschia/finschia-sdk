package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
)

var _ sdk.Msg = (*MsgRemoveSequencer)(nil)

func NewMsgRemoveSequencer() (*MsgRemoveSequencer, error) {
	return &MsgRemoveSequencer{}, nil
}

// ValidateBasic implements Msg.
func (m MsgRemoveSequencer) ValidateBasic() error {
	panic("implement me")
}

// GetSigners implements Msg
func (m MsgRemoveSequencer) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

// Type implements the LegacyMsg.Type method.
func (m MsgRemoveSequencer) Type() string {
	panic("implement me")
}

// Route implements the LegacyMsg.Route method.
func (m MsgRemoveSequencer) Route() string {
	panic("implement me")
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgRemoveSequencer) GetSignBytes() []byte {
	panic("implement me")
}
