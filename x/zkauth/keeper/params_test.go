package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	datest "github.com/Finschia/finschia-sdk/x/zkauth/testutil"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := datest.ZkAuthKeeper(t)
	params := types.DefaultParams()

	require.Nil(t, k.SetParams(ctx, params))

	savedParams, err := k.GetParams(ctx)

	require.Nil(t, err)
	require.EqualValues(t, params, savedParams)
}
