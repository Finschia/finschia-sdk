package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/lbm-sdk/v2/x/token/internal/keeper"
	"github.com/line/lbm-sdk/v2/x/token/internal/types"
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
