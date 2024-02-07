package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/Finschia/finschia-sdk/types"
	datest "github.com/Finschia/finschia-sdk/x/zkauth/testutil"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := datest.ZkAuthKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	testCases := map[string]struct {
		request any
		valid   bool
	}{
		"default genesis": {
			&types.QueryParamsRequest{},
			true,
		},
		"invalid request": {
			nil,
			false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			params := types.DefaultParams()
			err := keeper.SetParams(ctx, params)
			require.NoError(t, err)

			if tc.valid {
				response, err := keeper.Params(wctx, tc.request.(*types.QueryParamsRequest))
				require.NoError(t, err)
				require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
			} else {
				_, err := keeper.Params(wctx, nil)
				require.Error(t, err)
			}
		})
	}
}
