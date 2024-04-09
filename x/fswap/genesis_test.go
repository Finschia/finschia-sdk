// package fswap_test

// import (
// 	"testing"

// 	keepertest "github.com/Finschia/finschia-sdk/testutil/keeper"
// 	"github.com/Finschia/finschia-sdk/testutil/nullify"

// 	"github.com/Finschia/finschia-sdk/x/fswap"
// 	"github.com/Finschia/finschia-sdk/x/fswap/types"
// 	"github.com/stretchr/testify/require"
// )

// func TestGenesis(t *testing.T) {
// 	genesisState := types.GenesisState{
// 		Params: types.DefaultParams(),

// 		// this line is used by starport scaffolding # genesis/test/state
// 	}

// 	k, ctx := keepertest.FswapKeeper(t)
// 	fswap.InitGenesis(ctx, *k, genesisState)
// 	got := fswap.ExportGenesis(ctx, *k)
// 	require.NotNil(t, got)

// 	nullify.Fill(&genesisState)
// 	nullify.Fill(got)

// 	// this line is used by starport scaffolding # genesis/test/assert
// }
