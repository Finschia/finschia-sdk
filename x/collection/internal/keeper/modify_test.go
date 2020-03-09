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
		modifiedMeta = "modifiedMeta"
	)
	changes := linktype.NewChanges(
		linktype.NewChange("name", modifiedName),
		linktype.NewChange("base_img_uri", modifiedURI),
		linktype.NewChange("meta", modifiedMeta),
	)
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// Given collection and permission
	collection, err := keeper.GetCollection(ctx, defaultContractID)
	require.NoError(t, err)
	modifyPermission := types.NewModifyPermission(collection.GetContractID())
	keeper.AddPermission(ctx, addr1, modifyPermission)

	t.Log("Test to modify collection")
	{
		// When modify collection
		require.NoError(t, keeper.modifyCollection(ctx, addr1, defaultContractID, changes))

		// Then collection is modified
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.CollectionKey(collection.GetContractID()))
		actual := keeper.mustDecodeCollection(bz)
		require.Equal(t, modifiedName, actual.GetName())
		require.Equal(t, modifiedURI, actual.GetBaseImgURI())
		require.Equal(t, modifiedMeta, actual.GetMeta())
	}
	t.Log("Test with nonexistent contract")
	{
		// Given nonexistent contract, When modify collection name with invalid contract, Then error is occurred
		require.EqualError(t, keeper.modifyCollection(ctx, addr1, nonExistentID, changes),
			types.ErrCollectionNotExist(types.DefaultCodespace, nonExistentID).Error())
	}
	t.Log("Test without permission")
	{
		// Given user does not have permission
		invalidUser := addr2

		// When modify collection name with invalid permission, Then error is occurred
		require.EqualError(t, keeper.modifyCollection(ctx, invalidUser, collection.GetContractID(), changes),
			types.ErrTokenNoPermission(types.DefaultCodespace, invalidUser, modifyPermission).Error())
	}
}

func TestModifyTokenType(t *testing.T) {
	const modifiedName = "modifiedName"
	const modifiedURI = "modifiedURI"
	const modifiedMeta = "modifiedMeta"

	validChanges := linktype.NewChanges(
		linktype.NewChange("name", modifiedName),
		linktype.NewChange("meta", modifiedMeta),
	)
	invalidChanges := linktype.NewChanges(
		linktype.NewChange("base_img_uri", modifiedURI),
	)
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// Given collection and permission
	collection, err := keeper.GetCollection(ctx, defaultContractID)
	require.NoError(t, err)
	modifyPermission := types.NewModifyPermission(collection.GetContractID())
	keeper.AddPermission(ctx, addr1, modifyPermission)

	t.Log("Test to modify token type with valid fields")
	{
		// When modify token type name
		require.NoError(t, keeper.modifyTokenType(ctx, addr1, defaultContractID, defaultTokenType, validChanges))

		// Then collection name is modified
		actual, err := keeper.GetTokenType(ctx, defaultContractID, defaultTokenType)
		require.NoError(t, err)
		require.Equal(t, modifiedName, actual.GetName())
	}
	t.Log("Test to modify token type with invalid fields")
	{
		require.EqualError(t, keeper.modifyTokenType(ctx, addr1, defaultContractID, defaultTokenType, invalidChanges),
			types.ErrInvalidChangesField(types.DefaultCodespace, "base_img_uri").Error())
	}
	t.Log("Test with nonexistent contract")
	{
		// Given nonexistent token type, When modify token type name with invalid contract, Then error is occurred
		require.EqualError(t, keeper.modifyTokenType(ctx, addr1, defaultContractID, nonExistentID, validChanges),
			types.ErrTokenTypeNotExist(types.DefaultCodespace, defaultContractID, nonExistentID).Error())
	}
	t.Log("Test without permission")
	{
		// Given user does not have permission
		invalidUser := addr2

		// When modify token type name with invalid permission, Then error is occurred
		require.EqualError(t, keeper.modifyTokenType(ctx, invalidUser, defaultContractID, defaultTokenType, validChanges),
			types.ErrTokenNoPermission(types.DefaultCodespace, invalidUser, modifyPermission).Error())
	}
}

