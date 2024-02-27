package types

import (
	types "github.com/Finschia/finschia-sdk/codec/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
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
	anys := make([]*types.Any, len(msgs))
	for i, msg := range msgs {
		var err error
		anys[i], err = types.NewAnyWithValue(msg)
		if err != nil {
			return err
		}
	}
	e.Msgs = anys
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

func (e *MsgExecution) GetSigners() []sdk.AccAddress {
	// TODO:
	return []sdk.AccAddress{}
}

func (e *MsgExecution) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(e)
	return sdk.MustSortJSON(bz)
}

func (e *MsgExecution) Route() string { return RouterKey }
func (e *MsgExecution) Type() string  { return sdk.MsgTypeURL(e) }
