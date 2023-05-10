package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/Finschia/finschia-sdk/types"
	datest "github.com/Finschia/finschia-sdk/x/or/da/testutil"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := datest.DaKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
