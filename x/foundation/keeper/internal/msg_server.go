package internal

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

const gasCostPerIteration = uint64(20)

type msgServer struct {
	keeper Keeper
}

// NewMsgServer returns an implementation of the token MsgServer interface
// for the provided Keeper.
func NewMsgServer(keeper Keeper) foundation.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

var _ foundation.MsgServer = msgServer{}

// FundTreasury defines a method to fund the treasury.
func (s msgServer) FundTreasury(c context.Context, req *foundation.MsgFundTreasury) (*foundation.MsgFundTreasuryResponse, error) {
	from, err := s.keeper.addressCodec.StringToBytes(req.From)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", req.From)
	}

	if !req.Amount.IsValid() || !req.Amount.IsAllPositive() {
		return nil, sdkerrors.ErrInvalidCoins.Wrap(req.Amount.String())
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.FundTreasury(ctx, from, req.Amount); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventFundTreasury{
		From:   req.From,
		Amount: req.Amount,
	}); err != nil {
		panic(err)
	}

	return &foundation.MsgFundTreasuryResponse{}, nil
}

// WithdrawFromTreasury defines a method to withdraw coins from the treasury.
func (s msgServer) WithdrawFromTreasury(c context.Context, req *foundation.MsgWithdrawFromTreasury) (*foundation.MsgWithdrawFromTreasuryResponse, error) {
	if _, err := s.keeper.addressCodec.StringToBytes(req.Authority); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", req.Authority)
	}

	to, err := s.keeper.addressCodec.StringToBytes(req.To)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", req.To)
	}

	if !req.Amount.IsValid() || !req.Amount.IsAllPositive() {
		return nil, sdkerrors.ErrInvalidCoins.Wrap(req.Amount.String())
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.validateAuthority(req.Authority); err != nil {
		return nil, err
	}

	if err := s.keeper.Accept(ctx, to, req); err != nil {
		return nil, err
	}

	if err := s.keeper.WithdrawFromTreasury(ctx, to, req.Amount); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventWithdrawFromTreasury{
		To:     req.To,
		Amount: req.Amount,
	}); err != nil {
		panic(err)
	}

	return &foundation.MsgWithdrawFromTreasuryResponse{}, nil
}

func (s msgServer) UpdateMembers(c context.Context, req *foundation.MsgUpdateMembers) (*foundation.MsgUpdateMembersResponse, error) {
	if _, err := s.keeper.addressCodec.StringToBytes(req.Authority); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", req.Authority)
	}

	if len(req.MemberUpdates) == 0 {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty updates")
	}
	members := foundation.MemberRequests{Members: req.MemberUpdates}
	if err := members.ValidateBasic(s.keeper.addressCodec); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.validateAuthority(req.Authority); err != nil {
		return nil, err
	}

	if err := s.keeper.UpdateMembers(ctx, req.MemberUpdates); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventUpdateMembers{
		MemberUpdates: req.MemberUpdates,
	}); err != nil {
		panic(err)
	}

	return &foundation.MsgUpdateMembersResponse{}, nil
}

func (s msgServer) UpdateDecisionPolicy(c context.Context, req *foundation.MsgUpdateDecisionPolicy) (*foundation.MsgUpdateDecisionPolicyResponse, error) {
	if _, err := s.keeper.addressCodec.StringToBytes(req.Authority); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", req.Authority)
	}

	policy := req.GetDecisionPolicy()
	if policy == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("nil decision policy")
	}
	if err := req.GetDecisionPolicy().ValidateBasic(); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.validateAuthority(req.Authority); err != nil {
		return nil, err
	}

	if err := s.keeper.UpdateDecisionPolicy(ctx, policy); err != nil {
		return nil, err
	}

	event := &foundation.EventUpdateDecisionPolicy{}
	if err := event.SetDecisionPolicy(policy); err != nil {
		return nil, err
	}
	if err := ctx.EventManager().EmitTypedEvent(event); err != nil {
		panic(err)
	}

	return &foundation.MsgUpdateDecisionPolicyResponse{}, nil
}

