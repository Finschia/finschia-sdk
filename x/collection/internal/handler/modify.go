package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link-modules/x/collection/internal/keeper"
	"github.com/line/link-modules/x/collection/internal/types"
)

func handleMsgModify(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgModify) (*sdk.Result, error) {
	if err := keeper.Modify(ctx, msg.Owner, msg.TokenType, msg.TokenIndex, msg.Changes); err != nil {
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
