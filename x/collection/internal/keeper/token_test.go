package keeper

import (
	"strings"
	"testing"

	"github.com/line/link-modules/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper_GetToken(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Token")
	var expected types.Token
	expected = types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(defaultDecimals), true)
	{
		store := ctx.KVStore(keeper.storeKey)
		store.Set(types.TokenKey(defaultContractID, defaultTokenIDFT), keeper.cdc.MustMarshalBinaryBare(expected))
	}
	t.Log("Get Token")
	{
		actual, err := keeper.GetToken(ctx, defaultTokenIDFT)
		require.NoError(t, err)
		verifyTokenFunc(t, expected, actual)
	}
	t.Log("Prepare Token")
	expected = types.NewNFT(defaultContractID, defaultTokenID1, defaultName, defaultMeta, addr1)
	{
		store := ctx.KVStore(keeper.storeKey)
		store.Set(types.TokenKey(defaultContractID, defaultTokenID1), keeper.cdc.MustMarshalBinaryBare(expected))
	}
	t.Log("Get Token")
	{
		actual, err := keeper.GetToken(ctx, defaultTokenID1)
		require.NoError(t, err)
		verifyTokenFunc(t, expected, actual)
	}
}
func TestKeeper_SetToken(t *testing.T) {
	ctx := cacheKeeper()
	var expected types.Token
	t.Log("Set Token")
	expected = types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.SetToken(ctx, expected))
	}
	t.Log("Compare Token")
	{
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenKey(defaultContractID, defaultTokenIDFT))
		actual := keeper.mustDecodeToken(bz)
		verifyTokenFunc(t, expected, actual)
	}
	t.Log("Set Token")
	expected = types.NewNFT(defaultContractID, defaultTokenID1, defaultName, defaultMeta, addr1)
	{
		require.NoError(t, keeper.SetToken(ctx, expected))
	}
	t.Log("Compare Token")
	{
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenKey(defaultContractID, defaultTokenID1))
		actual := keeper.mustDecodeToken(bz)
		verifyTokenFunc(t, expected, actual)
	}
}

func TestKeeper_UpdateToken(t *testing.T) {
	ctx := cacheKeeper()
	var expected, token types.Token
	t.Log("Set Token")
	token = types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.SetToken(ctx, token))
	}
	t.Log("Update Token")
	expected = types.NewFT(defaultContractID, defaultTokenIDFT, "modifiedname", defaultMeta, sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.UpdateToken(ctx, expected))
	}
	t.Log("Compare Token")
	{
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenKey(defaultContractID, defaultTokenIDFT))
		actual := keeper.mustDecodeToken(bz)
		verifyTokenFunc(t, expected, actual)
	}
	t.Log("Set Token")
	token = types.NewNFT(defaultContractID, defaultTokenID1, defaultName, defaultMeta, addr1)
	{
		require.NoError(t, keeper.SetToken(ctx, token))
	}
	t.Log("Update Token")
	expected = types.NewFT(defaultContractID, defaultTokenID1, "modifiedname", defaultMeta, sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.UpdateToken(ctx, expected))
	}
	t.Log("Compare Token")
	{
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenKey(defaultContractID, defaultTokenID1))
		actual := keeper.mustDecodeToken(bz)
		verifyTokenFunc(t, expected, actual)
	}
}

func TestKeeper_GeTokens(t *testing.T) {
	ctx := cacheKeeper()
	var allTokens types.Tokens
	t.Log("Prepare collection")
	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, defaultName, defaultMeta, defaultImgURI), addr1))
	t.Log("Prepare FT Tokens")
	expected := types.Tokens{
		types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(defaultDecimals), true),
		types.NewFT(defaultContractID, defaultTokenIDFT2, defaultName, defaultMeta, sdk.NewInt(defaultDecimals), true),
		types.NewFT(defaultContractID, defaultTokenIDFT3, defaultName, defaultMeta, sdk.NewInt(defaultDecimals), true),
		types.NewFT(defaultContractID, defaultTokenIDFT4, defaultName, defaultMeta, sdk.NewInt(defaultDecimals), true),
	}
	allTokens = append(allTokens, expected...)
	{
		for _, to := range expected {
			require.NoError(t, keeper.IssueFT(ctx, addr1, addr1, to.(types.FT), sdk.NewInt(10)))
		}
	}
	t.Log("Compare FT Tokens")
	{
		actual, err := keeper.GetFTs(ctx)
		require.NoError(t, err)
		for index := range expected {
			verifyTokenFunc(t, expected[index], actual[index])
		}
	}
	t.Log("Prepare NFT Tokens")
	require.NoError(t, keeper.IssueNFT(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta), addr1))
	expected = types.Tokens{
		types.NewNFT(defaultContractID, defaultTokenID1, defaultName, defaultMeta, addr1),
		types.NewNFT(defaultContractID, defaultTokenID2, defaultName, defaultMeta, addr1),
		types.NewNFT(defaultContractID, defaultTokenID3, defaultName, defaultMeta, addr1),
		types.NewNFT(defaultContractID, defaultTokenID4, defaultName, defaultMeta, addr1),
		types.NewNFT(defaultContractID, defaultTokenID5, defaultName, defaultMeta, addr1),
	}
	allTokens = append(allTokens, expected...)
	{
		for _, to := range expected {
			require.NoError(t, keeper.MintNFT(ctx, addr1, to.(types.NFT)))
		}
	}
	t.Log("Compare NFT Tokens")
	{
		actual, err := keeper.GetNFTs(ctx, defaultTokenType)
		require.NoError(t, err)
		for index := range expected {
			verifyTokenFunc(t, expected[index], actual[index])
		}
	}
	t.Log("Compare NFT Tokens Count")
	{
		count, err := keeper.GetNFTCount(ctx, defaultTokenType)
		require.NoError(t, err)
		require.Equal(t, int64(5), count.Int64())
	}

	t.Log("Compare NFT Tokens Count Int")
	{
		count, err := keeper.GetNFTCountInt(ctx, defaultTokenType, types.QueryNFTCount)
		require.NoError(t, err)
		require.Equal(t, int64(5), count.Int64())
	}
	t.Log("Compare NFT Tokens Count Int")
	{
		count, err := keeper.GetNFTCountInt(ctx, defaultTokenType, types.QueryNFTMint)
		require.NoError(t, err)
		require.Equal(t, int64(5), count.Int64())
	}
	t.Log("Compare NFT Tokens Count Int")
	{
		count, err := keeper.GetNFTCountInt(ctx, defaultTokenType, types.QueryNFTBurn)
		require.NoError(t, err)
		require.Equal(t, int64(0), count.Int64())
	}

	t.Log("Compare All Tokens")
	{
		actual, err := keeper.GetTokens(ctx)
		require.NoError(t, err)
		for index := range allTokens {
			verifyTokenFunc(t, allTokens[index], actual[index])
		}
	}
}

