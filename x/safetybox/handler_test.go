package safetybox

import (
	"testing"

	"github.com/line/link/x/safetybox/internal/types"
	"github.com/line/link/x/token"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	testCommon "github.com/line/link/x/safetybox/internal/keeper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	SafetyBoxTestID   = "test_safety_box_id"
	TestContractID    = "9be17165"
	TestTokenName     = "name"
	TestTokenSymbol   = "BTC"
	TestTokenMeta     = "{}"
	TestTokenImageURI = "image-uri"
	TestTokenDecimals = 6
	TestTokenAmount   = 1000
)

//nolint:dupl
func TestHandler(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper := input.Ctx, input.Keeper
	ctx = ctx.WithEventManager(sdk.NewEventManager()) // to track emitted events

	h := NewHandler(keeper)

	_, err := h(ctx, sdk.NewTestMsg())
	require.Error(t, err)
	require.Error(t, sdkerrors.Wrapf(ErrSafetyBoxInvalidMsgType, "Type: %s", sdk.NewTestMsg().Type()))

	// issue token
	testToken := token.NewToken(TestContractID, TestTokenName, TestTokenSymbol, TestTokenMeta, TestTokenImageURI, sdk.NewInt(TestTokenDecimals), true)

	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	err = input.Tk.IssueToken(ctx, testToken, sdk.NewInt(TestTokenAmount), addr, addr)
	require.NoError(t, err)

	tok, err := input.Tk.GetToken(ctx, TestContractID)
	require.NoError(t, err)
	require.Equal(t, testToken.GetContractID(), tok.GetContractID())

	// create a box
	owner := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	msgSbCreate := MsgSafetyBoxCreate{
		SafetyBoxID:    SafetyBoxTestID,
		SafetyBoxOwner: owner,
		ContractID:     TestContractID,
	}
	res, err := h(ctx, msgSbCreate)
	require.NoError(t, err)

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
			sdk.NewAttribute(AttributeKeyContractID, msgSbCreate.ContractID),
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
	res, err = h(ctx, msgSbRegisterOperator)
	require.NoError(t, err)

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
	res, err = h(ctx, msgSbRegisterAllocator)
	require.NoError(t, err)

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
	res, err = h(ctx, msgSbRegisterIssuer)
	require.NoError(t, err)

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
	res, err = h(ctx, msgSbRegisterIssuer)
	require.NoError(t, err)

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
	res, err = h(ctx, msgSbRegisterReturner)
	require.NoError(t, err)

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

	// put amounts of token to all
	_, err = input.Tk.AddBalance(ctx, TestContractID, allocator, sdk.NewInt(10))
	require.NoError(t, err)
	_, err = input.Tk.AddBalance(ctx, TestContractID, issuer1, sdk.NewInt(10))
	require.NoError(t, err)
	_, err = input.Tk.AddBalance(ctx, TestContractID, issuer2, sdk.NewInt(10))
	require.NoError(t, err)
	_, err = input.Tk.AddBalance(ctx, TestContractID, returner, sdk.NewInt(10))
	require.NoError(t, err)

	// allocate, issue, return and recall
	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	msgSbAllocate := MsgSafetyBoxAllocateToken{
		SafetyBoxID:      SafetyBoxTestID,
		AllocatorAddress: allocator,
		ContractID:       TestContractID,
		Amount:           sdk.NewInt(1),
	}
	res, err = h(ctx, msgSbAllocate)
	require.NoError(t, err)

	// check emitted events including SendCoins of Bank module
	e = sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeySender, msgSbAllocate.AllocatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, safetyBoxAddress.String()),
			sdk.NewAttribute(types.AttributeKeyContractID, msgSbAllocate.ContractID),
			sdk.NewAttribute(types.AttributeKeyAmount, msgSbAllocate.Amount.String()),
		),
		sdk.NewEvent(
			EventSafetyBoxSendToken,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbAllocate.SafetyBoxID),
			sdk.NewAttribute(AttributeKeySafetyBoxAllocatorAddress, msgSbAllocate.AllocatorAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAction, ActionAllocate),
			sdk.NewAttribute(AttributeKeyContractID, msgSbAllocate.ContractID),
			sdk.NewAttribute(AttributeKeyAmount, msgSbAllocate.Amount.String()),
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
	msgSbIssue := MsgSafetyBoxIssueToken{
		SafetyBoxID: SafetyBoxTestID,
		FromAddress: issuer1,
		ToAddress:   issuer2,
		ContractID:  TestContractID,
		Amount:      sdk.NewInt(1),
	}
	res, err = h(ctx, msgSbIssue)
	require.NoError(t, err)

	// check emitted events including SendCoins of Bank module
	e = sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeySender, safetyBoxAddress.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, msgSbIssue.ToAddress.String()),
			sdk.NewAttribute(types.AttributeKeyContractID, msgSbIssue.ContractID),
			sdk.NewAttribute(types.AttributeKeyAmount, msgSbIssue.Amount.String()),
		),
		sdk.NewEvent(
			EventSafetyBoxSendToken,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbIssue.SafetyBoxID),
			sdk.NewAttribute(AttributeKeySafetyBoxIssueFromAddress, msgSbIssue.FromAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxIssueToAddress, msgSbIssue.ToAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAction, ActionIssue),
			sdk.NewAttribute(AttributeKeyContractID, msgSbIssue.ContractID),
			sdk.NewAttribute(AttributeKeyAmount, msgSbIssue.Amount.String()),
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
	msgSbReturn := MsgSafetyBoxReturnToken{
		SafetyBoxID:     SafetyBoxTestID,
		ReturnerAddress: returner,
		ContractID:      TestContractID,
		Amount:          sdk.NewInt(1),
	}
	res, err = h(ctx, msgSbReturn)
	require.NoError(t, err)

	// check emitted events including SendCoins of Bank module
	e = sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeySender, msgSbReturn.ReturnerAddress.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, safetyBoxAddress.String()),
			sdk.NewAttribute(types.AttributeKeyContractID, msgSbReturn.ContractID),
			sdk.NewAttribute(types.AttributeKeyAmount, msgSbReturn.Amount.String()),
		),
		sdk.NewEvent(
			EventSafetyBoxSendToken,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbReturn.SafetyBoxID),
			sdk.NewAttribute(AttributeKeySafetyBoxReturnerAddress, msgSbReturn.ReturnerAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAction, ActionReturn),
			sdk.NewAttribute(AttributeKeyContractID, msgSbReturn.ContractID),
			sdk.NewAttribute(AttributeKeyAmount, msgSbReturn.Amount.String()),
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
	msgSbRecall := MsgSafetyBoxRecallToken{
		SafetyBoxID:      SafetyBoxTestID,
		AllocatorAddress: allocator,
		ContractID:       TestContractID,
		Amount:           sdk.NewInt(1),
	}
	res, err = h(ctx, msgSbRecall)
	require.NoError(t, err)

	// check emitted events including SendCoins of Bank module
	e = sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeySender, safetyBoxAddress.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, msgSbRecall.AllocatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyContractID, msgSbRecall.ContractID),
			sdk.NewAttribute(types.AttributeKeyAmount, msgSbRecall.Amount.String()),
		),
		sdk.NewEvent(
			EventSafetyBoxSendToken,
			sdk.NewAttribute(AttributeKeySafetyBoxID, msgSbRecall.SafetyBoxID),
			sdk.NewAttribute(AttributeKeySafetyBoxAllocatorAddress, msgSbRecall.AllocatorAddress.String()),
			sdk.NewAttribute(AttributeKeySafetyBoxAction, ActionRecall),
			sdk.NewAttribute(AttributeKeyContractID, msgSbRecall.ContractID),
			sdk.NewAttribute(AttributeKeyAmount, msgSbRecall.Amount.String()),
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
	res, err = h(ctx, msgSbDeregisterAllocator)
	require.NoError(t, err)

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
	res, err = h(ctx, msgSbDeregisterIssuer)
	require.NoError(t, err)

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
	res, err = h(ctx, msgSbDeregisterIssuer)
	require.NoError(t, err)

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
	res, err = h(ctx, msgSbDeregisterReturner)
	require.NoError(t, err)

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
	res, err = h(ctx, msgSbDeregisterOperator)
	require.NoError(t, err)

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
