package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

var _ sdk.Msg = &MsgStartChallenge{}

func (m MsgStartChallenge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid challenger address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(m.To)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid defender address (%s)", err)
	}

	if len(m.RollupName) == 0 {
		return ErrInvalidRollupName
	}

	if m.BlockHeight == 0 {
		return ErrInvalidL2BlockHeight
	}

	// 0 < step_count
	if m.StepCount == 0 {
		return ErrInvalidStepCount
	}

	return nil
}

func (m MsgStartChallenge) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func (m MsgStartChallenge) Type() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgStartChallenge) Route() string {
	return RouterKey
}

func (m MsgStartChallenge) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = &MsgNsectChallenge{}

func (m MsgNsectChallenge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid challenger address (%s)", err)
	}

	if len(m.ChallengeId) != 64 {
		return ErrInvalidChallengeID
	}

	if len(m.StateHashes) == 0 {
		return ErrInvalidStateHashes
	}

	return nil
}

func (m MsgNsectChallenge) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func (m MsgNsectChallenge) Type() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgNsectChallenge) Route() string {
	return RouterKey
}

func (m MsgNsectChallenge) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = &MsgFinishChallenge{}

func (m MsgFinishChallenge) ValidateBasic() error {
	panic("implement me")
}

func (m MsgFinishChallenge) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

func (m MsgFinishChallenge) Type() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgFinishChallenge) Route() string {
	return RouterKey
}

func (m MsgFinishChallenge) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
