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

	savedParams, err := k.GetParams(ctx)
	require.Nil(t, err)
	savedCTCBatchMaxBytes, err := k.CTCBatchMaxBytes(ctx)
	require.Nil(t, err)
	savedSCCBatchMaxBytes, err := k.SCCBatchMaxBytes(ctx)
	require.Nil(t, err)
	require.EqualValues(t, params, savedParams)
	require.EqualValues(t, params.CTCBatchMaxBytes, savedCTCBatchMaxBytes)
	require.EqualValues(t, params.SCCBatchMaxBytes, savedSCCBatchMaxBytes)
}