func (s msgServer) SubmitProposal(c context.Context, req *foundation.MsgSubmitProposal) (*foundation.MsgSubmitProposalResponse, error) {
	if err := foundation.ValidateProposers(req.Proposers, s.keeper.addressCodec); err != nil {
		return nil, err
	}

	if len(req.GetMsgs()) == 0 {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("no msgs")
	}

	if _, ok := foundation.Exec_name[int32(req.Exec)]; !ok {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("invalid exec option")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.validateMembers(ctx, req.Proposers); err != nil {
		return nil, err
	}

	id, err := s.keeper.SubmitProposal(ctx, req.Proposers, req.Metadata, req.GetMsgs())
	if err != nil {
		return nil, err
	}

	proposal, err := s.keeper.GetProposal(ctx, *id)
	if err != nil {
		panic(err)
	}
	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventSubmitProposal{
		Proposal: *proposal,
	}); err != nil {
		panic(err)
	}

	// Try to execute proposal immediately
	if req.Exec == foundation.Exec_EXEC_TRY {
		// Consider proposers as Yes votes
		for _, proposer := range req.Proposers {
			ctx.GasMeter().ConsumeGas(gasCostPerIteration, "vote on proposal")

			vote := foundation.Vote{
				ProposalId: *id,
				Voter:      proposer,
				Option:     foundation.VOTE_OPTION_YES,
			}
			err = s.keeper.Vote(ctx, vote)
			if err != nil {
				return &foundation.MsgSubmitProposalResponse{ProposalId: *id}, errorsmod.Wrap(err, "The proposal was created but failed on vote")
			}
		}

		// Then try to execute the proposal
		if err = s.keeper.Exec(ctx, *id); err != nil {
			return &foundation.MsgSubmitProposalResponse{ProposalId: *id}, errorsmod.Wrap(err, "The proposal was created but failed on exec")
		}
	}

	return &foundation.MsgSubmitProposalResponse{ProposalId: *id}, nil
}

func (s msgServer) WithdrawProposal(c context.Context, req *foundation.MsgWithdrawProposal) (*foundation.MsgWithdrawProposalResponse, error) {
	id := req.ProposalId
	if id == 0 {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty proposal id")
	}

	if _, err := s.keeper.addressCodec.StringToBytes(req.Address); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid withdrawer address: %s", req.Address)
	}

	ctx := sdk.UnwrapSDKContext(c)

	proposal, err := s.keeper.GetProposal(ctx, id)
	if err != nil {
		return nil, err
	}

	// authority may withdraw any proposal.
	if err := s.keeper.validateAuthority(req.Address); err != nil {
		// check whether the address is in proposers list.
		if err := validateActorForProposal(req.Address, *proposal); err != nil {
			return nil, err
		}
	}

	if err := s.keeper.WithdrawProposal(ctx, id); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventWithdrawProposal{
		ProposalId: id,
	}); err != nil {
		panic(err)
	}

	return &foundation.MsgWithdrawProposalResponse{}, nil
}

func (s msgServer) Vote(c context.Context, req *foundation.MsgVote) (*foundation.MsgVoteResponse, error) {
	if req.ProposalId == 0 {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty proposal id")
	}

	if _, err := s.keeper.addressCodec.StringToBytes(req.Voter); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid voter address: %s", req.Voter)
	}

	if err := foundation.ValidateVoteOption(req.Option); err != nil {
		return nil, err
	}

	if _, ok := foundation.Exec_name[int32(req.Exec)]; !ok {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("invalid exec option")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.validateMembers(ctx, []string{req.Voter}); err != nil {
		return nil, err
	}

	vote := foundation.Vote{
		ProposalId: req.ProposalId,
		Voter:      req.Voter,
		Option:     req.Option,
		Metadata:   req.Metadata,
	}
	if err := s.keeper.Vote(ctx, vote); err != nil {
		return nil, err
	}

	// Try to execute proposal immediately
	if req.Exec == foundation.Exec_EXEC_TRY {
		if err := s.keeper.Exec(ctx, req.ProposalId); err != nil {
			return nil, err
		}
	}

	return &foundation.MsgVoteResponse{}, nil
}

