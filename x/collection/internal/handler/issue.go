package handler

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/collection/internal/keeper"
	"github.com/line/lbm-sdk/x/collection/internal/types"
)

func handleMsgIssueFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgIssueFT) (*sdk.Result, error) {
	_, err := keeper.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	perm := types.NewIssuePermission()
	if !keeper.HasPermission(ctx, msg.Owner, perm) {
		return nil, sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", msg.Owner.String(), perm.String())
	}

	tokenID, err := keeper.GetNextTokenIDFT(ctx)
	if err != nil {
		return nil, err
	}

	token := types.NewFT(msg.ContractID, tokenID, msg.Name, msg.Meta, msg.Decimals, msg.Mintable)
	err = keeper.IssueFT(ctx, msg.Owner, msg.To, token, msg.Amount)
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

func handleMsgIssueNFT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgIssueNFT) (*sdk.Result, error) {
	_, err := keeper.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	perm := types.NewIssuePermission()
	if !keeper.HasPermission(ctx, msg.Owner, perm) {
		return nil, sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", msg.Owner.String(), perm.String())
	}

	tokenTypeID, err := keeper.GetNextTokenType(ctx)
	if err != nil {
		return nil, err
	}

	tokenType := types.NewBaseTokenType(msg.ContractID, tokenTypeID, msg.Name, msg.Meta)
	err = keeper.IssueNFT(ctx, tokenType, msg.Owner)
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
