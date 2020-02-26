package keeper

import (
	"testing"

	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper_GetToken(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Token")
	var expected types.Token
	expected = types.NewFT(defaultSymbol, defaultTokenIDFT, defaultName, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
	{
		store := ctx.KVStore(keeper.storeKey)
		store.Set(types.TokenKey(defaultSymbol, defaultTokenIDFT), keeper.cdc.MustMarshalBinaryBare(expected))
	}
	t.Log("Get Token")
	{
		actual, err := keeper.GetToken(ctx, defaultSymbol, defaultTokenIDFT)
		require.NoError(t, err)
		verifyTokenFunc(t, expected, actual)
	}
	t.Log("Prepare Token")
	expected = types.NewNFT(defaultSymbol, defaultTokenID1, defaultName, defaultTokenURI, addr1)
	{
		store := ctx.KVStore(keeper.storeKey)
		store.Set(types.TokenKey(defaultSymbol, defaultTokenID1), keeper.cdc.MustMarshalBinaryBare(expected))
	}
	t.Log("Get Token")
	{
		actual, err := keeper.GetToken(ctx, defaultSymbol, defaultTokenID1)
		require.NoError(t, err)
		verifyTokenFunc(t, expected, actual)
	}
}
func TestKeeper_SetToken(t *testing.T) {
	ctx := cacheKeeper()
	var expected types.Token
	t.Log("Set Token")
	expected = types.NewFT(defaultSymbol, defaultTokenIDFT, defaultName, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.SetToken(ctx, defaultSymbol, expected))
	}
	t.Log("Compare Token")
	{
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenKey(defaultSymbol, defaultTokenIDFT))
		actual := keeper.mustDecodeToken(bz)
		verifyTokenFunc(t, expected, actual)
	}
	t.Log("Set Token")
	expected = types.NewNFT(defaultSymbol, defaultTokenID1, defaultName, defaultTokenURI, addr1)
	{
		require.NoError(t, keeper.SetToken(ctx, defaultSymbol, expected))
	}
	t.Log("Compare Token")
	{
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenKey(defaultSymbol, defaultTokenID1))
		actual := keeper.mustDecodeToken(bz)
		verifyTokenFunc(t, expected, actual)
	}
}

func TestKeeper_UpdateToken(t *testing.T) {
	ctx := cacheKeeper()
	var expected, token types.Token
	t.Log("Set Token")
	token = types.NewFT(defaultSymbol, defaultTokenIDFT, defaultName, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.SetToken(ctx, defaultSymbol, token))
	}
	t.Log("Update Token")
	expected = types.NewFT(defaultSymbol, defaultTokenIDFT, "modifiedname", "modifiedtokenuri", sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.UpdateToken(ctx, defaultSymbol, expected))
	}
	t.Log("Compare Token")
	{
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenKey(defaultSymbol, defaultTokenIDFT))
		actual := keeper.mustDecodeToken(bz)
		verifyTokenFunc(t, expected, actual)
	}
	t.Log("Set Token")
	token = types.NewNFT(defaultSymbol, defaultTokenID1, defaultName, defaultTokenURI, addr1)
	{
		require.NoError(t, keeper.SetToken(ctx, defaultSymbol, token))
	}
	t.Log("Update Token")
	expected = types.NewFT(defaultSymbol, defaultTokenID1, "modifiedname", "modifiedtokenuri", sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.UpdateToken(ctx, defaultSymbol, expected))
	}
	t.Log("Compare Token")
	{
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenKey(defaultSymbol, defaultTokenID1))
		actual := keeper.mustDecodeToken(bz)
		verifyTokenFunc(t, expected, actual)
	}
}

