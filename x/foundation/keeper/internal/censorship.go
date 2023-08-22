package internal

import (
	sdk "github.com/Finschia/finschia-rdk/types"
	sdkerrors "github.com/Finschia/finschia-rdk/types/errors"
	"github.com/Finschia/finschia-rdk/x/foundation"
)

func (k Keeper) GetCensorship(ctx sdk.Context, msgTypeURL string) (*foundation.Censorship, error) {
	store := ctx.KVStore(k.storeKey)
	key := censorshipKey(msgTypeURL)
	bz := store.Get(key)
	if bz == nil {
		return nil, sdkerrors.ErrNotFound.Wrap("censorship not found")
	}

	var censorship foundation.Censorship
	k.cdc.MustUnmarshal(bz, &censorship)

	return &censorship, nil
}

func (k Keeper) UpdateCensorship(ctx sdk.Context, censorship foundation.Censorship) error {
	url := censorship.MsgTypeUrl

	oldCensorship, err := k.GetCensorship(ctx, url)
	if err != nil {
		return err
	}

	newAuthority := censorship.Authority
	oldAuthority := oldCensorship.Authority
	if newAuthority >= oldAuthority {
		return sdkerrors.ErrInvalidRequest.Wrapf("bad transition; %s -> %s over %s", oldAuthority, newAuthority, url)
	}

	// clean up relevant authorizations
	if newAuthority == foundation.CensorshipAuthorityUnspecified {
		k.pruneAuthorizations(ctx, url)
	}

	k.SetCensorship(ctx, censorship)

	return nil
}

func (k Keeper) SetCensorship(ctx sdk.Context, censorship foundation.Censorship) {
	store := ctx.KVStore(k.storeKey)
	key := censorshipKey(censorship.MsgTypeUrl)

	if censorship.Authority == foundation.CensorshipAuthorityUnspecified {
		store.Delete(key)
		return
	}

	bz := k.cdc.MustMarshal(&censorship)
	store.Set(key, bz)
}

func (k Keeper) iterateCensorships(ctx sdk.Context, fn func(censorship foundation.Censorship) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, censorshipKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var censorship foundation.Censorship
		k.cdc.MustUnmarshal(iterator.Value(), &censorship)

		if stop := fn(censorship); stop {
			break
		}
	}
}

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

func (k Keeper) pruneAuthorizations(ctx sdk.Context, msgTypeURL string) {
	var pruning []foundation.GrantAuthorization
	k.iterateAuthorizations(ctx, func(grantee sdk.AccAddress, authorization foundation.Authorization) (stop bool) {
		if authorization.MsgTypeURL() == msgTypeURL {
			grant := foundation.GrantAuthorization{
				Grantee: grantee.String(),
			}.WithAuthorization(authorization)

			pruning = append(pruning, *grant)
		}
		return false
	})

	for _, grant := range pruning {
		k.deleteAuthorization(ctx, sdk.MustAccAddressFromBech32(grant.Grantee), grant.GetAuthorization().MsgTypeURL())
	}
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
