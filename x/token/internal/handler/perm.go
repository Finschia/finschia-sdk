// nolint:dupl
package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/lbm-sdk/v2/x/token/internal/keeper"
	"github.com/line/lbm-sdk/v2/x/token/internal/types"
)

func handleMsgGrant(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgGrantPermission) (*sdk.Result, error) {
	err := keeper.GrantPermission(ctx, msg.From, msg.To, msg.Permission)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRevoke(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgRevokePermission) (*sdk.Result, error) {
	err := keeper.RevokePermission(ctx, msg.From, msg.Permission)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
