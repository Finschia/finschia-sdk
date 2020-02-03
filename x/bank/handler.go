package bank

import (
	"fmt"

	"github.com/line/link/x/bank/internal/keeper"
	"github.com/line/link/x/bank/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgSend:
			return handleMsgSend(ctx, k, msg)

		case types.MsgMultiSend:
			return handleMsgMultiSend(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized bank message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
func handleMsgSend(ctx sdk.Context, k keeper.Keeper, msg types.MsgSend) sdk.Result {
	err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySender, msg.FromAddress.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// Handle MsgMultiSend.
func handleMsgMultiSend(ctx sdk.Context, k keeper.Keeper, msg types.MsgMultiSend) sdk.Result {
	err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	for _, in := range msg.Inputs {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(types.AttributeKeySender, in.Address.String()),
			),
		)
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}
