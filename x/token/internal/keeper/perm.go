package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token/internal/types"
)

func (k Keeper) AddPermission(ctx sdk.Context, addr sdk.AccAddress, perm types.Permission) {
	accPerm := k.getAccountPermission(ctx, addr)
	accPerm.AddPermission(perm)
	k.setAccountPermission(ctx, accPerm)
}

func (k Keeper) HasPermission(ctx sdk.Context, addr sdk.AccAddress, p types.Permission) bool {
	accPerm := k.getAccountPermission(ctx, addr)
	return accPerm.HasPermission(p)
}

func (k Keeper) GetPermissions(ctx sdk.Context, addr sdk.AccAddress) types.Permissions {
	accPerm := k.getAccountPermission(ctx, addr)
	return accPerm.GetPermissions()
}

func (k Keeper) RevokePermission(ctx sdk.Context, addr sdk.AccAddress, perm types.Permission) error {
	if !k.HasPermission(ctx, addr, perm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr.String(), perm.String())
	}
	accPerm := k.getAccountPermission(ctx, addr)
	accPerm.RemovePermission(perm)
	k.setAccountPermission(ctx, accPerm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRevokePermToken,
			sdk.NewAttribute(types.AttributeKeyFrom, addr.String()),
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyPerm, perm.String()),
		),
	})
	return nil
}

func (k Keeper) GrantPermission(ctx sdk.Context, from, to sdk.AccAddress, perm types.Permission) error {
	if !k.HasPermission(ctx, from, perm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", from.String(), perm.String())
	}
	k.AddPermission(ctx, to, perm)

	// Set Account if not exists yet
	account := k.accountKeeper.GetAccount(ctx, to)
	if account == nil {
		account = k.accountKeeper.NewAccountWithAddress(ctx, to)
		k.accountKeeper.SetAccount(ctx, account)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeGrantPermToken,
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyPerm, perm.String()),
		),
	})

	return nil
}

func (k Keeper) getAccountPermission(ctx sdk.Context, addr sdk.AccAddress) (accPerm types.AccountPermissionI) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PermKey(k.getContractID(ctx), addr))
	if bz != nil {
		accPerm = k.mustDecodeAccountPermission(bz)
		return accPerm
	}
	return types.NewAccountPermission(addr)
}

func (k Keeper) setAccountPermission(ctx sdk.Context, accPerm types.AccountPermissionI) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PermKey(k.getContractID(ctx), accPerm.GetAddress()), k.cdc.MustMarshalBinaryBare(accPerm))
}

func (k Keeper) mustDecodeAccountPermission(bz []byte) (accPerm types.AccountPermissionI) {
	err := k.cdc.UnmarshalBinaryBare(bz, &accPerm)
	if err != nil {
		panic(err)
	}
	return
}
