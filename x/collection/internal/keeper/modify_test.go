package keeper

import (
	"context"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/collection/internal/types"
	"github.com/line/link-modules/x/contract"
	"github.com/stretchr/testify/require"
)

const nonExistentID = "1234abcd"

func TestModifyCollection(t *testing.T) {
	const (
		modifiedName = "modifiedName"
		modifiedURI  = "modifiedURI"
		modifiedMeta = "modifiedMeta"
	)
	changes := types.NewChanges(
		types.NewChange("name", modifiedName),
		types.NewChange("base_img_uri", modifiedURI),
		types.NewChange("meta", modifiedMeta),
	)
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// Given collection and permission
	collection, err := keeper.GetCollection(ctx)
	require.NoError(t, err)
	modifyPermission := types.NewModifyPermission()
	keeper.AddPermission(ctx, addr1, modifyPermission)

	t.Log("Test to modify collection")
	{
		// When modify collection
		require.NoError(t, keeper.modifyCollection(ctx, addr1, changes))

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
		ctx2 := ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, nonExistentID))
		require.EqualError(t, keeper.modifyCollection(ctx2, addr1, changes),
			sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", nonExistentID).Error())
	}
	t.Log("Test without permission")
	{
		// Given user does not have permission
		invalidUser := addr2

		// When modify collection name with invalid permission, Then error is occurred
		require.EqualError(t, keeper.modifyCollection(ctx, invalidUser, changes),
			sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", invalidUser.String(), modifyPermission.String()).Error())
	}
}

func TestModifyTokenType(t *testing.T) {
	const modifiedName = "modifiedName"
	const modifiedURI = "modifiedURI"
	const modifiedMeta = "modifiedMeta"

	validChanges := types.NewChanges(
		types.NewChange("name", modifiedName),
		types.NewChange("meta", modifiedMeta),
	)
	invalidChanges := types.NewChanges(
		types.NewChange("base_img_uri", modifiedURI),
	)
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// Given collection and permission
	modifyPermission := types.NewModifyPermission()
	keeper.AddPermission(ctx, addr1, modifyPermission)

	t.Log("Test to modify token type with valid fields")
	{
		// When modify token type name
		require.NoError(t, keeper.modifyTokenType(ctx, addr1, defaultTokenType, validChanges))

		// Then collection name is modified
		actual, err := keeper.GetTokenType(ctx, defaultTokenType)
		require.NoError(t, err)
		require.Equal(t, modifiedName, actual.GetName())
	}
	t.Log("Test to modify token type with invalid fields")
	{
		require.EqualError(t, keeper.modifyTokenType(ctx, addr1, defaultTokenType, invalidChanges),
			sdkerrors.Wrap(types.ErrInvalidChangesField, "Field: base_img_uri").Error())
	}
	t.Log("Test with nonexistent contract")
	{
		// Given nonexistent token type, When modify token type name with invalid contract, Then error is occurred
		require.EqualError(t, keeper.modifyTokenType(ctx, addr1, nonExistentID, validChanges),
			sdkerrors.Wrapf(types.ErrTokenTypeNotExist, "ContractID: %s, TokenType: %s", defaultContractID, nonExistentID).Error())
	}
	t.Log("Test without permission")
	{
		// Given user does not have permission
		invalidUser := addr2

		// When modify token type name with invalid permission, Then error is occurred
		require.EqualError(t, keeper.modifyTokenType(ctx, invalidUser, defaultTokenType, validChanges),
			sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", invalidUser.String(), modifyPermission.String()).Error())
	}
}

func TestModifyToken(t *testing.T) {
	const modifiedName = "modifiedName"
	const modifiedURI = "modifiedURI"
	const modifiedMeta = "modifiedMeta"

	validChanges := types.NewChanges(
		types.NewChange("name", modifiedName),
		types.NewChange("meta", modifiedMeta),
	)
	invalidChanges := types.NewChanges(
		types.NewChange("base_img_uri", modifiedURI),
	)
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// Given collection and permission
	modifyPermission := types.NewModifyPermission()
	keeper.AddPermission(ctx, addr1, modifyPermission)
	// And token
	token, err := keeper.GetToken(ctx, defaultTokenID1)
	require.NoError(t, err)

	t.Log("Test to modify token with valid changes")
	{
		// When modify token name
		require.NoError(t, keeper.modifyToken(ctx, addr1, token.GetTokenID(), validChanges))

		// Then token name is modified
		actual, err := keeper.GetToken(ctx, token.GetTokenID())
		require.NoError(t, err)
		require.Equal(t, modifiedName, actual.GetName())
	}
	t.Log("Test to modify token with invalid changes")
	{
		require.EqualError(t, keeper.modifyToken(ctx, addr1, token.GetTokenID(), invalidChanges),
			sdkerrors.Wrap(types.ErrInvalidChangesField, "Field: base_img_uri").Error())
	}
	t.Log("Test with nonexistent contract")
	{
		// Given nonexistent token id, When modify token name with invalid contract, Then error is occurred
		require.EqualError(t, keeper.modifyToken(ctx, addr1, nonExistentID, validChanges),
			sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", token.GetContractID(), nonExistentID).Error())
	}
	t.Log("Test without permission")
	{
		// Given user does not have permission
		invalidUser := addr2

		// When modify token name with invalid permission, Then error is occurred
		require.EqualError(t, keeper.modifyToken(ctx, invalidUser, token.GetTokenID(), validChanges),
			sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", invalidUser.String(), modifyPermission.String()).Error())
	}
}

func TestModify(t *testing.T) {
	const modifiedName = "modifiedName"
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)
	changes := types.NewChanges(
		types.NewChange("name", modifiedName),
	)
	// Given permission
	modifyPermission := types.NewModifyPermission()
	keeper.AddPermission(ctx, addr1, modifyPermission)

	t.Logf("Test to modify name of collection to %s", modifiedName)
	{
		// When modify collection name
		require.NoError(t, keeper.Modify(ctx, addr1, "", "", changes))

		// Then collection name is modified
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.CollectionKey(defaultContractID))
		actual := keeper.mustDecodeCollection(bz)
		require.Equal(t, modifiedName, actual.GetName())
	}
	t.Logf("Test to modify name of token type to %s", modifiedName)
	{
		// When modify token type name
		require.NoError(t, keeper.Modify(ctx, addr1, defaultTokenType, "", changes))

		// Then token type name is modified
		actual, err := keeper.GetTokenType(ctx, defaultTokenType)
		require.NoError(t, err)
		require.Equal(t, modifiedName, actual.GetName())
	}
	t.Logf("Test to modify name of token to %s", modifiedName)
	{
		// When modify token name
		require.NoError(t, keeper.Modify(ctx, addr1, defaultTokenType, defaultTokenIndex, changes))

		// Then token name is modified
		actual, err := keeper.GetToken(ctx, defaultTokenID1)
		require.NoError(t, err)
		require.Equal(t, modifiedName, actual.GetName())
	}
	t.Log("Test with only token index not token type")
	{
		// When modify token name, Then error is occurred
		require.EqualError(t, keeper.Modify(ctx, addr1, "", defaultTokenIndex, changes), types.ErrTokenIndexWithoutType.Error())
	}
}
