package keeper

import (
	"context"
	"testing"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/contract"
	"github.com/line/lbm-sdk/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestModifyTokenName(t *testing.T) {
	const modifiedTokenName = "modifiedTokenName"
	const modifiedMeta = "modifiedMeta"
	const modifiedImageURI = "modifiedImageURI"
	changes := types.NewChanges(
		types.NewChange("name", modifiedTokenName),
		types.NewChange("meta", modifiedMeta),
		types.NewChange("img_uri", modifiedImageURI),
	)

	ctx := cacheKeeper()
	token := aToken(defaultContractID)
	tokenWithoutPerm := aToken(defaultContractID + "2")
	modifyPermission := types.NewModifyPermission()

	// Given Token And Permission
	require.NoError(t, keeper.SetToken(ctx, token))
	keeper.AddPermission(ctx, addr1, modifyPermission)

	t.Log("Test to modify token")
	{
		// When modify token name
		require.NoError(t, keeper.ModifyToken(ctx, addr1, changes))

		// Then token name is modified
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenKey(token.GetContractID()))
		actual := keeper.mustDecodeToken(bz)
		require.Equal(t, modifiedTokenName, actual.GetName())
	}
	t.Log("Test with nonexistent contract")
	{
		// Given nonexistent contractID
		nonExistentcontractID := "abcd1234"

		// When modify token name with invalid contractID, Then error is occurred
		ctx2 := ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, nonExistentcontractID))
		require.EqualError(t, keeper.ModifyToken(ctx2, addr1, changes),
			sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s", nonExistentcontractID).Error())
	}
	t.Log("Test without permission")
	{
		// Given Token without Permission
		ctx2 := ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, defaultContractID+"2"))
		require.NoError(t, keeper.SetToken(ctx2, tokenWithoutPerm))
		invalidPerm := types.NewModifyPermission()

		// When modify token name with invalid permission, Then error is occurred
		require.EqualError(t, keeper.ModifyToken(ctx2, addr1, changes),
			sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr1.String(), invalidPerm.String()).Error())
	}
}

func aToken(contractID string) types.Token {
	return types.NewToken(contractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
}
