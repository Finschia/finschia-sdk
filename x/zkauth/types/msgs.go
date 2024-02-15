package types

import (
	types "github.com/Finschia/finschia-sdk/codec/types"
	sdk "github.com/Finschia/finschia-sdk/types"
)

var (
	_ sdk.Msg = (*MsgExecution)(nil)
)

func NewMsgExecution(msgs []*types.Any, zkauthSignature ZKAuthSignature) *MsgExecution {
	return &MsgExecution{
		Msgs:            msgs,
		ZkAuthSignature: zkauthSignature,
	}
}

func (msg MsgExecution) GetSigners() []sdk.AccAddress {
	// TODO:
	return []sdk.AccAddress{}
}

func (msg MsgExecution) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgExecution) ValidateBasic() error { return nil }
func (msg MsgExecution) Route() string        { return RouterKey }
func (msg MsgExecution) Type() string         { return sdk.MsgTypeURL(&msg) }
