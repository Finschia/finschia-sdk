package keeper

import (
	"testing"

	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestApproveDisapproveScenario(t *testing.T) {
	ctx := cacheKeeper()
	const (
		defaultTokenIDFromSymbol2 = defaultTokenType2 + "00000001"
	)

	// prepare collection, FT, NFT
	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultSymbol, "name"), addr1))
	require.NoError(t, keeper.IssueFT(ctx, defaultSymbol, addr1, types.NewFT(defaultSymbol, defaultTokenIDFT, defaultName, defaultTokenURI, sdk.NewInt(defaultDecimals), true), sdk.NewInt(defaultAmount)))
	require.NoError(t, keeper.IssueNFT(ctx, defaultSymbol, types.NewBaseTokenType(defaultSymbol, defaultTokenType, defaultName), addr1))
	require.NoError(t, keeper.IssueNFT(ctx, defaultSymbol, types.NewBaseTokenType(defaultSymbol, defaultTokenType2, defaultName), addr1))
	require.NoError(t, keeper.MintNFT(ctx, defaultSymbol, addr1, types.NewNFT(defaultSymbol, defaultTokenID1, defaultName, defaultTokenURI, addr1)))
	require.NoError(t, keeper.MintNFT(ctx, defaultSymbol, addr1, types.NewNFT(defaultSymbol, defaultTokenType2+"00000001", defaultName, defaultTokenURI, addr1)))

	// approve test
	require.EqualError(t, keeper.SetApproved(ctx, addr3, addr1, defaultSymbol2), types.ErrCollectionNotExist(types.DefaultCodespace, defaultSymbol2).Error())
	require.NoError(t, keeper.SetApproved(ctx, addr3, addr1, defaultSymbol))
	require.EqualError(t, keeper.SetApproved(ctx, addr3, addr1, defaultSymbol), types.ErrCollectionAlreadyApproved(types.DefaultCodespace, addr3.String(), addr1.String(), defaultSymbol).Error())

	// attach_from/detach_from test
	require.EqualError(t, keeper.AttachFrom(ctx, addr2, addr1, defaultSymbol, defaultTokenID1, defaultTokenIDFromSymbol2), types.ErrCollectionNotApproved(types.DefaultCodespace, addr2.String(), addr1.String(), defaultSymbol).Error())
	require.NoError(t, keeper.AttachFrom(ctx, addr3, addr1, defaultSymbol, defaultTokenID1, defaultTokenIDFromSymbol2))
	require.EqualError(t, keeper.DetachFrom(ctx, addr2, addr1, defaultSymbol, defaultTokenIDFromSymbol2), types.ErrCollectionNotApproved(types.DefaultCodespace, addr2.String(), addr1.String(), defaultSymbol).Error())
	require.NoError(t, keeper.DetachFrom(ctx, addr3, addr1, defaultSymbol, defaultTokenIDFromSymbol2))

	// transfer_from test
	require.EqualError(t, keeper.TransferFTFrom(ctx, addr2, addr1, addr2, defaultSymbol, types.NewCoin(defaultTokenIDFT, sdk.NewInt(10))), types.ErrCollectionNotApproved(types.DefaultCodespace, addr2.String(), addr1.String(), defaultSymbol).Error())
	require.NoError(t, keeper.TransferFTFrom(ctx, addr3, addr1, addr2, defaultSymbol, types.NewCoin(defaultTokenIDFT, sdk.NewInt(10))))

	require.EqualError(t, keeper.TransferNFTFrom(ctx, addr2, addr1, addr2, defaultSymbol, defaultTokenID1), types.ErrCollectionNotApproved(types.DefaultCodespace, addr2.String(), addr1.String(), defaultSymbol).Error())
	require.NoError(t, keeper.TransferNFTFrom(ctx, addr3, addr1, addr2, defaultSymbol, defaultTokenID1))

	// disapprove test
	require.EqualError(t, keeper.DeleteApproved(ctx, addr3, addr1, defaultSymbol2), types.ErrCollectionNotExist(types.DefaultCodespace, defaultSymbol2).Error())
	require.NoError(t, keeper.DeleteApproved(ctx, addr3, addr1, defaultSymbol))
	require.EqualError(t, keeper.DeleteApproved(ctx, addr3, addr1, defaultSymbol), types.ErrCollectionNotApproved(types.DefaultCodespace, addr3.String(), addr1.String(), defaultSymbol).Error())
}
