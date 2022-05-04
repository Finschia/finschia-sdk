package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/authz"
	"github.com/line/lbm-sdk/x/foundation"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

func (k Keeper) Grant(ctx sdk.Context, granter foundation.Granter, grantee sdk.AccAddress, authorization authz.Authorization) error {
	return k.setAuthorization(ctx, granter, grantee, authorization)
}

func (k Keeper) Revoke(ctx sdk.Context, granter foundation.Granter, grantee sdk.AccAddress, authorization authz.Authorization) error {
	msgTypeURL := authorization.MsgTypeURL()
	if _, err := k.GetAuthorization(ctx, granter, grantee, msgTypeURL); err != nil {
		return err
	}
	k.deleteAuthorization(ctx, granter, grantee, msgTypeURL)

	return nil
}

func (k Keeper) GetAuthorization(ctx sdk.Context, granter foundation.Granter, grantee sdk.AccAddress, msgTypeURL string) (authz.Authorization, error) {
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

func (k Keeper) setAuthorization(ctx sdk.Context, granter foundation.Granter, grantee sdk.AccAddress, authorization authz.Authorization) error {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(granter, grantee, authorization.MsgTypeURL())

	bz, err := k.cdc.MarshalInterface(authorization)
	if err != nil {
		return err
	}
	store.Set(key, bz)

	return nil
}

func (k Keeper) deleteAuthorization(ctx sdk.Context, granter foundation.Granter, grantee sdk.AccAddress, msgTypeURL string) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(granter, grantee, msgTypeURL)
	store.Delete(key)
}

func (k Keeper) Accept(ctx sdk.Context, grantee sdk.AccAddress, msg sdk.Msg) error {
	granters := map[string]foundation.Granter{
		sdk.MsgTypeURL(&stakingtypes.MsgCreateValidator{}): foundation.GRANTER_GOVERNANCE,
		sdk.MsgTypeURL(&foundation.MsgWithdrawFromTreasury{}): foundation.GRANTER_FOUNDATION,
	}

	msgTypeURL := sdk.MsgTypeURL(msg)
	granter := granters[msgTypeURL]
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
