package settlement_test

import (
	"testing"

	"github.com/Finschia/finschia-sdk/x/or/settlement"
	keepertest "github.com/Finschia/finschia-sdk/x/or/settlement/testutil/keeper"
	"github.com/Finschia/finschia-sdk/x/or/settlement/testutil/nullify"
	"github.com/Finschia/finschia-sdk/x/or/settlement/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{}

	k, ctx := keepertest.SettlementKeeper(t)
	settlement.InitGenesis(ctx, *k, genesisState)
	got := settlement.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
