package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
)

var _ sdk.Msg = &MsgStartChallenge{}

func (m MsgStartChallenge) ValidateBasic() error {
	panic("implement me")
}

func (m MsgStartChallenge) GetSigners() []sdk.AccAddress {
	panic("implement me")
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
	panic("implement me")
}

func (m MsgNsectChallenge) GetSigners() []sdk.AccAddress {
	panic("implement me")
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