func TestModifyToken(t *testing.T) {
	const modifiedName = "modifiedName"
	const modifiedURI = "modifiedURI"
	const modifiedMeta = "modifiedMeta"

	validChanges := linktype.NewChanges(
		linktype.NewChange("name", modifiedName),
		linktype.NewChange("meta", modifiedMeta),
	)
	invalidChanges := linktype.NewChanges(
		linktype.NewChange("base_img_uri", modifiedURI),
	)
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

	t.Log("Test to modify token with valid changes")
	{
		// When modify token name
		require.NoError(t, keeper.modifyToken(ctx, addr1, defaultContractID, token.GetTokenID(), validChanges))

		// Then token name is modified
		actual, err := keeper.GetToken(ctx, defaultContractID, token.GetTokenID())
		require.NoError(t, err)
		require.Equal(t, modifiedName, actual.GetName())
	}
	t.Log("Test to modify token with invalid changes")
	{
		require.EqualError(t, keeper.modifyToken(ctx, addr1, defaultContractID, token.GetTokenID(), invalidChanges),
			types.ErrInvalidChangesField(types.DefaultCodespace, "base_img_uri").Error())
	}
	t.Log("Test with nonexistent contract")
	{
		// Given nonexistent token id, When modify token name with invalid contract, Then error is occurred
		require.EqualError(t, keeper.modifyToken(ctx, addr1, defaultContractID, nonExistentID, validChanges),
			types.ErrTokenNotExist(types.DefaultCodespace, token.GetContractID(), nonExistentID).Error())
	}
	t.Log("Test without permission")
	{
		// Given user does not have permission
		invalidUser := addr2

		// When modify token name with invalid permission, Then error is occurred
		require.EqualError(t, keeper.modifyToken(ctx, invalidUser, defaultContractID, token.GetTokenID(), validChanges),
			types.ErrTokenNoPermission(types.DefaultCodespace, invalidUser, modifyPermission).Error())
	}
}

func TestModify(t *testing.T) {
	const modifiedName = "modifiedName"
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)
	changes := linktype.NewChanges(
		linktype.NewChange("name", modifiedName),
	)
	// Given permission
	modifyPermission := types.NewModifyPermission(defaultContractID)
	keeper.AddPermission(ctx, addr1, modifyPermission)

	t.Logf("Test to modify name of collection to %s", modifiedName)
	{
		// When modify collection name
		require.NoError(t, keeper.Modify(ctx, addr1, defaultContractID, "", "", changes))

		// Then collection name is modified
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.CollectionKey(defaultContractID))
		actual := keeper.mustDecodeCollection(bz)
		require.Equal(t, modifiedName, actual.GetName())
	}
	t.Logf("Test to modify name of token type to %s", modifiedName)
	{
		// When modify token type name
		require.NoError(t, keeper.Modify(ctx, addr1, defaultContractID, defaultTokenType, "", changes))

		// Then token type name is modified
		actual, err := keeper.GetTokenType(ctx, defaultContractID, defaultTokenType)
		require.NoError(t, err)
		require.Equal(t, modifiedName, actual.GetName())
	}
	t.Logf("Test to modify name of token to %s", modifiedName)
	{
		// When modify token name
		require.NoError(t, keeper.Modify(ctx, addr1, defaultContractID, defaultTokenType, defaultTokenIndex, changes))

		// Then token name is modified
		actual, err := keeper.GetToken(ctx, defaultContractID, defaultTokenID1)
		require.NoError(t, err)
		require.Equal(t, modifiedName, actual.GetName())
	}
	t.Log("Test with only token index not token type")
	{
		// When modify token name, Then error is occurred
		require.EqualError(t, keeper.Modify(ctx, addr1, defaultContractID, "", defaultTokenIndex, changes), types.ErrTokenIndexWithoutType(types.DefaultCodespace).Error())
	}
}
