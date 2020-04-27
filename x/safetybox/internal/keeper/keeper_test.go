package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link/x/safetybox/internal/types"
	"github.com/line/link/x/token"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	safetyBoxID       = "test_safety_box_id"
	defaultContractID = "9be17165"
	tokenName         = "name"
	tokenSymbol       = "BTC"
	tokenMeta         = "{}"
	tokenImageURI     = "image-uri"
	tokenDecimals     = 6
	tokenAmount       = 1000
)

//nolint:dupl
func TestSafetyBox(t *testing.T) {
	input := SetupTestInput(t)
	_, ctx, keeper, tk := input.Cdc, input.Ctx, input.Keeper, input.Tk

	var sb types.SafetyBox
	var err error

	// issue token
	testToken := token.NewToken(defaultContractID, tokenName, tokenSymbol, tokenMeta, tokenImageURI, sdk.NewInt(tokenDecimals), true)

	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	err = tk.IssueToken(ctx, testToken, sdk.NewInt(tokenAmount), addr, addr)
	require.NoError(t, err)

	tok, err := tk.GetToken(ctx, defaultContractID)
	require.NoError(t, err)
	require.Equal(t, testToken.GetContractID(), tok.GetContractID())

	// owner
	owner := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc, err := tk.NewAccountWithAddress(ctx, defaultContractID, owner)
		require.NotNil(t, acc)
		require.NoError(t, err)
	}
	acc, err := tk.GetAccount(ctx, defaultContractID, owner)
	require.NotNil(t, acc)
	require.NoError(t, err)

	// create the safety box
	{
		// NewSafetyBox(ctx sdk.Context, msg types.MsgSafetyBoxCreate) (types.SafetyBox, sdk.Error)
		sb, err = keeper.NewSafetyBox(ctx, types.MsgSafetyBoxCreate{
			SafetyBoxID:    safetyBoxID,
			SafetyBoxOwner: owner,
			ContractID:     defaultContractID,
		})
		require.NoError(t, err)
		require.Equal(t, owner, sb.Owner)
		require.Equal(t, safetyBoxID, sb.ID)
		require.Equal(t, defaultContractID, sb.ContractID)
		require.Equal(t, sdk.ZeroInt(), sb.TotalAllocation)
		require.Equal(t, sdk.ZeroInt(), sb.CumulativeAllocation)
		require.Equal(t, sdk.ZeroInt(), sb.TotalIssuance)
	}

	// duplicated ID is not allowed
	{
		sb, err = keeper.NewSafetyBox(ctx, types.MsgSafetyBoxCreate{
			SafetyBoxID:    safetyBoxID,
			SafetyBoxOwner: owner,
			ContractID:     defaultContractID,
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxAccountExist, "ID: %s", safetyBoxID).Error())
	}

	// query non-exist safety box
	{
		sb, err = keeper.GetSafetyBox(ctx, "no-box")
		require.EqualError(t, err, sdkerrors.Wrap(types.ErrSafetyBoxNotExist, "ID: no-box").Error())
	}

	// query the safety box
	{
		// GetSafetyBox(ctx sdk.Context, safetyBoxId string) (types.SafetyBox, sdk.Error)
		sb, err = keeper.GetSafetyBox(ctx, safetyBoxID)
		require.NoError(t, err)
		require.Equal(t, owner, sb.Owner)
		require.Equal(t, safetyBoxID, sb.ID)
		require.Equal(t, sdk.ZeroInt(), sb.TotalAllocation)
		require.Equal(t, sdk.ZeroInt(), sb.CumulativeAllocation)
		require.Equal(t, sdk.ZeroInt(), sb.TotalIssuance)
	}

	// check permission of the owner
	{
		require.True(t, keeper.IsOwner(ctx, safetyBoxID, owner))
		require.False(t, keeper.IsOperator(ctx, safetyBoxID, owner))
		require.False(t, keeper.IsAllocator(ctx, safetyBoxID, owner))
		require.False(t, keeper.IsIssuer(ctx, safetyBoxID, owner))
		require.False(t, keeper.IsReturner(ctx, safetyBoxID, owner))
	}

	// operator
	operator := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc, err := tk.NewAccountWithAddress(ctx, defaultContractID, operator)
		require.NotNil(t, acc)
		require.NoError(t, err)
	}
	acc, err = tk.GetAccount(ctx, defaultContractID, operator)
	require.NotNil(t, acc)
	require.NoError(t, err)

	// the owner registers the operator
	{
		// GrantPermission(ctx sdk.Context, safetyBoxId string, by sdk.AccAddress, acc sdk.AccAddress, action string) sdk.Error
		err = keeper.GrantPermission(ctx, safetyBoxID, owner, operator, types.RoleOperator)
		require.NoError(t, err)
	}

	// check permission of the operator
	{
		require.False(t, keeper.IsOwner(ctx, safetyBoxID, operator))
		require.True(t, keeper.IsOperator(ctx, safetyBoxID, operator))
		require.False(t, keeper.IsAllocator(ctx, safetyBoxID, operator))
		require.False(t, keeper.IsIssuer(ctx, safetyBoxID, operator))
		require.False(t, keeper.IsReturner(ctx, safetyBoxID, operator))
	}

	// additional operator
	operator2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc, err := tk.NewAccountWithAddress(ctx, defaultContractID, operator2)
		require.NotNil(t, acc)
		require.NoError(t, err)
	}
	acc, err = tk.GetAccount(ctx, defaultContractID, operator2)
	require.NotNil(t, acc)
	require.NoError(t, err)

	// the owner registers second operator (multiple operators are allowed)
	{
		// GrantPermission(ctx sdk.Context, safetyBoxId string, by sdk.AccAddress, acc sdk.AccAddress, action string) sdk.Error
		err = keeper.GrantPermission(ctx, safetyBoxID, owner, operator2, types.RoleOperator)
		require.NoError(t, err)
	}

	// check permission of the operator2
	{
		require.False(t, keeper.IsOwner(ctx, safetyBoxID, operator2))
		require.True(t, keeper.IsOperator(ctx, safetyBoxID, operator2))
		require.False(t, keeper.IsAllocator(ctx, safetyBoxID, operator2))
		require.False(t, keeper.IsIssuer(ctx, safetyBoxID, operator2))
		require.False(t, keeper.IsReturner(ctx, safetyBoxID, operator2))
	}

	// the owner deregisters second operator
	{
		// RevokePermission(ctx sdk.Context, safetyBoxId string, by sdk.AccAddress, acc sdk.AccAddress, action string) sdk.Error
		err = keeper.RevokePermission(ctx, safetyBoxID, owner, operator2, types.RoleOperator)
		require.NoError(t, err)
	}

	// revoke error case
	{
		// RevokePermission(ctx sdk.Context, safetyBoxId string, by sdk.AccAddress, acc sdk.AccAddress, action string) sdk.Error
		err = keeper.RevokePermission(ctx, safetyBoxID, owner, owner, types.RoleOperator)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxSelfPermission, "Account: %s", owner.String()).Error())

		err = keeper.RevokePermission(ctx, safetyBoxID, operator2, owner, types.RoleAllocator)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionWhitelist, "Account: %s", operator2.String()).Error())

		err = keeper.RevokePermission(ctx, safetyBoxID, operator, operator2, types.RoleAllocator)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxDoesNotHavePermission, "Account: %s", operator2.String()).Error())

		err = keeper.RevokePermission(ctx, safetyBoxID, operator, operator2, types.RoleIssuer)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxDoesNotHavePermission, "Account: %s", operator2.String()).Error())

		err = keeper.RevokePermission(ctx, safetyBoxID, operator, operator2, types.RoleReturner)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxDoesNotHavePermission, "Account: %s", operator2.String()).Error())
	}

	// check permission of the operator2
	{
		require.False(t, keeper.IsOwner(ctx, safetyBoxID, operator2))
		require.False(t, keeper.IsOperator(ctx, safetyBoxID, operator2))
		require.False(t, keeper.IsAllocator(ctx, safetyBoxID, operator2))
		require.False(t, keeper.IsIssuer(ctx, safetyBoxID, operator2))
		require.False(t, keeper.IsReturner(ctx, safetyBoxID, operator2))
	}

	// allocator
	allocator := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc, err := tk.NewAccountWithAddress(ctx, defaultContractID, allocator)
		require.NotNil(t, acc)
		require.NoError(t, err)
	}
	acc, err = tk.GetAccount(ctx, defaultContractID, allocator)
	require.NotNil(t, acc)
	require.NoError(t, err)

	// issuer1
	issuer1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc, err := tk.NewAccountWithAddress(ctx, defaultContractID, issuer1)
		require.NotNil(t, acc)
		require.NoError(t, err)
	}
	acc, err = tk.GetAccount(ctx, defaultContractID, issuer1)
	require.NotNil(t, acc)
	require.NoError(t, err)

	// issuer2
	issuer2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc, err := tk.NewAccountWithAddress(ctx, defaultContractID, issuer2)
		require.NotNil(t, acc)
		require.NoError(t, err)
	}
	acc, err = tk.GetAccount(ctx, defaultContractID, issuer2)
	require.NotNil(t, acc)
	require.NoError(t, err)

	// returner
	returner := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc, err := tk.NewAccountWithAddress(ctx, defaultContractID, returner)
		require.NotNil(t, acc)
		require.NoError(t, err)
	}
	acc, err = tk.GetAccount(ctx, defaultContractID, returner)
	require.NotNil(t, acc)
	require.NoError(t, err)

	// the owner cannot register allocator, issuer or returner
	{
		err = keeper.GrantPermission(ctx, safetyBoxID, owner, allocator, types.RoleAllocator)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionWhitelist, "Account: %s", owner.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, owner, issuer1, types.RoleIssuer)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionWhitelist, "Account: %s", owner.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, owner, returner, types.RoleReturner)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionWhitelist, "Account: %s", owner.String()).Error())
	}

	// the operator registers allocator, issuers and returner
	{
		err = keeper.GrantPermission(ctx, safetyBoxID, operator, allocator, types.RoleAllocator)
		require.NoError(t, err)

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, issuer1, types.RoleIssuer)
		require.NoError(t, err)

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, issuer2, types.RoleIssuer)
		require.NoError(t, err)

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, returner, types.RoleReturner)
		require.NoError(t, err)
	}

	// put amounts of token to accounts
	_, err = input.Tk.AddBalance(ctx, defaultContractID, owner, sdk.NewInt(10))
	require.NoError(t, err)

	_, err = input.Tk.AddBalance(ctx, defaultContractID, operator, sdk.NewInt(10))
	require.NoError(t, err)

	_, err = input.Tk.AddBalance(ctx, defaultContractID, allocator, sdk.NewInt(10))
	require.NoError(t, err)

	_, err = input.Tk.AddBalance(ctx, defaultContractID, issuer1, sdk.NewInt(10))
	require.NoError(t, err)

	_, err = input.Tk.AddBalance(ctx, defaultContractID, returner, sdk.NewInt(10))
	require.NoError(t, err)

	// allocation
	// Allocate(ctx sdk.Context, msg types.MsgSafetyBoxAllocateToken) sdk.Error
	{
		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 0, 0, 0)

		// unable to allocate other token
		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateToken{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: allocator,
			ContractID:       "5ca12345",
			Amount:           sdk.NewInt(3),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxIncorrectContractID, "Expected: %s, Requested: %s", defaultContractID, "5ca12345").Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 0, 0, 0)

		// the allocator allocates to the safety box
		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateToken{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: allocator,
			ContractID:       defaultContractID,
			Amount:           sdk.NewInt(3),
		})
		require.NoError(t, err)

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 3, 3, 0)

		// the owner, operator, issuer, returner can not allocate
		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateToken{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: owner,
			ContractID:       defaultContractID,
			Amount:           sdk.NewInt(3),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionAllocate, "Account: %s", owner.String()).Error())

		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateToken{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: operator,
			ContractID:       defaultContractID,
			Amount:           sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionAllocate, "Account: %s", operator.String()).Error())

		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateToken{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: issuer1,
			ContractID:       defaultContractID,
			Amount:           sdk.NewInt(3),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionAllocate, "Account: %s", issuer1.String()).Error())

		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateToken{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: returner,
			ContractID:       defaultContractID,
			Amount:           sdk.NewInt(3),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionAllocate, "Account: %s", returner.String()).Error())

		// the safety box balance check - no changes
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 3, 3, 0)

		// insufficient fund
		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateToken{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: allocator,
			ContractID:       defaultContractID,
			Amount:           sdk.NewInt(1000),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(token.ErrInsufficientBalance, "insufficient account funds for token [%s]; 7 < 1000", defaultContractID).Error())

		// the safety box balance check - no changes
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 3, 3, 0)
	}

	// recall
	// Recall(ctx sdk.Context, msg types.MsgSafetyBoxRecallToken) sdk.Error
	{
		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 3, 3, 0)

		// unable to recalls other token
		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallToken{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: allocator,
			ContractID:       "5ca12345",
			Amount:           sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxIncorrectContractID, "Expected: %s, Requested: %s", defaultContractID, "5ca12345").Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 3, 3, 0)

		// the allocator recalls from the safety box
		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallToken{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: allocator,
			ContractID:       defaultContractID,
			Amount:           sdk.NewInt(1),
		})
		require.NoError(t, err)

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 0)

		// the owner, operator, issuer, returner can not recall
		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallToken{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: owner,
			ContractID:       defaultContractID,
			Amount:           sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionRecall, "Account: %s", owner.String()).Error())

		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallToken{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: operator,
			ContractID:       defaultContractID,
			Amount:           sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionRecall, "Account: %s", operator.String()).Error())

		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallToken{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: issuer1,
			ContractID:       defaultContractID,
			Amount:           sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionRecall, "Account: %s", issuer1.String()).Error())

		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallToken{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: returner,
			ContractID:       defaultContractID,
			Amount:           sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionRecall, "Account: %s", returner.String()).Error())

		// the safety box balance check - no changes
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 0)

		// insufficient fund
		sb, err := keeper.GetSafetyBox(ctx, safetyBoxID)
		require.NoError(t, err)
		toRecall := sdk.NewInt(1000)
		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallToken{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: allocator,
			ContractID:       defaultContractID,
			Amount:           toRecall,
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxRecallMoreThanAllocated, "Has: %v, Requested: %v", sb.TotalAllocation, toRecall).Error())

		// the safety box balance check - no changes
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 0)
	}

	// issue
	// Issue(ctx sdk.Context, msg types.MsgSafetyBoxIssueToken) sdk.Error
	{
		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 0)

		// unable to issue other token
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueToken{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer1,
			FromAddress: issuer1,
			ContractID:  "5ca12345",
			Amount:      sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxIncorrectContractID, "Expected: %s, Requested: %s", defaultContractID, "5ca12345").Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 0)

		// issuer1 request issuance to issuer1
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueToken{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer1,
			FromAddress: issuer1,
			ContractID:  defaultContractID,
			Amount:      sdk.NewInt(1),
		})
		require.NoError(t, err)

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 1)

		// unable to issue other token
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueToken{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer2,
			FromAddress: issuer1,
			ContractID:  "5ca12345",
			Amount:      sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxIncorrectContractID, "Expected: %s, Requested: %s", defaultContractID, "5ca12345").Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 1)

		// issuer1 request issuance to issuer2
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueToken{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer2,
			FromAddress: issuer1,
			ContractID:  defaultContractID,
			Amount:      sdk.NewInt(1),
		})
		require.NoError(t, err)

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 2)

		// the owner, operator, allocator, returner can not issue
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueToken{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer2,
			FromAddress: owner,
			ContractID:  defaultContractID,
			Amount:      sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionIssue, "Account: %s", owner.String()).Error())

		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueToken{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer2,
			FromAddress: operator,
			ContractID:  defaultContractID,
			Amount:      sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionIssue, "Account: %s", operator.String()).Error())

		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueToken{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer2,
			FromAddress: allocator,
			ContractID:  defaultContractID,
			Amount:      sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionIssue, "Account: %s", allocator.String()).Error())

		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueToken{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer2,
			FromAddress: returner,
			ContractID:  defaultContractID,
			Amount:      sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionIssue, "Account: %s", returner.String()).Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 2)

		// insufficient fund
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueToken{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer2,
			FromAddress: issuer1,
			ContractID:  defaultContractID,
			Amount:      sdk.NewInt(1000),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(token.ErrInsufficientBalance, "insufficient account funds for token [%s]; 0 < 1000", defaultContractID).Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 2)
	}

	// return
	// Return(ctx sdk.Context, msg types.MsgSafetyBoxReturnToken) sdk.Error
	{
		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 2)

		// unable to return other token
		err = keeper.Return(ctx, types.MsgSafetyBoxReturnToken{
			SafetyBoxID:     safetyBoxID,
			ReturnerAddress: returner,
			ContractID:      "5ca12345",
			Amount:          sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxIncorrectContractID, "Expected: %s, Requested: %s", defaultContractID, "5ca12345").Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 2)

		// the returner returns to the safety box
		err = keeper.Return(ctx, types.MsgSafetyBoxReturnToken{
			SafetyBoxID:     safetyBoxID,
			ReturnerAddress: returner,
			ContractID:      defaultContractID,
			Amount:          sdk.NewInt(1),
		})
		require.NoError(t, err)

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 1)

		// the owner, operator, issuer, allocator can not return
		err = keeper.Return(ctx, types.MsgSafetyBoxReturnToken{
			SafetyBoxID:     safetyBoxID,
			ReturnerAddress: owner,
			ContractID:      defaultContractID,
			Amount:          sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionReturn, "Account: %s", owner.String()).Error())

		err = keeper.Return(ctx, types.MsgSafetyBoxReturnToken{
			SafetyBoxID:     safetyBoxID,
			ReturnerAddress: operator,
			ContractID:      defaultContractID,
			Amount:          sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionReturn, "Account: %s", operator.String()).Error())

		err = keeper.Return(ctx, types.MsgSafetyBoxReturnToken{
			SafetyBoxID:     safetyBoxID,
			ReturnerAddress: allocator,
			ContractID:      defaultContractID,
			Amount:          sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionReturn, "Account: %s", allocator.String()).Error())

		err = keeper.Return(ctx, types.MsgSafetyBoxReturnToken{
			SafetyBoxID:     safetyBoxID,
			ReturnerAddress: issuer1,
			ContractID:      defaultContractID,
			Amount:          sdk.NewInt(1),
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxPermissionReturn, "Account: %s", issuer1.String()).Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 1)

		// insufficient fund
		sb, err := keeper.GetSafetyBox(ctx, safetyBoxID)
		require.NoError(t, err)
		toReturn := sdk.NewInt(1000)
		err = keeper.Return(ctx, types.MsgSafetyBoxReturnToken{
			SafetyBoxID:     safetyBoxID,
			ReturnerAddress: returner,
			ContractID:      defaultContractID,
			Amount:          toReturn,
		})
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxReturnMoreThanIssued, "Has: %v, Requested: %v", sb.TotalIssuance, toReturn).Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 1)
	}

	// check permissions via GetPermissions
	// GetPermissions(ctx sdk.Context, safetyBoxId, role string, acc sdk.AccAddress) (types.MsgSafetyBoxRoleResponse, sdk.Error)
	var pms types.MsgSafetyBoxRoleResponse
	{
		// owner to be only owner, not others
		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleOwner, owner)
		require.NoError(t, err)
		require.True(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleOperator, owner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleAllocator, owner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleIssuer, owner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleReturner, owner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		// operator to be only operator, not others
		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleOwner, operator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleOperator, operator)
		require.NoError(t, err)
		require.True(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleAllocator, operator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleIssuer, operator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleReturner, operator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		// allocator to be only allocator, not others
		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleOwner, allocator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleOperator, allocator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleAllocator, allocator)
		require.NoError(t, err)
		require.True(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleIssuer, allocator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleReturner, allocator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		// issuer to be only issuer, not others
		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleOwner, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleOperator, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleAllocator, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleIssuer, issuer1)
		require.NoError(t, err)
		require.True(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleReturner, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		// issuer to be only issuer, not others
		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleOwner, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleOperator, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleAllocator, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleIssuer, issuer1)
		require.NoError(t, err)
		require.True(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleReturner, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		// returner to be only returner, not others
		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleOwner, returner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleOperator, returner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleAllocator, returner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleIssuer, returner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleReturner, returner)
		require.NoError(t, err)
		require.True(t, pms.HasRole)
	}

	// try get invalid role permission -> should fail w/ ErrSafetyBoxInvalidRole
	{
		invalidRole := "invalidRole"
		pms, err = keeper.GetPermissions(ctx, safetyBoxID, invalidRole, returner)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxInvalidRole, "Role: %s", invalidRole).Error())
		require.False(t, pms.HasRole)
	}

	// try grant the same permission again -> all should fail w/ ErrSafetyBoxHasPermissionAlready
	{
		err = keeper.GrantPermission(ctx, safetyBoxID, owner, operator, types.RoleOperator)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxHasPermissionAlready, "Account: %s", operator.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, allocator, types.RoleAllocator)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxHasPermissionAlready, "Account: %s", allocator.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, issuer1, types.RoleIssuer)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxHasPermissionAlready, "Account: %s", issuer1.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, returner, types.RoleReturner)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxHasPermissionAlready, "Account: %s", returner.String()).Error())
	}

	// try grant other permissions again -> all should fail w/ ErrSafetyBoxHasOtherPermission
	{
		err = keeper.GrantPermission(ctx, safetyBoxID, owner, operator2, types.RoleOperator)
		require.NoError(t, err)

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, operator2, types.RoleAllocator)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxHasOtherPermission, "Account: %s", operator2.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, operator2, types.RoleIssuer)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxHasOtherPermission, "Account: %s", operator2.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, operator2, types.RoleReturner)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxHasOtherPermission, "Account: %s", operator2.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, allocator, types.RoleIssuer)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxHasOtherPermission, "Account: %s", allocator.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, allocator, types.RoleReturner)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxHasOtherPermission, "Account: %s", allocator.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, owner, allocator, types.RoleOperator)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxHasOtherPermission, "Account: %s", allocator.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, issuer1, types.RoleAllocator)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxHasOtherPermission, "Account: %s", issuer1.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, issuer1, types.RoleReturner)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxHasOtherPermission, "Account: %s", issuer1.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, owner, issuer1, types.RoleOperator)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxHasOtherPermission, "Account: %s", issuer1.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, returner, types.RoleAllocator)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxHasOtherPermission, "Account: %s", returner.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, returner, types.RoleIssuer)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxHasOtherPermission, "Account: %s", returner.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, owner, returner, types.RoleOperator)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrSafetyBoxHasOtherPermission, "Account: %s", returner.String()).Error())
	}

	// revoke permissions
	{
		err = keeper.RevokePermission(ctx, safetyBoxID, operator, allocator, types.RoleAllocator)
		require.NoError(t, err)

		err = keeper.RevokePermission(ctx, safetyBoxID, operator, issuer1, types.RoleIssuer)
		require.NoError(t, err)

		err = keeper.RevokePermission(ctx, safetyBoxID, operator, returner, types.RoleReturner)
		require.NoError(t, err)

		err = keeper.RevokePermission(ctx, safetyBoxID, owner, operator, types.RoleOperator)
		require.NoError(t, err)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleOperator, operator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleAllocator, allocator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleIssuer, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxID, types.RoleReturner, returner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)
	}
}
