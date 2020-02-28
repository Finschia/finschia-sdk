package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) BurnToken(ctx sdk.Context, contractID string, amount sdk.Int, from sdk.AccAddress) sdk.Error {
	err := k.isBurnable(ctx, contractID, from, from, amount)
	if err != nil {
		return err
	}

	err = k.BurnSupply(ctx, contractID, from, amount)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnToken,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
		),
	})
	return nil
}

func (k Keeper) isBurnable(ctx sdk.Context, contractID string, permissionOwner, tokenOwner sdk.AccAddress, amount sdk.Int) sdk.Error {
	if !k.HasBalance(ctx, contractID, tokenOwner, amount) {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("%v has not enough coins for %v", tokenOwner, amount))
	}

	perm := types.NewBurnPermission(contractID)
	if !k.HasPermission(ctx, permissionOwner, perm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, permissionOwner, perm)
	}
	return nil
}
