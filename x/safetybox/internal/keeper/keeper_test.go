package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/safetybox/internal/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	safetyBoxID = "test_safety_box_id"
)

//nolint:dupl
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
			SafetyBoxID:     safetyBoxID,
			SafetyBoxOwner:  owner,
			SafetyBoxDenoms: []string{"link"},
		})
		require.NoError(t, err)
		require.Equal(t, owner, sb.Owner)
		require.Equal(t, safetyBoxID, sb.ID)
		require.Equal(t, sdk.Coins{}, sb.TotalAllocation)
		require.Equal(t, sdk.Coins{}, sb.CumulativeAllocation)
		require.Equal(t, sdk.Coins{}, sb.TotalIssuance)
	}

	// duplicated ID is not allowed
	{
		sb, err = keeper.NewSafetyBox(ctx, types.MsgSafetyBoxCreate{
			SafetyBoxID:     safetyBoxID,
			SafetyBoxOwner:  owner,
			SafetyBoxDenoms: []string{"link"},
		})
		require.EqualError(t, err, types.ErrSafetyBoxAccountExist(types.DefaultCodespace, safetyBoxID).Error())
	}

	// multiple denoms not allowed
	{
		tooManyDenoms := []string{"link", "stake"}
		sb, err = keeper.NewSafetyBox(ctx, types.MsgSafetyBoxCreate{
			SafetyBoxID:     "new_id",
			SafetyBoxOwner:  owner,
			SafetyBoxDenoms: tooManyDenoms,
		})
		require.EqualError(t, err, types.ErrSafetyBoxTooManyCoinDenoms(types.DefaultCodespace, tooManyDenoms).Error())
	}

	// query non-exist safety box
	{
		sb, err = keeper.GetSafetyBox(ctx, "no-box")
		require.EqualError(t, err, types.ErrSafetyBoxNotExist(types.DefaultCodespace, "no-box").Error())
	}

	// query the safety box
	{
		// GetSafetyBox(ctx sdk.Context, safetyBoxId string) (types.SafetyBox, sdk.Error)
		sb, err = keeper.GetSafetyBox(ctx, safetyBoxID)
		require.NoError(t, err)
		require.Equal(t, owner, sb.Owner)
		require.Equal(t, safetyBoxID, sb.ID)
		require.Equal(t, sdk.Coins(nil), sb.TotalAllocation)
		require.Equal(t, sdk.Coins(nil), sb.CumulativeAllocation)
		require.Equal(t, sdk.Coins(nil), sb.TotalIssuance)
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
		acc := ak.NewAccountWithAddress(ctx, operator)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, operator))

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
		acc := ak.NewAccountWithAddress(ctx, operator2)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, operator2))

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
		require.EqualError(t, err, types.ErrSafetyBoxSelfPermission(types.DefaultCodespace, owner.String()).Error())

		err = keeper.RevokePermission(ctx, safetyBoxID, operator2, owner, types.RoleAllocator)
		require.EqualError(t, err, types.ErrSafetyBoxPermissionWhitelist(types.DefaultCodespace, operator2.String()).Error())

		err = keeper.RevokePermission(ctx, safetyBoxID, operator, operator2, types.RoleAllocator)
		require.EqualError(t, err, types.ErrSafetyBoxDoesNotHavePermission(types.DefaultCodespace, operator2.String()).Error())

		err = keeper.RevokePermission(ctx, safetyBoxID, operator, operator2, types.RoleIssuer)
		require.EqualError(t, err, types.ErrSafetyBoxDoesNotHavePermission(types.DefaultCodespace, operator2.String()).Error())

		err = keeper.RevokePermission(ctx, safetyBoxID, operator, operator2, types.RoleReturner)
		require.EqualError(t, err, types.ErrSafetyBoxDoesNotHavePermission(types.DefaultCodespace, operator2.String()).Error())
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
		err = keeper.GrantPermission(ctx, safetyBoxID, owner, allocator, types.RoleAllocator)
		require.EqualError(t, err, types.ErrSafetyBoxPermissionWhitelist(types.DefaultCodespace, owner.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, owner, issuer1, types.RoleIssuer)
		require.EqualError(t, err, types.ErrSafetyBoxPermissionWhitelist(types.DefaultCodespace, owner.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, owner, returner, types.RoleReturner)
		require.EqualError(t, err, types.ErrSafetyBoxPermissionWhitelist(types.DefaultCodespace, owner.String()).Error())
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
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 0, 0, 0)

		// unable to allocate other coins
		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateCoins{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: allocator,
			Coins:            sdk.Coins{sdk.Coin{Denom: "stake2", Amount: sdk.NewInt(3)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxIncorrectDenom(types.DefaultCodespace, "link", "stake2").Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 0, 0, 0)

		// the allocator allocates to the safety box
		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateCoins{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: allocator,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(3)}},
		})
		require.NoError(t, err)

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 3, 3, 0)

		// the owner, operator, issuer, returner can not allocate
		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateCoins{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: owner,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(3)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionAllocate(types.DefaultCodespace, owner.String()).Error())

		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateCoins{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: operator,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionAllocate(types.DefaultCodespace, operator.String()).Error())

		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateCoins{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: issuer1,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(3)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionAllocate(types.DefaultCodespace, issuer1.String()).Error())

		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateCoins{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: returner,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(3)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionAllocate(types.DefaultCodespace, returner.String()).Error())

		// the safety box balance check - no changes
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 3, 3, 0)

		// insufficient fund
		err = keeper.Allocate(ctx, types.MsgSafetyBoxAllocateCoins{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: allocator,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1000)}},
		})
		require.EqualError(t, err, sdk.ErrInsufficientCoins("insufficient account funds; 7link < 1000link").Error())

		// the safety box balance check - no changes
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 3, 3, 0)
	}

	// recall
	// Recall(ctx sdk.Context, msg types.MsgSafetyBoxRecallCoins) sdk.Error
	{
		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 3, 3, 0)

		// unable to recalls other coins
		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallCoins{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: allocator,
			Coins:            sdk.Coins{sdk.Coin{Denom: "stake2", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxIncorrectDenom(types.DefaultCodespace, "link", "stake2").Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 3, 3, 0)

		// the allocator recalls from the safety box
		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallCoins{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: allocator,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.NoError(t, err)

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 0)

		// the owner, operator, issuer, returner can not recall
		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallCoins{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: owner,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionRecall(types.DefaultCodespace, owner.String()).Error())

		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallCoins{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: operator,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionRecall(types.DefaultCodespace, operator.String()).Error())

		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallCoins{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: issuer1,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionRecall(types.DefaultCodespace, issuer1.String()).Error())

		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallCoins{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: returner,
			Coins:            sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionRecall(types.DefaultCodespace, returner.String()).Error())

		// the safety box balance check - no changes
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 0)

		// insufficient fund
		sb, err := keeper.GetSafetyBox(ctx, safetyBoxID)
		require.NoError(t, err)
		toRecall := sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1000)}}
		err = keeper.Recall(ctx, types.MsgSafetyBoxRecallCoins{
			SafetyBoxID:      safetyBoxID,
			AllocatorAddress: allocator,
			Coins:            toRecall,
		})
		require.EqualError(t, err, types.ErrSafetyBoxRecallMoreThanAllocated(types.DefaultCodespace, sb.TotalAllocation, toRecall).Error())

		// the safety box balance check - no changes
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 0)
	}

	// issue
	// Issue(ctx sdk.Context, msg types.MsgSafetyBoxIssueCoins) sdk.Error
	{
		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 0)

		// unable to issue other coins
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer1,
			FromAddress: issuer1,
			Coins:       sdk.Coins{sdk.Coin{Denom: "stake2", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxIncorrectDenom(types.DefaultCodespace, "link", "stake2").Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 0)

		// issuer1 request issuance to issuer1
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer1,
			FromAddress: issuer1,
			Coins:       sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.NoError(t, err)

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 1)

		// unable to issue other coins
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer2,
			FromAddress: issuer1,
			Coins:       sdk.Coins{sdk.Coin{Denom: "stake2", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxIncorrectDenom(types.DefaultCodespace, "link", "stake2").Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 1)

		// issuer1 request issuance to issuer2
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer2,
			FromAddress: issuer1,
			Coins:       sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.NoError(t, err)

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 2)

		// the owner, operator, allocator, returner can not issue
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer2,
			FromAddress: owner,
			Coins:       sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionIssue(types.DefaultCodespace, owner.String()).Error())

		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer2,
			FromAddress: operator,
			Coins:       sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionIssue(types.DefaultCodespace, operator.String()).Error())

		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer2,
			FromAddress: allocator,
			Coins:       sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionIssue(types.DefaultCodespace, allocator.String()).Error())

		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer2,
			FromAddress: returner,
			Coins:       sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionIssue(types.DefaultCodespace, returner.String()).Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 2)

		// insufficient fund
		err = keeper.Issue(ctx, types.MsgSafetyBoxIssueCoins{
			SafetyBoxID: safetyBoxID,
			ToAddress:   issuer2,
			FromAddress: issuer1,
			Coins:       sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1000)}},
		})
		require.EqualError(t, err, sdk.ErrInsufficientCoins("insufficient account funds;  < 1000link").Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 2)
	}

	// return
	// Return(ctx sdk.Context, msg types.MsgSafetyBoxReturnCoins) sdk.Error
	{
		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 2)

		// unable to return other coins
		err = keeper.Return(ctx, types.MsgSafetyBoxReturnCoins{
			SafetyBoxID:     safetyBoxID,
			ReturnerAddress: returner,
			Coins:           sdk.Coins{sdk.Coin{Denom: "stake2", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxIncorrectDenom(types.DefaultCodespace, "link", "stake2").Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 2)

		// the returner returns to the safety box
		err = keeper.Return(ctx, types.MsgSafetyBoxReturnCoins{
			SafetyBoxID:     safetyBoxID,
			ReturnerAddress: returner,
			Coins:           sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.NoError(t, err)

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 1)

		// the owner, operator, issuer, allocator can not return
		err = keeper.Return(ctx, types.MsgSafetyBoxReturnCoins{
			SafetyBoxID:     safetyBoxID,
			ReturnerAddress: owner,
			Coins:           sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionReturn(types.DefaultCodespace, owner.String()).Error())

		err = keeper.Return(ctx, types.MsgSafetyBoxReturnCoins{
			SafetyBoxID:     safetyBoxID,
			ReturnerAddress: operator,
			Coins:           sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionReturn(types.DefaultCodespace, operator.String()).Error())

		err = keeper.Return(ctx, types.MsgSafetyBoxReturnCoins{
			SafetyBoxID:     safetyBoxID,
			ReturnerAddress: allocator,
			Coins:           sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionReturn(types.DefaultCodespace, allocator.String()).Error())

		err = keeper.Return(ctx, types.MsgSafetyBoxReturnCoins{
			SafetyBoxID:     safetyBoxID,
			ReturnerAddress: issuer1,
			Coins:           sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1)}},
		})
		require.EqualError(t, err, types.ErrSafetyBoxPermissionReturn(types.DefaultCodespace, issuer1.String()).Error())

		// the safety box balance check
		checkSafetyBoxBalance(t, keeper, ctx, safetyBoxID, 2, 3, 1)

		// insufficient fund
		sb, err := keeper.GetSafetyBox(ctx, safetyBoxID)
		require.NoError(t, err)
		toReturn := sdk.Coins{sdk.Coin{Denom: "link", Amount: sdk.NewInt(1000)}}
		err = keeper.Return(ctx, types.MsgSafetyBoxReturnCoins{
			SafetyBoxID:     safetyBoxID,
			ReturnerAddress: returner,
			Coins:           toReturn,
		})
		require.EqualError(t, err, types.ErrSafetyBoxReturnMoreThanIssued(types.DefaultCodespace, sb.TotalIssuance, toReturn).Error())

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
		require.EqualError(t, err, types.ErrSafetyBoxInvalidRole(types.DefaultCodespace, invalidRole).Error())
		require.False(t, pms.HasRole)
	}

	// try grant the same permission again -> all should fail w/ ErrSafetyBoxHasPermissionAlready
	{
		err = keeper.GrantPermission(ctx, safetyBoxID, owner, operator, types.RoleOperator)
		require.EqualError(t, err, types.ErrSafetyBoxHasPermissionAlready(types.DefaultCodespace, operator.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, allocator, types.RoleAllocator)
		require.EqualError(t, err, types.ErrSafetyBoxHasPermissionAlready(types.DefaultCodespace, allocator.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, issuer1, types.RoleIssuer)
		require.EqualError(t, err, types.ErrSafetyBoxHasPermissionAlready(types.DefaultCodespace, issuer1.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, returner, types.RoleReturner)
		require.EqualError(t, err, types.ErrSafetyBoxHasPermissionAlready(types.DefaultCodespace, returner.String()).Error())
	}

	// try grant other permissions again -> all should fail w/ ErrSafetyBoxHasOtherPermission
	{
		err = keeper.GrantPermission(ctx, safetyBoxID, owner, operator2, types.RoleOperator)
		require.NoError(t, err)

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, operator2, types.RoleAllocator)
		require.EqualError(t, err, types.ErrSafetyBoxHasOtherPermission(types.DefaultCodespace, operator2.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, operator2, types.RoleIssuer)
		require.EqualError(t, err, types.ErrSafetyBoxHasOtherPermission(types.DefaultCodespace, operator2.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, operator2, types.RoleReturner)
		require.EqualError(t, err, types.ErrSafetyBoxHasOtherPermission(types.DefaultCodespace, operator2.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, allocator, types.RoleIssuer)
		require.EqualError(t, err, types.ErrSafetyBoxHasOtherPermission(types.DefaultCodespace, allocator.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, allocator, types.RoleReturner)
		require.EqualError(t, err, types.ErrSafetyBoxHasOtherPermission(types.DefaultCodespace, allocator.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, owner, allocator, types.RoleOperator)
		require.EqualError(t, err, types.ErrSafetyBoxHasOtherPermission(types.DefaultCodespace, allocator.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, issuer1, types.RoleAllocator)
		require.EqualError(t, err, types.ErrSafetyBoxHasOtherPermission(types.DefaultCodespace, issuer1.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, issuer1, types.RoleReturner)
		require.EqualError(t, err, types.ErrSafetyBoxHasOtherPermission(types.DefaultCodespace, issuer1.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, owner, issuer1, types.RoleOperator)
		require.EqualError(t, err, types.ErrSafetyBoxHasOtherPermission(types.DefaultCodespace, issuer1.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, returner, types.RoleAllocator)
		require.EqualError(t, err, types.ErrSafetyBoxHasOtherPermission(types.DefaultCodespace, returner.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, operator, returner, types.RoleIssuer)
		require.EqualError(t, err, types.ErrSafetyBoxHasOtherPermission(types.DefaultCodespace, returner.String()).Error())

		err = keeper.GrantPermission(ctx, safetyBoxID, owner, returner, types.RoleOperator)
		require.EqualError(t, err, types.ErrSafetyBoxHasOtherPermission(types.DefaultCodespace, returner.String()).Error())
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
