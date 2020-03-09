package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestModifyTokenName(t *testing.T) {
	const modifiedTokenName = "modifiedTokenName"
	const modifiedMeta = "modifiedMeta"
	const modifiedImageURI = "modifiedImageURI"
	changes := linktype.NewChanges(
		linktype.NewChange("name", modifiedTokenName),
		linktype.NewChange("meta", modifiedMeta),
		linktype.NewChange("img_uri", modifiedImageURI),
	)

	ctx := cacheKeeper()
	token := aToken(defaultContractID)
	tokenWithoutPerm := aToken(defaultContractID + "2")
	modifyPermission := types.NewModifyPermission(token.GetContractID())

	// Given Token And Permission
	require.NoError(t, keeper.SetToken(ctx, token))
	keeper.AddPermission(ctx, addr1, modifyPermission)

	t.Log("Test to modify token")
	{
		// When modify token name
		require.NoError(t, keeper.ModifyToken(ctx, addr1, token.GetContractID(), changes))

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
		require.EqualError(t, keeper.ModifyToken(ctx, addr1, nonExistentcontractID, changes),
			types.ErrTokenNotExist(types.DefaultCodespace, nonExistentcontractID).Error())
	}
	t.Log("Test without permission")
	{
		// Given Token without Permission
		require.NoError(t, keeper.SetToken(ctx, tokenWithoutPerm))
		invalidPerm := types.NewModifyPermission(tokenWithoutPerm.GetContractID())

		// When modify token name with invalid permission, Then error is occurred
		require.EqualError(t, keeper.ModifyToken(ctx, addr1, tokenWithoutPerm.GetContractID(), changes),
			types.ErrTokenNoPermission(types.DefaultCodespace, addr1, invalidPerm).Error())
	}
}

func aToken(contractID string) types.Token {
	return types.NewToken(contractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
}
