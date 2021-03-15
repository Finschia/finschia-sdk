package types

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/line/lbm-sdk/v2/codec"
	sdk "github.com/line/lbm-sdk/v2/types"
	auth "github.com/line/lbm-sdk/v2/x/auth/types"
	bankexported "github.com/line/lbm-sdk/v2/x/bank/exported"
)

// StakingKeeper defines the expected staking keeper (noalias)
type StakingKeeper interface {
	ApplyAndReturnValidatorSetUpdates(sdk.Context) (updates []abci.ValidatorUpdate, err error)
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	NewAccount(sdk.Context, auth.AccountI) auth.AccountI
	SetAccount(sdk.Context, auth.AccountI)
	IterateAccounts(ctx sdk.Context, process func(auth.AccountI) (stop bool))
}

// GenesisAccountsIterator defines the expected iterating genesis accounts object (noalias)
type GenesisAccountsIterator interface {
	IterateGenesisAccounts(
		cdc *codec.LegacyAmino,
		appGenesis map[string]json.RawMessage,
		cb func(auth.AccountI) (stop bool),
	)
}

// GenesisAccountsIterator defines the expected iterating genesis accounts object (noalias)
type GenesisBalancesIterator interface {
	IterateGenesisBalances(
		cdc codec.JSONMarshaler,
		appGenesis map[string]json.RawMessage,
		cb func(bankexported.GenesisBalance) (stop bool),
	)
}
