package proxy

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link/x/proxy/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgProxyApproveCoins:
			return handleMsgProxyApproveCoins(ctx, keeper, msg)
		case types.MsgProxyDisapproveCoins:
			return handleMsgProxyDisapproveCoins(ctx, keeper, msg)
		case types.MsgProxySendCoinsFrom:
			return handleMsgProxySendCoinsFrom(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrapf(types.ErrProxyInvalidMsgType, "Type: %s", msg.Type())
		}
	}
}

func handleMsgProxyApproveCoins(ctx sdk.Context, keeper Keeper, msg types.MsgProxyApproveCoins) (*sdk.Result, error) {
	err := keeper.ApproveCoins(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventProxyApproveCoins,
			sdk.NewAttribute(AttributeKeyProxyAddress, msg.Proxy.String()),
			sdk.NewAttribute(AttributeKeyProxyOnBehalfOfAddress, msg.OnBehalfOf.String()),
			sdk.NewAttribute(AttributeKeyProxyDenom, msg.Denom),
			sdk.NewAttribute(AttributeKeyProxyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OnBehalfOf.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgProxyDisapproveCoins(ctx sdk.Context, keeper Keeper, msg types.MsgProxyDisapproveCoins) (*sdk.Result, error) {
	err := keeper.DisapproveCoins(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventProxyDisapproveCoins,
			sdk.NewAttribute(AttributeKeyProxyAddress, msg.Proxy.String()),
			sdk.NewAttribute(AttributeKeyProxyOnBehalfOfAddress, msg.OnBehalfOf.String()),
			sdk.NewAttribute(AttributeKeyProxyDenom, msg.Denom),
			sdk.NewAttribute(AttributeKeyProxyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OnBehalfOf.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgProxySendCoinsFrom(ctx sdk.Context, keeper Keeper, msg types.MsgProxySendCoinsFrom) (*sdk.Result, error) {
	err := keeper.SendCoinsFrom(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventProxySendCoinsFrom,
			sdk.NewAttribute(AttributeKeyProxyAddress, msg.Proxy.String()),
			sdk.NewAttribute(AttributeKeyProxyOnBehalfOfAddress, msg.OnBehalfOf.String()),
			sdk.NewAttribute(AttributeKeyProxyToAddress, msg.ToAddress.String()),
			sdk.NewAttribute(AttributeKeyProxyDenom, msg.Denom),
			sdk.NewAttribute(AttributeKeyProxyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OnBehalfOf.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
