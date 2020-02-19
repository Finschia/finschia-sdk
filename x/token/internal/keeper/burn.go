package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) BurnTokens(ctx sdk.Context, amount sdk.Coins, from sdk.AccAddress) sdk.Error {
	err := k.isBurnable(ctx, from, from, amount)
	if err != nil {
		return err
	}

	err = k.burnTokens(ctx, from, amount)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnToken,
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
		),
	})
	return nil
}

func (k Keeper) isBurnable(ctx sdk.Context, permissionOwner, tokenOwner sdk.AccAddress, amount sdk.Coins) sdk.Error {
	if !k.hasEnoughCoins(ctx, tokenOwner, amount) {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("%v has not enough coins for %v", tokenOwner, amount))
	}

	for _, coin := range amount {
		perm := types.NewBurnPermission(coin.Denom)
		if !k.HasPermission(ctx, permissionOwner, perm) {
			return types.ErrTokenNoPermission(types.DefaultCodespace, permissionOwner, perm)
		}
	}
	return nil
}

func (k Keeper) hasEnoughCoins(ctx sdk.Context, from sdk.AccAddress, amount sdk.Coins) bool {
	return k.bankKeeper.GetCoins(ctx, from).IsAllGTE(amount)
}

func (k Keeper) burnTokens(ctx sdk.Context, from sdk.AccAddress, amount sdk.Coins) sdk.Error {
	moduleAddr := k.supplyKeeper.GetModuleAddress(types.ModuleName)
	if moduleAddr == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", types.ModuleName))
	}

	_, err := k.bankKeeper.SubtractCoins(ctx, from, amount)
	if err != nil {
		return err
	}

	_, err = k.bankKeeper.AddCoins(ctx, moduleAddr, amount)
	if err != nil {
		return err
	}

	err = k.supplyKeeper.BurnCoins(ctx, types.ModuleName, amount)
	if err != nil {
		return err
	}
	return nil
}
