package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	rutypes "github.com/Finschia/finschia-sdk/x/or/rollup/types"
)

type AccountKeeper interface {
	GetParams(ctx sdk.Context) (params authtypes.Params)
}

type RollupKeeper interface {
	GetRollup(ctx sdk.Context, rollupName string) (val rutypes.Rollup, found bool)
	GetAllRollup(ctx sdk.Context) (list []rutypes.Rollup)
	GetSequencersByRollupName(ctx sdk.Context, rollupName string) (val rutypes.SequencersByRollup, found bool)
}
