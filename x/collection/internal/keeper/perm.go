package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link/x/collection/internal/types"
)

func (k Keeper) AddPermission(ctx sdk.Context, contractID string, addr sdk.AccAddress, perm types.Permission) {
	accPerm := k.getAccountPermission(ctx, contractID, addr)
	accPerm.AddPermission(perm)
	k.setAccountPermission(ctx, contractID, accPerm)
}

func (k Keeper) HasPermission(ctx sdk.Context, contractID string, addr sdk.AccAddress, p types.Permission) bool {
	accPerm := k.getAccountPermission(ctx, contractID, addr)
	return accPerm.HasPermission(p)
}

func (k Keeper) GetPermissions(ctx sdk.Context, contractID string, addr sdk.AccAddress) types.Permissions {
	accPerm := k.getAccountPermission(ctx, contractID, addr)
	return accPerm.GetPermissions()
}

func (k Keeper) RevokePermission(ctx sdk.Context, contractID string, addr sdk.AccAddress, perm types.Permission) error {
	if !k.HasPermission(ctx, contractID, addr, perm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr.String(), perm.String())
	}
	accPerm := k.getAccountPermission(ctx, contractID, addr)
	accPerm.RemovePermission(perm)
	k.setAccountPermission(ctx, contractID, accPerm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRevokePermToken,
			sdk.NewAttribute(types.AttributeKeyFrom, addr.String()),
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyPerm, perm.String()),
		),
	})
	return nil
}

func (k Keeper) GrantPermission(ctx sdk.Context, contractID string, from, to sdk.AccAddress, perm types.Permission) error {
	if !k.HasPermission(ctx, contractID, from, perm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", from.String(), perm.String())
	}
	k.AddPermission(ctx, contractID, to, perm)

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
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyPerm, perm.String()),
		),
	})

	return nil
}

func (k Keeper) getAccountPermission(ctx sdk.Context, contractID string, addr sdk.AccAddress) (accPerm types.AccountPermissionI) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PermKey(contractID, addr))
	if bz != nil {
		accPerm = k.mustDecodeAccountPermission(bz)
		return accPerm
	}
	return types.NewAccountPermission(addr)
}

func (k Keeper) setAccountPermission(ctx sdk.Context, contractID string, accPerm types.AccountPermissionI) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PermKey(contractID, accPerm.GetAddress()), k.cdc.MustMarshalBinaryBare(accPerm))
}

func (k Keeper) mustDecodeAccountPermission(bz []byte) (accPerm types.AccountPermissionI) {
	err := k.cdc.UnmarshalBinaryBare(bz, &accPerm)
	if err != nil {
		panic(err)
	}
	return
}
