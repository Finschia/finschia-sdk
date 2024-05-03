package keeper

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(ctx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	panic("implement me")
}

func (k Keeper) NextSeqSend(goCtx context.Context, req *types.QueryNextSeqSendRequest) (*types.QueryNextSeqSendResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	seq := k.GetNextSequence(ctx)

	return &types.QueryNextSeqSendResponse{Seq: seq}, nil
}

func (k Keeper) SeqToBlocknums(goCtx context.Context, req *types.QuerySeqToBlocknumsRequest) (*types.QuerySeqToBlocknumsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	const lowerBound int = 1
	const upperBound int = 1000
	if len(req.Seqs) < lowerBound || len(req.Seqs) > upperBound {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("The number of seqs requested must be between %d and %d", lowerBound, upperBound))
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	bhList := make([]uint64, len(req.Seqs))
	for i, seq := range req.Seqs {
		bh, err := k.GetSeqToBlocknum(ctx, seq)
		if err != nil {
			return nil, err
		}
		bhList[i] = bh
	}

	return &types.QuerySeqToBlocknumsResponse{Blocknums: bhList}, nil
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

func (k Keeper) Guardians(ctx context.Context, request *types.QueryGuardiansRequest) (*types.QueryGuardiansResponse, error) {
	panic("implement me")
}

func (k Keeper) Operators(ctx context.Context, request *types.QueryOperatorsRequest) (*types.QueryOperatorsResponse, error) {
	panic("implement me")
}

func (k Keeper) Judges(ctx context.Context, request *types.QueryJudgesRequest) (*types.QueryJudgesResponse, error) {
	panic("implement me")
}

func (k Keeper) Proposals(ctx context.Context, request *types.QueryProposalsRequest) (*types.QueryProposalsResponse, error) {
	panic("implement me")
}

func (k Keeper) Proposal(ctx context.Context, request *types.QueryProposalRequest) (*types.QueryProposalResponse, error) {
	panic("implement me")
}
