package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/lbm-sdk/v2/x/collection/internal/keeper"
	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
)

func handleMsgApprove(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgApprove) (*sdk.Result, error) {
	err := keeper.SetApproved(ctx, msg.Proxy, msg.Approver)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Approver.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDisapprove(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgDisapprove) (*sdk.Result, error) {
	err := keeper.DeleteApproved(ctx, msg.Proxy, msg.Approver)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Approver.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
