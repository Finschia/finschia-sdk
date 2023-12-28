package foundation

import (
	"github.com/gogo/protobuf/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Finschia/finschia-sdk/x/foundation/codec"
)

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgFundTreasury) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgWithdrawFromTreasury) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgUpdateMembers) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgUpdateDecisionPolicy) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgUpdateDecisionPolicy) GetDecisionPolicy() DecisionPolicy {
	if m.DecisionPolicy == nil {
		return nil
	}

	policy, ok := m.DecisionPolicy.GetCachedValue().(DecisionPolicy)
	if !ok {
		return nil
	}
	return policy
}

func (m *MsgUpdateDecisionPolicy) SetDecisionPolicy(policy DecisionPolicy) error {
	msg, ok := policy.(proto.Message)
	if !ok {
		return sdkerrors.ErrInvalidType.Wrapf("can't proto marshal %T", msg)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return err
	}
	m.DecisionPolicy = any

	return nil
}

func (m MsgUpdateDecisionPolicy) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var policy DecisionPolicy
	return unpacker.UnpackAny(m.DecisionPolicy, &policy)
}

func (m MsgSubmitProposal) GetMsgs() []sdk.Msg {
	msgs, err := GetMsgs(m.Messages, "proposal")
	if err != nil {
		panic(err)
	}
	return msgs
}

func (m *MsgSubmitProposal) SetMsgs(msgs []sdk.Msg) error {
	anys, err := SetMsgs(msgs)
	if err != nil {
		return err
	}
	m.Messages = anys
	return nil
}

func (m MsgSubmitProposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return UnpackInterfaces(unpacker, m.Messages)
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgSubmitProposal) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgWithdrawProposal) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgVote) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgExec) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgLeaveFoundation) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgUpdateCensorship) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgGrant) GetAuthorization() Authorization {
	if m.Authorization == nil {
		return nil
	}

	a, err := GetAuthorization(m.Authorization, "grant")
	if err != nil {
		return nil
	}
	return a
}

func (m *MsgGrant) SetAuthorization(a Authorization) error {
	any, err := SetAuthorization(a)
	if err != nil {
		return err
	}
	m.Authorization = any

	return nil
}

func (m MsgGrant) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var authorization Authorization
	return unpacker.UnpackAny(m.Authorization, &authorization)
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgGrant) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgRevoke) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}
