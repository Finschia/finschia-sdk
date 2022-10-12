package foundation

import (
	"github.com/gogo/protobuf/proto"

	codectypes "github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

var _ sdk.Msg = (*MsgFundTreasury)(nil)

// ValidateBasic implements Msg.
func (m MsgFundTreasury) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.From)
	}

	if !m.Amount.IsValid() || !m.Amount.IsAllPositive() {
		return sdkerrors.ErrInvalidCoins.Wrap(m.Amount.String())
	}

	return nil
}

// GetSigners implements Msg.
func (m MsgFundTreasury) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgWithdrawFromTreasury)(nil)

// ValidateBasic implements Msg.
func (m MsgWithdrawFromTreasury) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}

	if _, err := sdk.AccAddressFromBech32(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", m.To)
	}

	if !m.Amount.IsValid() || !m.Amount.IsAllPositive() {
		return sdkerrors.ErrInvalidCoins.Wrap(m.Amount.String())
	}

	return nil
}

// GetSigners implements Msg.
func (m MsgWithdrawFromTreasury) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Operator)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgUpdateMembers)(nil)

// ValidateBasic implements Msg.
func (m MsgUpdateMembers) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Operator); err != nil {
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

// GetSigners implements Msg.
func (m MsgUpdateMembers) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Operator)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgUpdateDecisionPolicy)(nil)

// ValidateBasic implements Msg.
func (m MsgUpdateDecisionPolicy) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Operator); err != nil {
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

// GetSigners implements Msg.
func (m MsgUpdateDecisionPolicy) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Operator)
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

// GetSigners implements Msg.
func (m MsgSubmitProposal) GetSigners() []sdk.AccAddress {
	signers := make([]sdk.AccAddress, len(m.Proposers))
	for i, p := range m.Proposers {
		proposer, _ := sdk.AccAddressFromBech32(p)
		signers[i] = proposer
	}
	return signers
}

var _ sdk.Msg = (*MsgWithdrawProposal)(nil)

// ValidateBasic implements Msg.
func (m MsgWithdrawProposal) ValidateBasic() error {
	if m.ProposalId == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty proposal id")
	}

	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid withdrawer address: %s", m.Address)
	}

	return nil
}

// GetSigners implements Msg.
func (m MsgWithdrawProposal) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Address)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgVote)(nil)

// ValidateBasic implements Msg.
func (m MsgVote) ValidateBasic() error {
	if m.ProposalId == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty proposal id")
	}

	if _, err := sdk.AccAddressFromBech32(m.Voter); err != nil {
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

// GetSigners implements Msg.
func (m MsgVote) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Voter)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgExec)(nil)

// ValidateBasic implements Msg.
func (m MsgExec) ValidateBasic() error {
	if m.ProposalId == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty proposal id")
	}

	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid approver address: %s", m.Signer)
	}

	return nil
}

// GetSigners implements Msg.
func (m MsgExec) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Signer)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgLeaveFoundation)(nil)

// ValidateBasic implements Msg.
func (m MsgLeaveFoundation) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid member address: %s", m.Address)
	}

	return nil
}

// GetSigners implements Msg.
func (m MsgLeaveFoundation) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Address)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgGrant)(nil)

// ValidateBasic implements Msg.
func (m MsgGrant) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}

	if _, err := sdk.AccAddressFromBech32(m.Grantee); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid grantee address: %s", m.Grantee)
	}

	if a := m.GetAuthorization(); a != nil {
		if err := a.ValidateBasic(); err != nil {
			return err
		}
	} else {
		return sdkerrors.ErrInvalidType.Wrap("invalid authorization")
	}

	return nil
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

// GetSigners implements Msg.
func (m MsgGrant) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Operator)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgRevoke)(nil)

// ValidateBasic implements Msg.
func (m MsgRevoke) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}

	if _, err := sdk.AccAddressFromBech32(m.Grantee); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid grantee address: %s", m.Grantee)
	}

	if len(m.MsgTypeUrl) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf("empty url")
	}

	return nil
}

// GetSigners implements Msg.
func (m MsgRevoke) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Operator)
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = (*MsgOneTimeMint)(nil)

// ValidateBasic implements Msg.
func (m MsgOneTimeMint) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Operator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", m.Operator)
	}

	if !m.Amount.IsValid() || !m.Amount.IsAllPositive() {
		return sdkerrors.ErrInvalidCoins.Wrap(m.Amount.String())
	}

	return nil
}

// GetSigners implements Msg.
func (m MsgOneTimeMint) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Operator)
	return []sdk.AccAddress{signer}
}
