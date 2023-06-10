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

	require.Nil(t, k.SetParams(ctx, params))

	savedParams := k.GetParams(ctx)
	savedCCBatchMaxBytes := k.CCBatchMaxBytes(ctx)
	savedSCCBatchMaxBytes := k.SCCBatchMaxBytes(ctx)
	require.EqualValues(t, params, savedParams)
	require.EqualValues(t, params.CCBatchMaxBytes, savedCCBatchMaxBytes)
	require.EqualValues(t, params.SCCBatchMaxBytes, savedSCCBatchMaxBytes)
}
