package keeper

import (
	"testing"

	"github.com/line/lbm-sdk/x/collection/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_AddNFTOwner(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Add Owner")
	{
		keeper.AddNFTOwner(ctx, addr1, defaultTokenID1)
	}
	t.Log("Add Owner Again")
	{
		require.Panics(t, func() { keeper.AddNFTOwner(ctx, addr1, defaultTokenID1) }, "")
	}
	t.Log("Get The Data")
	{
		store := ctx.KVStore(keeper.storeKey)
		tokenOwnerKey := types.AccountOwnNFTKey(defaultContractID, addr1, defaultTokenID1)
		require.True(t, store.Has(tokenOwnerKey))
	}
	t.Log("Get The Wrong Data")
	{
		store := ctx.KVStore(keeper.storeKey)
		tokenOwnerKey := types.AccountOwnNFTKey(defaultContractID, addr1, defaultTokenID2)
		require.False(t, store.Has(tokenOwnerKey))
	}
}

func TestKeeper_DeleteNFTOwner(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Owner")
	{
		store := ctx.KVStore(keeper.storeKey)
		tokenOwnerKey := types.AccountOwnNFTKey(defaultContractID, addr1, defaultTokenID1)
		store.Set(tokenOwnerKey, []byte(defaultTokenID1))
	}
	t.Log("Delete the Data")
	{
		keeper.DeleteNFTOwner(ctx, addr1, defaultTokenID1)
	}
	t.Log("Is deleted")
	{
		store := ctx.KVStore(keeper.storeKey)
		tokenOwnerKey := types.AccountOwnNFTKey(defaultContractID, addr1, defaultTokenID1)
		require.False(t, store.Has(tokenOwnerKey))
	}
	t.Log("Delete Wrong Data")
	{
		require.Panics(t, func() { keeper.DeleteNFTOwner(ctx, addr1, defaultTokenID1) }, "")
	}
}

func TestKeeepr_HasNFTOwner(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Has Not Data")
	{
		require.False(t, keeper.HasNFTOwner(ctx, addr1, defaultTokenID1))
	}
	t.Log("Prepare Owner")
	{
		store := ctx.KVStore(keeper.storeKey)
		tokenOwnerKey := types.AccountOwnNFTKey(defaultContractID, addr1, defaultTokenID1)
		store.Set(tokenOwnerKey, []byte(defaultTokenID1))
	}
	t.Log("Has Data")
	{
		require.True(t, keeper.HasNFTOwner(ctx, addr1, defaultTokenID1))
	}
}

func TestKeeper_ChangeNFTOwner(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Owner")
	{
		store := ctx.KVStore(keeper.storeKey)
		tokenOwnerKey := types.AccountOwnNFTKey(defaultContractID, addr1, defaultTokenID1)
		store.Set(tokenOwnerKey, []byte(defaultTokenID1))
	}
	t.Log("transfer")
	{
		require.NoError(t, keeper.ChangeNFTOwner(ctx, addr1, addr2, defaultTokenID1))
	}
	t.Log("transfer again")
	{
		require.EqualError(t, keeper.ChangeNFTOwner(ctx, addr1, addr2, defaultTokenID1), "insufficient token: insufficient account funds[abcdef01]; account has no coin")
	}
}

func TestKeeper_GetNFTsOwner(t *testing.T) {
	ctx := cacheKeeper()
	{
		tokenIDs := keeper.GetNFTsOwner(ctx, addr1)
		require.Empty(t, tokenIDs)
	}
	t.Log("Prepare Owner")
	{
		store := ctx.KVStore(keeper.storeKey)
		tokenOwnerKey := types.AccountOwnNFTKey(defaultContractID, addr1, defaultTokenID1)
		store.Set(tokenOwnerKey, []byte(defaultTokenID1))

		tokenOwnerKey = types.AccountOwnNFTKey(defaultContractID, addr1, defaultTokenID2)
		store.Set(tokenOwnerKey, []byte(defaultTokenID2))

		tokenOwnerKey = types.AccountOwnNFTKey(defaultContractID, addr1, defaultTokenID3)
		store.Set(tokenOwnerKey, []byte(defaultTokenID3))
	}
	t.Log("Get the data")
	{
		tokenIDs := keeper.GetNFTsOwner(ctx, addr1)
		require.NotEmpty(t, tokenIDs)
		require.Equal(t, defaultTokenID1, tokenIDs[0])
		require.Equal(t, defaultTokenID2, tokenIDs[1])
		require.Equal(t, defaultTokenID3, tokenIDs[2])
	}
}
