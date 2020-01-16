package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) MintTokens(ctx sdk.Context, amount sdk.Coins, to sdk.AccAddress) sdk.Error {
	for _, coin := range amount {
		token, err := k.GetToken(ctx, coin.Denom)
		if err != nil {
			return err
		}

		ft, ok := token.(types.FT)
		if !ok {
			return types.ErrTokenNotMintable(types.DefaultCodespace, token.GetDenom())
		}

		if !ft.GetMintable() {
			return types.ErrTokenNotMintable(types.DefaultCodespace, ft.GetDenom())
		}
		perm := types.NewMintPermission(coin.Denom)
		if !k.HasPermission(ctx, to, perm) {
			return types.ErrTokenPermission(types.DefaultCodespace, to, perm)
		}
	}
	return k.mintTokens(ctx, amount, to)
}
func (k Keeper) mintTokens(ctx sdk.Context, amount sdk.Coins, to sdk.AccAddress) sdk.Error {
	err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, amount)
	if err != nil {
		return err
	}

	err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, amount)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintToken,
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
		),
	})

	return nil
}

func (k Keeper) BurnTokens(ctx sdk.Context, amount sdk.Coins, from sdk.AccAddress) sdk.Error {
	if !k.hasEnoughCoins(ctx, amount, from) {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("%v has not enough coins for %v", from, amount))
	}

	return k.burnTokens(ctx, amount, from)
}

func (k Keeper) hasEnoughCoins(ctx sdk.Context, amount sdk.Coins, from sdk.AccAddress) bool {
	return k.accountKeeper.GetAccount(ctx, from).GetCoins().IsAllGTE(amount)
}

func (k Keeper) burnTokens(ctx sdk.Context, amount sdk.Coins, from sdk.AccAddress) sdk.Error {
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, amount)
	if err != nil {
		return err
	}

	err = k.supplyKeeper.BurnCoins(ctx, types.ModuleName, amount)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnToken,
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyTo, from.String()),
		),
	})

	return nil
}
