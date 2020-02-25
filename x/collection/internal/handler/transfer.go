package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/keeper"
	"github.com/line/link/x/collection/internal/types"
)

func handleMsgTransferFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTransferFT) sdk.Result {
	err := keeper.TransferFT(ctx, msg.From, msg.To, msg.Symbol, msg.Amount...)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgTransferNFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTransferNFT) sdk.Result {
	err := keeper.TransferNFT(ctx, msg.From, msg.To, msg.Symbol, msg.TokenIDs...)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgTransferFTFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTransferFTFrom) sdk.Result {
	err := keeper.TransferFTFrom(ctx, msg.Proxy, msg.From, msg.To, msg.Symbol, msg.Amount...)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Proxy.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgTransferNFTFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTransferNFTFrom) sdk.Result {
	err := keeper.TransferNFTFrom(ctx, msg.Proxy, msg.From, msg.To, msg.Symbol, msg.TokenIDs...)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Proxy.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}
