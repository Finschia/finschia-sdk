package keeper

import (
	"context"

	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(ctx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	panic("implement me")
}

func (k Keeper) NextSeqSend(ctx context.Context, request *types.NextSeqSendRequest) (*types.NextSeqSendResponse, error) {
	panic("implement me")
}

func (k Keeper) GreatestSeqByOperator(ctx context.Context, request *types.GreatestSeqByOperatorRequest) (*types.GreatestSeqByOperatorResponse, error) {
	panic("implement me")
}

func (k Keeper) GreatestConsecutiveConfirmedSeq(ctx context.Context, request *types.GreatestConsecutiveConfirmedSeqRequest) (*types.GreatestConsecutiveConfirmedSeqResponse, error) {
	panic("implement me")
}

func (k Keeper) SubmittedProvision(ctx context.Context, request *types.SubmittedProvisionRequest) (*types.SubmittedProvisionResponse, error) {
	panic("implement me")
}

func (k Keeper) ConfirmedProvision(ctx context.Context, request *types.ConfirmedProvisionRequest) (*types.ConfirmedProvisionResponse, error) {
	panic("implement me")
}

func (k Keeper) NeededSubmissionSeqs(ctx context.Context, request *types.NeededSubmissionSeqsRequest) (*types.NeededSubmissionSeqsResponse, error) {
	panic("implement me")
}

func (k Keeper) Commitments(ctx context.Context, request *types.CommitmentsRequest) (*types.CommitmentsResponse, error) {
	panic("implement me")
}

func (k Keeper) Guardians(ctx context.Context, request *types.GuardiansRequest) (*types.GuardiansResponse, error) {
	panic("implement me")
}

func (k Keeper) Operators(ctx context.Context, request *types.OperatorsRequest) (*types.OperatorsResponse, error) {
	panic("implement me")
}

func (k Keeper) Judges(ctx context.Context, request *types.JudgesRequest) (*types.JudgesResponse, error) {
	panic("implement me")
}
