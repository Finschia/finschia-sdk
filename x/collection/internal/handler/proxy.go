package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/keeper"
	"github.com/line/link/x/collection/internal/types"
)

func handleMsgApprove(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgApprove) sdk.Result {
	err := keeper.SetApproved(ctx, msg.ContractID, msg.Proxy, msg.Approver)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Approver.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgDisapprove(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgDisapprove) sdk.Result {
	err := keeper.DeleteApproved(ctx, msg.ContractID, msg.Proxy, msg.Approver)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Approver.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}
