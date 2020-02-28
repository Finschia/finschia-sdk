package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/keeper"
	"github.com/line/link/x/token/internal/types"
)

func handleMsgModify(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgModify) sdk.Result {
	for _, change := range msg.Changes {
		err := keeper.ModifyToken(ctx, msg.Owner, msg.ContractID, change)
		if err != nil {
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
