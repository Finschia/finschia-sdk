package keeper

import (
	"testing"

	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestCollectionAndPermission(t *testing.T) {
	ctx := cacheKeeper()

	issuePerm := types.NewIssuePermission(defaultSymbol)
	{
		require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultSymbol, defaultName), addr1))
		require.True(t, keeper.HasPermission(ctx, addr1, issuePerm))
		require.Error(t, keeper.CreateCollection(ctx, types.NewCollection(defaultSymbol, defaultName), addr1))
		collection, err := keeper.GetCollection(ctx, defaultSymbol)
		require.NoError(t, err)
		require.Equal(t, defaultSymbol, collection.GetSymbol())

		{
			require.NoError(t, keeper.IssueFT(ctx, defaultSymbol, addr1, types.NewFT(defaultSymbol, defaultTokenIDFT, defaultName, defaultTokenURI, sdk.NewInt(defaultDecimals), true), sdk.NewInt(defaultAmount)))
			token, err := keeper.GetToken(ctx, defaultSymbol, defaultTokenIDFT)
			require.NoError(t, err)
			require.Equal(t, defaultSymbol, token.GetSymbol())
			require.Equal(t, defaultTokenIDFT, token.GetTokenID())
		}
		{
			require.NoError(t, keeper.IssueNFT(ctx, defaultSymbol, types.NewBaseTokenType(defaultSymbol, defaultTokenType, defaultName), addr1))
			require.NoError(t, keeper.MintNFT(ctx, defaultSymbol, addr1, types.NewNFT(defaultSymbol, defaultTokenID1, defaultName, defaultTokenURI, addr1)))

			token, err := keeper.GetToken(ctx, defaultSymbol, defaultTokenID1)
			require.NoError(t, err)
			require.Equal(t, defaultSymbol, token.GetSymbol())
			require.Equal(t, defaultTokenID1, token.GetTokenID())

			require.NoError(t, keeper.MintNFT(ctx, defaultSymbol, addr1, types.NewNFT(defaultSymbol, defaultTokenID2, defaultName, defaultTokenURI, addr1)))
			token, err = keeper.GetToken(ctx, defaultSymbol, defaultTokenID2)
			require.NoError(t, err)
			require.Equal(t, defaultSymbol, token.GetSymbol())
			require.Equal(t, defaultTokenID2, token.GetTokenID())

			count, err := keeper.GetNFTCount(ctx, defaultSymbol, defaultTokenType)
			require.NoError(t, err)
			require.Equal(t, int64(2), count.Int64())

			require.NoError(t, keeper.IssueNFT(ctx, defaultSymbol, types.NewBaseTokenType(defaultSymbol, defaultTokenType2, defaultName), addr1))
			require.NoError(t, keeper.MintNFT(ctx, defaultSymbol, addr1, types.NewNFT(defaultSymbol, defaultTokenType2+"00000001", defaultName, defaultTokenURI, addr1)))
			token, err = keeper.GetToken(ctx, defaultSymbol, defaultTokenType2+"00000001")
			require.NoError(t, err)
			require.Equal(t, defaultSymbol, token.GetSymbol())
			require.Equal(t, defaultTokenType2+"00000001", token.GetTokenID())
		}
	}
	{
		require.NoError(t, keeper.GrantPermission(ctx, addr1, addr2, issuePerm))
		require.True(t, keeper.HasPermission(ctx, addr1, issuePerm))
		require.True(t, keeper.HasPermission(ctx, addr2, issuePerm))
	}

	issuePerm2 := types.NewIssuePermission(defaultSymbol2)
	{
		require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultSymbol2, defaultName), addr1))
		require.True(t, keeper.HasPermission(ctx, addr1, issuePerm2))
		require.Error(t, keeper.CreateCollection(ctx, types.NewCollection(defaultSymbol2, defaultName), addr1))
		collection, err := keeper.GetCollection(ctx, defaultSymbol2)
		require.NoError(t, err)
		require.Equal(t, defaultSymbol2, collection.GetSymbol())
	}
	{
		collections := keeper.GetAllCollections(ctx)
		require.Equal(t, 2, len(collections))
		require.Equal(t, defaultSymbol, collections[0].GetSymbol())
		require.Equal(t, defaultSymbol2, collections[1].GetSymbol())
	}
}

func TestPermission(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	require.EqualError(t, keeper.RevokePermission(ctx, addr3, types.NewMintPermission(defaultSymbol)), types.ErrTokenNoPermission(types.DefaultCodespace, addr3, types.NewMintPermission(defaultSymbol)).Error())
	require.NoError(t, keeper.RevokePermission(ctx, addr1, types.NewMintPermission(defaultSymbol)))
	require.EqualError(t, keeper.GrantPermission(ctx, addr3, addr1, types.NewMintPermission(defaultSymbol)), types.ErrTokenNoPermission(types.DefaultCodespace, addr3, types.NewMintPermission(defaultSymbol)).Error())
}
