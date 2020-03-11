package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_MintFT(t *testing.T) {
	ctx := cacheKeeper()

	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, "name", defaultMeta, defaultImgURI), addr1))
	require.NoError(t, keeper.IssueFT(ctx, addr1, addr1, types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(1), true), sdk.NewInt(defaultAmount)))
	require.NoError(t, keeper.IssueFT(ctx, addr1, addr1, types.NewFT(defaultContractID, defaultTokenIDFT2, defaultName, defaultMeta, sdk.NewInt(1), true), sdk.NewInt(defaultAmount)))
	require.NoError(t, keeper.IssueFT(ctx, addr1, addr1, types.NewFT(defaultContractID, defaultTokenIDFT3, defaultName, defaultMeta, sdk.NewInt(1), false), sdk.NewInt(defaultAmount)))
	require.NoError(t, keeper.IssueNFT(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta), addr1))

	require.EqualError(t, keeper.MintFT(ctx, wrongContractID, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(10)))), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", wrongContractID, defaultTokenIDFT).Error())
	require.EqualError(t, keeper.MintFT(ctx, defaultContractID, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT4, sdk.NewInt(10)))), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenIDFT4).Error())
	require.EqualError(t, keeper.MintFT(ctx, defaultContractID, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT3, sdk.NewInt(10)))), sdkerrors.Wrapf(types.ErrTokenNotMintable, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenIDFT3).Error())
	require.EqualError(t, keeper.MintFT(ctx, defaultContractID, addr2, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(10)))), sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr2.String(), types.NewMintPermission(defaultContractID)).Error())
	require.EqualError(t, keeper.MintFT(ctx, defaultContractID, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenID1, sdk.NewInt(10)))), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenID1).Error())

	require.NoError(t, keeper.MintFT(ctx, defaultContractID, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(10)))))
	require.NoError(t, keeper.MintFT(ctx, defaultContractID, addr1, addr2, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(10)), types.NewCoin(defaultTokenIDFT2, sdk.NewInt(20)))))
}

func TestKeeper_MintNFT(t *testing.T) {
	ctx := cacheKeeper()

	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, "name", defaultMeta, defaultImgURI), addr1))
	require.NoError(t, keeper.IssueNFT(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta), addr1))

	require.EqualError(t, keeper.MintNFT(ctx, addr2, types.NewNFT(defaultContractID, defaultTokenID1, "sword", defaultMeta, addr1)), sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr2.String(), types.NewMintPermission(defaultContractID)).Error())
	require.EqualError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(wrongContractID, defaultTokenID1, "sword", defaultMeta, addr1)), sdkerrors.Wrapf(types.ErrTokenTypeNotExist, "ContractID: %s, TokenType: %s", wrongContractID, defaultTokenType).Error())
	require.EqualError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, wrongTokenID, "sword", defaultMeta, addr1)), sdkerrors.Wrapf(types.ErrTokenTypeNotExist, "ContractID: %s, TokenType: %s", defaultContractID, defaultTokenType2).Error())
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID1, "sword", defaultMeta, addr1)))
}
