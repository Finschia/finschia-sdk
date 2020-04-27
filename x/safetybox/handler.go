package safetybox

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgSafetyBoxCreate:
			return handleMsgSafetyBoxCreate(ctx, keeper, msg)
		case MsgSafetyBoxAllocateToken:
			return handleMsgSafetyBoxAllocateToken(ctx, keeper, msg)
		case MsgSafetyBoxRecallToken:
			return handleMsgSafetyBoxRecallToken(ctx, keeper, msg)
		case MsgSafetyBoxIssueToken:
			return handleMsgSafetyBoxIssueToken(ctx, keeper, msg)
		case MsgSafetyBoxReturnToken:
			return handleMsgSafetyBoxReturnToken(ctx, keeper, msg)
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
			return nil, sdkerrors.Wrapf(ErrSafetyBoxInvalidMsgType, "Type: %s", msg.Type())
		}
	}
}

func handleMsgSafetyBoxCreate(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxCreate) (*sdk.Result, error) {
	sb, err := keeper.NewSafetyBox(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxCreate,
			sdk.NewAttribute(AttributeKeySafetyBoxID, sb.ID),
			sdk.NewAttribute(AttributeKeySafetyBoxOwner, sb.Owner.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAddress, sb.Address.String()),
			sdk.NewAttribute(AttributeKeyContractID, sb.ContractID),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, sb.Owner.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSafetyBoxAllocateToken(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxAllocateToken) (*sdk.Result, error) {
	err := keeper.Allocate(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxSendToken,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msg.SafetyBoxID),
			sdk.NewAttribute(AttributeKeySafetyBoxAllocatorAddress, msg.AllocatorAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAction, ActionAllocate),
			sdk.NewAttribute(AttributeKeyContractID, msg.ContractID),
			sdk.NewAttribute(AttributeKeyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.AllocatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSafetyBoxRecallToken(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxRecallToken) (*sdk.Result, error) {
	err := keeper.Recall(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxSendToken,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msg.SafetyBoxID),
			sdk.NewAttribute(AttributeKeySafetyBoxAllocatorAddress, msg.AllocatorAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAction, ActionRecall),
			sdk.NewAttribute(AttributeKeyContractID, msg.ContractID),
			sdk.NewAttribute(AttributeKeyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.AllocatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSafetyBoxIssueToken(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxIssueToken) (*sdk.Result, error) {
	err := keeper.Issue(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxSendToken,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msg.SafetyBoxID),
			sdk.NewAttribute(AttributeKeySafetyBoxIssueFromAddress, msg.FromAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxIssueToAddress, msg.ToAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAction, ActionIssue),
			sdk.NewAttribute(AttributeKeyContractID, msg.ContractID),
			sdk.NewAttribute(AttributeKeyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSafetyBoxReturnToken(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxReturnToken) (*sdk.Result, error) {
	err := keeper.Return(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxSendToken,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msg.SafetyBoxID),
			sdk.NewAttribute(AttributeKeySafetyBoxReturnerAddress, msg.ReturnerAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAction, ActionReturn),
			sdk.NewAttribute(AttributeKeyContractID, msg.ContractID),
			sdk.NewAttribute(AttributeKeyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ReturnerAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSafetyBoxRegisterAllocator(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxRegisterAllocator) (*sdk.Result, error) {
	err := keeper.GrantPermission(ctx, msg.SafetyBoxID, msg.Operator, msg.Address, RoleAllocator)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msg.SafetyBoxID),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msg.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msg.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxGrantAllocatorPermission, RoleAllocator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSafetyBoxDeregisterAllocator(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxDeregisterAllocator) (*sdk.Result, error) {
	err := keeper.RevokePermission(ctx, msg.SafetyBoxID, msg.Operator, msg.Address, RoleAllocator)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msg.SafetyBoxID),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msg.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msg.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxRevokeAllocatorPermission, RoleAllocator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSafetyBoxRegisterOperator(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxRegisterOperator) (*sdk.Result, error) {
	err := keeper.GrantPermission(ctx, msg.SafetyBoxID, msg.SafetyBoxOwner, msg.Address, RoleOperator)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msg.SafetyBoxID),
			sdk.NewAttribute(AttributeKeySafetyBoxOwner, msg.SafetyBoxOwner.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msg.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxGrantOperatorPermission, RoleOperator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SafetyBoxOwner.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSafetyBoxDeregisterOperator(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxDeregisterOperator) (*sdk.Result, error) {
	err := keeper.RevokePermission(ctx, msg.SafetyBoxID, msg.SafetyBoxOwner, msg.Address, RoleOperator)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msg.SafetyBoxID),
			sdk.NewAttribute(AttributeKeySafetyBoxOwner, msg.SafetyBoxOwner.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msg.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxRevokeOperatorPermission, RoleOperator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SafetyBoxOwner.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSafetyBoxRegisterIssuer(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxRegisterIssuer) (*sdk.Result, error) {
	err := keeper.GrantPermission(ctx, msg.SafetyBoxID, msg.Operator, msg.Address, RoleIssuer)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msg.SafetyBoxID),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msg.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msg.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxGrantIssuerPermission, RoleIssuer),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSafetyBoxDeregisterIssuer(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxDeregisterIssuer) (*sdk.Result, error) {
	err := keeper.RevokePermission(ctx, msg.SafetyBoxID, msg.Operator, msg.Address, RoleIssuer)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msg.SafetyBoxID),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msg.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msg.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxRevokeIssuerPermission, RoleIssuer),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSafetyBoxRegisterReturner(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxRegisterReturner) (*sdk.Result, error) {
	err := keeper.GrantPermission(ctx, msg.SafetyBoxID, msg.Operator, msg.Address, RoleReturner)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msg.SafetyBoxID),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msg.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msg.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxGrantReturnerPermission, RoleReturner),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSafetyBoxDeregisterReturner(ctx sdk.Context, keeper Keeper, msg MsgSafetyBoxDeregisterReturner) (*sdk.Result, error) {
	err := keeper.RevokePermission(ctx, msg.SafetyBoxID, msg.Operator, msg.Address, RoleReturner)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msg.SafetyBoxID),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msg.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msg.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxRevokeReturnerPermission, RoleReturner),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
