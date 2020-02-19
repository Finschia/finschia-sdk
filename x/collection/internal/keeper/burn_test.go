package keeper

import (
	"fmt"
	"testing"

	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
)

func TestKeeper_BurnCFT(t *testing.T) {
	t.Log("implement me - ", t.Name())
}

func TestKeeper_BurnCFTFrom(t *testing.T) {
	t.Log("implement me - ", t.Name())
}

func TestKeeper_BurnCNFT(t *testing.T) {
	t.Log("implement me - ", t.Name())
}

func TestKeeper_BurnCNFTFrom(t *testing.T) {
	t.Log("implement me - ", t.Name())
}

func TestMintBurn(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	const (
		wrongTokenID = "12345678"
	)
	require.EqualError(t, keeper.MintCNFT(ctx, addr1, types.NewNFT(types.NewCollection(defaultSymbol2, defaultName), defaultName, wrongTokenID, defaultTokenURI, addr1)), types.ErrCollectionTokenTypeNotExist(types.DefaultCodespace, defaultSymbol2, wrongTokenID[:4]).Error())
	require.EqualError(t, keeper.MintCNFT(ctx, addr3, types.NewNFT(types.NewCollection(defaultSymbol, defaultName), defaultName, defaultTokenID1, defaultTokenURI, addr1)), types.ErrTokenNoPermission(types.DefaultCodespace, addr3, types.NewMintPermission(defaultSymbol+defaultTokenID1[:4])).Error())

	require.NoError(t, keeper.BurnCFT(ctx, addr1, linktype.NewCoinWithTokenIDs(linktype.NewCoinWithTokenID(defaultSymbol, defaultTokenIDFT, sdk.NewInt(1)))))
	require.EqualError(t, keeper.BurnCNFT(ctx, addr1, defaultSymbol, wrongTokenID), types.ErrCollectionTokenNotExist(types.DefaultCodespace, defaultSymbol, wrongTokenID).Error())
	require.EqualError(t, keeper.BurnCNFT(ctx, addr3, defaultSymbol, defaultTokenID1), types.ErrTokenNoPermission(types.DefaultCodespace, addr3, types.NewBurnPermission(defaultSymbol+defaultTokenID1[:4])).Error())
}

func TestBurnCNFTScenario(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// attach token1 <- token2 (basic case) : success
	require.NoError(t, keeper.Attach(ctx, addr1, defaultSymbol, defaultTokenID1, defaultTokenID2))
	// attach token2 <- token3 (attach to a child): success
	require.NoError(t, keeper.Attach(ctx, addr1, defaultSymbol, defaultTokenID2, defaultTokenID3))
	// attach token1 <- token4 (attach to a root): success
	require.NoError(t, keeper.Attach(ctx, addr1, defaultSymbol, defaultTokenID1, defaultTokenID4))

	require.NoError(t, keeper.BurnCNFT(ctx, addr1, defaultSymbol, defaultTokenID1))

	_, err := keeper.GetNFT(ctx, defaultSymbol, defaultTokenID1)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultSymbol, defaultTokenID2)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultSymbol, defaultTokenID3)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultSymbol, defaultTokenID4)
	require.Error(t, err)

	require.Equal(t, int64(0), keeper.GetAccountBalance(ctx, defaultSymbol, defaultTokenID1, addr1).Int64())
	require.Equal(t, int64(0), keeper.GetAccountBalance(ctx, defaultSymbol, defaultTokenID2, addr1).Int64())
	require.Equal(t, int64(0), keeper.GetAccountBalance(ctx, defaultSymbol, defaultTokenID3, addr1).Int64())
	require.Equal(t, int64(0), keeper.GetAccountBalance(ctx, defaultSymbol, defaultTokenID4, addr1).Int64())
}

