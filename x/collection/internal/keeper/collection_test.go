package keeper

import (
	"context"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/collection/internal/types"
	"github.com/line/link-modules/x/contract"
	"github.com/stretchr/testify/require"
)

func TestKeeper_NewContractID(t *testing.T) {
	ctx := cacheKeeper()
	contractID := keeper.NewContractID(ctx)
	require.Equal(t, len(contractID), 8)
	contractID2 := keeper.NewContractID(ctx)
	require.NotEqual(t, contractID, contractID2)
}

func TestKeeper_CreateCollection(t *testing.T) {
	ctx := cacheKeeper()
	contractID := keeper.NewContractID(ctx)

	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(contractID, "MyCollection", "meta", "base url"), addr1))
	require.EqualError(t, keeper.CreateCollection(ctx, types.NewCollection(contractID, "MyCollection", "meta", "base url"), addr1), sdkerrors.Wrapf(types.ErrCollectionExist, "ContractID: %s", contractID).Error())
}

func TestKeeper_ExistCollection(t *testing.T) {
	ctx := cacheKeeper()
	contractID := keeper.NewContractID(ctx)
	ctx = ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, contractID))

	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(contractID, "MyCollection", "meta", "base url"), addr1))
	require.True(t, keeper.ExistCollection(ctx))
	ctx = ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, "abcd1234"))
	require.False(t, keeper.ExistCollection(ctx))
}

func TestKeeper_GetCollection(t *testing.T) {
	ctx := cacheKeeper()
	contractID := keeper.NewContractID(ctx)
	ctx = ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, contractID))

	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(contractID, "MyCollection", "meta", "base url"), addr1))
	collection, err := keeper.GetCollection(ctx)
	require.NoError(t, err)
	require.Equal(t, collection.GetContractID(), contractID)
	require.Equal(t, collection.GetName(), "MyCollection")
	require.Equal(t, collection.GetBaseImgURI(), "base url")

	ctx = ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, "abcd1234"))
	_, err = keeper.GetCollection(ctx)
	require.EqualError(t, err, sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: abcd1234").Error())
}

func TestKeeper_SetCollection(t *testing.T) {
	ctx := cacheKeeper()
	contractID := keeper.NewContractID(ctx)

	require.NoError(t, keeper.SetCollection(ctx, types.NewCollection(contractID, "MyCollection", "meta", "base url")))
	require.EqualError(t, keeper.SetCollection(ctx, types.NewCollection(contractID, "MyCollection", "meta", "base url")), sdkerrors.Wrapf(types.ErrCollectionExist, "ContractID: %s", contractID).Error())
}

func TestKeeper_UpdateCollection(t *testing.T) {
	ctx := cacheKeeper()
	contractID := keeper.NewContractID(ctx)
	ctx = ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, contractID))

	require.EqualError(t, keeper.UpdateCollection(ctx, types.NewCollection(contractID, "MyCollection", "meta", "base url")), sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", contractID).Error())
	require.NoError(t, keeper.SetCollection(ctx, types.NewCollection(contractID, "MyCollection", "meta", "base url")))
	require.NoError(t, keeper.UpdateCollection(ctx, types.NewCollection(contractID, "MyCollection2", "meta", "base url2")))

	collection, err := keeper.GetCollection(ctx)
	require.NoError(t, err)
	require.Equal(t, collection.GetContractID(), contractID)
	require.Equal(t, collection.GetName(), "MyCollection2")
	require.Equal(t, collection.GetBaseImgURI(), "base url2")
}

func TestKeeper_GetAllCollections(t *testing.T) {
	ctx := cacheKeeper()
	contractID1 := keeper.NewContractID(ctx)
	contractID2 := keeper.NewContractID(ctx)
	contractID3 := keeper.NewContractID(ctx)

	collections := keeper.GetAllCollections(ctx)
	require.Equal(t, len(collections), 0)

	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(contractID1, "MyCollection1", "meta1", "base url1"), addr1))
	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(contractID2, "MyCollection2", "meta2", "base url2"), addr1))
	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(contractID3, "MyCollection3", "meta3", "base url3"), addr1))

	collections = keeper.GetAllCollections(ctx)
	require.Equal(t, len(collections), 3)
	require.Equal(t, collections[2].GetContractID(), contractID1)
	require.Equal(t, collections[2].GetName(), "MyCollection1")
	require.Equal(t, collections[2].GetMeta(), "meta1")
	require.Equal(t, collections[2].GetBaseImgURI(), "base url1")
	require.Equal(t, collections[1].GetContractID(), contractID2)
	require.Equal(t, collections[1].GetName(), "MyCollection2")
	require.Equal(t, collections[1].GetMeta(), "meta2")
	require.Equal(t, collections[1].GetBaseImgURI(), "base url2")
	require.Equal(t, collections[0].GetContractID(), contractID3)
	require.Equal(t, collections[0].GetName(), "MyCollection3")
	require.Equal(t, collections[0].GetMeta(), "meta3")
	require.Equal(t, collections[0].GetBaseImgURI(), "base url3")
}
