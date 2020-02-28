package keeper

import (
	"testing"

	linktype "github.com/line/link/types"
	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"
)

const nonExistentID = "1234abcd"

func TestModifyCollection(t *testing.T) {
	const (
		modifiedName = "modifiedName"
		modifiedURI  = "modifiedURI"
	)
	nameChange := linktype.NewChange("name", modifiedName)
	imgURIChange := linktype.NewChange("base_img_uri", modifiedURI)
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// Given collection and permission
	collection, err := keeper.GetCollection(ctx, defaultContractID)
	require.NoError(t, err)
	modifyPermission := types.NewModifyPermission(collection.GetContractID())
	keeper.AddPermission(ctx, addr1, modifyPermission)

	t.Logf("Test to modify name of collection to %s", modifiedName)
	{
		// When modify collection name
		require.NoError(t, keeper.modifyCollection(ctx, addr1, defaultContractID, nameChange))

		// Then collection name is modified
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.CollectionKey(collection.GetContractID()))
		actual := keeper.mustDecodeCollection(bz)
		require.Equal(t, modifiedName, actual.GetName())
	}
	t.Logf("Test to modify img uri of collection to %s", modifiedURI)
	{
		// When modify img uri
		require.NoError(t, keeper.modifyCollection(ctx, addr1, defaultContractID, imgURIChange))

		// Then img uri is modified
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.CollectionKey(collection.GetContractID()))
		actual := keeper.mustDecodeCollection(bz)
		require.Equal(t, modifiedURI, actual.GetBaseImgURI())
	}
	t.Log("Test with nonexistent contract")
	{
		// Given nonexistent contract, When modify collection name with invalid contract, Then error is occurred
		require.EqualError(t, keeper.modifyCollection(ctx, addr1, nonExistentID, nameChange),
			types.ErrCollectionNotExist(types.DefaultCodespace, nonExistentID).Error())
	}
	t.Log("Test without permission")
	{
		// Given user does not have permission
		invalidUser := addr2

		// When modify collection name with invalid permission, Then error is occurred
		require.EqualError(t, keeper.modifyCollection(ctx, invalidUser, collection.GetContractID(), nameChange),
			types.ErrTokenNoPermission(types.DefaultCodespace, invalidUser, modifyPermission).Error())
	}
}

func TestModifyTokenType(t *testing.T) {
	const modifiedName = "modifiedName"
	const modifiedURI = "modifiedURI"
	nameChange := linktype.NewChange("name", modifiedName)
	imgURIChange := linktype.NewChange("base_img_uri", modifiedURI)
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// Given collection and permission
	collection, err := keeper.GetCollection(ctx, defaultContractID)
	require.NoError(t, err)
	modifyPermission := types.NewModifyPermission(collection.GetContractID())
	keeper.AddPermission(ctx, addr1, modifyPermission)

	t.Logf("Test to modify name of token type to %s", modifiedName)
	{
		// When modify token type name
		require.NoError(t, keeper.modifyTokenType(ctx, addr1, defaultContractID, defaultTokenType, nameChange))

		// Then collection name is modified
		actual, err := keeper.GetTokenType(ctx, defaultContractID, defaultTokenType)
		require.NoError(t, err)
		require.Equal(t, modifiedName, actual.GetName())
	}
	t.Log("Test to modify img uri of token type")
	{
		require.EqualError(t, keeper.modifyTokenType(ctx, addr1, defaultContractID, defaultTokenType, imgURIChange),
			types.ErrInvalidChangesField(types.DefaultCodespace, imgURIChange.Field).Error())
	}
	t.Log("Test with nonexistent contract")
	{
		// Given nonexistent token type, When modify token type name with invalid contract, Then error is occurred
		require.EqualError(t, keeper.modifyTokenType(ctx, addr1, defaultContractID, nonExistentID, nameChange),
			types.ErrTokenTypeNotExist(types.DefaultCodespace, defaultContractID, nonExistentID).Error())
	}
	t.Log("Test without permission")
	{
		// Given user does not have permission
		invalidUser := addr2

		// When modify token type name with invalid permission, Then error is occurred
		require.EqualError(t, keeper.modifyTokenType(ctx, invalidUser, defaultContractID, defaultTokenType, nameChange),
			types.ErrTokenNoPermission(types.DefaultCodespace, invalidUser, modifyPermission).Error())
	}
}

