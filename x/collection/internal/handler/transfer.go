package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/keeper"
	"github.com/line/link/x/collection/internal/types"
)

func handleMsgTransferCFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTransferCFT) sdk.Result {
	err := keeper.TransferCFT(ctx, msg.From, msg.To, msg.Symbol, msg.TokenID, msg.Amount)
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

func handleMsgTransferCNFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTransferCNFT) sdk.Result {
	err := keeper.TransferCNFT(ctx, msg.From, msg.To, msg.Symbol, msg.TokenID)
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

func handleMsgTransferCFTFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTransferCFTFrom) sdk.Result {
	err := keeper.TransferCFTFrom(ctx, msg.Proxy, msg.From, msg.To, msg.Symbol, msg.TokenID, msg.Amount)
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

func handleMsgTransferCNFTFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTransferCNFTFrom) sdk.Result {
	err := keeper.TransferCNFTFrom(ctx, msg.Proxy, msg.From, msg.To, msg.Symbol, msg.TokenID)
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
