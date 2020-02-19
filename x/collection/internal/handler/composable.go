package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/keeper"
	"github.com/line/link/x/collection/internal/types"
)

func handleMsgAttach(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgAttach) sdk.Result {
	err := keeper.Attach(ctx, msg.From, msg.Symbol, msg.ToTokenID, msg.TokenID)
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

func handleMsgDetach(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgDetach) sdk.Result {
	err := keeper.Detach(ctx, msg.From, msg.To, msg.Symbol, msg.TokenID)
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

func handleMsgAttachFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgAttachFrom) sdk.Result {
	err := keeper.AttachFrom(ctx, msg.Proxy, msg.From, msg.Symbol, msg.ToTokenID, msg.TokenID)
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

func handleMsgDetachFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgDetachFrom) sdk.Result {
	err := keeper.DetachFrom(ctx, msg.Proxy, msg.From, msg.To, msg.Symbol, msg.TokenID)
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
