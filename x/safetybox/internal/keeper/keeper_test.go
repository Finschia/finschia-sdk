package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/safetybox/internal/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

const (
	safetyBoxId = "test_safety_box_id"
)

// ToDo: check emitted events in the unit tests https://github.com/line/link/issues/270
func TestSafetyBox(t *testing.T) {
	input := SetupTestInput(t)
	_, ctx, keeper, ak := input.Cdc, input.Ctx, input.Keeper, input.Ak

	var sb types.SafetyBox
	var err sdk.Error

	// owner
	owner := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, owner)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, owner))
	require.Equal(t, uint64(0), ak.GetAccount(ctx, owner).GetAccountNumber())

	// create the safety box
	{
		// NewSafetyBox(ctx sdk.Context, msg types.MsgSafetyBoxCreate) (types.SafetyBox, sdk.Error)
		sb, err = keeper.NewSafetyBox(ctx, types.MsgSafetyBoxCreate{
			SafetyBoxId:     safetyBoxId,
			SafetyBoxOwner:  owner,
			SafetyBoxDenoms: []string{"link"},
		})
		require.NoError(t, err)
		require.Equal(t, owner, sb.Owner)
		require.Equal(t, safetyBoxId, sb.ID)
		require.Equal(t, sdk.Coins{}, sb.TotalAllocation)
		require.Equal(t, sdk.Coins{}, sb.CumulativeAllocation)
		require.Equal(t, sdk.Coins{}, sb.TotalIssuance)
	}

	// duplicated ID is not allowed
	{
		sb, err = keeper.NewSafetyBox(ctx, types.MsgSafetyBoxCreate{
			SafetyBoxId:     safetyBoxId,
			SafetyBoxOwner:  owner,
			SafetyBoxDenoms: []string{"link"},
		})
		require.EqualError(t, err, types.ErrSafetyBoxAccountExist(types.DefaultCodespace).Error())
	}

	// multiple denoms not allowed
	{
		sb, err = keeper.NewSafetyBox(ctx, types.MsgSafetyBoxCreate{
			SafetyBoxId:     "new_id",
			SafetyBoxOwner:  owner,
			SafetyBoxDenoms: []string{"link", "stake"},
		})
		require.EqualError(t, err, types.ErrSafetyBoxTooManyCoinDenoms(types.DefaultCodespace).Error())
	}

	// query the safety box
	{
		// GetSafetyBox(ctx sdk.Context, safetyBoxId string) (types.SafetyBox, sdk.Error)
		sb, err = keeper.GetSafetyBox(ctx, safetyBoxId)
		require.NoError(t, err)
		require.Equal(t, owner, sb.Owner)
		require.Equal(t, safetyBoxId, sb.ID)
		require.Equal(t, sdk.Coins(nil), sb.TotalAllocation)
		require.Equal(t, sdk.Coins(nil), sb.CumulativeAllocation)
		require.Equal(t, sdk.Coins(nil), sb.TotalIssuance)
	}

	// check permission of the owner
	{
		require.True(t, keeper.IsOwner(ctx, safetyBoxId, owner))
		require.False(t, keeper.IsOperator(ctx, safetyBoxId, owner))
		require.False(t, keeper.IsAllocator(ctx, safetyBoxId, owner))
		require.False(t, keeper.IsIssuer(ctx, safetyBoxId, owner))
		require.False(t, keeper.IsReturner(ctx, safetyBoxId, owner))
	}

	// operator
	operator := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, operator)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, operator))

	// the owner registers the operator
	{
		// GrantPermission(ctx sdk.Context, safetyBoxId string, by sdk.AccAddress, acc sdk.AccAddress, action string) sdk.Error
		err = keeper.GrantPermission(ctx, safetyBoxId, owner, operator, types.RoleOperator)
		require.NoError(t, err)
	}

	// check permission of the operator
	{
		require.False(t, keeper.IsOwner(ctx, safetyBoxId, operator))
		require.True(t, keeper.IsOperator(ctx, safetyBoxId, operator))
		require.False(t, keeper.IsAllocator(ctx, safetyBoxId, operator))
		require.False(t, keeper.IsIssuer(ctx, safetyBoxId, operator))
		require.False(t, keeper.IsReturner(ctx, safetyBoxId, operator))
	}

	// additional operator
	operator2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, operator2)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, operator2))

	// the owner registers second operator (multiple operators are allowed)
	{
		// GrantPermission(ctx sdk.Context, safetyBoxId string, by sdk.AccAddress, acc sdk.AccAddress, action string) sdk.Error
		err = keeper.GrantPermission(ctx, safetyBoxId, owner, operator2, types.RoleOperator)
		require.NoError(t, err)
	}

	// check permission of the operator2
	{
		require.False(t, keeper.IsOwner(ctx, safetyBoxId, operator2))
		require.True(t, keeper.IsOperator(ctx, safetyBoxId, operator2))
		require.False(t, keeper.IsAllocator(ctx, safetyBoxId, operator2))
		require.False(t, keeper.IsIssuer(ctx, safetyBoxId, operator2))
		require.False(t, keeper.IsReturner(ctx, safetyBoxId, operator2))
	}

	// the owner deregisters second operator
	{
		// RevokePermission(ctx sdk.Context, safetyBoxId string, by sdk.AccAddress, acc sdk.AccAddress, action string) sdk.Error
		err = keeper.RevokePermission(ctx, safetyBoxId, owner, operator2, types.RoleOperator)
		require.NoError(t, err)
	}

	// check permission of the operator2
	{
		require.False(t, keeper.IsOwner(ctx, safetyBoxId, operator2))
		require.False(t, keeper.IsOperator(ctx, safetyBoxId, operator2))
		require.False(t, keeper.IsAllocator(ctx, safetyBoxId, operator2))
		require.False(t, keeper.IsIssuer(ctx, safetyBoxId, operator2))
		require.False(t, keeper.IsReturner(ctx, safetyBoxId, operator2))
	}

	// allocator
	allocator := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, allocator)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, allocator))

	// issuer1
	issuer1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, issuer1)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, issuer1))

	// issuer2
	issuer2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, issuer2)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, issuer2))

	// returner
	returner := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, returner)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, returner))

	// the owner cannot register allocator, issuer or returner
	{
		err = keeper.GrantPermission(ctx, safetyBoxId, owner, allocator, types.RoleAllocator)
		require.EqualError(t, err, types.ErrSafetyBoxPermissionWhitelist(types.DefaultCodespace).Error())

		err = keeper.GrantPermission(ctx, safetyBoxId, owner, issuer1, types.RoleIssuer)
		require.EqualError(t, err, types.ErrSafetyBoxPermissionWhitelist(types.DefaultCodespace).Error())

		err = keeper.GrantPermission(ctx, safetyBoxId, owner, returner, types.RoleReturner)
		require.EqualError(t, err, types.ErrSafetyBoxPermissionWhitelist(types.DefaultCodespace).Error())
	}

	// the operator registers allocator, issuers and returner
	{
		err = keeper.GrantPermission(ctx, safetyBoxId, operator, allocator, types.RoleAllocator)
		require.NoError(t, err)

		err = keeper.GrantPermission(ctx, safetyBoxId, operator, issuer1, types.RoleIssuer)
		require.NoError(t, err)

		err = keeper.GrantPermission(ctx, safetyBoxId, operator, issuer2, types.RoleIssuer)
		require.NoError(t, err)

		err = keeper.GrantPermission(ctx, safetyBoxId, operator, returner, types.RoleReturner)
		require.NoError(t, err)
	}

	// put some coins to accounts
	_, err = input.Bk.AddCoins(ctx, owner, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(10))))
	require.NoError(t, err)

	_, err = input.Bk.AddCoins(ctx, operator, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(10))))
	require.NoError(t, err)

	_, err = input.Bk.AddCoins(ctx, allocator, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(10))))
	require.NoError(t, err)

	_, err = input.Bk.AddCoins(ctx, issuer1, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(10))))
	require.NoError(t, err)

	_, err = input.Bk.AddCoins(ctx, returner, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(10))))
	require.NoError(t, err)

	// allocation
	// Allocate(ctx sdk.Context, msg types.MsgSafetyBoxAllocateCoins) sdk.Error
	{
		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 0, 0, 0)

		// unable to allocate other coins
		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateCoins{
			SafetyBoxId:      safetyBoxId,
			AllocatorAddress: allocator,
			Coins:            sdk.Coins{sdk.Coin{Denom: "stake2", Amount: sdk.NewInt(3)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxIncorrectDenom(types.DefaultCodespace).Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 0, 0, 0)

		// the allocator allocates to the safety box
		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateCoins{
			SafetyBoxId:      safetyBoxId,
			AllocatorAddress: allocator,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(3)}},
		})
		require.NoError(t, err)

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 3, 3, 0)

		// the owner, operator, issuer, returner can not allocate
		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateCoins{
			SafetyBoxId:      safetyBoxId,
			AllocatorAddress: owner,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(3)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionAllocate(types.DefaultCodespace).Error())

		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateCoins{
			SafetyBoxId:      safetyBoxId,
			AllocatorAddress: operator,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionAllocate(types.DefaultCodespace).Error())

		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateCoins{
			SafetyBoxId:      safetyBoxId,
			AllocatorAddress: issuer1,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(3)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionAllocate(types.DefaultCodespace).Error())

		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateCoins{
			SafetyBoxId:      safetyBoxId,
			AllocatorAddress: returner,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(3)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionAllocate(types.DefaultCodespace).Error())

		// the safety box balance check - no changes
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 3, 3, 0)

		// insufficient fund
		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateCoins{
			SafetyBoxId:      safetyBoxId,
			AllocatorAddress: allocator,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1000)}},
		})
		require.EqualError(t, err, sdk.ErrInsufficientCoins("insufficient account funds; 7link < 1000link").Error())

		// the safety box balance check - no changes
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 3, 3, 0)
	}

	// recall
	// Recall(ctx sdk.Context, msg types.MsgSafetyBoxRecallCoins) sdk.Error
	{
		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 3, 3, 0)

		// unable to recalls other coins
		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallCoins{
			SafetyBoxId:      safetyBoxId,
			AllocatorAddress: allocator,
			Coins:            sdk.Coins{sdk.Coin{Denom: "stake2", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxIncorrectDenom(types.DefaultCodespace).Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 3, 3, 0)

		// the allocator recalls from the safety box
		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallCoins{
			SafetyBoxId:      safetyBoxId,
			AllocatorAddress: allocator,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.NoError(t, err)

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 2, 3, 0)

		// the owner, operator, issuer, returner can not recall
		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallCoins{
			SafetyBoxId:      safetyBoxId,
			AllocatorAddress: owner,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionRecall(types.DefaultCodespace).Error())

		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallCoins{
			SafetyBoxId:      safetyBoxId,
			AllocatorAddress: operator,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionRecall(types.DefaultCodespace).Error())

		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallCoins{
			SafetyBoxId:      safetyBoxId,
			AllocatorAddress: issuer1,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionRecall(types.DefaultCodespace).Error())

		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallCoins{
			SafetyBoxId:      safetyBoxId,
			AllocatorAddress: returner,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionRecall(types.DefaultCodespace).Error())

		// the safety box balance check - no changes
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 2, 3, 0)

		// insufficient fund
		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallCoins{
			SafetyBoxId:      safetyBoxId,
			AllocatorAddress: allocator,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1000)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxRecallMoreThanAllocated(types.DefaultCodespace).Error())

		// the safety box balance check - no changes
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 2, 3, 0)
	}

	// issue
	// Issue(ctx sdk.Context, msg types.MsgSafetyBoxIssueCoins) sdk.Error
	{
		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 2, 3, 0)

		// unable to issue other coins
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxId: safetyBoxId,
			ToAddress:   issuer1,
			FromAddress: issuer1,
			Coins:       sdk.Coins{sdk.Coin{Denom: "stake2", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxIncorrectDenom(types.DefaultCodespace).Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 2, 3, 0)

		// issuer1 request issuance to issuer1
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxId: safetyBoxId,
			ToAddress:   issuer1,
			FromAddress: issuer1,
			Coins:       sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.NoError(t, err)

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 2, 3, 1)

		// unable to issue other coins
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxId: safetyBoxId,
			ToAddress:   issuer2,
			FromAddress: issuer1,
			Coins:       sdk.Coins{sdk.Coin{Denom: "stake2", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxIncorrectDenom(types.DefaultCodespace).Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 2, 3, 1)

		// issuer1 request issuance to issuer2
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxId: safetyBoxId,
			ToAddress:   issuer2,
			FromAddress: issuer1,
			Coins:       sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.NoError(t, err)

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 2, 3, 2)

		// the owner, operator, allocator, returner can not issue
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxId: safetyBoxId,
			ToAddress:   issuer2,
			FromAddress: owner,
			Coins:       sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionIssue(types.DefaultCodespace).Error())

		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxId: safetyBoxId,
			ToAddress:   issuer2,
			FromAddress: operator,
			Coins:       sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionIssue(types.DefaultCodespace).Error())

		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxId: safetyBoxId,
			ToAddress:   issuer2,
			FromAddress: allocator,
			Coins:       sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionIssue(types.DefaultCodespace).Error())

		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxId: safetyBoxId,
			ToAddress:   issuer2,
			FromAddress: returner,
			Coins:       sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionIssue(types.DefaultCodespace).Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 2, 3, 2)

		// insufficient fund
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxId: safetyBoxId,
			ToAddress:   issuer2,
			FromAddress: issuer1,
			Coins:       sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1000)}},
		})
		require.EqualError(t, err, sdk.ErrInsufficientCoins("insufficient account funds;  < 1000link").Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 2, 3, 2)
	}

	// return
	// Return(ctx sdk.Context, msg types.MsgSafetyBoxReturnCoins) sdk.Error
	{
		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 2, 3, 2)

		// unable to return other coins
		err = keeper.Return(ctx, types.MsgSafetyBoxReturnCoins{
			SafetyBoxId:     safetyBoxId,
			ReturnerAddress: returner,
			Coins:           sdk.Coins{sdk.Coin{Denom: "stake2", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxIncorrectDenom(types.DefaultCodespace).Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 2, 3, 2)

		// the returner returns to the safety box
		err = keeper.Return(ctx, types.MsgSafetyBoxReturnCoins{
			SafetyBoxId:     safetyBoxId,
			ReturnerAddress: returner,
			Coins:           sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.NoError(t, err)

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 2, 3, 1)

		// the owner, operator, issuer, allocator can not return
		err = keeper.Return(ctx, types.MsgSafetyBoxReturnCoins{
			SafetyBoxId:     safetyBoxId,
			ReturnerAddress: owner,
			Coins:           sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionReturn(types.DefaultCodespace).Error())

		err = keeper.Return(ctx, types.MsgSafetyBoxReturnCoins{
			SafetyBoxId:     safetyBoxId,
			ReturnerAddress: operator,
			Coins:           sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionReturn(types.DefaultCodespace).Error())

		err = keeper.Return(ctx, types.MsgSafetyBoxReturnCoins{
			SafetyBoxId:     safetyBoxId,
			ReturnerAddress: allocator,
			Coins:           sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionReturn(types.DefaultCodespace).Error())

		err = keeper.Return(ctx, types.MsgSafetyBoxReturnCoins{
			SafetyBoxId:     safetyBoxId,
			ReturnerAddress: issuer1,
			Coins:           sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionReturn(types.DefaultCodespace).Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 2, 3, 1)

		// insufficient fund
		err = keeper.Return(ctx, types.MsgSafetyBoxReturnCoins{
			SafetyBoxId:     safetyBoxId,
			ReturnerAddress: returner,
			Coins:           sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1000)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxReturnMoreThanIssued(types.DefaultCodespace).Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxId, 2, 3, 1)
	}

	// check permissions via GetPermissions
	// GetPermissions(ctx sdk.Context, safetyBoxId, role string, acc sdk.AccAddress) (types.MsgSafetyBoxRoleResponse, sdk.Error)
	var pms types.MsgSafetyBoxRoleResponse
	{
		// owner to be only owner, not others
		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleOwner, owner)
		require.NoError(t, err)
		require.True(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleOperator, owner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleAllocator, owner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleIssuer, owner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleReturner, owner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		// operator to be only operator, not others
		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleOwner, operator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleOperator, operator)
		require.NoError(t, err)
		require.True(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleAllocator, operator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleIssuer, operator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleReturner, operator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		// allocator to be only allocator, not others
		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleOwner, allocator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleOperator, allocator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleAllocator, allocator)
		require.NoError(t, err)
		require.True(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleIssuer, allocator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleReturner, allocator)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		// issuer to be only issuer, not others
		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleOwner, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleOperator, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleAllocator, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleIssuer, issuer1)
		require.NoError(t, err)
		require.True(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleReturner, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		// issuer to be only issuer, not others
		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleOwner, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleOperator, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleAllocator, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleIssuer, issuer1)
		require.NoError(t, err)
		require.True(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleReturner, issuer1)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		// returner to be only returner, not others
		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleOwner, returner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleOperator, returner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleAllocator, returner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleIssuer, returner)
		require.NoError(t, err)
		require.False(t, pms.HasRole)

		pms, err = keeper.GetPermissions(ctx, safetyBoxId, types.RoleReturner, returner)
		require.NoError(t, err)
		require.True(t, pms.HasRole)
	}
}
