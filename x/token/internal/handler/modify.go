package handler

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token/internal/keeper"
	"github.com/line/lbm-sdk/x/token/internal/types"
)

func handleMsgModify(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgModify) (*sdk.Result, error) {
	err := keeper.ModifyToken(ctx, msg.Owner, msg.Changes)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
