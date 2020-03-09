package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/keeper"
	"github.com/line/link/x/collection/internal/types"
)

func handleMsgCreateCollection(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgCreateCollection) sdk.Result {
	collection := types.NewCollection(keeper.NewContractID(ctx), msg.Name, msg.Meta, msg.BaseImgURI)
	err := keeper.CreateCollection(ctx, collection, msg.Owner)
	if err != nil {
		return err.Result()
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
