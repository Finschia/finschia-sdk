package keeper

import (
	"context"
	"testing"

	"github.com/line/link-modules/x/collection/internal/types"
	"github.com/line/link-modules/x/contract"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestCollectionAndPermission(t *testing.T) {
	ctx := cacheKeeper()

	issuePerm := types.NewIssuePermission()
	{
		require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, defaultName,
			defaultMeta, defaultImgURI), addr1))
		require.True(t, keeper.HasPermission(ctx, addr1, issuePerm))
		require.Error(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, defaultName,
			defaultMeta, defaultImgURI), addr1))
		collection, err := keeper.GetCollection(ctx)
		require.NoError(t, err)
		require.Equal(t, defaultContractID, collection.GetContractID())

		{
			require.NoError(t, keeper.IssueFT(ctx, addr1, addr1, types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(defaultDecimals), true), sdk.NewInt(defaultAmount)))
			token, err := keeper.GetToken(ctx, defaultTokenIDFT)
			require.NoError(t, err)
			require.Equal(t, defaultContractID, token.GetContractID())
			require.Equal(t, defaultTokenIDFT, token.GetTokenID())
		}
		{
			require.NoError(t, keeper.IssueNFT(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta), addr1))
			require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID1, defaultName, defaultMeta, addr1)))

			token, err := keeper.GetToken(ctx, defaultTokenID1)
			require.NoError(t, err)
			require.Equal(t, defaultContractID, token.GetContractID())
			require.Equal(t, defaultTokenID1, token.GetTokenID())

			require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID2, defaultName, defaultMeta, addr1)))
			token, err = keeper.GetToken(ctx, defaultTokenID2)
			require.NoError(t, err)
			require.Equal(t, defaultContractID, token.GetContractID())
			require.Equal(t, defaultTokenID2, token.GetTokenID())

			count, err := keeper.GetNFTCount(ctx, defaultTokenType)
			require.NoError(t, err)
			require.Equal(t, int64(2), count.Int64())

			require.NoError(t, keeper.IssueNFT(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType2, defaultName, defaultMeta), addr1))
			require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenType2+"00000001", defaultName, defaultMeta, addr1)))
			token, err = keeper.GetToken(ctx, defaultTokenType2+"00000001")
			require.NoError(t, err)
			require.Equal(t, defaultContractID, token.GetContractID())
			require.Equal(t, defaultTokenType2+"00000001", token.GetTokenID())
		}
	}
	{
		require.NoError(t, keeper.GrantPermission(ctx, addr1, addr2, issuePerm))
		require.True(t, keeper.HasPermission(ctx, addr1, issuePerm))
		require.True(t, keeper.HasPermission(ctx, addr2, issuePerm))
	}

	issuePerm2 := types.NewIssuePermission()
	ctx = ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, defaultContractID2))
	{
		require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID2, defaultName,
			defaultMeta, defaultImgURI), addr1))
		require.True(t, keeper.HasPermission(ctx, addr1, issuePerm2))
		require.Error(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID2, defaultName,
			defaultMeta, defaultImgURI), addr1))
		collection, err := keeper.GetCollection(ctx)
		require.NoError(t, err)
		require.Equal(t, defaultContractID2, collection.GetContractID())
	}
	{
		collections := keeper.GetAllCollections(ctx)
		require.Equal(t, 2, len(collections))
		require.Equal(t, defaultContractID, collections[0].GetContractID())
		require.Equal(t, defaultContractID2, collections[1].GetContractID())
	}
}

func TestPermission(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	require.EqualError(t, keeper.RevokePermission(ctx, addr3, types.NewMintPermission()), sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr3.String(), types.NewMintPermission().String()).Error())
	require.NoError(t, keeper.RevokePermission(ctx, addr1, types.NewMintPermission()))
	require.EqualError(t, keeper.GrantPermission(ctx, addr3, addr1, types.NewMintPermission()), sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr3.String(), types.NewMintPermission().String()).Error())
}
