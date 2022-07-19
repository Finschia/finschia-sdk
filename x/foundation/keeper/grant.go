package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/foundation"
)

func (k Keeper) Grant(ctx sdk.Context, granter string, grantee sdk.AccAddress, authorization foundation.Authorization) error {
	if _, err := k.GetAuthorization(ctx, granter, grantee, authorization.MsgTypeURL()); err == nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("authorization for %s already exists", authorization.MsgTypeURL())
	}

	k.setAuthorization(ctx, granter, grantee, authorization)

	any, err := foundation.SetAuthorization(authorization)
	if err != nil {
		return err
	}
	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventGrant{
		Granter:       granter,
		Grantee:       grantee.String(),
		Authorization: any,
	}); err != nil {
		panic(err)
	}

	return nil
}

func (k Keeper) Revoke(ctx sdk.Context, granter string, grantee sdk.AccAddress, msgTypeURL string) error {
	if _, err := k.GetAuthorization(ctx, granter, grantee, msgTypeURL); err != nil {
		return err
	}
	k.deleteAuthorization(ctx, granter, grantee, msgTypeURL)

	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventRevoke{
		Granter:    granter,
		Grantee:    grantee.String(),
		MsgTypeUrl: msgTypeURL,
	}); err != nil {
		panic(err)
	}

	return nil
}

func (k Keeper) GetAuthorization(ctx sdk.Context, granter string, grantee sdk.AccAddress, msgTypeURL string) (foundation.Authorization, error) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(grantee, msgTypeURL, granter)
	bz := store.Get(key)
	if bz == nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap("authorization not found")
	}

	var auth foundation.Authorization
	if err := k.cdc.UnmarshalInterface(bz, &auth); err != nil {
		panic(err)
	}

	return auth, nil
}

func (k Keeper) setAuthorization(ctx sdk.Context, granter string, grantee sdk.AccAddress, authorization foundation.Authorization) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(grantee, authorization.MsgTypeURL(), granter)

	bz, err := k.cdc.MarshalInterface(authorization)
	if err != nil {
		panic(err)
	}
	store.Set(key, bz)
}

func (k Keeper) deleteAuthorization(ctx sdk.Context, granter string, grantee sdk.AccAddress, msgTypeURL string) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(grantee, msgTypeURL, granter)
	store.Delete(key)
}

func (k Keeper) Accept(ctx sdk.Context, granter string, grantee sdk.AccAddress, msg sdk.Msg) error {
	msgTypeURL := sdk.MsgTypeURL(msg)
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
		k.setAuthorization(ctx, granter, grantee, resp.Updated)
	}

	if !resp.Accept {
		return sdkerrors.ErrUnauthorized
	}

	return nil
}

func (k Keeper) iterateAuthorizations(ctx sdk.Context, fn func(granter string, grantee sdk.AccAddress, authorization foundation.Authorization) (stop bool)) {
	k.iterateAuthorizationsImpl(ctx, grantKeyPrefix, fn)
}

func (k Keeper) iterateAuthorizationsImpl(ctx sdk.Context, prefix []byte, fn func(granter string, grantee sdk.AccAddress, authorization foundation.Authorization) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var authorization foundation.Authorization
		if err := k.cdc.UnmarshalInterface(iterator.Value(), &authorization); err != nil {
			panic(err)
		}

		grantee, _, granter := splitGrantKey(iterator.Key())
		if stop := fn(granter, grantee, authorization); stop {
			break
		}
	}
}
