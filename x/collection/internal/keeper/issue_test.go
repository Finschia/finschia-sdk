package keeper

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/collection/internal/types"
	"github.com/line/link-modules/x/contract"
	"github.com/stretchr/testify/require"
)

func TestKeeper_IssueFT(t *testing.T) {
	ctx := cacheKeeper()

	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, "name", defaultMeta, defaultImgURI), addr1))

	ctx2 := ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, wrongContractID))
	require.EqualError(t, keeper.IssueFT(ctx2, addr1, addr1, types.NewFT(wrongContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(1), true), sdk.NewInt(defaultAmount)), sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", wrongContractID).Error())
	require.NoError(t, keeper.IssueFT(ctx, addr1, addr1, types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(1), true), sdk.NewInt(defaultAmount)))
	require.EqualError(t, keeper.IssueFT(ctx, addr1, addr1, types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(1), true), sdk.NewInt(defaultAmount)), sdkerrors.Wrapf(types.ErrTokenExist, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenIDFT).Error())
}

func TestKeeper_IssueNFT(t *testing.T) {
	ctx := cacheKeeper()

	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, "name", defaultMeta, defaultImgURI), addr1))

	ctx2 := ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, wrongContractID))
	require.EqualError(t, keeper.IssueNFT(ctx2, types.NewBaseTokenType(wrongContractID, defaultTokenType, defaultName, defaultMeta), addr1), sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", wrongContractID).Error())
	require.NoError(t, keeper.IssueNFT(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta), addr1))
	require.EqualError(t, keeper.IssueNFT(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta), addr1), sdkerrors.Wrapf(types.ErrTokenTypeExist, "ContractID: %s, TokenType: %s", defaultContractID, defaultTokenType).Error())
}
