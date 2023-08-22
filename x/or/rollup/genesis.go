package rollup

import (
	sdk "github.com/Finschia/finschia-rdk/types"
	"github.com/Finschia/finschia-rdk/x/or/rollup/keeper"
	"github.com/Finschia/finschia-rdk/x/or/rollup/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, accoutKeeper types.AccountKeeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
