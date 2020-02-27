package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/keeper"
	"github.com/line/link/x/collection/internal/types"
)

func handleMsgIssueFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgIssueFT) sdk.Result {
	_, err := keeper.GetCollection(ctx, msg.Symbol)
	if err != nil {
		return err.Result()
	}
	perm := types.NewIssuePermission(msg.Symbol)
	if !keeper.HasPermission(ctx, msg.Owner, perm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, msg.Owner, perm).Result()
	}

	tokenID, err := keeper.GetNextTokenIDFT(ctx, msg.Symbol)
	if err != nil {
		return err.Result()
	}

	token := types.NewFT(msg.Symbol, tokenID, msg.Name, msg.Decimals, msg.Mintable)
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

	tokenTypeID, err := keeper.GetNextTokenType(ctx, msg.Symbol)
	if err != nil {
		return err.Result()
	}

	tokenType := types.NewBaseTokenType(msg.Symbol, tokenTypeID, msg.Name)
	err = keeper.IssueNFT(ctx, msg.Symbol, tokenType, msg.Owner)
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