func TestKeeper_GetNextTokenIDFT(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare collection")
	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, defaultName, defaultMeta, defaultImgURI), addr1))
	t.Log("Get Next Token ID FT")
	{
		tokenID, err := keeper.GetNextTokenIDFT(ctx)
		require.NoError(t, err)
		require.Equal(t, defaultTokenIDFT, tokenID)
	}
	t.Log("Issue a token and get next token id")
	{
		require.NoError(t, keeper.SetToken(ctx, types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(defaultDecimals), true)))
		tokenID, err := keeper.GetNextTokenIDFT(ctx)
		require.NoError(t, err)
		require.Equal(t, defaultTokenIDFT2, tokenID)
	}
	t.Log("Set Full")
	{
		keeper.setNextTokenTypeFT(ctx, "0fffffff")
		_, err := keeper.GetNextTokenIDFT(ctx)
		require.Error(t, err)
	}
}
func TestKeeper_GetNextTokenIDNFT(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare collection")
	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, defaultName, defaultMeta, defaultImgURI), addr1))
	t.Log("Prepare Token Type")
	expected := types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta)
	{
		require.NoError(t, keeper.SetTokenType(ctx, expected))
	}
	t.Log("Get Next Token ID NFT")
	{
		tokenID, err := keeper.GetNextTokenIDNFT(ctx, defaultTokenType)
		require.NoError(t, err)
		require.Equal(t, defaultTokenID1, tokenID)
	}
	t.Log("Issue a token and get next token id")
	{
		require.NoError(t, keeper.SetToken(ctx, types.NewNFT(defaultContractID, defaultTokenID1, defaultName, defaultMeta, addr1)))
		tokenID, err := keeper.GetNextTokenIDNFT(ctx, defaultTokenType)
		require.NoError(t, err)
		require.Equal(t, defaultTokenID2, tokenID)
	}
	t.Log("Set Full")
	{
		keeper.setNextTokenIndexNFT(ctx, defaultTokenType, "ffffffff")
		_, err := keeper.GetNextTokenIDNFT(ctx, defaultTokenType)
		require.Error(t, err)
	}
}

func TestKeeper_getNFTCountMint(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare collection")
	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, defaultName, defaultMeta, defaultImgURI), addr1))
	t.Log("Prepare Token Type")
	require.NoError(t, keeper.IssueNFT(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta), addr1))

	keeper.setNextTokenIndexNFT(ctx, defaultTokenType, strings.Repeat("f", len(types.ReservedEmpty)))
	require.Equal(t, int64(0), keeper.getNFTCountMint(ctx, defaultTokenType).Int64())
}

func TestNextTokenID(t *testing.T) {
	require.Panics(t, func() { nextID("") })
	require.Equal(t, "b", nextID("a"))
	require.Equal(t, "00", nextID("ff"))
	require.Equal(t, "0001", nextID("0000"))
	require.Equal(t, "000a", nextID("0009"))
	require.Equal(t, "0010", nextID("000f"))
	require.Equal(t, "0000", nextID("ffff"))
	require.Equal(t, "00000000", nextID("ffffffff"))
	require.Equal(t, "abce0000", nextID("abcdffff"))
	require.Equal(t, "abcdabc1", nextID("abcdabc0"))
	require.Equal(t, "abcd0001", nextID("abcd0000"))
	require.Equal(t, "abcd0010", nextID("abcd000f"))
	require.Equal(t, "abcdeef0", nextID("abcdeeef"))
	require.Equal(t, "abcdef00", nextID("abcdeeff"))
	require.Equal(t, "abcd999a", nextID("abcd9999"))
	require.Equal(t, "abcd99a0", nextID("abcd999f"))

	next := "0000"
	for idx := 0; idx < 16*16*16*16; idx++ {
		next = nextID(next)
	}
	require.Equal(t, "0000", next)
}
