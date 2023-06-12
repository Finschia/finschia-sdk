package rollup

import (
	"fmt"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/or/rollup/keeper"
	"github.com/Finschia/finschia-sdk/x/or/rollup/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, accoutKeeper types.AccountKeeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	rollupModuleAddress := k.GetRollupAccount(ctx)
	if rollupModuleAddress == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	accoutKeeper.SetModuleAccount(ctx, rollupModuleAddress)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
