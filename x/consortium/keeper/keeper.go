package keeper

import (
	"github.com/line/ostracon/libs/log"

	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/consortium/types"
)

// Keeper defines the consortium module Keeper
type Keeper struct {
	stakingKeeper  types.StakingKeeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec for binary encoding/decoding.
	cdc codec.BinaryMarshaler
}

// NewKeeper returns a consortium keeper. It handles:
// - editing allowed validators list
//
// CONTRACT: the parameter Subspace must have the param key table already initialized
func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey,
	stakingKeeper types.StakingKeeper,
) Keeper {
	return Keeper{
		storeKey:       key,
		stakingKeeper:  stakingKeeper,
		cdc:            cdc,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// Clean up the states
func (k Keeper) Cleanup(ctx sdk.Context) {
	valAddrs := []sdk.ValAddress{}
	k.IterateValidatorAuths(ctx, func(auth types.ValidatorAuth) (stop bool) {
		addr := sdk.ValAddress(auth.OperatorAddress)
		valAddrs = append(valAddrs, addr)
		return false
	})

	store := ctx.KVStore(k.storeKey)
	for _, addr := range valAddrs {
		key := types.ValidatorAuthKey(addr)
		store.Delete(key)
	}
}
