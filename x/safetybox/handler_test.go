package safetybox

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	testCommon "github.com/line/link/x/safetybox/internal/keeper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

const (
	SafetyBoxTestId = "test_safety_box_id"
)

func TestHandler(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper := input.Ctx, input.Keeper
	ctx = ctx.WithEventManager(sdk.NewEventManager()) // to track emitted events

	h := NewHandler(keeper)

	res := h(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.Error(t, ErrSafetyBoxInvalidMsgType(DefaultCodespace, sdk.NewTestMsg().Type()))

	// create a box
	owner := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	msgSbCreate := MsgSafetyBoxCreate{
		SafetyBoxId:     SafetyBoxTestId,
		SafetyBoxOwner:  owner,
		SafetyBoxDenoms: []string{"link"},
	}
	res = h(ctx, msgSbCreate)
	require.True(t, res.IsOK())

	sb, _ := keeper.GetSafetyBox(ctx, SafetyBoxTestId)
	safetyBoxAddress := sb.Address

	// check emitted events
	e := sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxCreate,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msgSbCreate.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOwner, msgSbCreate.SafetyBoxOwner.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAddress, safetyBoxAddress.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbCreate.SafetyBoxOwner.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	}
	testCommon.VerifyEventFunc(t, e, res.Events)

	// the owner registers an operator
	operator := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	msgSbRegisterOperator := MsgSafetyBoxRegisterOperator{
		SafetyBoxId:    SafetyBoxTestId,
		SafetyBoxOwner: owner,
		Address:        operator,
	}
	res = h(ctx, msgSbRegisterOperator)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msgSbRegisterOperator.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOwner, msgSbRegisterOperator.SafetyBoxOwner.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msgSbRegisterOperator.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxGrantOperatorPermission, RoleOperator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbRegisterOperator.SafetyBoxOwner.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	}
	testCommon.VerifyEventFunc(t, e, res.Events)

	// the operator registers allocator, issuers, and returner
	allocator := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	msgSbRegisterAllocator := MsgSafetyBoxRegisterAllocator{
		SafetyBoxId: SafetyBoxTestId,
		Operator:    operator,
		Address:     allocator,
	}
	res = h(ctx, msgSbRegisterAllocator)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msgSbRegisterAllocator.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msgSbRegisterAllocator.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msgSbRegisterAllocator.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxGrantAllocatorPermission, RoleAllocator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbRegisterAllocator.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	}
	testCommon.VerifyEventFunc(t, e, res.Events)

	issuer1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	msgSbRegisterIssuer := MsgSafetyBoxRegisterIssuer{
		SafetyBoxId: SafetyBoxTestId,
		Operator:    operator,
		Address:     issuer1,
	}
	res = h(ctx, msgSbRegisterIssuer)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msgSbRegisterIssuer.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msgSbRegisterIssuer.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msgSbRegisterIssuer.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxGrantIssuerPermission, RoleIssuer),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbRegisterIssuer.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	}
	testCommon.VerifyEventFunc(t, e, res.Events)

	issuer2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	msgSbRegisterIssuer = MsgSafetyBoxRegisterIssuer{
		SafetyBoxId: SafetyBoxTestId,
		Operator:    operator,
		Address:     issuer2,
	}
	res = h(ctx, msgSbRegisterIssuer)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msgSbRegisterIssuer.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msgSbRegisterIssuer.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msgSbRegisterIssuer.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxGrantIssuerPermission, RoleIssuer),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbRegisterIssuer.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	}
	testCommon.VerifyEventFunc(t, e, res.Events)

	returner := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	msgSbRegisterReturner := MsgSafetyBoxRegisterReturner{
		SafetyBoxId: SafetyBoxTestId,
		Operator:    operator,
		Address:     returner,
	}
	res = h(ctx, msgSbRegisterReturner)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msgSbRegisterReturner.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msgSbRegisterReturner.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msgSbRegisterReturner.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxGrantReturnerPermission, RoleReturner),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbRegisterReturner.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	}
	testCommon.VerifyEventFunc(t, e, res.Events)

	// put some coins to all
	_, err := input.Bk.AddCoins(ctx, allocator, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(10))))
	require.NoError(t, err)
	_, err = input.Bk.AddCoins(ctx, issuer1, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(10))))
	require.NoError(t, err)
	_, err = input.Bk.AddCoins(ctx, issuer2, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(10))))
	require.NoError(t, err)
	_, err = input.Bk.AddCoins(ctx, returner, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(10))))
	require.NoError(t, err)

	// allocate, issue, return and recall
	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	msgSbAllocate := MsgSafetyBoxAllocateCoins{
		SafetyBoxId:      SafetyBoxTestId,
		AllocatorAddress: allocator,
		Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
	}
	res = h(ctx, msgSbAllocate)
	require.True(t, res.IsOK())

	// check emitted events including SendCoins of Bank module
	e = sdk.Events{
		sdk.NewEvent(
			"transfer",
			sdk.NewAttribute("recipient", safetyBoxAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msgSbAllocate.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbAllocate.AllocatorAddress.String()),
		),
		sdk.NewEvent(
			EventSafetyBoxSendCoin,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msgSbAllocate.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxAllocatorAddress, msgSbAllocate.AllocatorAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAction, ActionAllocate),
			sdk.NewAttribute(AttributeKeySafetyBoxCoins, msgSbAllocate.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbAllocate.AllocatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	}
	testCommon.VerifyEventFunc(t, e, res.Events)

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	msgSbIssue := MsgSafetyBoxIssueCoins{
		SafetyBoxId: SafetyBoxTestId,
		FromAddress: issuer1,
		ToAddress:   issuer2,
		Coins:       sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
	}
	res = h(ctx, msgSbIssue)
	require.True(t, res.IsOK())

	// check emitted events including SendCoins of Bank module
	e = sdk.Events{
		sdk.NewEvent(
			"transfer",
			sdk.NewAttribute("recipient", msgSbIssue.ToAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msgSbIssue.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, safetyBoxAddress.String()),
		),
		sdk.NewEvent(
			EventSafetyBoxSendCoin,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msgSbIssue.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxIssueFromAddress, msgSbIssue.FromAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxIssueToAddress, msgSbIssue.ToAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAction, ActionIssue),
			sdk.NewAttribute(AttributeKeySafetyBoxCoins, msgSbIssue.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbIssue.FromAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	}
	testCommon.VerifyEventFunc(t, e, res.Events)

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	msgSbReturn := MsgSafetyBoxReturnCoins{
		SafetyBoxId:     SafetyBoxTestId,
		ReturnerAddress: returner,
		Coins:           sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
	}
	res = h(ctx, msgSbReturn)
	require.True(t, res.IsOK())

	// check emitted events including SendCoins of Bank module
	e = sdk.Events{
		sdk.NewEvent(
			"transfer",
			sdk.NewAttribute("recipient", safetyBoxAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msgSbReturn.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbReturn.ReturnerAddress.String()),
		),
		sdk.NewEvent(
			EventSafetyBoxSendCoin,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msgSbReturn.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxReturnerAddress, msgSbReturn.ReturnerAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAction, ActionReturn),
			sdk.NewAttribute(AttributeKeySafetyBoxCoins, msgSbReturn.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbReturn.ReturnerAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	}
	testCommon.VerifyEventFunc(t, e, res.Events)

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	msgSbRecall := MsgSafetyBoxRecallCoins{
		SafetyBoxId:      SafetyBoxTestId,
		AllocatorAddress: allocator,
		Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
	}
	res = h(ctx, msgSbRecall)
	require.True(t, res.IsOK())

	// check emitted events including SendCoins of Bank module
	e = sdk.Events{
		sdk.NewEvent(
			"transfer",
			sdk.NewAttribute("recipient", msgSbRecall.AllocatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msgSbReturn.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, safetyBoxAddress.String()),
		),
		sdk.NewEvent(
			EventSafetyBoxSendCoin,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msgSbRecall.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxAllocatorAddress, msgSbRecall.AllocatorAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAction, ActionRecall),
			sdk.NewAttribute(AttributeKeySafetyBoxCoins, msgSbRecall.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbRecall.AllocatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	}
	testCommon.VerifyEventFunc(t, e, res.Events)

	// the operator deregisters allocator, issuers and returner
	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	msgSbDeregisterAllocator := MsgSafetyBoxDeregisterAllocator{
		SafetyBoxId: SafetyBoxTestId,
		Operator:    operator,
		Address:     allocator,
	}
	res = h(ctx, msgSbDeregisterAllocator)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msgSbDeregisterAllocator.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msgSbDeregisterAllocator.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msgSbDeregisterAllocator.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxRevokeAllocatorPermission, RoleAllocator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbDeregisterAllocator.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	}
	testCommon.VerifyEventFunc(t, e, res.Events)

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	msgSbDeregisterIssuer := MsgSafetyBoxDeregisterIssuer{
		SafetyBoxId: SafetyBoxTestId,
		Operator:    operator,
		Address:     issuer1,
	}
	res = h(ctx, msgSbDeregisterIssuer)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msgSbDeregisterIssuer.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msgSbDeregisterIssuer.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msgSbDeregisterIssuer.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxRevokeIssuerPermission, RoleIssuer),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbDeregisterIssuer.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	}
	testCommon.VerifyEventFunc(t, e, res.Events)

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	msgSbDeregisterIssuer = MsgSafetyBoxDeregisterIssuer{
		SafetyBoxId: SafetyBoxTestId,
		Operator:    operator,
		Address:     issuer2,
	}
	res = h(ctx, msgSbDeregisterIssuer)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msgSbDeregisterIssuer.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msgSbDeregisterIssuer.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msgSbDeregisterIssuer.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxRevokeIssuerPermission, RoleIssuer),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbDeregisterIssuer.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	}
	testCommon.VerifyEventFunc(t, e, res.Events)

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	msgSbDeregisterReturner := MsgSafetyBoxDeregisterReturner{
		SafetyBoxId: SafetyBoxTestId,
		Operator:    operator,
		Address:     returner,
	}
	res = h(ctx, msgSbDeregisterReturner)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msgSbDeregisterReturner.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOperator, msgSbDeregisterReturner.Operator.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msgSbDeregisterReturner.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxRevokeReturnerPermission, RoleReturner),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbDeregisterReturner.Operator.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	}
	testCommon.VerifyEventFunc(t, e, res.Events)

	// the owner deregisters an operator
	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	msgSbDeregisterOperator := MsgSafetyBoxDeregisterOperator{
		SafetyBoxId:    SafetyBoxTestId,
		SafetyBoxOwner: owner,
		Address:        operator,
	}
	res = h(ctx, msgSbDeregisterOperator)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxId, msgSbDeregisterOperator.SafetyBoxId),
			sdk.NewAttribute(AttributeKeySafetyBoxOwner, msgSbDeregisterOperator.SafetyBoxOwner.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxTarget, msgSbDeregisterOperator.Address.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxRevokeOperatorPermission, RoleOperator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msgSbDeregisterOperator.SafetyBoxOwner.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	}
	testCommon.VerifyEventFunc(t, e, res.Events)
}