func TestModifyToken(t *testing.T) {
	const modifiedName = "modifiedName"
	const modifiedURI = "modifiedURI"
	nameChange := linktype.NewChange("name", modifiedName)
	imgURIChange := linktype.NewChange("base_img_uri", modifiedURI)
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// Given collection and permission
	collection, err := keeper.GetCollection(ctx, defaultContractID)
	require.NoError(t, err)
	modifyPermission := types.NewModifyPermission(collection.GetContractID())
	keeper.AddPermission(ctx, addr1, modifyPermission)
	// And token
	token, err := keeper.GetToken(ctx, defaultContractID, defaultTokenID1)
	require.NoError(t, err)

	t.Logf("Test to modify name of token to %s", modifiedName)
	{
		// When modify token name
		require.NoError(t, keeper.modifyToken(ctx, addr1, defaultContractID, token.GetTokenID(), nameChange))

		// Then token name is modified
		actual, err := keeper.GetToken(ctx, defaultContractID, token.GetTokenID())
		require.NoError(t, err)
		require.Equal(t, modifiedName, actual.GetName())
	}
	t.Log("Test to modify img uri")
	{
		require.EqualError(t, keeper.modifyToken(ctx, addr1, defaultContractID, token.GetTokenID(), imgURIChange),
			types.ErrInvalidChangesField(types.DefaultCodespace, imgURIChange.Field).Error())
	}
	t.Log("Test with nonexistent contract")
	{
		// Given nonexistent token id, When modify token name with invalid contract, Then error is occurred
		require.EqualError(t, keeper.modifyToken(ctx, addr1, defaultContractID, nonExistentID, nameChange),
			types.ErrTokenNotExist(types.DefaultCodespace, token.GetContractID(), nonExistentID).Error())
	}
	t.Log("Test without permission")
	{
		// Given user does not have permission
		invalidUser := addr2

		// When modify token name with invalid permission, Then error is occurred
		require.EqualError(t, keeper.modifyToken(ctx, invalidUser, defaultContractID, token.GetTokenID(),
			nameChange), types.ErrTokenNoPermission(types.DefaultCodespace, invalidUser, modifyPermission).Error())
	}
}

func TestModify(t *testing.T) {
	const modifiedName = "modifiedName"
	nameChange := linktype.NewChange("name", modifiedName)
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// Given permission
	modifyPermission := types.NewModifyPermission(defaultContractID)
	keeper.AddPermission(ctx, addr1, modifyPermission)

	t.Logf("Test to modify name of collection to %s", modifiedName)
	{
		// Given contract of collection, When modify collection name
		require.NoError(t, keeper.Modify(ctx, addr1, defaultContractID, "", "", nameChange))

		// Then collection name is modified
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.CollectionKey(defaultContractID))
		actual := keeper.mustDecodeCollection(bz)
		require.Equal(t, modifiedName, actual.GetName())
	}
	t.Logf("Test to modify name of token type to %s", modifiedName)
	{
		// Given contract of token type, When modify token type name
		require.NoError(t, keeper.Modify(ctx, addr1, defaultContractID, defaultTokenType, "", nameChange))

		// Then token type name is modified
		actual, err := keeper.GetTokenType(ctx, defaultContractID, defaultTokenType)
		require.NoError(t, err)
		require.Equal(t, modifiedName, actual.GetName())
	}
	t.Logf("Test to modify name of token to %s", modifiedName)
	{
		// Given contract of token, When modify token name
		require.NoError(t, keeper.Modify(ctx, addr1, defaultContractID, defaultTokenType, defaultTokenIndex, nameChange))

		// Then token name is modified
		actual, err := keeper.GetToken(ctx, defaultContractID, defaultTokenID1)
		require.NoError(t, err)
		require.Equal(t, modifiedName, actual.GetName())
	}
	t.Log("Test with only token index not token type")
	{
		// When modify token name
		require.EqualError(t, keeper.Modify(ctx, addr1, defaultContractID, "", defaultTokenIndex, nameChange),
			types.ErrTokenIndexWithoutType(types.DefaultCodespace).Error())
	}
}
