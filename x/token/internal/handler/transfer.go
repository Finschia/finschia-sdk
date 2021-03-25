package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link-modules/x/token/internal/keeper"
	"github.com/line/link-modules/x/token/internal/types"
)

func handleMsgTransfer(ctx sdk.Context, k keeper.Keeper, msg types.MsgTransfer) (*sdk.Result, error) {
	err := k.Transfer(ctx, msg.From, msg.To, msg.Amount)
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

func handleMsgTransferFrom(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTransferFrom) (*sdk.Result, error) {
	err := keeper.TransferFrom(ctx, msg.Proxy, msg.From, msg.To, msg.Amount)
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
