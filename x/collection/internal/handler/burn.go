package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/keeper"
	"github.com/line/link/x/collection/internal/types"
)

func handleMsgBurnCNFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBurnCNFT) sdk.Result {
	err := keeper.BurnCNFT(ctx, msg.From, msg.Symbol, msg.TokenID)
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

func handleMsgBurnCNFTFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBurnCNFTFrom) sdk.Result {
	err := keeper.BurnCNFTFrom(ctx, msg.Proxy, msg.From, msg.Symbol, msg.TokenID)
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

func handleMsgBurnCFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBurnCFT) sdk.Result {
	err := keeper.BurnCFT(ctx, msg.From, msg.Amount)
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

func handleMsgBurnCFTFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBurnCFTFrom) sdk.Result {
	err := keeper.BurnCFTFrom(ctx, msg.Proxy, msg.From, msg.Amount)
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
