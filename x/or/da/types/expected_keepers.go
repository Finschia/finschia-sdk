package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	rolluptypes "github.com/Finschia/finschia-sdk/x/or/rollup/types"
)

type AccountKeeper interface {
	GetParams(ctx sdk.Context) (params authtypes.Params)
}

type RollupKeeper interface {
	GetAllRollup(ctx sdk.Context) (list []rolluptypes.Rollup)
	GetRollup(ctx sdk.Context, rollupName string) (val rolluptypes.Rollup, found bool)
}

type RollupInfo struct {
	ID string

	// The ratio between the cost of gas on L1 and L2.
	// This is a positive integer.
	L1ToL2GasRatio uint64
}
