package types

import (
	cdctypes "github.com/Finschia/finschia-sdk/codec/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgExecution)(nil)
)

func NewMsgExecution(msgs []sdk.Msg, zkauthSignature ZKAuthSignature) *MsgExecution {
	msgsAny := make([]*cdctypes.Any, len(msgs))
	for i, msg := range msgs {
		any, err := cdctypes.NewAnyWithValue(msg)
		if err != nil {
			panic(err)
		}

		msgsAny[i] = any
	}

	return &MsgExecution{
		Msgs:            msgsAny,
		ZkAuthSignature: zkauthSignature,
	}
}

func (msg MsgExecution) GetMessages() ([]sdk.Msg, error) {
	msgs := make([]sdk.Msg, len(msg.Msgs))
	for i, msgAny := range msg.Msgs {
		msg, ok := msgAny.GetCachedValue().(sdk.Msg)
		if !ok {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "messages contains %T which is not a sdk.MsgRequest", msgAny)
		}
		msgs[i] = msg
	}

	return msgs, nil
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
