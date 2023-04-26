package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	datest "github.com/Finschia/finschia-sdk/x/or/da/testutil"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := datest.DaKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
	require.EqualValues(t, params.Placeholder, k.Placeholder(ctx))
}
