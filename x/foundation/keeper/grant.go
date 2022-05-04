package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/authz"
	"github.com/line/lbm-sdk/x/foundation"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
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
		sdk.MsgTypeURL(&stakingtypes.MsgCreateValidator{}): govtypes.ModuleName,
		sdk.MsgTypeURL(&foundation.MsgWithdrawFromTreasury{}): foundation.ModuleName,
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
