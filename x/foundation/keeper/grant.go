package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/foundation"
)

func (k Keeper) Grant(ctx sdk.Context, grantee sdk.AccAddress, authorization foundation.Authorization) error {
	msgTypeURL := authorization.MsgTypeURL()
	if !k.IsCensoredMessage(ctx, msgTypeURL) {
		return sdkerrors.ErrInvalidRequest.Wrapf("%s is not being censored", msgTypeURL)
	}

	if _, err := k.GetAuthorization(ctx, grantee, msgTypeURL); err == nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("authorization for %s already exists", msgTypeURL)
	}

	k.setAuthorization(ctx, grantee, authorization)

	any, err := foundation.SetAuthorization(authorization)
	if err != nil {
		return err
	}

	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventGrant{
		Grantee:       grantee.String(),
		Authorization: any,
	}); err != nil {
		panic(err)
	}

	return nil
}

func (k Keeper) Revoke(ctx sdk.Context, grantee sdk.AccAddress, msgTypeURL string) error {
	if _, err := k.GetAuthorization(ctx, grantee, msgTypeURL); err != nil {
		return err
	}
	k.deleteAuthorization(ctx, grantee, msgTypeURL)

	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventRevoke{
		Grantee:    grantee.String(),
		MsgTypeUrl: msgTypeURL,
	}); err != nil {
		panic(err)
	}

	return nil
}

func (k Keeper) GetAuthorization(ctx sdk.Context, grantee sdk.AccAddress, msgTypeURL string) (foundation.Authorization, error) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(grantee, msgTypeURL)
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

func (k Keeper) setAuthorization(ctx sdk.Context, grantee sdk.AccAddress, authorization foundation.Authorization) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(grantee, authorization.MsgTypeURL())

	bz, err := k.cdc.MarshalInterface(authorization)
	if err != nil {
		panic(err)
	}
	store.Set(key, bz)
}

func (k Keeper) deleteAuthorization(ctx sdk.Context, grantee sdk.AccAddress, msgTypeURL string) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(grantee, msgTypeURL)
	store.Delete(key)
}

func (k Keeper) Accept(ctx sdk.Context, grantee sdk.AccAddress, msg sdk.Msg) error {
	msgTypeURL := sdk.MsgTypeURL(msg)

	// check whether the msg is being censored
	if !k.IsCensoredMessage(ctx, msgTypeURL) {
		return nil
	}

	authorization, err := k.GetAuthorization(ctx, grantee, msgTypeURL)
	if err != nil {
		return err
	}

	resp, err := authorization.Accept(ctx, msg)
	if err != nil {
		return err
	}

	if resp.Delete {
		k.deleteAuthorization(ctx, grantee, msgTypeURL)
	} else if resp.Updated != nil {
		k.setAuthorization(ctx, grantee, resp.Updated)
	}

	if !resp.Accept {
		return sdkerrors.ErrUnauthorized
	}

	return nil
}

func (k Keeper) iterateAuthorizations(ctx sdk.Context, fn func(grantee sdk.AccAddress, authorization foundation.Authorization) (stop bool)) {
	k.iterateAuthorizationsImpl(ctx, grantKeyPrefix, fn)
}

func (k Keeper) iterateAuthorizationsImpl(ctx sdk.Context, prefix []byte, fn func(grantee sdk.AccAddress, authorization foundation.Authorization) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var authorization foundation.Authorization
		if err := k.cdc.UnmarshalInterface(iterator.Value(), &authorization); err != nil {
			panic(err)
		}

		grantee, _ := splitGrantKey(iterator.Key())
		if stop := fn(grantee, authorization); stop {
			break
		}
	}
}
