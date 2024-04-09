package fswap_test

import (
	"testing"

	keepertest "fswap/testutil/keeper"
	"fswap/testutil/nullify"
	"fswap/x/fswap"
	"fswap/x/fswap/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.FswapKeeper(t)
	fswap.InitGenesis(ctx, *k, genesisState)
	got := fswap.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
