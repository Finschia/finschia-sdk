package keeper

import (
	"context"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/foundation"
)

const gasCostPerIteration = uint64(20)

func canFoundationAuthorize(msgTypeURL string) bool {
	urls := map[string]bool{
		foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(): true,
	}
	return urls[msgTypeURL]
}

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
	from, err := sdk.AccAddressFromBech32(req.From)
	if err != nil {
		return nil, err
	}
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
	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.validateOperator(ctx, req.Operator); err != nil {
		return nil, err
	}

	to, err := sdk.AccAddressFromBech32(req.To)
	if err != nil {
		return nil, err
	}
	if err := s.keeper.Accept(ctx, foundation.ModuleName, to, req); err != nil {
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
	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.validateOperator(ctx, req.Operator); err != nil {
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
	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.validateOperator(ctx, req.Operator); err != nil {
		return nil, err
	}

	policy := req.GetDecisionPolicy()
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
	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.validateMembers(ctx, req.Proposers); err != nil {
		return nil, err
	}

	id, err := s.keeper.SubmitProposal(ctx, req.Proposers, req.Metadata, req.GetMsgs())
	if err != nil {
		return nil, err
	}

	proposal, err := s.keeper.GetProposal(ctx, id)
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
				ProposalId: id,
				Voter:      proposer,
				Option:     foundation.VOTE_OPTION_YES,
			}
			err = s.keeper.Vote(ctx, vote)
			if err != nil {
				return &foundation.MsgSubmitProposalResponse{ProposalId: id}, sdkerrors.Wrap(err, "The proposal was created but failed on vote")
			}
		}

		// Then try to execute the proposal
		// We consider the first proposer as the MsgExecRequest signer
		if err = s.keeper.Exec(ctx, id); err != nil {
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

	// operator may withdraw any proposal.
	if req.Address != s.keeper.GetOperator(ctx).String() {
		// check whether the address is in proposers list.
		if err = validateActorForProposal(req.Address, *proposal); err != nil {
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
	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.validateMembers(ctx, []string{req.Address}); err != nil {
		return nil, err
	}

	update := foundation.Member{
		Address: req.Address,
	}
	if err := s.keeper.UpdateMembers(ctx, []foundation.Member{update}); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventLeaveFoundation{
		Address: req.Address,
	}); err != nil {
		panic(err)
	}

	return &foundation.MsgLeaveFoundationResponse{}, nil
}

func (s msgServer) Grant(c context.Context, req *foundation.MsgGrant) (*foundation.MsgGrantResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.validateOperator(ctx, req.Operator); err != nil {
		return nil, err
	}

	msgTypeURL := req.GetAuthorization().MsgTypeURL()
	if !canFoundationAuthorize(msgTypeURL) {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("foundation cannot grant %s", msgTypeURL)
	}

	authorization := req.GetAuthorization()
	grantee, err := sdk.AccAddressFromBech32(req.Grantee)
	if err != nil {
		return nil, err
	}
	if err := s.keeper.Grant(ctx, foundation.ModuleName, grantee, authorization); err != nil {
		return nil, err
	}

	return &foundation.MsgGrantResponse{}, nil
}

func (s msgServer) Revoke(c context.Context, req *foundation.MsgRevoke) (*foundation.MsgRevokeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := s.keeper.validateOperator(ctx, req.Operator); err != nil {
		return nil, err
	}

	if !canFoundationAuthorize(req.MsgTypeUrl) {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("foundation cannot revoke %s", req.MsgTypeUrl)
	}

	grantee, err := sdk.AccAddressFromBech32(req.Grantee)
	if err != nil {
		return nil, err
	}
	if err := s.keeper.Revoke(ctx, foundation.ModuleName, grantee, req.MsgTypeUrl); err != nil {
		return nil, err
	}

	return &foundation.MsgRevokeResponse{}, nil
}
