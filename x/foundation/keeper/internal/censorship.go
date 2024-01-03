package internal

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (k Keeper) GetCensorship(ctx sdk.Context, msgTypeURL string) (*foundation.Censorship, error) {
	store := k.storeService.OpenKVStore(ctx)
	key := censorshipKey(msgTypeURL)
	bz, err := store.Get(key)
	if err != nil {
		return nil, err
	}
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
	store := k.storeService.OpenKVStore(ctx)
	key := censorshipKey(censorship.MsgTypeUrl)

	if censorship.Authority == foundation.CensorshipAuthorityUnspecified {
		store.Delete(key)
		return
	}

	bz := k.cdc.MustMarshal(&censorship)
	store.Set(key, bz)
}

func (k Keeper) iterateCensorships(ctx sdk.Context, fn func(censorship foundation.Censorship) (stop bool)) {
	store := k.storeService.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(store)
	iterator := storetypes.KVStorePrefixIterator(adapter, censorshipKeyPrefix)
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

	granteeStr, err := k.addressCodec().BytesToString(grantee)
	if err != nil {
		panic(err)
	}
	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventGrant{
		Grantee:       granteeStr,
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

	granteeStr, err := k.addressCodec().BytesToString(grantee)
	if err != nil {
		panic(err)
	}
	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventRevoke{
		Grantee:    granteeStr,
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
			granteeStr, err := k.addressCodec().BytesToString(grantee)
			if err != nil {
				panic(err)
			}

			grant := foundation.GrantAuthorization{
				Grantee: granteeStr,
			}.WithAuthorization(authorization)

			pruning = append(pruning, *grant)
		}
		return false
	})

	for _, grant := range pruning {
		grantee, err := k.addressCodec().StringToBytes(grant.Grantee)
		if err != nil {
			panic(err)
		}

		k.deleteAuthorization(ctx, grantee, grant.GetAuthorization().MsgTypeURL())
	}
}

func (k Keeper) GetAuthorization(ctx sdk.Context, grantee sdk.AccAddress, msgTypeURL string) (foundation.Authorization, error) {
	store := k.storeService.OpenKVStore(ctx)
	key := grantKey(grantee, msgTypeURL)
	bz, err := store.Get(key)
	if err != nil {
		return nil, err
	}
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
	store := k.storeService.OpenKVStore(ctx)
	key := grantKey(grantee, authorization.MsgTypeURL())

	bz, err := k.cdc.MarshalInterface(authorization)
	if err != nil {
		panic(err)
	}
	store.Set(key, bz)
}

func (k Keeper) deleteAuthorization(ctx sdk.Context, grantee sdk.AccAddress, msgTypeURL string) {
	store := k.storeService.OpenKVStore(ctx)
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
	store := k.storeService.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(store)
	iterator := storetypes.KVStorePrefixIterator(adapter, prefix)
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
