package proxy

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/proxy/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgProxyApproveCoins:
			return handleMsgProxyApproveCoins(ctx, keeper, msg)
		case types.MsgProxyDisapproveCoins:
			return handleMsgProxyDisapproveCoins(ctx, keeper, msg)
		case types.MsgProxySendCoinsFrom:
			return handleMsgProxySendCoinsFrom(ctx, keeper, msg)
		default:
			return types.ErrProxyInvalidMsgType(types.DefaultCodespace, msg.Type()).Result()
		}
	}
}

func handleMsgProxyApproveCoins(ctx sdk.Context, keeper Keeper, msg types.MsgProxyApproveCoins) sdk.Result {
	err := keeper.ApproveCoins(ctx, msg)
	if err != nil {
		return err.Result()
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
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgProxyDisapproveCoins(ctx sdk.Context, keeper Keeper, msg types.MsgProxyDisapproveCoins) sdk.Result {
	err := keeper.DisapproveCoins(ctx, msg)
	if err != nil {
		return err.Result()
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
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgProxySendCoinsFrom(ctx sdk.Context, keeper Keeper, msg types.MsgProxySendCoinsFrom) sdk.Result {
	err := keeper.SendCoinsFrom(ctx, msg)
	if err != nil {
		return err.Result()
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
	return sdk.Result{Events: ctx.EventManager().Events()}
}
