// nolint:dupl
package handler

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token/internal/keeper"
	"github.com/line/lbm-sdk/x/token/internal/types"
)

func handleMsgMint(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgMint) (*sdk.Result, error) {
	err := keeper.MintToken(ctx, msg.Amount, msg.From, msg.To)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgBurn(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBurn) (*sdk.Result, error) {
	err := keeper.BurnToken(ctx, msg.Amount, msg.From)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgBurnFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBurnFrom) (*sdk.Result, error) {
	err := keeper.BurnTokenFrom(ctx, msg.Proxy, msg.From, msg.Amount)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Proxy.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
