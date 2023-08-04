package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/or/settlement/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Challenge(c context.Context, req *types.QueryChallengeRequest) (*types.QueryChallengeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	challenge, err := k.GetChallenge(ctx, req.ChallengeId)

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &types.QueryChallengeResponse{Challenge: challenge}, err
}
