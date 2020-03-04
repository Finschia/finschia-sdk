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
	require.EqualError(t, keeper.MintNFT(ctx, defaultContractID2, addr1, types.NewNFT(defaultContractID2, wrongTokenID, defaultName, addr1)), types.ErrTokenTypeNotExist(types.DefaultCodespace, defaultContractID2, wrongTokenID[:types.TokenTypeLength]).Error())
	require.EqualError(t, keeper.MintNFT(ctx, defaultContractID, addr3, types.NewNFT(defaultContractID, defaultTokenID1, defaultName, addr1)), types.ErrTokenNoPermission(types.DefaultCodespace, addr3, types.NewMintPermission(defaultContractID)).Error())

	require.NoError(t, keeper.BurnFT(ctx, defaultContractID, addr1, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))))
	require.EqualError(t, keeper.BurnNFT(ctx, defaultContractID, addr1, wrongTokenID), types.ErrTokenNotExist(types.DefaultCodespace, defaultContractID, wrongTokenID).Error())
	require.EqualError(t, keeper.BurnNFT(ctx, defaultContractID, addr3, defaultTokenID1), types.ErrTokenNoPermission(types.DefaultCodespace, addr3, types.NewBurnPermission(defaultContractID)).Error())
}

func TestBurnNFTScenario(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// attach token1 <- token2 (basic case) : success
	require.NoError(t, keeper.Attach(ctx, defaultContractID, addr1, defaultTokenID1, defaultTokenID2))
	// attach token2 <- token3 (attach to a child): success
	require.NoError(t, keeper.Attach(ctx, defaultContractID, addr1, defaultTokenID2, defaultTokenID3))
	// attach token1 <- token4 (attach to a root): success
	require.NoError(t, keeper.Attach(ctx, defaultContractID, addr1, defaultTokenID1, defaultTokenID4))

	require.NoError(t, keeper.BurnNFT(ctx, defaultContractID, addr1, defaultTokenID1))

	_, err := keeper.GetNFT(ctx, defaultContractID, defaultTokenID1)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultContractID, defaultTokenID2)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultContractID, defaultTokenID3)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultContractID, defaultTokenID4)
	require.Error(t, err)

	balance, err := keeper.GetBalance(ctx, defaultContractID, defaultTokenID1, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultContractID, defaultTokenID2, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultContractID, defaultTokenID3, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultContractID, defaultTokenID4, addr1)
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
	require.NoError(t, keeper.Attach(ctx, defaultContractID, addr1, defaultTokenID1, defaultTokenID2))
	// attach token2 <- token3 (attach to a child): success
	require.NoError(t, keeper.Attach(ctx, defaultContractID, addr1, defaultTokenID2, defaultTokenID3))
	// attach token1 <- token4 (attach to a root): success
	require.NoError(t, keeper.Attach(ctx, defaultContractID, addr1, defaultTokenID1, defaultTokenID4))

	// transfer tokens to addr2
	require.NoError(t, keeper.TransferNFT(ctx, defaultContractID, addr1, addr2, defaultTokenID1))
	require.NoError(t, keeper.TransferFT(ctx, defaultContractID, addr1, addr2, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount))))

	// addr2 approves addr1 for the contractID
	require.NoError(t, keeper.SetApproved(ctx, defaultContractID, addr1, addr2))

	// test burnNFTFrom
	require.NoError(t, keeper.BurnNFTFrom(ctx, defaultContractID, addr1, addr2, defaultTokenID1))
	require.NoError(t, keeper.BurnFTFrom(ctx, defaultContractID, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))))

	_, err := keeper.GetNFT(ctx, defaultContractID, defaultTokenID1)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultContractID, defaultTokenID2)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultContractID, defaultTokenID3)
	require.Error(t, err)
	_, err = keeper.GetNFT(ctx, defaultContractID, defaultTokenID4)
	require.Error(t, err)

	balance, err := keeper.GetBalance(ctx, defaultContractID, defaultTokenID1, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultContractID, defaultTokenID2, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultContractID, defaultTokenID3, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultContractID, defaultTokenID4, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
	balance, err = keeper.GetBalance(ctx, defaultContractID, defaultTokenIDFT, addr1)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())
}

func TestBurnNFTFromFailure1(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// failure case1
	// addr1 has: burn permission, approved, token
	// addr2 has: nothing

	// addr2 approves addr1 for the contractID
	require.NoError(t, keeper.SetApproved(ctx, defaultContractID, addr1, addr2))

	// test burnNFTFrom, burnFTFrom fail
	require.EqualError(t, keeper.BurnNFTFrom(ctx, defaultContractID, addr1, addr2, defaultTokenID1), types.ErrTokenNotOwnedBy(types.DefaultCodespace, defaultTokenID1, addr2).Error())
	require.EqualError(t, keeper.BurnFTFrom(ctx, defaultContractID, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), types.ErrInsufficientToken(types.DefaultCodespace, fmt.Sprintf("%v has not enough coins for %v", addr2, types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)).String())).Error())
}

func TestBurnNFTFromFailure2(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// failure case2
	// addr1 has: burn permission (not approved)
	// addr2 has: token

	// transfer tokens to addr2
	require.NoError(t, keeper.TransferNFT(ctx, defaultContractID, addr1, addr2, defaultTokenID1))
	require.NoError(t, keeper.TransferFT(ctx, defaultContractID, addr1, addr2, types.NewCoin(defaultTokenIDFT, sdk.NewInt(1))))

	// test burnNFTFrom fail
	require.EqualError(t, keeper.BurnNFTFrom(ctx, defaultContractID, addr1, addr2, defaultTokenID1), types.ErrCollectionNotApproved(types.DefaultCodespace, addr1.String(), addr2.String(), defaultContractID).Error())
	require.EqualError(t, keeper.BurnFTFrom(ctx, defaultContractID, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), types.ErrCollectionNotApproved(types.DefaultCodespace, addr1.String(), addr2.String(), defaultContractID).Error())
}

func TestBurnNFTFromFailure3(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	// failure case3
	// addr2 has: approved (no permission)
	// addr3 has: token

	// transfer tokens to addr2
	require.NoError(t, keeper.TransferNFT(ctx, defaultContractID, addr1, addr3, defaultTokenID1))
	require.NoError(t, keeper.TransferFT(ctx, defaultContractID, addr1, addr3, types.NewCoin(defaultTokenIDFT, sdk.NewInt(1))))
	// addr3 approves addr2 for the contractID
	require.NoError(t, keeper.SetApproved(ctx, defaultContractID, addr2, addr3))

	// test burnNFTFrom fail
	require.EqualError(t, keeper.BurnNFTFrom(ctx, defaultContractID, addr2, addr3, defaultTokenID1), types.ErrTokenNoPermission(types.DefaultCodespace, addr2, types.NewBurnPermission(defaultContractID)).Error())
	require.EqualError(t, keeper.BurnFTFrom(ctx, defaultContractID, addr2, addr3, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), types.ErrTokenNoPermission(types.DefaultCodespace, addr2, types.NewBurnPermission(defaultContractID)).Error())
}
