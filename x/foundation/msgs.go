package foundation

import (
	"github.com/line/lbm-sdk/codec/legacy"
	"github.com/line/lbm-sdk/codec/types"

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

	any, err := types.NewAnyWithValue(msg)
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
	if len(m.Proposers) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("no proposers")
	}

	proposers := map[string]bool{}
	for _, proposer := range m.Proposers {
		if proposers[proposer] {
			return sdkerrors.ErrInvalidRequest.Wrapf("duplicated proposer: %s", proposer)
		}
		proposers[proposer] = true

		if err := sdk.ValidateAccAddress(proposer); err != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid proposer address: %s", proposer)
		}
	}

	msgs := m.GetMsgs()
	if len(msgs) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("no msgs")
	}
	for i, msg := range msgs {
		if err := msg.ValidateBasic(); err != nil {
			return sdkerrors.Wrapf(err, "msg %d", i)
		}
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
	var signers []sdk.AccAddress
	for _, proposer := range m.Proposers {
		signers = append(signers, sdk.AccAddress(proposer))
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

	if m.Option == VOTE_OPTION_UNSPECIFIED {
		return sdkerrors.ErrInvalidRequest.Wrap("empty vote option")
	}
	if _, ok := VoteOption_name[int32(m.Option)]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid vote option")
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

// var _ sdk.Msg = (*MsgXxx)(nil)

// // Route implements Msg.
// func (m MsgXxx) Route() string { return RouterKey }

// // Type implements Msg.
// func (m MsgXxx) Type() string { return sdk.MsgTypeURL(&m) }

// // ValidateBasic implements Msg.
// func (m MsgXxx) ValidateBasic() error {
// 	if err := class.ValidateID(m.ClassId); err != nil {
// 		return err
// 	}
// 	if err := sdk.ValidateAccAddress(m.Approver); err != nil {
// 		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid approver address: %s", m.Approver)
// 	}

// 	if err := sdk.ValidateAccAddress(m.Proxy); err != nil {
// 		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid proxy address: %s", m.Proxy)
// 	}

// 	return nil
// }

// // GetSignBytes implements Msg.
// func (m MsgXxx) GetSignBytes() []byte {
// 	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
// }

// // GetSigners implements Msg.
// func (m MsgXxx) GetSigners() []sdk.AccAddress {
// 	signer := sdk.AccAddress(m.Approver)
// 	return []sdk.AccAddress{signer}
// }
