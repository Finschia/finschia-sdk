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

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

type SequencerKeeper interface {
	GetRollupInfo(ctx sdk.Context, rollupID string) (RollupInfo, error)
}

type RollupInfo struct {
	ID             string
	L1ToL2GasRatio uint64
}
