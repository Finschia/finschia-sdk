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
		defaultTokenIDFromSymbol2 = defaultTokenType2 + "0001"
	)

	// prepare collection, CFT, CNFT
	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultSymbol, "name"), addr1))
	collection, err := keeper.GetCollection(ctx, defaultSymbol)
	require.NoError(t, err)
	err = keeper.IssueCFT(ctx, addr1, types.NewFT(collection, defaultName, defaultTokenURI, sdk.NewInt(defaultDecimals), true), sdk.NewInt(defaultAmount))
	require.NoError(t, err)
	err = keeper.IssueCNFT(ctx, addr1, defaultSymbol)
	require.NoError(t, err)
	err = keeper.IssueCNFT(ctx, addr1, defaultSymbol)
	require.NoError(t, err)
	require.NoError(t, keeper.MintCNFT(ctx, addr1, types.NewNFT(collection, defaultName, defaultTokenType, defaultTokenURI, addr1)))
	require.NoError(t, keeper.MintCNFT(ctx, addr1, types.NewNFT(collection, defaultName, defaultTokenType2, defaultTokenURI, addr1)))

	// approve test
	require.EqualError(t, keeper.SetApproved(ctx, addr3, addr1, defaultSymbol2), types.ErrCollectionNotExist(types.DefaultCodespace, defaultSymbol2).Error())
	require.NoError(t, keeper.SetApproved(ctx, addr3, addr1, defaultSymbol))
	require.EqualError(t, keeper.SetApproved(ctx, addr3, addr1, defaultSymbol), types.ErrCollectionAlreadyApproved(types.DefaultCodespace, addr3.String(), addr1.String(), defaultSymbol).Error())

	// attach_from/detach_from test
	require.EqualError(t, keeper.AttachFrom(ctx, addr2, addr1, defaultSymbol, defaultTokenID1, defaultTokenIDFromSymbol2), types.ErrCollectionNotApproved(types.DefaultCodespace, addr2.String(), addr1.String(), defaultSymbol).Error())
	require.NoError(t, keeper.AttachFrom(ctx, addr3, addr1, defaultSymbol, defaultTokenID1, defaultTokenIDFromSymbol2))
	require.EqualError(t, keeper.DetachFrom(ctx, addr2, addr1, addr2, defaultSymbol, defaultTokenIDFromSymbol2), types.ErrCollectionNotApproved(types.DefaultCodespace, addr2.String(), addr1.String(), defaultSymbol).Error())
	require.NoError(t, keeper.DetachFrom(ctx, addr3, addr1, addr1, defaultSymbol, defaultTokenIDFromSymbol2))

	// transfer_from test
	require.EqualError(t, keeper.TransferCFTFrom(ctx, addr2, addr1, addr2, defaultSymbol, defaultTokenIDFT, sdk.NewInt(10)), types.ErrCollectionNotApproved(types.DefaultCodespace, addr2.String(), addr1.String(), defaultSymbol).Error())
	require.NoError(t, keeper.TransferCFTFrom(ctx, addr3, addr1, addr2, defaultSymbol, defaultTokenIDFT, sdk.NewInt(10)))

	require.EqualError(t, keeper.TransferCNFTFrom(ctx, addr2, addr1, addr2, defaultSymbol, defaultTokenID1), types.ErrCollectionNotApproved(types.DefaultCodespace, addr2.String(), addr1.String(), defaultSymbol).Error())
	require.NoError(t, keeper.TransferCNFTFrom(ctx, addr3, addr1, addr2, defaultSymbol, defaultTokenID1))

	// disapprove test
	require.EqualError(t, keeper.DeleteApproved(ctx, addr3, addr1, defaultSymbol2), types.ErrCollectionNotExist(types.DefaultCodespace, defaultSymbol2).Error())
	require.NoError(t, keeper.DeleteApproved(ctx, addr3, addr1, defaultSymbol))
	require.EqualError(t, keeper.DeleteApproved(ctx, addr3, addr1, defaultSymbol), types.ErrCollectionNotApproved(types.DefaultCodespace, addr3.String(), addr1.String(), defaultSymbol).Error())
}
