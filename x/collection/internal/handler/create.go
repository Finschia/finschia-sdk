package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link-modules/x/collection/internal/keeper"
	"github.com/line/link-modules/x/collection/internal/types"
	"github.com/line/link-modules/x/contract"
)

func handleMsgCreateCollection(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgCreateCollection) (*sdk.Result, error) {
	contractI := ctx.Context().Value(contract.CtxKey{})
	if contractI == nil {
		panic("contract id does not set")
	}
	collection := types.NewCollection(contractI.(string), msg.Name, msg.Meta, msg.BaseImgURI)
	err := keeper.CreateCollection(ctx, collection, msg.Owner)
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
