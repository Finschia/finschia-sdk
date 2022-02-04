package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

// InitGenesis new token genesis
func (k Keeper) InitGenesis(ctx sdk.Context, data *token.GenesisState) {
	panic("Not implemented")
}

// ExportGenesis returns a GenesisState for a given context.
func (k Keeper) ExportGenesis(ctx sdk.Context) *token.GenesisState {
	panic("Not implemented")
}
