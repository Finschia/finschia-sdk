package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/line/lbm-sdk/store/prefix"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/query"
	"github.com/line/lbm-sdk/x/foundation"
)

type queryServer struct {
	keeper Keeper
}

func NewQueryServer(keeper Keeper) foundation.QueryServer {
	return &queryServer{
		keeper: keeper,
	}
}

var _ foundation.QueryServer = (*queryServer)(nil)

func (s queryServer) Params(c context.Context, req *foundation.QueryParamsRequest) (*foundation.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	return &foundation.QueryParamsResponse{Params: s.keeper.GetParams(ctx)}, nil
}

func (s queryServer) ValidatorAuth(c context.Context, req *foundation.QueryValidatorAuthRequest) (*foundation.QueryValidatorAuthResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty validator address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	addr := sdk.ValAddress(req.ValidatorAddress)
	auth, err := s.keeper.GetValidatorAuth(ctx, addr)
	if err != nil {
		return nil, err
	}

	return &foundation.QueryValidatorAuthResponse{Auth: auth}, nil
}

func (s queryServer) ValidatorAuths(c context.Context, req *foundation.QueryValidatorAuthsRequest) (*foundation.QueryValidatorAuthsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var auths []foundation.ValidatorAuth
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	validatorStore := prefix.NewStore(store, validatorAuthKeyPrefix)
	pageRes, err := query.Paginate(validatorStore, req.Pagination, func(key []byte, value []byte) error {
		var auth foundation.ValidatorAuth
		s.keeper.cdc.MustUnmarshal(value, &auth)
		auths = append(auths, auth)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &foundation.QueryValidatorAuthsResponse{Auths: auths, Pagination: pageRes}, nil
}

func (s queryServer) Treasury(c context.Context, req *foundation.QueryTreasuryRequest) (*foundation.QueryTreasuryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	amount := s.keeper.GetTreasury(ctx)

	return &foundation.QueryTreasuryResponse{Amount: amount}, nil
}

func (s queryServer) FoundationInfo(c context.Context, req *foundation.QueryFoundationInfoRequest) (*foundation.QueryFoundationInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	info := s.keeper.GetFoundationInfo(ctx)

	return &foundation.QueryFoundationInfoResponse{Info: info}, nil
}

func (s queryServer) FoundationMember(c context.Context, req *foundation.QueryFoundationMemberRequest) (*foundation.QueryFoundationMemberResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	member, err := s.keeper.GetMember(ctx, sdk.AccAddress(req.Address))
	if err != nil {
		return nil, err
	}

	return &foundation.QueryFoundationMemberResponse{Member: member}, nil
}

func (s queryServer) FoundationMembers(c context.Context, req *foundation.QueryFoundationMembersRequest) (*foundation.QueryFoundationMembersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var members []foundation.Member
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	memberStore := prefix.NewStore(store, memberKeyPrefix)
	pageRes, err := query.Paginate(memberStore, req.Pagination, func(key []byte, value []byte) error {
		var member foundation.Member
		s.keeper.cdc.MustUnmarshal(value, &member)
		members = append(members, member)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &foundation.QueryFoundationMembersResponse{Members: members, Pagination: pageRes}, nil
}

func (s queryServer) Proposal(c context.Context, req *foundation.QueryProposalRequest) (*foundation.QueryProposalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	proposal, err := s.keeper.GetProposal(ctx, req.ProposalId)
	if err != nil {
		return nil, err
	}

	return &foundation.QueryProposalResponse{Proposal: proposal}, nil
}

func (s queryServer) Proposals(c context.Context, req *foundation.QueryProposalsRequest) (*foundation.QueryProposalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var proposals []foundation.Proposal
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	proposalStore := prefix.NewStore(store, proposalKeyPrefix)
	pageRes, err := query.Paginate(proposalStore, req.Pagination, func(key []byte, value []byte) error {
		var proposal foundation.Proposal
		s.keeper.cdc.MustUnmarshal(value, &proposal)
		proposals = append(proposals, proposal)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &foundation.QueryProposalsResponse{Proposals: proposals, Pagination: pageRes}, nil
}

func (s queryServer) Vote(c context.Context, req *foundation.QueryVoteRequest) (*foundation.QueryVoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	vote, err := s.keeper.GetVote(ctx, req.ProposalId, sdk.AccAddress(req.Voter))
	if err != nil {
		return nil, err
	}

	return &foundation.QueryVoteResponse{Vote: vote}, nil
}

func (s queryServer) Votes(c context.Context, req *foundation.QueryVotesRequest) (*foundation.QueryVotesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var votes []foundation.Vote
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	voteStore := prefix.NewStore(store, append(voteKeyPrefix, Uint64ToBytes(req.ProposalId)...))
	pageRes, err := query.Paginate(voteStore, req.Pagination, func(key []byte, value []byte) error {
		var vote foundation.Vote
		s.keeper.cdc.MustUnmarshal(value, &vote)
		votes = append(votes, vote)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &foundation.QueryVotesResponse{Votes: votes, Pagination: pageRes}, nil
}

func (s queryServer) TallyResult(c context.Context, req *foundation.QueryTallyResultRequest) (*foundation.QueryTallyResultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	proposal, err := s.keeper.GetProposal(ctx, req.ProposalId)
	if err != nil {
		return nil, err
	}

	tally, err := s.keeper.tally(ctx, *proposal)
	if err != nil {
		return nil, err
	}

	return &foundation.QueryTallyResultResponse{Tally: tally}, nil
}
