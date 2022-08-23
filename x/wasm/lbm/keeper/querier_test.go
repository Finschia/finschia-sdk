package keeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
	wasmkeeper "github.com/line/lbm-sdk/x/wasm/keeper"
	"github.com/line/lbm-sdk/x/wasm/keeper/wasmtesting"
	lbmwasmtypes "github.com/line/lbm-sdk/x/wasm/lbm/types"
)

func TestQueryInactiveContracts(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, SupportedFeatures, nil, nil)
	keeper := keepers.WasmKeeper

	var mock wasmtesting.MockWasmer
	wasmtesting.MakeInstantiable(&mock)
	example1 := wasmkeeper.SeedNewContractInstance(t, ctx, keepers.TestKeepers, &mock)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	example2 := wasmkeeper.SeedNewContractInstance(t, ctx, keepers.TestKeepers, &mock)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// set inactive
	err := keeper.deactivateContract(ctx, example1.Contract)
	require.NoError(t, err)
	err = keeper.deactivateContract(ctx, example2.Contract)
	require.NoError(t, err)

	q := Querier(keeper)
	rq := lbmwasmtypes.QueryInactiveContractsRequest{}
	res, err := q.InactiveContracts(sdk.WrapSDKContext(ctx), &rq)
	require.NoError(t, err)
	expect := []string{example1.Contract.String(), example2.Contract.String()}
	for _, exp := range expect {
		assert.Contains(t, res.Addresses, exp)
	}
}
