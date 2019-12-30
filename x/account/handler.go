package account

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/line/link/x/account/internal/types"
)

// NewHandler returns a handler for "account" type messages.
func NewHandler(k auth.AccountKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCreateAccount:
			return handleMsgCreateAccount(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized account message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateAccount(ctx sdk.Context, keeper auth.AccountKeeper, msg types.MsgCreateAccount) sdk.Result {
	if keeper.GetAccount(ctx, msg.TargetAddress) != nil {
		return types.ErrAccountAlreadyExist(types.DefaultCodespace).Result()
	}

	acc := keeper.NewAccountWithAddress(ctx, msg.TargetAddress)

	keeper.SetAccount(ctx, acc)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventCreateAccount,
			sdk.NewAttribute(types.AttributeKeyCreateAccountFrom, msg.FromAddress.String()),
			sdk.NewAttribute(types.AttributeKeyCreateAccountTarget, msg.TargetAddress.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
