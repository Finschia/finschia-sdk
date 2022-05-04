package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/authz"
	"github.com/line/lbm-sdk/x/foundation"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

func (k Keeper) Grant(ctx sdk.Context, granter string, grantee sdk.AccAddress, authorization authz.Authorization) error {
	msgTypeURL := authorization.MsgTypeURL()
	if granter != getGranter(msgTypeURL) {
		return sdkerrors.ErrInvalidRequest.Wrapf("granter %s cannot grant a msg type of %s", granter, msgTypeURL)
	}

	return k.setAuthorization(ctx, granter, grantee, authorization)
}

func (k Keeper) Revoke(ctx sdk.Context, granter string, grantee sdk.AccAddress, msgTypeURL string) error {
	if granter != getGranter(msgTypeURL) {
		return sdkerrors.ErrInvalidRequest.Wrapf("granter %s cannot revoke a msg type of %s", granter, msgTypeURL)
	}

	if _, err := k.GetAuthorization(ctx, granter, grantee, msgTypeURL); err != nil {
		return err
	}
	k.deleteAuthorization(ctx, granter, grantee, msgTypeURL)

	return nil
}

func (k Keeper) GetAuthorization(ctx sdk.Context, granter string, grantee sdk.AccAddress, msgTypeURL string) (authz.Authorization, error) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(granter, grantee, msgTypeURL)
	bz := store.Get(key)
	if bz == nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap("authorization not found")
	}

	var auth authz.Authorization
	if err := k.cdc.UnmarshalInterface(bz, auth); err != nil {
		return nil, err
	}

	return auth, nil
}

func (k Keeper) setAuthorization(ctx sdk.Context, granter string, grantee sdk.AccAddress, authorization authz.Authorization) error {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(granter, grantee, authorization.MsgTypeURL())

	bz, err := k.cdc.MarshalInterface(authorization)
	if err != nil {
		return err
	}
	store.Set(key, bz)

	return nil
}

func (k Keeper) deleteAuthorization(ctx sdk.Context, granter string, grantee sdk.AccAddress, msgTypeURL string) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(granter, grantee, msgTypeURL)
	store.Delete(key)
}

func getGranter(msgTypeURL string) string {
	granters := map[string]string{
		foundation.CreateValidatorAuthorization{}.MsgTypeURL(): govtypes.ModuleName,
		foundation.WithdrawFromTreasuryAuthorization{}.MsgTypeURL(): foundation.ModuleName,
	}
	return granters[msgTypeURL]
}

func (k Keeper) Accept(ctx sdk.Context, grantee sdk.AccAddress, msg sdk.Msg) error {
	msgTypeURL := sdk.MsgTypeURL(msg)
	granter := getGranter(msgTypeURL)
	authorization, err := k.GetAuthorization(ctx, granter, grantee, msgTypeURL)
	if err != nil {
		return err
	}

	resp, err := authorization.Accept(ctx, msg)
	if err != nil {
		return err
	}

	if resp.Delete {
		k.deleteAuthorization(ctx, granter, grantee, msgTypeURL)
	} else if resp.Updated != nil {
		if err := k.setAuthorization(ctx, granter, grantee, resp.Updated); err != nil {
			return err
		}
	}

	if !resp.Accept {
		return sdkerrors.ErrUnauthorized
	}

	return nil
}

func (k Keeper) iterateAuthorizations(ctx sdk.Context, grantee string, fn func(granter string, grantee sdk.AccAddress, authorization authz.Authorization) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefix := append(grantKeyPrefix, grantee...)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var authorization authz.Authorization
		if err := k.cdc.UnmarshalInterface(iterator.Value(), &authorization); err != nil {
			panic(err)
		}

		granter, grantee, _ := splitGrantKey(iterator.Key())
		if stop := fn(granter, grantee, authorization); stop {
			break
		}
	}
}
