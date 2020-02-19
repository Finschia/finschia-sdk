package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
	"github.com/line/link/x/iam/exported"
)

func (k Keeper) AddPermission(ctx sdk.Context, addr sdk.AccAddress, perm types.PermissionI) {
	k.iamKeeper.GrantPermission(ctx, addr, perm)
}

func (k Keeper) RevokePermission(ctx sdk.Context, addr sdk.AccAddress, perm types.PermissionI) sdk.Error {
	if !k.HasPermission(ctx, addr, perm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, addr, perm)
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

func (k Keeper) GrantPermission(ctx sdk.Context, from, to sdk.AccAddress, perm types.PermissionI) sdk.Error {
	if !k.HasPermission(ctx, from, perm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, from, perm)
	}
	k.AddPermission(ctx, to, perm)

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
