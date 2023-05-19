package keeper

import (
	"context"

	"github.com/Finschia/finschia-sdk/x/or/settlement/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Challenge(ctx context.Context, req *types.QueryChallengeRequest) (*types.QueryChallengeResponse, error) {
	panic("implement me")
}