func (s msgServer) Exec(c context.Context, req *foundation.MsgExec) (*foundation.MsgExecResponse, error) {
	if req.ProposalId == 0 {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty proposal id")
	}

	if _, err := s.keeper.addressCodec.StringToBytes(req.Signer); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid signer address: %s", req.Signer)
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.validateMembers(ctx, []string{req.Signer}); err != nil {
		return nil, err
	}

	if err := s.keeper.Exec(ctx, req.ProposalId); err != nil {
		return nil, err
	}

	return &foundation.MsgExecResponse{}, nil
}

func (s msgServer) LeaveFoundation(c context.Context, req *foundation.MsgLeaveFoundation) (*foundation.MsgLeaveFoundationResponse, error) {
	if _, err := s.keeper.addressCodec.StringToBytes(req.Address); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid member address: %s", req.Address)
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.validateMembers(ctx, []string{req.Address}); err != nil {
		return nil, err
	}

	update := foundation.MemberRequest{
		Address: req.Address,
		Remove:  true,
	}
	if err := s.keeper.UpdateMembers(ctx, []foundation.MemberRequest{update}); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventLeaveFoundation{
		Address: req.Address,
	}); err != nil {
		panic(err)
	}

	return &foundation.MsgLeaveFoundationResponse{}, nil
}

func (s msgServer) UpdateCensorship(c context.Context, req *foundation.MsgUpdateCensorship) (*foundation.MsgUpdateCensorshipResponse, error) {
	if _, err := s.keeper.addressCodec.StringToBytes(req.Authority); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", req.Authority)
	}

	if err := req.Censorship.ValidateBasic(); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	url := req.Censorship.MsgTypeUrl
	if err := s.keeper.validateCensorshipAuthority(ctx, url, req.Authority); err != nil {
		return nil, err
	}

	if err := s.keeper.UpdateCensorship(ctx, req.Censorship); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventUpdateCensorship{
		Censorship: req.Censorship,
	}); err != nil {
		panic(err)
	}

	return &foundation.MsgUpdateCensorshipResponse{}, nil
}

func (s msgServer) Grant(c context.Context, req *foundation.MsgGrant) (*foundation.MsgGrantResponse, error) {
	if _, err := s.keeper.addressCodec.StringToBytes(req.Authority); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", req.Authority)
	}

	grantee, err := s.keeper.addressCodec.StringToBytes(req.Grantee)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid grantee address: %s", req.Grantee)
	}

	authorization := req.GetAuthorization()
	if authorization != nil {
		if err := authorization.ValidateBasic(); err != nil {
			return nil, err
		}
	} else {
		return nil, sdkerrors.ErrInvalidType.Wrap("invalid authorization")
	}

	ctx := sdk.UnwrapSDKContext(c)

	url := authorization.MsgTypeURL()
	if err := s.keeper.validateCensorshipAuthority(ctx, url, req.Authority); err != nil {
		return nil, err
	}

	if err := s.keeper.Grant(ctx, grantee, authorization); err != nil {
		return nil, err
	}

	return &foundation.MsgGrantResponse{}, nil
}

func (s msgServer) Revoke(c context.Context, req *foundation.MsgRevoke) (*foundation.MsgRevokeResponse, error) {
	if _, err := s.keeper.addressCodec.StringToBytes(req.Authority); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", req.Authority)
	}

	grantee, err := s.keeper.addressCodec.StringToBytes(req.Grantee)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid grantee address: %s", req.Grantee)
	}

	url := req.MsgTypeUrl
	if len(url) == 0 {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("empty url")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.validateCensorshipAuthority(ctx, url, req.Authority); err != nil {
		return nil, err
	}

	if err := s.keeper.Revoke(ctx, grantee, req.MsgTypeUrl); err != nil {
		return nil, err
	}

	return &foundation.MsgRevokeResponse{}, nil
}
