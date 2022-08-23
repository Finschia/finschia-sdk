package keeper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
	wasmkeeper "github.com/line/lbm-sdk/x/wasm/keeper"
	"github.com/line/lbm-sdk/x/wasm/keeper/wasmtesting"
)

const SupportedFeatures = "iterator,staking,stargate"

func TestActivateContract(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, SupportedFeatures, nil, nil)

	k := keepers.WasmKeeper
	var mock wasmtesting.MockWasmer
	wasmtesting.MakeInstantiable(&mock)
	example := wasmkeeper.SeedNewContractInstance(t, ctx, keepers.TestKeepers, &mock)
	em := sdk.NewEventManager()

	// request no contract address -> fail
	err := k.activateContract(ctx, example.CreatorAddr)
	require.Error(t, err, fmt.Sprintf("no contract %s", example.CreatorAddr))

	// first activate -> fail
	err = k.activateContract(ctx.WithEventManager(em), example.Contract)
	require.Error(t, err, fmt.Sprintf("no inactivate contract %s", example.Contract))

	// add to inactive contract
	err = k.deactivateContract(ctx, example.Contract)
	require.NoError(t, err)

	// second activate -> success
	err = k.activateContract(ctx, example.Contract)
	require.NoError(t, err)
}

func TestDeactivateContract(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, SupportedFeatures, nil, nil)

	k := keepers.WasmKeeper
	var mock wasmtesting.MockWasmer
	wasmtesting.MakeInstantiable(&mock)
	example := wasmkeeper.SeedNewContractInstance(t, ctx, keepers.TestKeepers, &mock)
	em := sdk.NewEventManager()

	// request no contract address -> fail
	err := k.deactivateContract(ctx, example.CreatorAddr)
	require.Error(t, err, fmt.Sprintf("no contract %s", example.CreatorAddr))

	// success case
	err = k.deactivateContract(ctx, example.Contract)
	require.NoError(t, err)

	// already inactivate contract -> fail
	err = k.deactivateContract(ctx.WithEventManager(em), example.Contract)
	require.Error(t, err, fmt.Sprintf("already inactivate contract %s", example.Contract))
}

func TestIterateInactiveContracts(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, SupportedFeatures, nil, nil)
	k := keepers.WasmKeeper

	var mock wasmtesting.MockWasmer
	wasmtesting.MakeInstantiable(&mock)
	example1 := wasmkeeper.SeedNewContractInstance(t, ctx, keepers.TestKeepers, &mock)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	example2 := wasmkeeper.SeedNewContractInstance(t, ctx, keepers.TestKeepers, &mock)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	err := k.deactivateContract(ctx, example1.Contract)
	require.NoError(t, err)
	err = k.deactivateContract(ctx, example2.Contract)
	require.NoError(t, err)

	var inactiveContracts []sdk.AccAddress
	k.IterateInactiveContracts(ctx, func(contractAddress sdk.AccAddress) (stop bool) {
		inactiveContracts = append(inactiveContracts, contractAddress)
		return false
	})
	expectList := []sdk.AccAddress{example2.Contract, example1.Contract}
	assert.Equal(t, expectList, inactiveContracts)
}
