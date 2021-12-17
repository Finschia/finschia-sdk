package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/consortium/types"
)

func (k Keeper) GetValidatorAuth(ctx sdk.Context, valAddr sdk.ValAddress) (*types.ValidatorAuth, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.ValidatorAuthKey(valAddr)
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "validator auth not found")
	}

	var auth types.ValidatorAuth
	if err := k.cdc.UnmarshalBinaryBare(bz, &auth); err != nil {
		return nil, err
	}

	return &auth, nil
}

func (k Keeper) SetValidatorAuth(ctx sdk.Context, auth *types.ValidatorAuth) error {
	store := ctx.KVStore(k.storeKey)
	key := types.ValidatorAuthKey(sdk.ValAddress(auth.OperatorAddress))

	bz, err := k.cdc.MarshalBinaryBare(auth)
	if err != nil {
		return err
	}
	store.Set(key, bz)

	return nil
}

func (k Keeper) addPendingRejectedDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.PendingRejectedDelegationKey(delAddr, valAddr)
	store.Set(key, []byte(""))
}

func (k Keeper) deletePendingRejectedDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.PendingRejectedDelegationKey(delAddr, valAddr)
	store.Delete(key)
}

// Iterators

// IterateValidatorAuths iterates over the allowed validators
// and performs a callback function
func (k Keeper) IterateValidatorAuths(ctx sdk.Context, cb func(auth types.ValidatorAuth) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorAuthKeyPrefix)

	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var auth types.ValidatorAuth
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &auth)
		if cb(auth) {
			break
		}
	}
}

// IteratePendingRejectedDelegations iterates over the pending rejected delegations
// and performs a callback function
func (k Keeper) IteratePendingRejectedDelegations(ctx sdk.Context, cb func(delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.PendingRejectedDelegationKeyPrefix)

	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		delAddr, valAddr := types.SplitPendingRejectedDelegationKey(iter.Key())
		if cb(delAddr, valAddr) {
			break
		}
	}
}

func (k Keeper) RejectDelegations(ctx sdk.Context) {
	type DVPair struct {
		Delegator sdk.AccAddress
		Validator sdk.ValAddress
	}

	dvpairs := []DVPair{}
	k.IteratePendingRejectedDelegations(ctx, func(delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stop bool) {
		dvpairs = append(dvpairs, DVPair{delAddr, valAddr})
		return false
	})

	for _, dvpair := range dvpairs {
		delAddr, valAddr := dvpair.Delegator, dvpair.Validator

		delegation, found := k.stakingKeeper.GetDelegation(ctx, delAddr, valAddr)
		if !found {
			panic("delegation not found")
		}
		share := delegation.GetShares()
		if _, err := k.stakingKeeper.Undelegate(ctx, delAddr, valAddr, share); err != nil {
			// TODO: prevent DoS attack
			if _ /*validator*/, found := k.stakingKeeper.GetValidator(ctx, valAddr); !found {
				panic("validator must exist")
			} else {
				// k.stakingKeeper.jailValidator(ctx, validator)
			}
		} else {
			k.deletePendingRejectedDelegation(ctx, delAddr, valAddr)
		}
	}
}

// utility functions
func (k Keeper) GetValidatorAuths(ctx sdk.Context) []*types.ValidatorAuth {
	auths := []*types.ValidatorAuth{}
	k.IterateValidatorAuths(ctx, func(auth types.ValidatorAuth) (stop bool) {
		auths = append(auths, &auth)
		return false
	})
	return auths
}
