package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/keeper"
	"github.com/line/link/x/collection/internal/types"
)

func handleMsgIssueFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgIssueFT) sdk.Result {
	collection, err := keeper.GetCollection(ctx, msg.Symbol)
	if err != nil {
		return err.Result()
	}
	perm := types.NewIssuePermission(collection.GetSymbol())
	if !keeper.HasPermission(ctx, msg.Owner, perm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, msg.Owner, perm).Result()
	}

	token := types.NewFT(collection, msg.Name, msg.TokenURI, msg.Decimals, msg.Mintable)
	err = keeper.IssueFT(ctx, msg.Symbol, msg.Owner, token, msg.Amount)
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

func handleMsgIssueNFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgIssueNFT) sdk.Result {
	_, err := keeper.GetCollection(ctx, msg.Symbol)
	if err != nil {
		return err.Result()
	}

	perm := types.NewIssuePermission(msg.Symbol)
	if !keeper.HasPermission(ctx, msg.Owner, perm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, msg.Owner, perm).Result()
	}

	err = keeper.IssueNFT(ctx, msg.Symbol, msg.Owner)
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
