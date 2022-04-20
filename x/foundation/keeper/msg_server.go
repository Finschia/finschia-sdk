package keeper

import (
	"context"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/foundation"
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
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.fundTreasury(ctx, sdk.AccAddress(req.From), req.Amount); err != nil {
		return nil, err
	}

	return &foundation.MsgFundTreasuryResponse{}, nil
}

// WithdrawFromTreasury defines a method to withdraw coins from the treasury.
func (s msgServer) WithdrawFromTreasury(c context.Context, req *foundation.MsgWithdrawFromTreasury) (*foundation.MsgWithdrawFromTreasuryResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if sdk.AccAddress(req.Operator) != s.keeper.GetOperator(ctx) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "Not authorized to %s: %s", sdk.MsgTypeURL(req), req.Operator)
	}

	if err := s.keeper.withdrawFromTreasury(ctx, sdk.AccAddress(req.To), req.Amount); err != nil {
		return nil, err
	}

	return &foundation.MsgWithdrawFromTreasuryResponse{}, nil
}

func (s msgServer) UpdateMembers(c context.Context, req *foundation.MsgUpdateMembers) (*foundation.MsgUpdateMembersResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if sdk.AccAddress(req.Operator) != s.keeper.GetOperator(ctx) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "Not authorized to %s: %s", sdk.MsgTypeURL(req), req.Operator)
	}

	if err := s.keeper.updateMembers(ctx, req.MemberUpdates); err != nil {
		return nil, err
	}

	return &foundation.MsgUpdateMembersResponse{}, nil
}

func (s msgServer) UpdateDecisionPolicy(c context.Context, req *foundation.MsgUpdateDecisionPolicy) (*foundation.MsgUpdateDecisionPolicyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if sdk.AccAddress(req.Operator) != s.keeper.GetOperator(ctx) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "Not authorized to %s: %s", sdk.MsgTypeURL(req), req.Operator)
	}

	policy := req.GetDecisionPolicy()
	if err := policy.Validate(s.keeper.GetFoundationInfo(ctx), s.keeper.config); err != nil {
		return nil, err
	}

	if err := s.keeper.updateDecisionPolicy(ctx, policy); err != nil {
		return nil, err
	}

	return &foundation.MsgUpdateDecisionPolicyResponse{}, nil
}

func (s msgServer) SubmitProposal(c context.Context, req *foundation.MsgSubmitProposal) (*foundation.MsgSubmitProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	id, err := s.keeper.submitProposal(ctx, req.Proposers, req.Metadata, req.GetMsgs())
	if err != nil {
		return nil, err
	}

	// Try to execute proposal immediately
	if req.Exec == foundation.Exec_EXEC_TRY {
		// Consider proposers as Yes votes
		for _, proposer := range req.Proposers {
			ctx.GasMeter().ConsumeGas(gasCostPerIteration, "vote on proposal")
			err = s.keeper.vote(ctx, id, proposer, foundation.VOTE_OPTION_YES, "")
			if err != nil {
				return &foundation.MsgSubmitProposalResponse{ProposalId: id}, sdkerrors.Wrap(err, "The proposal was created but failed on vote")
			}
		}

		// Then try to execute the proposal
		// We consider the first proposer as the MsgExecRequest signer
		signer := req.Proposers[0]
		if err = s.keeper.exec(ctx, id, signer); err != nil {
			return &foundation.MsgSubmitProposalResponse{ProposalId: id}, sdkerrors.Wrap(err, "The proposal was created but failed on exec")
		}
	}

	return &foundation.MsgSubmitProposalResponse{ProposalId: id}, nil
}

func (s msgServer) WithdrawProposal(c context.Context, req *foundation.MsgWithdrawProposal) (*foundation.MsgWithdrawProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	id := req.ProposalId

	proposal, err := s.keeper.GetProposal(ctx, id)
	if err != nil {
		return nil, err
	}

	// check whether the address is in proposers list.
	if err = validateActorForProposal(req.Address, *proposal); err != nil {
		return nil, err
	}

	err = ctx.EventManager().EmitTypedEvent(&foundation.EventWithdrawProposal{ProposalId: id})
	if err != nil {
		return nil, err
	}

	if err = s.keeper.withdrawProposal(ctx, *proposal); err != nil {
		return nil, err
	}

	return &foundation.MsgWithdrawProposalResponse{}, nil
}

func (s msgServer) Vote(c context.Context, req *foundation.MsgVote) (*foundation.MsgVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.vote(ctx, req.ProposalId, req.Voter, req.Option, req.Metadata); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventVote{ProposalId: req.ProposalId}); err != nil {
		return nil, err
	}

	// Try to execute proposal immediately
	if req.Exec == foundation.Exec_EXEC_TRY {
		if err := s.keeper.exec(ctx, req.ProposalId, req.Voter); err != nil {
			return nil, err
		}
	}

	return &foundation.MsgVoteResponse{}, nil
}

func (s msgServer) Exec(c context.Context, req *foundation.MsgExec) (*foundation.MsgExecResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	id := req.ProposalId

	proposal, err := s.keeper.GetProposal(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.keeper.exec(ctx, req.ProposalId, req.Signer); err != nil {
		return nil, err
	}

	if err = ctx.EventManager().EmitTypedEvent(&foundation.EventExec{
		ProposalId: id,
		Result:     proposal.ExecutorResult,
	}); err != nil {
		return nil, err
	}

	return &foundation.MsgExecResponse{}, nil
}

func (s msgServer) LeaveFoundation(c context.Context, req *foundation.MsgLeaveFoundation) (*foundation.MsgLeaveFoundationResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	update := foundation.Member{
		Address: req.Address,
		Weight: sdk.ZeroDec(),
	}
	if err := s.keeper.updateMembers(ctx, []foundation.Member{update}); err != nil {
		return nil, err
	}

	return &foundation.MsgLeaveFoundationResponse{}, nil
}