func TestKeeper_GeTokens(t *testing.T) {
	ctx := cacheKeeper()
	var allTokens types.Tokens
	t.Log("Prepare collection")
	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultSymbol, defaultName), addr1))
	t.Log("Prepare FT Tokens")
	expected := types.Tokens{
		types.NewFT(defaultSymbol, defaultTokenIDFT, defaultName, defaultTokenURI, sdk.NewInt(defaultDecimals), true),
		types.NewFT(defaultSymbol, defaultTokenIDFT2, defaultName, defaultTokenURI, sdk.NewInt(defaultDecimals), true),
		types.NewFT(defaultSymbol, defaultTokenIDFT3, defaultName, defaultTokenURI, sdk.NewInt(defaultDecimals), true),
		types.NewFT(defaultSymbol, defaultTokenIDFT4, defaultName, defaultTokenURI, sdk.NewInt(defaultDecimals), true),
	}
	allTokens = append(allTokens, expected...)
	{
		store := ctx.KVStore(keeper.storeKey)
		for _, t := range expected {
			store.Set(types.TokenKey(defaultSymbol, t.GetTokenID()), keeper.cdc.MustMarshalBinaryBare(t))
		}
	}
	t.Log("Compare FT Tokens")
	{
		actual, err := keeper.GetFTs(ctx, defaultSymbol)
		require.NoError(t, err)
		for index := range expected {
			verifyTokenFunc(t, expected[index], actual[index])
		}
	}
	t.Log("Prepare NFT Tokens")
	expected = types.Tokens{
		types.NewNFT(defaultSymbol, defaultTokenID1, defaultName, defaultTokenURI, addr1),
		types.NewNFT(defaultSymbol, defaultTokenID2, defaultName, defaultTokenURI, addr1),
		types.NewNFT(defaultSymbol, defaultTokenID3, defaultName, defaultTokenURI, addr1),
		types.NewNFT(defaultSymbol, defaultTokenID4, defaultName, defaultTokenURI, addr1),
		types.NewNFT(defaultSymbol, defaultTokenID5, defaultName, defaultTokenURI, addr1),
	}
	allTokens = append(allTokens, expected...)
	{
		store := ctx.KVStore(keeper.storeKey)
		for _, t := range expected {
			store.Set(types.TokenKey(defaultSymbol, t.GetTokenID()), keeper.cdc.MustMarshalBinaryBare(t))
		}
	}
	t.Log("Compare NFT Tokens")
	{
		actual, err := keeper.GetNFTs(ctx, defaultSymbol, defaultTokenType)
		require.NoError(t, err)
		for index := range expected {
			verifyTokenFunc(t, expected[index], actual[index])
		}
	}
	t.Log("Compare NFT Tokens Count")
	{
		count, err := keeper.GetNFTCount(ctx, defaultSymbol, defaultTokenType)
		require.NoError(t, err)
		require.Equal(t, int64(5), count.Int64())
	}

	t.Log("Compare All Tokens")
	{
		actual, err := keeper.GetTokens(ctx, defaultSymbol)
		require.NoError(t, err)
		for index := range allTokens {
			verifyTokenFunc(t, allTokens[index], actual[index])
		}
	}
	t.Log("Get Next Token ID FT")
	{
		tokenID, err := keeper.GetNextTokenIDFT(ctx, defaultSymbol)
		require.NoError(t, err)
		require.Equal(t, defaultTokenIDFT5, tokenID)
	}
	t.Log("Get Next Token ID NFT")
	{
		tokenID, err := keeper.GetNextTokenIDNFT(ctx, defaultSymbol, defaultTokenType)
		require.NoError(t, err)
		require.Equal(t, defaultTokenID6, tokenID)
	}
}

func TestKeeper_ModifyTokenURI(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Collection and Token. Add Modify Permission to Addr1")
	{
		store := ctx.KVStore(keeper.storeKey)
		collection := types.NewCollection(defaultSymbol, defaultName)
		store.Set(types.CollectionKey(collection.GetSymbol()), keeper.cdc.MustMarshalBinaryBare(collection))
		token := types.NewFT(defaultSymbol, defaultTokenIDFT, defaultName, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
		err := keeper.SetToken(ctx, defaultSymbol, token)
		require.NoError(t, err)
		keeper.AddPermission(ctx, addr1, types.NewModifyTokenURIPermission(defaultSymbol, defaultTokenIDFT))
	}

	const modifiedTokenURI = "modifiedtokenuri"
	t.Log("Modify Token URI")
	{
		require.NoError(t, keeper.ModifyTokenURI(ctx, addr1, defaultSymbol, defaultTokenIDFT, modifiedTokenURI))
	}
	t.Log("Compare Token")
	{
		actual, err := keeper.GetFT(ctx, defaultSymbol, defaultTokenIDFT)
		require.NoError(t, err)
		require.Equal(t, modifiedTokenURI, actual.GetTokenURI())
	}
}
func TestNextTokenID(t *testing.T) {
	require.Equal(t, "b", nextID("a", ""))
	require.Equal(t, "0001", nextID("0000", ""))
	require.Equal(t, "000a", nextID("0009", ""))
	require.Equal(t, "0010", nextID("000z", ""))
	require.Equal(t, "0000", nextID("zzzz", ""))
	require.Equal(t, "00000000", nextID("zzzzzzzz", ""))
	require.Equal(t, "abce0000", nextID("abcdzzzz", ""))
	require.Equal(t, "abcdabc1", nextID("abcdabc0", ""))

	require.Equal(t, "", nextID("", ""))
	require.Equal(t, "", nextID("", "zzzzz"))
	require.Equal(t, "z0", nextID("zz", "z"))
	require.Equal(t, "item0001", nextID("item0000", "item"))
	require.Equal(t, "item0010", nextID("item000z", "item"))
	require.Equal(t, "itemyyz0", nextID("itemyyyz", "item"))
	require.Equal(t, "itemyz00", nextID("itemyyzz", "item"))
	require.Equal(t, "item999a", nextID("item9999", "item"))
	require.Equal(t, "item99a0", nextID("item999z", "item"))
	require.Equal(t, "z0000000", nextID("zzzzzzzz", "z"))
	require.Equal(t, "zz000000", nextID("zzzzzzzz", "zz"))
	require.Equal(t, "zzzzzzz0", nextID("zzzzzzzz", "zzzzzzz"))
	require.Equal(t, "zzzzzzzz", nextID("zzzzzzzz", "zzzzzzzz"))
	require.Equal(t, "item0000", nextID("itemzzzz", "item"))
	require.Equal(t, "itemz000", nextID("itemyzzz", "item"))
	require.Equal(t, "item0000", nextID("itezzzzz", "item"))

	next := "0000"
	for idx := 0; idx < 36*36*36*36; idx++ {
		next = nextID(next, "")
	}
	require.Equal(t, "0000", next)
}
