package keeper

import (
	"context"

	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(ctx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GreatestSeqByOperator(ctx context.Context, request *types.GreatestSeqByOperatorRequest) (*types.GreatestSeqByOperatorResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) NextSeqToConfirm(ctx context.Context, request *types.NextSeqToConfirmRequest) (*types.NextSeqToConfirmResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) SubmittedProvision(ctx context.Context, request *types.SubmittedProvisionRequest) (*types.SubmittedProvisionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) ConfirmedProvision(ctx context.Context, request *types.ConfirmedProvisionRequest) (*types.ConfirmedProvisionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) Commitments(ctx context.Context, request *types.CommitmentsRequest) (*types.CommitmentsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) Guardians(ctx context.Context, request *types.GuardiansRequest) (*types.GuardiansResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) Operators(ctx context.Context, request *types.OperatorsRequest) (*types.OperatorsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) Judges(ctx context.Context, request *types.JudgesRequest) (*types.JudgesResponse, error) {
	//TODO implement me
	panic("implement me")
}
