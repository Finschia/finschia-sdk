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
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "no validator auth found for: %s", valAddr)
	}

	var auth types.ValidatorAuth
	if err := k.cdc.Unmarshal(bz, &auth); err != nil {
		return nil, err
	}

	return &auth, nil
}

func (k Keeper) SetValidatorAuth(ctx sdk.Context, auth *types.ValidatorAuth) error {
	store := ctx.KVStore(k.storeKey)
	key := types.ValidatorAuthKey(sdk.ValAddress(auth.OperatorAddress))

	bz, err := k.cdc.Marshal(auth)
	if err != nil {
		return err
	}
	store.Set(key, bz)

	return nil
}

// Iterators

// IterateValidatorAuths iterates over the validator auths
// and performs a callback function
func (k Keeper) IterateValidatorAuths(ctx sdk.Context, cb func(auth types.ValidatorAuth) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorAuthKeyPrefix)

	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var auth types.ValidatorAuth
		k.cdc.MustUnmarshal(iter.Value(), &auth)
		if cb(auth) {
			break
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
