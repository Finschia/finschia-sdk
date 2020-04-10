package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link/x/iam/exported"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) AddPermission(ctx sdk.Context, addr sdk.AccAddress, perm types.PermissionI) {
	k.iamKeeper.GrantPermission(ctx, addr, perm)
}

func (k Keeper) RevokePermission(ctx sdk.Context, addr sdk.AccAddress, perm types.PermissionI) error {
	if !k.HasPermission(ctx, addr, perm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr.String(), perm.String())
	}
	k.iamKeeper.RevokePermission(ctx, addr, perm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRevokePermToken,
			sdk.NewAttribute(types.AttributeKeyFrom, addr.String()),
			sdk.NewAttribute(types.AttributeKeyResource, perm.GetResource()),
			sdk.NewAttribute(types.AttributeKeyAction, perm.GetAction()),
		),
	})
	return nil
}

func (k Keeper) HasPermission(ctx sdk.Context, addr sdk.AccAddress, perm types.PermissionI) bool {
	return k.iamKeeper.HasPermission(ctx, addr, perm)
}

func (k Keeper) GrantPermission(ctx sdk.Context, from, to sdk.AccAddress, perm types.PermissionI) error {
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
			sdk.NewAttribute(types.AttributeKeyResource, perm.GetResource()),
			sdk.NewAttribute(types.AttributeKeyAction, perm.GetAction()),
		),
	})

	return nil
}

func (k Keeper) GetPermissions(ctx sdk.Context, addr sdk.AccAddress) []exported.PermissionI {
	return k.iamKeeper.GetPermissions(ctx, addr)
}
