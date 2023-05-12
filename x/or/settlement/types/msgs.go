package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
)

var _ sdk.Msg = &MsgAddTrieNode{}

func (m MsgAddTrieNode) ValidateBasic() error {
	panic("implement me")
}

func (m MsgAddTrieNode) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

func (m MsgAddTrieNode) Type() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgAddTrieNode) Route() string {
	return RouterKey
}

func (m MsgAddTrieNode) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = &MsgConfirmStateTransition{}

func (m MsgConfirmStateTransition) ValidateBasic() error {
	panic("implement me")
}

func (m MsgConfirmStateTransition) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

func (m MsgConfirmStateTransition) Type() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgConfirmStateTransition) Route() string {
	return RouterKey
}

func (m MsgConfirmStateTransition) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = &MsgDenyStateTransition{}

func (m MsgDenyStateTransition) ValidateBasic() error {
	panic("implement me")
}

func (m MsgDenyStateTransition) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

func (m MsgDenyStateTransition) Type() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgDenyStateTransition) Route() string {
	return RouterKey
}

func (m MsgDenyStateTransition) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = &MsgInitiateChallenge{}

func (m MsgInitiateChallenge) ValidateBasic() error {
	panic("implement me")
}

func (m MsgInitiateChallenge) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

func (m MsgInitiateChallenge) Type() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgInitiateChallenge) Route() string {
	return RouterKey
}

func (m MsgInitiateChallenge) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = &MsgProposeState{}

func (m MsgProposeState) ValidateBasic() error {
	panic("implement me")
}

func (m MsgProposeState) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

func (m MsgProposeState) Type() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgProposeState) Route() string {
	return RouterKey
}

func (m MsgProposeState) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

var _ sdk.Msg = &MsgRespondState{}

func (m MsgRespondState) ValidateBasic() error {
	panic("implement me")
}

func (m MsgRespondState) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

func (m MsgRespondState) Type() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgRespondState) Route() string {
	return RouterKey
}

func (m MsgRespondState) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
