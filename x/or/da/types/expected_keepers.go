package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
)

type AccountKeeper interface {
	GetParams(ctx sdk.Context) (params authtypes.Params)
}

type RollupKeeper interface {
	GetRollupInfo(ctx sdk.Context, name string) (RollupInfo, error)
	GetRegisteredRollups(ctx sdk.Context) []string
}

type RollupInfo struct {
	ID string

	// The ratio between the cost of gas on L1 and L2.
	// This is a positive integer.
	L1ToL2GasRatio uint64
}
