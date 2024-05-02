package keeper

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Finschia/finschia-sdk/store/prefix"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/query"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) NextSeqSend(goCtx context.Context, req *types.QueryNextSeqSendRequest) (*types.QueryNextSeqSendResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	seq := k.GetNextSequence(ctx)

	return &types.QueryNextSeqSendResponse{Seq: seq}, nil
}

func (k Keeper) GreatestSeqByOperator(ctx context.Context, request *types.QueryGreatestSeqByOperatorRequest) (*types.QueryGreatestSeqByOperatorResponse, error) {
	panic("implement me")
}

func (k Keeper) GreatestConsecutiveConfirmedSeq(ctx context.Context, request *types.QueryGreatestConsecutiveConfirmedSeqRequest) (*types.QueryGreatestConsecutiveConfirmedSeqResponse, error) {
	panic("implement me")
}

func (k Keeper) SubmittedProvision(ctx context.Context, request *types.QuerySubmittedProvisionRequest) (*types.QuerySubmittedProvisionResponse, error) {
	panic("implement me")
}

func (k Keeper) ConfirmedProvision(ctx context.Context, request *types.QueryConfirmedProvisionRequest) (*types.QueryConfirmedProvisionResponse, error) {
	panic("implement me")
}

func (k Keeper) NeededSubmissionSeqs(ctx context.Context, request *types.QueryNeededSubmissionSeqsRequest) (*types.QueryNeededSubmissionSeqsResponse, error) {
	panic("implement me")
}

func (k Keeper) Commitments(ctx context.Context, request *types.QueryCommitmentsRequest) (*types.QueryCommitmentsResponse, error) {
	panic("implement me")
}

func (k Keeper) Guardians(goCtx context.Context, req *types.QueryGuardiansRequest) (*types.QueryGuardiansResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	guardians := make([]string, 0)
	roles := k.GetRolePairs(ctx)
	for _, pair := range roles {
		if pair.Role == types.RoleGuardian {
			guardians = append(guardians, pair.Address)
		}
	}

	return &types.QueryGuardiansResponse{Guardians: guardians}, nil
}

func (k Keeper) Operators(goCtx context.Context, req *types.QueryOperatorsRequest) (*types.QueryOperatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	operators := make([]string, 0)
	roles := k.GetRolePairs(ctx)
	for _, pair := range roles {
		if pair.Role == types.RoleOperator {
			operators = append(operators, pair.Address)
		}
	}

	return &types.QueryOperatorsResponse{Operators: operators}, nil
}

func (k Keeper) Judges(goCtx context.Context, req *types.QueryJudgesRequest) (*types.QueryJudgesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	judges := make([]string, 0)
	roles := k.GetRolePairs(ctx)
	for _, pair := range roles {
		if pair.Role == types.RoleJudge {
			judges = append(judges, pair.Address)
		}
	}

	return &types.QueryJudgesResponse{Judges: judges}, nil
}

func (k Keeper) Proposals(goCtx context.Context, req *types.QueryProposalsRequest) (*types.QueryProposalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyProposalPrefix)
	proposals := make([]types.RoleProposal, 0)
	pageRes, err := query.Paginate(store, req.Pagination, func(_, value []byte) error {
		var proposal types.RoleProposal
		k.cdc.MustUnmarshal(value, &proposal)
		proposals = append(proposals, proposal)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryProposalsResponse{Proposals: proposals, Pagination: pageRes}, nil
}

func (k Keeper) Proposal(goCtx context.Context, req *types.QueryProposalRequest) (*types.QueryProposalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	proposal, found := k.GetRoleProposal(ctx, req.ProposalId)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("role proposal %d", req.ProposalId))
	}

	return &types.QueryProposalResponse{Proposal: proposal}, nil
}

func (k Keeper) Votes(goCtx context.Context, req *types.QueryVotesRequest) (*types.QueryVotesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	votes := k.GetProposalVotes(ctx, req.ProposalId)
	return &types.QueryVotesResponse{Votes: votes}, nil
}

func (k Keeper) Vote(goCtx context.Context, req *types.QueryVoteRequest) (*types.QueryVoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	opt, err := k.GetVote(ctx, req.ProposalId, sdk.MustAccAddressFromBech32(req.Voter))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &types.QueryVoteResponse{Vote: types.Vote{ProposalId: req.ProposalId, Voter: req.Voter, Option: opt}}, nil
}
