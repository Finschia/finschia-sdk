package keeper

import (
	"fmt"
	"testing"

	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper_BurnFT(t *testing.T) {
	t.Log("implement me - ", t.Name())
}

func TestKeeper_BurnFTFrom(t *testing.T) {
	t.Log("implement me - ", t.Name())
}

func TestKeeper_BurnNFT(t *testing.T) {
	t.Log("implement me - ", t.Name())
}

func TestKeeper_BurnNFTFrom(t *testing.T) {
	t.Log("implement me - ", t.Name())
}

func TestMintBurn(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	const (
		wrongTokenID = "12345678"
	)
	require.EqualError(t, keeper.MintNFT(ctx, defaultSymbol2, addr1, types.NewNFT(defaultSymbol2, wrongTokenID, defaultName, defaultTokenURI, addr1)), types.ErrTokenTypeNotExist(types.DefaultCodespace, defaultSymbol2, wrongTokenID[:types.TokenTypeLength]).Error())
	require.EqualError(t, keeper.MintNFT(ctx, defaultSymbol, addr3, types.NewNFT(defaultSymbol, defaultTokenID1, defaultName, defaultTokenURI, addr1)), types.ErrTokenNoPermission(types.DefaultCodespace, addr3, types.NewMintPermission(defaultSymbol)).Error())

	require.NoError(t, keeper.BurnFT(ctx, defaultSymbol, addr1, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))))
	require.EqualError(t, keeper.BurnNFT(ctx, defaultSymbol, addr1, wrongTokenID), types.ErrTokenNotExist(types.DefaultCodespace, defaultSymbol, wrongTokenID).Error())
	require.EqualError(t, keeper.BurnNFT(ctx, defaultSymbol, addr3, defaultTokenID1), types.ErrTokenNoPermission(types.DefaultCodespace, addr3, types.NewBurnPermission(defaultSymbol)).Error())
}

func TestBurnNFTScenario(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// attach token1 <- token2 (basic case) : success
	require.NoError(t, keeper.Attach(ctx, addr1, defaultSymbol, defaultTokenID1, defaultTokenID2))
	// attach token2 <- token3 (attach to a child): success
	require.NoError(t, keeper.Attach(ctx, addr1, defaultSymbol, defaultTokenID2, defaultTokenID3))
	// attach token1 <- token4 (attach to a root): success
	require.NoError(t, keeper.Attach(ctx, addr1, defaultSymbol, defaultTokenID1, defaultTokenID4))

	require.NoError(t, keeper.BurnNFT(ctx, defaultSymbol, addr1, defaultTokenID1))

	_, err := keeper.GetNFT(ctx, defaultSymbol, defaultTokenID1)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultSymbol, defaultTokenID2)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultSymbol, defaultTokenID3)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultSymbol, defaultTokenID4)
	require.Error(t, err)

	balance, err := keeper.GetBalance(ctx, defaultSymbol, defaultTokenID1, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultSymbol, defaultTokenID2, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultSymbol, defaultTokenID3, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultSymbol, defaultTokenID4, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
}

func TestBurnNFTFromSuccess(t *testing.T) {
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
	require.NoError(t, keeper.TransferNFT(ctx, addr1, addr2, defaultSymbol, defaultTokenID1))
	require.NoError(t, keeper.TransferFT(ctx, addr1, addr2, defaultSymbol, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount))))

	// addr2 approves addr1 for the symbol
	require.NoError(t, keeper.SetApproved(ctx, addr1, addr2, defaultSymbol))

	// test burnNFTFrom
	require.NoError(t, keeper.BurnNFTFrom(ctx, defaultSymbol, addr1, addr2, defaultTokenID1))
	require.NoError(t, keeper.BurnFTFrom(ctx, defaultSymbol, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))))

	_, err := keeper.GetNFT(ctx, defaultSymbol, defaultTokenID1)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultSymbol, defaultTokenID2)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultSymbol, defaultTokenID3)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultSymbol, defaultTokenID4)
	require.Error(t, err)

	balance, err := keeper.GetBalance(ctx, defaultSymbol, defaultTokenID1, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultSymbol, defaultTokenID2, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultSymbol, defaultTokenID3, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultSymbol, defaultTokenID4, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultSymbol, defaultTokenIDFT, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
}

func TestBurnNFTFromFailure1(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// failure case1
	// addr1 has: burn permission, approved, token
	// addr2 has: nothing

	// addr2 approves addr1 for the symbol
	require.NoError(t, keeper.SetApproved(ctx, addr1, addr2, defaultSymbol))

	// test burnNFTFrom, burnFTFrom fail
	require.EqualError(t, keeper.BurnNFTFrom(ctx, defaultSymbol, addr1, addr2, defaultTokenID1), types.ErrTokenNotOwnedBy(types.DefaultCodespace, defaultTokenID1, addr2).Error())
	require.EqualError(t, keeper.BurnFTFrom(ctx, defaultSymbol, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), sdk.ErrInsufficientCoins(fmt.Sprintf("%v has not enough coins for %v", addr2, types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)).String())).Error())
}

func TestBurnNFTFromFailure2(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// failure case2
	// addr1 has: burn permission (not approved)
	// addr2 has: token

	// transfer tokens to addr2
	require.NoError(t, keeper.TransferNFT(ctx, addr1, addr2, defaultSymbol, defaultTokenID1))
	require.NoError(t, keeper.TransferFT(ctx, addr1, addr2, defaultSymbol, types.NewCoin(defaultTokenIDFT, sdk.NewInt(1))))

	// test burnNFTFrom fail
	require.EqualError(t, keeper.BurnNFTFrom(ctx, defaultSymbol, addr1, addr2, defaultTokenID1), types.ErrCollectionNotApproved(types.DefaultCodespace, addr1.String(), addr2.String(), defaultSymbol).Error())
	require.EqualError(t, keeper.BurnFTFrom(ctx, defaultSymbol, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), types.ErrCollectionNotApproved(types.DefaultCodespace, addr1.String(), addr2.String(), defaultSymbol).Error())
}

func TestBurnNFTFromFailure3(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// failure case3
	// addr2 has: approved (no permission)
	// addr3 has: token

	// transfer tokens to addr2
	require.NoError(t, keeper.TransferNFT(ctx, addr1, addr3, defaultSymbol, defaultTokenID1))
	require.NoError(t, keeper.TransferFT(ctx, addr1, addr3, defaultSymbol, types.NewCoin(defaultTokenIDFT, sdk.NewInt(1))))
	// addr3 approves addr2 for the symbol
	require.NoError(t, keeper.SetApproved(ctx, addr2, addr3, defaultSymbol))

	// test burnNFTFrom fail
	require.EqualError(t, keeper.BurnNFTFrom(ctx, defaultSymbol, addr2, addr3, defaultTokenID1), types.ErrTokenNoPermission(types.DefaultCodespace, addr2, types.NewBurnPermission(defaultSymbol)).Error())
	require.EqualError(t, keeper.BurnFTFrom(ctx, defaultSymbol, addr2, addr3, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), types.ErrTokenNoPermission(types.DefaultCodespace, addr2, types.NewBurnPermission(defaultSymbol)).Error())
}
