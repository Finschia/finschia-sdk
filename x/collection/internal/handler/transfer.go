package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link-modules/x/collection/internal/keeper"
	"github.com/line/link-modules/x/collection/internal/types"
)

func handleMsgTransferFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTransferFT) (*sdk.Result, error) {
	err := keeper.TransferFT(ctx, msg.From, msg.To, msg.Amount...)
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

func handleMsgTransferNFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTransferNFT) (*sdk.Result, error) {
	err := keeper.TransferNFT(ctx, msg.From, msg.To, msg.TokenIDs...)
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

func handleMsgTransferFTFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTransferFTFrom) (*sdk.Result, error) {
	err := keeper.TransferFTFrom(ctx, msg.Proxy, msg.From, msg.To, msg.Amount...)
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

func handleMsgTransferNFTFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTransferNFTFrom) (*sdk.Result, error) {
	err := keeper.TransferNFTFrom(ctx, msg.Proxy, msg.From, msg.To, msg.TokenIDs...)
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
