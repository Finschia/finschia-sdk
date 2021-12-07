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
	slashingKeeper types.SlashingKeeper

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
	stakingKeeper types.StakingKeeper, slashingKeeper types.SlashingKeeper,
) Keeper {
	return Keeper{
		storeKey:       key,
		stakingKeeper:  stakingKeeper,
		slashingKeeper: slashingKeeper,
		cdc:            cdc,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) GetEnabled(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.EnabledKeyPrefix
	bz := store.Get(key)

	return bz != nil
}

func (k Keeper) SetEnabled(ctx sdk.Context, enable bool) error {
	store := ctx.KVStore(k.storeKey)
	key := types.EnabledKeyPrefix

	if enable {
		// you can only enable consortium module at genesis
		if ctx.BlockHeight() != int64(0) {
			panic("Consortium module cannot enabled after genesis")
		}
		store.Set(key, []byte(""))
	} else {
		store.Delete(key)

		k.Cleanup(ctx)
	}

	return nil
}

func (k Keeper) GetAllowedValidator(ctx sdk.Context, valAddr sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.AllowedValidatorKey(valAddr)
	bz := store.Get(key)

	return bz != nil
}

func (k Keeper) SetAllowedValidator(ctx sdk.Context, valAddr sdk.ValAddress, set bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.AllowedValidatorKey(valAddr)

	if set {
		store.Set(key, []byte(""))
	} else {
		store.Delete(key)
	}
}

func (k Keeper) GetDeniedValidator(ctx sdk.Context, valAddr sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.DeniedValidatorKey(valAddr)
	bz := store.Get(key)

	return bz != nil
}

func (k Keeper) SetDeniedValidator(ctx sdk.Context, valAddr sdk.ValAddress, set bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.DeniedValidatorKey(valAddr)

	if set {
		store.Set(key, []byte(""))
	} else {
		store.Delete(key)
	}
}

// Iterators

// IterateAllowedValidators iterates over the allowed validators
// and performs a callback function
func (k Keeper) IterateAllowedValidators(ctx sdk.Context, cb func(valAddr sdk.ValAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AllowedValidatorKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		addr := types.SplitValidatorKey(iterator.Key())
		if cb(addr) {
			break
		}
	}
}

// IterateDeniedValidators iterates over the denied validators
// and performs a callback function
func (k Keeper) IterateDeniedValidators(ctx sdk.Context, cb func(valAddr sdk.ValAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DeniedValidatorKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		addr := types.SplitValidatorKey(iterator.Key())
		if cb(addr) {
			break
		}
	}
}

func (k Keeper) TombstoneValidators(ctx sdk.Context) {
	for _, addr := range k.GetDeniedValidators(ctx) {
		validator := k.stakingKeeper.Validator(ctx, addr)
		consAddr, err := validator.GetConsAddr()
		if err == nil {
			k.slashingKeeper.Tombstone(ctx, consAddr)
		}

		k.SetDeniedValidator(ctx, addr, false)
	}
}

// Clean up the states
func (k Keeper) Cleanup(ctx sdk.Context) {
	for _, addr := range k.GetAllowedValidators(ctx) {
		k.SetAllowedValidator(ctx, addr, false)
	}
}

// utility functions
func (k Keeper) GetAllowedValidators(ctx sdk.Context) []sdk.ValAddress {
	addrs := []sdk.ValAddress{}
	k.IterateAllowedValidators(ctx, func(addr sdk.ValAddress) (stop bool) {
		addrs = append(addrs, addr)
		return false
	})
	return addrs
}

func (k Keeper) GetDeniedValidators(ctx sdk.Context) []sdk.ValAddress {
	addrs := []sdk.ValAddress{}
	k.IterateDeniedValidators(ctx, func(addr sdk.ValAddress) (stop bool) {
		addrs = append(addrs, addr)
		return false
	})
	return addrs
}