func TestBurnCNFTFromSuccess(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// success case
	// addr1 has: burn permission, approved
	// addr2 has: token

	// attach token1 <- token2 (basic case) : success
	require.NoError(t, keeper.Attach(ctx, addr1, defaultSymbol, defaultTokenID1, defaultTokenID2))
	// attach token2 <- token3 (attach to a child): success
	require.NoError(t, keeper.Attach(ctx, addr1, defaultSymbol, defaultTokenID2, defaultTokenID3))
	// attach token1 <- token4 (attach to a root): success
	require.NoError(t, keeper.Attach(ctx, addr1, defaultSymbol, defaultTokenID1, defaultTokenID4))

	// transfer tokens to addr2
	require.NoError(t, keeper.TransferCNFT(ctx, addr1, addr2, defaultSymbol, defaultTokenID1))
	require.NoError(t, keeper.TransferCFT(ctx, addr1, addr2, defaultSymbol, defaultTokenIDFT, sdk.NewInt(defaultAmount)))

	// addr2 approves addr1 for the symbol
	require.NoError(t, keeper.SetApproved(ctx, addr1, addr2, defaultSymbol))

	// test burnCNFTFrom
	require.NoError(t, keeper.BurnCNFTFrom(ctx, addr1, addr2, defaultSymbol, defaultTokenID1))
	require.NoError(t, keeper.BurnCFTFrom(ctx, addr1, addr2, linktype.NewCoinWithTokenIDs(linktype.NewCoinWithTokenID(defaultSymbol, defaultTokenIDFT, sdk.NewInt(defaultAmount)))))

	_, err := keeper.GetNFT(ctx, defaultSymbol, defaultTokenID1)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultSymbol, defaultTokenID2)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultSymbol, defaultTokenID3)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultSymbol, defaultTokenID4)
	require.Error(t, err)

	require.Equal(t, int64(0), keeper.GetAccountBalance(ctx, defaultSymbol, defaultTokenID1, addr1).Int64())
	require.Equal(t, int64(0), keeper.GetAccountBalance(ctx, defaultSymbol, defaultTokenID2, addr1).Int64())
	require.Equal(t, int64(0), keeper.GetAccountBalance(ctx, defaultSymbol, defaultTokenID3, addr1).Int64())
	require.Equal(t, int64(0), keeper.GetAccountBalance(ctx, defaultSymbol, defaultTokenID4, addr1).Int64())
	require.Equal(t, int64(0), keeper.GetAccountBalance(ctx, defaultSymbol, defaultTokenIDFT, addr1).Int64())
}

func TestBurnCNFTFromFailure1(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// failure case1
	// addr1 has: burn permission, approved, token
	// addr2 has: nothing

	// addr2 approves addr1 for the symbol
	require.NoError(t, keeper.SetApproved(ctx, addr1, addr2, defaultSymbol))

	// test burnCNFTFrom, burnCFTFrom fail
	require.EqualError(t, keeper.BurnCNFTFrom(ctx, addr1, addr2, defaultSymbol, defaultTokenID1), types.ErrTokenNotOwnedBy(types.DefaultCodespace, defaultSymbol+defaultTokenID1, addr2).Error())
	require.EqualError(t, keeper.BurnCFTFrom(ctx, addr1, addr2, linktype.NewCoinWithTokenIDs(linktype.NewCoinWithTokenID(defaultSymbol, defaultTokenIDFT, sdk.NewInt(1)))), sdk.ErrInsufficientCoins(fmt.Sprintf("%v has not enough coins for %v", addr2, linktype.NewCoinWithTokenID(defaultSymbol, defaultTokenIDFT, sdk.NewInt(1)).ToCoin().String())).Error())
}

func TestBurnCNFTFromFailure2(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// failure case2
	// addr1 has: burn permission (not approved)
	// addr2 has: token

	// transfer tokens to addr2
	require.NoError(t, keeper.TransferCNFT(ctx, addr1, addr2, defaultSymbol, defaultTokenID1))
	require.NoError(t, keeper.TransferCFT(ctx, addr1, addr2, defaultSymbol, defaultTokenIDFT, sdk.NewInt(1)))

	// test burnCNFTFrom fail
	require.EqualError(t, keeper.BurnCNFTFrom(ctx, addr1, addr2, defaultSymbol, defaultTokenID1), types.ErrCollectionNotApproved(types.DefaultCodespace, addr1.String(), addr2.String(), defaultSymbol).Error())
	require.EqualError(t, keeper.BurnCFTFrom(ctx, addr1, addr2, linktype.NewCoinWithTokenIDs(linktype.NewCoinWithTokenID(defaultSymbol, defaultTokenIDFT, sdk.NewInt(1)))), types.ErrCollectionNotApproved(types.DefaultCodespace, addr1.String(), addr2.String(), defaultSymbol).Error())
}

func TestBurnCNFTFromFailure3(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// failure case3
	// addr2 has: approved (no permission)
	// addr3 has: token

	// transfer tokens to addr2
	require.NoError(t, keeper.TransferCNFT(ctx, addr1, addr3, defaultSymbol, defaultTokenID1))
	require.NoError(t, keeper.TransferCFT(ctx, addr1, addr3, defaultSymbol, defaultTokenIDFT, sdk.NewInt(1)))

	// addr3 approves addr2 for the symbol
	require.NoError(t, keeper.SetApproved(ctx, addr2, addr3, defaultSymbol))

	// test burnCNFTFrom fail
	require.EqualError(t, keeper.BurnCNFTFrom(ctx, addr2, addr3, defaultSymbol, defaultTokenID1), types.ErrTokenNoPermission(types.DefaultCodespace, addr2, types.NewBurnPermission(defaultSymbol+defaultTokenID1[:4])).Error())
	require.EqualError(t, keeper.BurnCFTFrom(ctx, addr2, addr3, linktype.NewCoinWithTokenIDs(linktype.NewCoinWithTokenID(defaultSymbol, defaultTokenIDFT, sdk.NewInt(1)))), types.ErrTokenNoPermission(types.DefaultCodespace, addr2, types.NewBurnPermission(defaultSymbol+defaultTokenIDFT)).Error())
}
