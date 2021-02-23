package handler

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection/internal/keeper"
	"github.com/line/lbm-sdk/x/collection/internal/types"
)

func handleMsgAttach(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgAttach) (*sdk.Result, error) {
	err := keeper.Attach(ctx, msg.From, msg.ToTokenID, msg.TokenID)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDetach(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgDetach) (*sdk.Result, error) {
	err := keeper.Detach(ctx, msg.From, msg.TokenID)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgAttachFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgAttachFrom) (*sdk.Result, error) {
	err := keeper.AttachFrom(ctx, msg.Proxy, msg.From, msg.ToTokenID, msg.TokenID)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Proxy.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDetachFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgDetachFrom) (*sdk.Result, error) {
	err := keeper.DetachFrom(ctx, msg.Proxy, msg.From, msg.TokenID)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Proxy.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
