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
	if err := signature.ZkAuthInputs.Validate(); err != nil {
		return err
	}

	if signature.MaxBlockHeight == 0 {
		return sdkerrors.Wrapf(ErrInvalidZKAuthSignature, "invalid max_block_height %d", signature.MaxBlockHeight)
	}

	return nil
}

func (msg *MsgExecution) ValidateBasic() error {
	// validate msg
	for _, msg := range msg.Msgs {
		message, ok := msg.GetCachedValue().(sdk.Msg)
		if !ok {
			return ErrInvalidMessage
		}
		if err := message.ValidateBasic(); err != nil {
			return sdkerrors.Wrap(ErrInvalidMessage, err.Error())
		}
	}

	// validate signature
	if err := ValidateZkAuthSignature(msg.ZkAuthSignature); err != nil {
		return err
	}

	return nil
}

func (msg *MsgExecution) GetSigners() []sdk.AccAddress {
	// TODO:
	return []sdk.AccAddress{}
}

func (msg *MsgExecution) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgExecution) Route() string { return RouterKey }
func (msg *MsgExecution) Type() string  { return sdk.MsgTypeURL(msg) }
