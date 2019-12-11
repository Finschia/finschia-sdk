package safetybox

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/link-chain/link/x/safetybox/internal/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSafetyBoxCreate:
			return handleMsgSafetyBoxCreate(ctx, keeper, msg)
		case MsgSafetyBoxAllocateCoins:
			return handleMsgSafetyBoxAllocateCoins(ctx, keeper, msg)
		case MsgSafetyBoxRecallCoins:
			return handleMsgSafetyBoxRecallCoins(ctx, keeper, msg)
		case MsgSafetyBoxIssueCoins:
			return handleMsgSafetyBoxIssueCoins(ctx, keeper, msg)
		case MsgSafetyBoxReturnCoins:
			return handleMsgSafetyBoxReturnCoins(ctx, keeper, msg)
		case MsgSafetyBoxRegisterAllocator:
			return handleMsgSafetyBoxRegisterAllocator(ctx, keeper, msg)
		case MsgSafetyBoxDeregisterAllocator:
			return handleMsgSafetyBoxDeregisterAllocator(ctx, keeper, msg)
		case MsgSafetyBoxRegisterOperator:
			return handleMsgSafetyBoxRegisterOperator(ctx, keeper, msg)
		case MsgSafetyBoxDeregisterOperator:
			return handleMsgSafetyBoxDeregisterOperator(ctx, keeper, msg)
		case MsgSafetyBoxRegisterIssuer:
			return handleMsgSafetyBoxRegisterIssuer(ctx, keeper, msg)
		case MsgSafetyBoxDeregisterIssuer:
			return handleMsgSafetyBoxDeregisterIssuer(ctx, keeper, msg)
		case MsgSafetyBoxRegisterReturner:
			return handleMsgSafetyBoxRegisterReturner(ctx, keeper, msg)
		case MsgSafetyBoxDeregisterReturner:
			return handleMsgSafetyBoxDeregisterReturner(ctx, keeper, msg)
		default:
			return ErrSafetyBoxInvalidMsgType(types.DefaultCodespace).Result()
		}
	}
}

func handleMsgSafetyBoxCreate(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxCreate) sdk.Result {
	sb, err := keeper.NewSafetyBox(ctx, msg)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxCreate,
			sdk.NewAttribute(AttributeKeySafetyBoxId, sb.ID),
			sdk.NewAttribute(AttributeKeySafetyBoxOwner, sb.Owner.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAddress, sb.Address.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, sb.Owner.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSafetyBoxAllocateCoins(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxAllocateCoins) sdk.Result {
	err := keeper.Allocate(ctx, msg)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxSendCoin,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msg.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxAllocatorAddress, msg.AllocatorAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAction, ActionAllocate),
			sdk.NewAttribute(AttributeKeySafetyBoxCoins, msg.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.AllocatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSafetyBoxRecallCoins(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxRecallCoins) sdk.Result {
	err := keeper.Recall(ctx, msg)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxSendCoin,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msg.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxAllocatorAddress, msg.AllocatorAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAction, ActionRecall),
			sdk.NewAttribute(AttributeKeySafetyBoxCoins, msg.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.AllocatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSafetyBoxIssueCoins(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxIssueCoins) sdk.Result {
	err := keeper.Issue(ctx, msg)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxSendCoin,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msg.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxIssueFromAddress, msg.FromAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxIssueToAddress, msg.ToAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAction, ActionIssue),
			sdk.NewAttribute(AttributeKeySafetyBoxCoins, msg.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSafetyBoxReturnCoins(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxReturnCoins) sdk.Result {
	err := keeper.Return(ctx, msg)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxSendCoin,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msg.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxReturnerAddress, msg.ReturnerAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAction, ActionReturn),
			sdk.NewAttribute(AttributeKeySafetyBoxCoins, msg.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ReturnerAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSafetyBoxRegisterAllocator(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxRegisterAllocator) sdk.Result {
	err := keeper.GrantPermission(ctx, msg.SafetyBoxId, msg.Operator, msg.Address, RoleAllocator)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msg.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msg.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msg.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxGrantAllocatePermission, RoleAllocator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSafetyBoxDeregisterAllocator(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxDeregisterAllocator) sdk.Result {
	err := keeper.RevokePermission(ctx, msg.SafetyBoxId, msg.Operator, msg.Address, RoleAllocator)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msg.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msg.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msg.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxRevokeAllocatePermission, RoleAllocator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSafetyBoxRegisterOperator(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxRegisterOperator) sdk.Result {
	err := keeper.GrantPermission(ctx, msg.SafetyBoxId, msg.SafetyBoxOwner, msg.Address, RoleOperator)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msg.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOwner, msg.SafetyBoxOwner.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msg.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxGrantRecallPermission, RoleOperator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SafetyBoxOwner.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSafetyBoxDeregisterOperator(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxDeregisterOperator) sdk.Result {
	err := keeper.RevokePermission(ctx, msg.SafetyBoxId, msg.SafetyBoxOwner, msg.Address, RoleOperator)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msg.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOwner, msg.SafetyBoxOwner.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msg.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxRevokeRecallPermission, RoleOperator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SafetyBoxOwner.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSafetyBoxRegisterIssuer(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxRegisterIssuer) sdk.Result {
	err := keeper.GrantPermission(ctx, msg.SafetyBoxId, msg.Operator, msg.Address, RoleIssuer)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msg.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msg.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msg.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxGrantIssuePermission, RoleIssuer),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSafetyBoxDeregisterIssuer(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxDeregisterIssuer) sdk.Result {
	err := keeper.RevokePermission(ctx, msg.SafetyBoxId, msg.Operator, msg.Address, RoleIssuer)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msg.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msg.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msg.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxRevokeIssuePermission, RoleIssuer),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSafetyBoxRegisterReturner(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxRegisterReturner) sdk.Result {
	err := keeper.GrantPermission(ctx, msg.SafetyBoxId, msg.Operator, msg.Address, RoleReturner)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msg.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msg.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msg.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxGrantReturnPermission, RoleReturner),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSafetyBoxDeregisterReturner(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxDeregisterReturner) sdk.Result {
	err := keeper.RevokePermission(ctx, msg.SafetyBoxId, msg.Operator, msg.Address, RoleReturner)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msg.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msg.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msg.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxRevokeReturnPermission, RoleReturner),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
