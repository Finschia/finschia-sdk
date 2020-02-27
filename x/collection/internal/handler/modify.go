package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/keeper"
	"github.com/line/link/x/collection/internal/types"
)

func handleMsgModify(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgModify) sdk.Result {
	for _, change := range msg.Changes {
		if err := keeper.Modify(ctx, msg.Owner, msg.Symbol, msg.TokenType, msg.TokenIndex, change); err != nil {
			return err.Result()
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
