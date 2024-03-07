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

func ValidateZkAuthSignature(signature ZKAuthSignature) error {
	if signature.ZkAuthInputs == nil {
		return sdkerrors.Wrap(ErrInvalidZKAuthSignature, "ZkAuthInputs is empty")
	}

	if err := signature.ZkAuthInputs.Validate(); err != nil {
		return err
	}

	if signature.MaxBlockHeight == 0 {
		return sdkerrors.Wrapf(ErrInvalidZKAuthSignature, "invalid max_block_height %d", signature.MaxBlockHeight)
	}

	return nil
}

func (e *MsgExecution) SetMsgs(msgs []sdk.Msg) error {
	anys := make([]*cdctypes.Any, len(msgs))
	for i, msg := range msgs {
		var err error
		anys[i], err = cdctypes.NewAnyWithValue(msg)
		if err != nil {
			return err
		}
	}
	e.Msgs = anys
	return nil
}

func (e *MsgExecution) UnpackInterfaces(unpacker cdctypes.AnyUnpacker) error {
	for _, any := range e.GetMsgs() {
		var msg sdk.Msg
		err := unpacker.UnpackAny(any, &msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *MsgExecution) ValidateBasic() error {
	if len(e.GetMsgs()) == 0 {
		return sdkerrors.Wrap(ErrInvalidMessage, "message is empty")
	}

	// validate msg
	for _, m := range e.GetMsgs() {
		message, ok := m.GetCachedValue().(sdk.Msg)
		if !ok {
			return sdkerrors.Wrapf(ErrInvalidMessage, "message contains %T which is not a sdk.MsgRequest", m)
		}
		if err := message.ValidateBasic(); err != nil {
			return sdkerrors.Wrap(ErrInvalidMessage, err.Error())
		}
	}

	// validate signature
	if err := ValidateZkAuthSignature(e.ZkAuthSignature); err != nil {
		return err
	}

	return nil
}

func (e *MsgExecution) GetMessages() ([]sdk.Msg, error) {
	msgs := make([]sdk.Msg, len(e.Msgs))
	for i, msgAny := range e.Msgs {
		msg, ok := msgAny.GetCachedValue().(sdk.Msg)
		if !ok {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "messages contains %T which is not a sdk.MsgRequest", msgAny)
		}
		msgs[i] = msg
	}

	return msgs, nil
}

func (e *MsgExecution) GetSigners() []sdk.AccAddress {
	addr, err := e.ZkAuthSignature.ZkAuthInputs.AccAddress()
	if err != nil {
		return nil
	}

	return []sdk.AccAddress{addr}
}

func (e *MsgExecution) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(e)
	return sdk.MustSortJSON(bz)
}

func (e *MsgExecution) Route() string { return RouterKey }
func (e *MsgExecution) Type() string  { return sdk.MsgTypeURL(e) }
