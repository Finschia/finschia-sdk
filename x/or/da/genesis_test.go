package da_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/x/or/da"
	datest "github.com/Finschia/finschia-sdk/x/or/da/testutil"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	k, ctx := datest.DaKeeper(t)
	da.InitGenesis(ctx, k, genesisState)
	got := da.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	// TODO: compare got & genesisState for each field
}
