package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
	wasmkeeper "github.com/line/lbm-sdk/x/wasm/keeper"
	"github.com/line/lbm-sdk/x/wasm/keeper/wasmtesting"
	lbmwasmtypes "github.com/line/lbm-sdk/x/wasm/lbm/types"
)

func TestValidateDeactivateContractProposal(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, "staking", nil, nil)
	govKeeper, wasmKeeper := keepers.GovKeeper, keepers.WasmKeeper

	var mock wasmtesting.MockWasmer
	wasmtesting.MakeInstantiable(&mock)
	example := wasmkeeper.SeedNewContractInstance(t, ctx, keepers.TestKeepers, &mock)

	src := lbmwasmtypes.DeactivateContract{
		Title:       "Foo",
		Description: "Bar",
		Contract:    example.Contract.String(),
	}

	em := sdk.NewEventManager()

	// when stored
	storedProposal, err := govKeeper.SubmitProposal(ctx, &src)
	require.NoError(t, err)

	// proposal execute
	handler := govKeeper.Router().GetRoute(storedProposal.ProposalRoute())
	err = handler(ctx.WithEventManager(em), storedProposal.GetContent())
	require.NoError(t, err)

	// then
	isInactive := wasmKeeper.IsInactiveContract(ctx, example.Contract)
	require.True(t, isInactive)
}

func TestActivateContractProposal(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, "staking", nil, nil)
	govKeeper, wasmKeeper := keepers.GovKeeper, keepers.WasmKeeper

	var mock wasmtesting.MockWasmer
	wasmtesting.MakeInstantiable(&mock)
	example := wasmkeeper.SeedNewContractInstance(t, ctx, keepers.TestKeepers, &mock)
	// set deactivate
	err := wasmKeeper.deactivateContract(ctx, example.Contract)
	require.NoError(t, err)

	src := lbmwasmtypes.ActivateContract{
		Title:       "Foo",
		Description: "Bar",
		Contract:    example.Contract.String(),
	}

	em := sdk.NewEventManager()

	// when stored
	storedProposal, err := govKeeper.SubmitProposal(ctx, &src)
	require.NoError(t, err)

	// proposal execute
	handler := govKeeper.Router().GetRoute(storedProposal.ProposalRoute())
	err = handler(ctx.WithEventManager(em), storedProposal.GetContent())
	require.NoError(t, err)

	// then
	isInactive := wasmKeeper.IsInactiveContract(ctx, example.Contract)
	require.False(t, isInactive)
}
