package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/keeper"
	"github.com/line/link/x/collection/internal/types"
)

func handleMsgBurnNFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBurnNFT) sdk.Result {
	err := keeper.BurnNFT(ctx, msg.ContractID, msg.From, msg.TokenIDs...)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgBurnNFTFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBurnNFTFrom) sdk.Result {
	err := keeper.BurnNFTFrom(ctx, msg.ContractID, msg.Proxy, msg.From, msg.TokenIDs...)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Proxy.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgBurnFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBurnFT) sdk.Result {
	err := keeper.BurnFT(ctx, msg.ContractID, msg.From, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgBurnFTFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBurnFTFrom) sdk.Result {
	err := keeper.BurnFTFrom(ctx, msg.ContractID, msg.Proxy, msg.From, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Proxy.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
