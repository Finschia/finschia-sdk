package safetybox

import (
	"testing"

	"github.com/line/link/x/safetybox/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testCommon "github.com/line/link/x/safetybox/internal/keeper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	SafetyBoxTestID = "test_safety_box_id"
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
		SafetyBoxID:     SafetyBoxTestID,
		SafetyBoxOwner:  owner,
		SafetyBoxDenoms: []string{"link"},
	}
	res = h(ctx, msgSbCreate)
	require.True(t, res.IsOK())

	sb, err := keeper.GetSafetyBox(ctx, SafetyBoxTestID)
	require.NoError(t, err)
	safetyBoxAddress := sb.Address

	// check emitted events
	e := sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxCreate,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbCreate.SafetyBoxID),
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
		SafetyBoxID:    SafetyBoxTestID,
		SafetyBoxOwner: owner,
		Address:        operator,
	}
	res = h(ctx, msgSbRegisterOperator)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbRegisterOperator.SafetyBoxID),
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
		SafetyBoxID: SafetyBoxTestID,
		Operator:    operator,
		Address:     allocator,
	}
	res = h(ctx, msgSbRegisterAllocator)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbRegisterAllocator.SafetyBoxID),
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
		SafetyBoxID: SafetyBoxTestID,
		Operator:    operator,
		Address:     issuer1,
	}
	res = h(ctx, msgSbRegisterIssuer)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbRegisterIssuer.SafetyBoxID),
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
		SafetyBoxID: SafetyBoxTestID,
		Operator:    operator,
		Address:     issuer2,
	}
	res = h(ctx, msgSbRegisterIssuer)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbRegisterIssuer.SafetyBoxID),
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
		SafetyBoxID: SafetyBoxTestID,
		Operator:    operator,
		Address:     returner,
	}
	res = h(ctx, msgSbRegisterReturner)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbRegisterReturner.SafetyBoxID),
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
	_, err = input.Bk.AddCoins(ctx, allocator, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(10))))
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
		SafetyBoxID:      SafetyBoxTestID,
		AllocatorAddress: allocator,
		Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
	}
	res = h(ctx, msgSbAllocate)
	require.True(t, res.IsOK())

	// check emitted events including SendCoins of Bank module
	e = sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeySender, msgSbAllocate.AllocatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, safetyBoxAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msgSbAllocate.Coins.String()),
		),
		sdk.NewEvent(
			EventSafetyBoxSendCoin,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbAllocate.SafetyBoxID),
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
		SafetyBoxID: SafetyBoxTestID,
		FromAddress: issuer1,
		ToAddress:   issuer2,
		Coins:       sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
	}
	res = h(ctx, msgSbIssue)
	require.True(t, res.IsOK())

	// check emitted events including SendCoins of Bank module
	e = sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeySender, safetyBoxAddress.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, msgSbIssue.ToAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msgSbIssue.Coins.String()),
		),
		sdk.NewEvent(
			EventSafetyBoxSendCoin,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbIssue.SafetyBoxID),
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
		SafetyBoxID:     SafetyBoxTestID,
		ReturnerAddress: returner,
		Coins:           sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
	}
	res = h(ctx, msgSbReturn)
	require.True(t, res.IsOK())

	// check emitted events including SendCoins of Bank module
	e = sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeySender, msgSbReturn.ReturnerAddress.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, safetyBoxAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msgSbReturn.Coins.String()),
		),
		sdk.NewEvent(
			EventSafetyBoxSendCoin,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbReturn.SafetyBoxID),
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
		SafetyBoxID:      SafetyBoxTestID,
		AllocatorAddress: allocator,
		Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
	}
	res = h(ctx, msgSbRecall)
	require.True(t, res.IsOK())

	// check emitted events including SendCoins of Bank module
	e = sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeySender, safetyBoxAddress.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, msgSbRecall.AllocatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msgSbReturn.Coins.String()),
		),
		sdk.NewEvent(
			EventSafetyBoxSendCoin,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbRecall.SafetyBoxID),
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
		SafetyBoxID: SafetyBoxTestID,
		Operator:    operator,
		Address:     allocator,
	}
	res = h(ctx, msgSbDeregisterAllocator)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbDeregisterAllocator.SafetyBoxID),
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
		SafetyBoxID: SafetyBoxTestID,
		Operator:    operator,
		Address:     issuer1,
	}
	res = h(ctx, msgSbDeregisterIssuer)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbDeregisterIssuer.SafetyBoxID),
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
		SafetyBoxID: SafetyBoxTestID,
		Operator:    operator,
		Address:     issuer2,
	}
	res = h(ctx, msgSbDeregisterIssuer)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbDeregisterIssuer.SafetyBoxID),
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
		SafetyBoxID: SafetyBoxTestID,
		Operator:    operator,
		Address:     returner,
	}
	res = h(ctx, msgSbDeregisterReturner)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbDeregisterReturner.SafetyBoxID),
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
		SafetyBoxID:    SafetyBoxTestID,
		SafetyBoxOwner: owner,
		Address:        operator,
	}
	res = h(ctx, msgSbDeregisterOperator)
	require.True(t, res.IsOK())

	// check emitted events
	e = sdk.Events{
		sdk.NewEvent(
			EventSafetyBoxPermission,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbDeregisterOperator.SafetyBoxID),
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
