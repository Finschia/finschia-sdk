package safetybox

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	testCommon "github.com/line/link/x/safetybox/internal/keeper"
	types "github.com/line/link/x/safetybox/internal/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

const (
	SafetyBoxTestId = "test_safety_box_id"
)

// ToDo: MsgSafetyBoxAllocateCoins, MsgSafetyBoxRecallCoins, MsgSafetyBoxIssueCoins, MsgSafetyBoxReturnCoins
func TestHandler(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper := input.Ctx, input.Keeper

	h := NewHandler(keeper)

	res := h(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.Error(t, ErrSafetyBoxInvalidMsgType(types.DefaultCodespace))

	// create a box
	owner := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	msgSbCreate := MsgSafetyBoxCreate{
		SafetyBoxId:    SafetyBoxTestId,
		SafetyBoxOwner: owner,
	}
	res = h(ctx, msgSbCreate)
	require.True(t, res.IsOK())

	// the owner registers an operator
	operator := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	msgSbRegisterOperator := MsgSafetyBoxRegisterOperator{
		SafetyBoxId:    SafetyBoxTestId,
		SafetyBoxOwner: owner,
		Address:        operator,
	}
	res = h(ctx, msgSbRegisterOperator)
	require.True(t, res.IsOK())

	// the operator registers allocator, issuers, and returner
	allocator := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	msgSbRegisterAllocator := MsgSafetyBoxRegisterAllocator{
		SafetyBoxId: SafetyBoxTestId,
		Operator:    operator,
		Address:     allocator,
	}
	res = h(ctx, msgSbRegisterAllocator)
	require.True(t, res.IsOK())

	issuer1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	msgSbRegisterIssuer := MsgSafetyBoxRegisterIssuer{
		SafetyBoxId: SafetyBoxTestId,
		Operator:    operator,
		Address:     issuer1,
	}
	res = h(ctx, msgSbRegisterIssuer)
	require.True(t, res.IsOK())

	issuer2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	msgSbRegisterIssuer = MsgSafetyBoxRegisterIssuer{
		SafetyBoxId: SafetyBoxTestId,
		Operator:    operator,
		Address:     issuer2,
	}
	res = h(ctx, msgSbRegisterIssuer)
	require.True(t, res.IsOK())

	returner := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	msgSbRegisterReturner := MsgSafetyBoxRegisterReturner{
		SafetyBoxId: SafetyBoxTestId,
		Operator:    operator,
		Address:     returner,
	}
	res = h(ctx, msgSbRegisterReturner)
	require.True(t, res.IsOK())

	// the operator deregisters allocator, issuers and returner
	msgSbDeregisterAllocator := MsgSafetyBoxDeregisterAllocator{
		SafetyBoxId: SafetyBoxTestId,
		Operator:    operator,
		Address:     allocator,
	}
	res = h(ctx, msgSbDeregisterAllocator)
	require.True(t, res.IsOK())

	msgSbDeregisterIssuer := MsgSafetyBoxDeregisterIssuer{
		SafetyBoxId: SafetyBoxTestId,
		Operator:    operator,
		Address:     issuer1,
	}
	res = h(ctx, msgSbDeregisterIssuer)
	require.True(t, res.IsOK())

	msgSbDeregisterIssuer = MsgSafetyBoxDeregisterIssuer{
		SafetyBoxId: SafetyBoxTestId,
		Operator:    operator,
		Address:     issuer2,
	}
	res = h(ctx, msgSbDeregisterIssuer)
	require.True(t, res.IsOK())

	msgSbDeregisterReturner := MsgSafetyBoxDeregisterReturner{
		SafetyBoxId: SafetyBoxTestId,
		Operator:    operator,
		Address:     returner,
	}
	res = h(ctx, msgSbDeregisterReturner)
	require.True(t, res.IsOK())

	// the owner deregisters an operator
	msgSbDeregisterOperator := MsgSafetyBoxDeregisterOperator{
		SafetyBoxId:    SafetyBoxTestId,
		SafetyBoxOwner: owner,
		Address:        operator,
	}
	res = h(ctx, msgSbDeregisterOperator)
	require.True(t, res.IsOK())
}
