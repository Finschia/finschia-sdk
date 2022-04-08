package keeper

import (
	"github.com/line/ostracon/libs/log"

	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

// Keeper defines the foundation module Keeper
type Keeper struct {
	stakingKeeper foundation.StakingKeeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec for binary encoding/decoding.
	cdc codec.Codec
}

// NewKeeper returns a foundation keeper. It handles:
// - updating validator auths.
//
// CONTRACT: the parameter Subspace must have the param key table already initialized
func NewKeeper(
	cdc codec.Codec,
	key sdk.StoreKey,
	stakingKeeper foundation.StakingKeeper,
) Keeper {
	return Keeper{
		storeKey:      key,
		stakingKeeper: stakingKeeper,
		cdc:           cdc,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+foundation.ModuleName)
}

// Cleaning up the states
func (k Keeper) Cleanup(ctx sdk.Context) {
	valAddrs := []sdk.ValAddress{}
	k.IterateValidatorAuths(ctx, func(auth foundation.ValidatorAuth) (stop bool) {
		addr := sdk.ValAddress(auth.OperatorAddress)
		valAddrs = append(valAddrs, addr)
		return false
	})

	store := ctx.KVStore(k.storeKey)
	for _, addr := range valAddrs {
		key := validatorAuthKey(addr)
		store.Delete(key)
	}
}
