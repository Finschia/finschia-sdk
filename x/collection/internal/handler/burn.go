package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/lbm-sdk/v2/x/collection/internal/keeper"
	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
)

func handleMsgBurnNFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBurnNFT) (*sdk.Result, error) {
	err := keeper.BurnNFT(ctx, msg.From, msg.TokenIDs...)
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

func handleMsgBurnNFTFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBurnNFTFrom) (*sdk.Result, error) {
	err := keeper.BurnNFTFrom(ctx, msg.Proxy, msg.From, msg.TokenIDs...)
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

func handleMsgBurnFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBurnFT) (*sdk.Result, error) {
	err := keeper.BurnFT(ctx, msg.From, msg.Amount)
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

func handleMsgBurnFTFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBurnFTFrom) (*sdk.Result, error) {
	err := keeper.BurnFTFrom(ctx, msg.Proxy, msg.From, msg.Amount)
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
