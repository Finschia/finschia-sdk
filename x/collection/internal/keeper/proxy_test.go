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

func TestApproveDisapproveScenario(t *testing.T) {
	ctx := cacheKeeper()
	const (
		defaultTokenIDFromContractID2 = defaultTokenType2 + "00000001"
	)

	// prepare collection, FT, NFT
	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, "name", defaultMeta, defaultImgURI), addr1))
	require.NoError(t, keeper.IssueFT(ctx, addr1, addr1, types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(defaultDecimals), true), sdk.NewInt(defaultAmount)))
	require.NoError(t, keeper.IssueNFT(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta), addr1))
	require.NoError(t, keeper.IssueNFT(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType2, defaultName, defaultMeta), addr1))
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID1, defaultName, defaultMeta, addr1)))
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenType2+"00000001", defaultName, defaultMeta, addr1)))

	// approve test
	ctx2 := ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, defaultContractID2))
	require.EqualError(t, keeper.SetApproved(ctx2, addr3, addr1), sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", defaultContractID2).Error())
	require.NoError(t, keeper.SetApproved(ctx, addr3, addr1))
	require.EqualError(t, keeper.SetApproved(ctx, addr3, addr1), sdkerrors.Wrapf(types.ErrCollectionAlreadyApproved, "Proxy: %s, Approver: %s, ContractID: %s", addr3.String(), addr1.String(), defaultContractID).Error())

	// attach_from/detach_from test
	require.EqualError(t, keeper.AttachFrom(ctx, addr2, addr1, defaultTokenID1, defaultTokenIDFromContractID2), sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", addr2.String(), addr1.String(), defaultContractID).Error())
	require.NoError(t, keeper.AttachFrom(ctx, addr3, addr1, defaultTokenID1, defaultTokenIDFromContractID2))
	require.EqualError(t, keeper.DetachFrom(ctx, addr2, addr1, defaultTokenIDFromContractID2), sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", addr2.String(), addr1.String(), defaultContractID).Error())
	require.NoError(t, keeper.DetachFrom(ctx, addr3, addr1, defaultTokenIDFromContractID2))

	// transfer_from test
	require.EqualError(t, keeper.TransferFTFrom(ctx, addr2, addr1, addr2, types.NewCoin(defaultTokenIDFT, sdk.NewInt(10))), sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", addr2.String(), addr1.String(), defaultContractID).Error())
	require.NoError(t, keeper.TransferFTFrom(ctx, addr3, addr1, addr2, types.NewCoin(defaultTokenIDFT, sdk.NewInt(10))))

	require.EqualError(t, keeper.TransferNFTFrom(ctx, addr2, addr1, addr2, defaultTokenID1), sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", addr2.String(), addr1.String(), defaultContractID).Error())
	require.NoError(t, keeper.TransferNFTFrom(ctx, addr3, addr1, addr2, defaultTokenID1))

	// disapprove test
	ctx2 = ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, defaultContractID2))
	require.EqualError(t, keeper.DeleteApproved(ctx2, addr3, addr1), sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", defaultContractID2).Error())
	require.NoError(t, keeper.DeleteApproved(ctx, addr3, addr1))
	require.EqualError(t, keeper.DeleteApproved(ctx, addr3, addr1), sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", addr3.String(), addr1.String(), defaultContractID).Error())
}
