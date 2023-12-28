package internal

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"

	"github.com/cosmos/cosmos-sdk/runtime"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/Finschia/finschia-sdk/x/foundation"
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

func (s queryServer) Member(c context.Context, req *foundation.QueryMemberRequest) (*foundation.QueryMemberResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	addr, err := s.keeper.addressCodec.StringToBytes(req.Address)
	if err != nil {
		return nil, err
	}
	member, err := s.keeper.GetMember(ctx, addr)
	if err != nil {
		return nil, err
	}

	return &foundation.QueryMemberResponse{Member: member}, nil
}

func (s queryServer) Members(c context.Context, req *foundation.QueryMembersRequest) (*foundation.QueryMembersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var members []foundation.Member
	ctx := sdk.UnwrapSDKContext(c)
	store := s.keeper.storeService.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(store)
	memberStore := prefix.NewStore(adapter, memberKeyPrefix)
	pageRes, err := query.Paginate(memberStore, req.Pagination, func(key, value []byte) error {
		var member foundation.Member
		s.keeper.cdc.MustUnmarshal(value, &member)
		members = append(members, member)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &foundation.QueryMembersResponse{Members: members, Pagination: pageRes}, nil
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
	store := s.keeper.storeService.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(store)
	proposalStore := prefix.NewStore(adapter, proposalKeyPrefix)
	pageRes, err := query.Paginate(proposalStore, req.Pagination, func(key, value []byte) error {
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
	voter, err := s.keeper.addressCodec.StringToBytes(req.Voter)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid voter address")
	}
	vote, err := s.keeper.GetVote(ctx, req.ProposalId, voter)
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
	store := s.keeper.storeService.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(store)
	voteStore := prefix.NewStore(adapter, append(voteKeyPrefix, Uint64ToBytes(req.ProposalId)...))
	pageRes, err := query.Paginate(voteStore, req.Pagination, func(key, value []byte) error {
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

func (s queryServer) Censorships(c context.Context, req *foundation.QueryCensorshipsRequest) (*foundation.QueryCensorshipsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var censorships []foundation.Censorship
	ctx := sdk.UnwrapSDKContext(c)
	store := s.keeper.storeService.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(store)
	censorshipStore := prefix.NewStore(adapter, censorshipKeyPrefix)
	pageRes, err := query.Paginate(censorshipStore, req.Pagination, func(key, value []byte) error {
		var censorship foundation.Censorship
		s.keeper.cdc.MustUnmarshal(value, &censorship)
		censorships = append(censorships, censorship)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &foundation.QueryCensorshipsResponse{Censorships: censorships, Pagination: pageRes}, nil
}

func (s queryServer) Grants(c context.Context, req *foundation.QueryGrantsRequest) (*foundation.QueryGrantsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	grantee, err := s.keeper.addressCodec.StringToBytes(req.Grantee)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := s.keeper.storeService.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(store)

	if req.MsgTypeUrl != "" {
		keyPrefix := grantKey(grantee, req.MsgTypeUrl)
		grantStore := prefix.NewStore(adapter, keyPrefix)

		var authorizations []*codectypes.Any
		_, err = query.Paginate(grantStore, req.Pagination, func(key, value []byte) error {
			var authorization foundation.Authorization
			if err := s.keeper.cdc.UnmarshalInterface(value, &authorization); err != nil {
				panic(err)
			}

			msg, ok := authorization.(proto.Message)
			if !ok {
				panic(sdkerrors.ErrInvalidType.Wrapf("can't proto marshal %T", msg))
			}
			any, err := codectypes.NewAnyWithValue(msg)
			if err != nil {
				panic(err)
			}
			authorizations = append(authorizations, any)

			return nil
		})
		if err != nil {
			return nil, err
		}

		return &foundation.QueryGrantsResponse{Authorizations: authorizations}, nil

	}

	keyPrefix := grantKeyPrefixByGrantee(grantee)
	grantStore := prefix.NewStore(adapter, keyPrefix)

	var authorizations []*codectypes.Any
	pageRes, err := query.Paginate(grantStore, req.Pagination, func(key, value []byte) error {
		var authorization foundation.Authorization
		if err := s.keeper.cdc.UnmarshalInterface(value, &authorization); err != nil {
			panic(err)
		}

		msg, ok := authorization.(proto.Message)
		if !ok {
			panic(sdkerrors.ErrInvalidType.Wrapf("can't proto marshal %T", msg))
		}
		any, err := codectypes.NewAnyWithValue(msg)
		if err != nil {
			panic(err)
		}
		authorizations = append(authorizations, any)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &foundation.QueryGrantsResponse{Authorizations: authorizations, Pagination: pageRes}, nil
}
