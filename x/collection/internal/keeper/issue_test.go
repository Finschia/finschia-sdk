package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_IssueFT(t *testing.T) {
	ctx := cacheKeeper()

	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, "name", defaultMeta, defaultImgURI), addr1))

	require.EqualError(t, keeper.IssueFT(ctx, addr1, addr1, types.NewFT(wrongContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(1), true), sdk.NewInt(defaultAmount)), types.ErrCollectionNotExist(types.DefaultCodespace, wrongContractID).Error())
	require.NoError(t, keeper.IssueFT(ctx, addr1, addr1, types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(1), true), sdk.NewInt(defaultAmount)))
	require.EqualError(t, keeper.IssueFT(ctx, addr1, addr1, types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(1), true), sdk.NewInt(defaultAmount)), types.ErrTokenExist(types.DefaultCodespace, defaultContractID, defaultTokenIDFT).Error())
}

func TestKeeper_IssueNFT(t *testing.T) {
	ctx := cacheKeeper()

	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, "name", defaultMeta, defaultImgURI), addr1))

	require.EqualError(t, keeper.IssueNFT(ctx, types.NewBaseTokenType(wrongContractID, defaultTokenType, defaultName, defaultMeta), addr1), types.ErrCollectionNotExist(types.DefaultCodespace, wrongContractID).Error())
	require.NoError(t, keeper.IssueNFT(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta), addr1))
	require.EqualError(t, keeper.IssueNFT(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta), addr1), types.ErrTokenTypeExist(types.DefaultCodespace, defaultContractID, defaultTokenType).Error())
}
