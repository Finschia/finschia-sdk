package keeper_test

import (
	"testing"

	testkeeper "github.com/Finschia/finschia-sdk/testutil/keeper"

	"github.com/Finschia/finschia-sdk/x/fswap/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.FswapKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
