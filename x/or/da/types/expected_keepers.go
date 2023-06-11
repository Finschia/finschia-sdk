package types

import (
	"github.com/Finschia/finschia-sdk/x/auth/types"

	sdk "github.com/Finschia/finschia-sdk/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

type RollupKeeper interface {
	GetRollupInfo(ctx sdk.Context, rollupID string) (RollupInfo, error)
	GetRegisteredRollups(ctx sdk.Context) []string
}

type RollupInfo struct {
	ID string

	// The ratio between the cost of gas on L1 and L2.
	// This is a positive integer.
	L1ToL2GasRatio uint64
}
