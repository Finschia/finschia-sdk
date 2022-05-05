package foundation

import (
	"github.com/line/lbm-sdk/codec/legacy"
	"github.com/line/lbm-sdk/x/authz"

	"github.com/gogo/protobuf/proto"

	codectypes "github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

var _ sdk.Msg = (*MsgFundTreasury)(nil)

// Route implements Msg.
func (m MsgFundTreasury) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgFundTreasury) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgFundTreasury) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if !m.Amount.IsAllPositive() {
		return sdkerrors.ErrInvalidCoins.Wrap(m.Amount.String())
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgFundTreasury) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg.
func (m MsgFundTreasury) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.From)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgWithdrawFromTreasury)(nil)

// Route implements Msg.
func (m MsgWithdrawFromTreasury) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgWithdrawFromTreasury) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgWithdrawFromTreasury) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}

	if err := sdk.ValidateAccAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if !m.Amount.IsAllPositive() {
		return sdkerrors.ErrInvalidCoins.Wrap(m.Amount.String())
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgWithdrawFromTreasury) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg.
func (m MsgWithdrawFromTreasury) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Operator)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgUpdateMembers)(nil)

// Route implements Msg.
func (m MsgUpdateMembers) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgUpdateMembers) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgUpdateMembers) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}

	if len(m.MemberUpdates) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty updates")
	}
	if err := validateMembers(m.MemberUpdates); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgUpdateMembers) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg.
func (m MsgUpdateMembers) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Operator)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgUpdateDecisionPolicy)(nil)

// Route implements Msg.
func (m MsgUpdateDecisionPolicy) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgUpdateDecisionPolicy) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgUpdateDecisionPolicy) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}

	if m.GetDecisionPolicy() == nil {
		return sdkerrors.ErrInvalidRequest.Wrap("nil decision policy")
	}

	if err := m.GetDecisionPolicy().ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgUpdateDecisionPolicy) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg.
func (m MsgUpdateDecisionPolicy) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Operator)
	return []sdk.AccAddress{signer}
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

var _ sdk.Msg = (*MsgSubmitProposal)(nil)

// Route implements Msg.
func (m MsgSubmitProposal) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgSubmitProposal) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgSubmitProposal) ValidateBasic() error {
	if err := validateProposers(m.Proposers); err != nil {
		return err
	}

	if err := validateMsgs(m.GetMsgs()); err != nil {
		return err
	}

	if _, ok := Exec_name[int32(m.Exec)]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid exec option")
	}

	return nil
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

// GetSignBytes implements Msg.
func (m MsgSubmitProposal) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg.
func (m MsgSubmitProposal) GetSigners() []sdk.AccAddress {
	signers := make([]sdk.AccAddress, len(m.Proposers))
	for i, proposer := range m.Proposers {
		signers[i] = sdk.AccAddress(proposer)
	}
	return signers
}

var _ sdk.Msg = (*MsgWithdrawProposal)(nil)

// Route implements Msg.
func (m MsgWithdrawProposal) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgWithdrawProposal) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgWithdrawProposal) ValidateBasic() error {
	if m.ProposalId == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty proposal id")
	}

	if err := sdk.ValidateAccAddress(m.Address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid withdrawer address: %s", m.Address)
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgWithdrawProposal) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg.
func (m MsgWithdrawProposal) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Address)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgVote)(nil)

// Route implements Msg.
func (m MsgVote) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgVote) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgVote) ValidateBasic() error {
	if m.ProposalId == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty proposal id")
	}

	if err := sdk.ValidateAccAddress(m.Voter); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid voter address: %s", m.Voter)
	}

	if err := validateVoteOption(m.Option); err != nil {
		return err
	}

	if _, ok := Exec_name[int32(m.Exec)]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid exec option")
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgVote) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg.
func (m MsgVote) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Voter)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgExec)(nil)

// Route implements Msg.
func (m MsgExec) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgExec) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgExec) ValidateBasic() error {
	if m.ProposalId == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty proposal id")
	}

	if err := sdk.ValidateAccAddress(m.Signer); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid approver address: %s", m.Signer)
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgExec) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg.
func (m MsgExec) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Signer)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgLeaveFoundation)(nil)

// Route implements Msg.
func (m MsgLeaveFoundation) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgLeaveFoundation) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgLeaveFoundation) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid member address: %s", m.Address)
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgLeaveFoundation) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg.
func (m MsgLeaveFoundation) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Address)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgGrant)(nil)

// Route implements Msg.
func (m MsgGrant) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgGrant) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgGrant) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}

	if err := sdk.ValidateAccAddress(m.Grantee); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid grantee address: %s", m.Grantee)
	}

	if a := m.GetAuthorization(); a == nil {
		return sdkerrors.ErrInvalidType.Wrap("invalid authorization")
	}
	if err := m.GetAuthorization().ValidateBasic(); err != nil {
		return err
	}

	return nil
}

func (m MsgGrant) GetAuthorization() authz.Authorization {
	if m.Authorization == nil {
		return nil
	}

	a, ok := m.Authorization.GetCachedValue().(authz.Authorization)
	if !ok {
		return nil
	}
	return a
}

func (m *MsgGrant) SetAuthorization(a authz.Authorization) error {
	msg, ok := a.(proto.Message)
	if !ok {
		return sdkerrors.ErrInvalidType.Wrapf("can't proto marshal %T", msg)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return err
	}
	m.Authorization = any

	return nil
}

func (m MsgGrant) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var authorization authz.Authorization
	return unpacker.UnpackAny(m.Authorization, &authorization)
}

// GetSignBytes implements Msg.
func (m MsgGrant) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg.
func (m MsgGrant) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Operator)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgRevoke)(nil)

// Route implements Msg.
func (m MsgRevoke) Route() string { return RouterKey }

// Type implements Msg.
func (m MsgRevoke) Type() string { return sdk.MsgTypeURL(&m) }

// ValidateBasic implements Msg.
func (m MsgRevoke) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}

	if err := sdk.ValidateAccAddress(m.Grantee); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid grantee address: %s", m.Grantee)
	}

	return nil
}

// GetSignBytes implements Msg.
func (m MsgRevoke) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// GetSigners implements Msg.
func (m MsgRevoke) GetSigners() []sdk.AccAddress {
	signer := sdk.AccAddress(m.Operator)
	return []sdk.AccAddress{signer}
}
