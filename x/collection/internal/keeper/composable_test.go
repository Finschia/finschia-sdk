package keeper

import (
	"context"
	"testing"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection/internal/types"
	"github.com/line/lbm-sdk/x/contract"
	"github.com/stretchr/testify/require"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

func TestKeeper_Attach(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	require.EqualError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID1), sdkerrors.Wrapf(types.ErrCannotAttachToItself, "TokenID: %s", defaultTokenID1).Error())
	require.EqualError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID6), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenID6).Error())
	ctx2 := ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, wrongContractID))
	require.EqualError(t, keeper.Attach(ctx2, addr1, defaultTokenID1, defaultTokenID2), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", wrongContractID, defaultTokenID2).Error())
	require.EqualError(t, keeper.Attach(ctx, addr2, defaultTokenID1, defaultTokenID2), sdkerrors.Wrapf(types.ErrTokenNotOwnedBy, "TokenID: %s, Owner: %v", defaultTokenID2, addr2.String()).Error())

	require.EqualError(t, keeper.Attach(ctx, addr1, defaultTokenID6, defaultTokenID1), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenID6).Error())
	require.EqualError(t, keeper.Attach(ctx, addr2, defaultTokenID1, defaultTokenID5), sdkerrors.Wrapf(types.ErrTokenNotOwnedBy, "TokenID: %s, Owner: %s", defaultTokenID1, addr2.String()).Error())

	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID2))
	require.EqualError(t, keeper.Attach(ctx, addr1, defaultTokenID3, defaultTokenID2), sdkerrors.Wrapf(types.ErrTokenAlreadyAChild, "TokenID: %s", defaultTokenID2).Error())
	require.EqualError(t, keeper.Attach(ctx, addr1, defaultTokenID2, defaultTokenID1), sdkerrors.Wrapf(types.ErrCannotAttachToADescendant, "TokenID: %s, ToTokenID: %s", defaultTokenID1, defaultTokenID2).Error())
}

func TestKeeper_AttachFrom(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	require.EqualError(t, keeper.AttachFrom(ctx, addr1, addr2, defaultTokenID1, defaultTokenID2), sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", addr1.String(), addr2.String(), defaultContractID).Error())
	prepareProxy(ctx, t)
	require.NoError(t, keeper.AttachFrom(ctx, addr1, addr2, defaultTokenID1, defaultTokenID2))
}

func TestKeeper_Detach(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	require.EqualError(t, keeper.Detach(ctx, addr1, defaultTokenID6), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenID6).Error())
	ctx2 := ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, wrongContractID))
	require.EqualError(t, keeper.Detach(ctx2, addr1, defaultTokenID1), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", wrongContractID, defaultTokenID1).Error())
	require.EqualError(t, keeper.Detach(ctx, addr2, defaultTokenID1), sdkerrors.Wrapf(types.ErrTokenNotOwnedBy, "TokenID: %s, Owner: %s", defaultTokenID1, addr2.String()).Error())
	require.EqualError(t, keeper.Detach(ctx, addr1, defaultTokenID1), sdkerrors.Wrapf(types.ErrTokenNotAChild, "TokenID: %s", defaultTokenID1).Error())
	require.EqualError(t, keeper.Detach(ctx, addr1, defaultTokenIDFT), sdkerrors.Wrapf(types.ErrTokenNotNFT, "TokenID: %s", defaultTokenIDFT).Error())

	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID2))
	require.NoError(t, keeper.Detach(ctx, addr1, defaultTokenID2))
}

func TestKeeper_DetachFrom(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	require.EqualError(t, keeper.DetachFrom(ctx, addr1, addr2, defaultTokenID2), sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", addr1.String(), addr2.String(), defaultContractID).Error())
	prepareProxy(ctx, t)
	require.NoError(t, keeper.Attach(ctx, addr2, defaultTokenID1, defaultTokenID2))
	require.NoError(t, keeper.DetachFrom(ctx, addr1, addr2, defaultTokenID2))
}

func TestKeeper_RootOf(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	_, err := keeper.RootOf(ctx, defaultTokenID6)
	require.EqualError(t, err, sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenID6).Error())

	ctx2 := ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, wrongContractID))
	_, err = keeper.RootOf(ctx2, defaultTokenID1)
	require.EqualError(t, err, sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", wrongContractID, defaultTokenID1).Error())

	_, err = keeper.RootOf(ctx, defaultTokenIDFT)
	require.EqualError(t, err, sdkerrors.Wrapf(types.ErrTokenNotNFT, "TokenID: %s", defaultTokenIDFT).Error())

	nft, err := keeper.RootOf(ctx, defaultTokenID1)
	require.NoError(t, err)
	require.Equal(t, nft.GetContractID(), defaultContractID)
	require.Equal(t, nft.GetTokenID(), defaultTokenID1)

	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID2))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID2, defaultTokenID3))

	nft, err = keeper.RootOf(ctx, defaultTokenID3)
	require.NoError(t, err)
	require.Equal(t, nft.GetContractID(), defaultContractID)
	require.Equal(t, nft.GetTokenID(), defaultTokenID1)
}

func TestKeeper_ParentOf(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	_, err := keeper.ParentOf(ctx, defaultTokenID6)
	require.EqualError(t, err, sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenID6).Error())

	ctx2 := ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, wrongContractID))
	_, err = keeper.ParentOf(ctx2, defaultTokenID1)
	require.EqualError(t, err, sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", wrongContractID, defaultTokenID1).Error())

	_, err = keeper.ParentOf(ctx, defaultTokenIDFT)
	require.EqualError(t, err, sdkerrors.Wrapf(types.ErrTokenNotNFT, "TokenID: %s", defaultTokenIDFT).Error())

	nft, err := keeper.ParentOf(ctx, defaultTokenID1)
	require.NoError(t, err)
	require.Equal(t, nft, nil)

	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID2))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID2, defaultTokenID3))

	nft, err = keeper.ParentOf(ctx, defaultTokenID3)
	require.NoError(t, err)
	require.Equal(t, nft.GetContractID(), defaultContractID)
	require.Equal(t, nft.GetTokenID(), defaultTokenID2)
}

func TestKeeper_ChildrenOf(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	_, err := keeper.ChildrenOf(ctx, defaultTokenID6)
	require.EqualError(t, err, sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenID6).Error())

	ctx2 := ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, wrongContractID))
	_, err = keeper.ChildrenOf(ctx2, defaultTokenID1)
	require.EqualError(t, err, sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", wrongContractID, defaultTokenID1).Error())

	_, err = keeper.ChildrenOf(ctx, defaultTokenIDFT)
	require.EqualError(t, err, sdkerrors.Wrapf(types.ErrTokenNotNFT, "TokenID: %s", defaultTokenIDFT).Error())

	tokens, err := keeper.ChildrenOf(ctx, defaultTokenID1)
	require.NoError(t, err)
	require.Equal(t, len(tokens), 0)

	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID2))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID2, defaultTokenID3))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID4))

	tokens, err = keeper.ChildrenOf(ctx, defaultTokenID1)
	require.NoError(t, err)
	require.Equal(t, len(tokens), 2)
	require.Equal(t, tokens[0].GetTokenID(), defaultTokenID2)
	require.Equal(t, tokens[1].GetTokenID(), defaultTokenID4)

	tokens, err = keeper.ChildrenOf(ctx, defaultTokenID2)
	require.NoError(t, err)
	require.Equal(t, len(tokens), 1)
	require.Equal(t, tokens[0].GetTokenID(), defaultTokenID3)
}

func TestAttachDetachScenario(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	//
	// attach success cases
	//

	// attach token1 <- token2 (basic case) : success
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID2))

	// attach token2 <- token3 (attach to a child): success
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID2, defaultTokenID3))

	// attach token1 <- token4 (attach to a root): success
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID4))

	// verify the relations

	// root of token1 is token1
	rootOfToken1, err1 := keeper.RootOf(ctx, defaultTokenID1)
	require.NoError(t, err1)
	require.Equal(t, rootOfToken1.GetTokenID(), defaultTokenID1)

	// root of token2 is token1
	rootOfToken2, err2 := keeper.RootOf(ctx, defaultTokenID2)
	require.NoError(t, err2)
	require.Equal(t, rootOfToken2.GetTokenID(), defaultTokenID1)

	// root of token3 is token1
	rootOfToken3, err3 := keeper.RootOf(ctx, defaultTokenID3)
	require.NoError(t, err3)
	require.Equal(t, rootOfToken3.GetTokenID(), defaultTokenID1)

	// root of token4 is token1
	rootOfToken4, err4 := keeper.RootOf(ctx, defaultTokenID4)
	require.NoError(t, err4)
	require.Equal(t, rootOfToken4.GetTokenID(), defaultTokenID1)

	// parent of token1 is nil
	parentOfToken1, err5 := keeper.ParentOf(ctx, defaultTokenID1)
	require.NoError(t, err5)
	require.Nil(t, parentOfToken1)

	// parent of token2 is token1
	parentOfToken2, err6 := keeper.ParentOf(ctx, defaultTokenID2)
	require.NoError(t, err6)
	require.Equal(t, parentOfToken2.GetTokenID(), defaultTokenID1)

	// parent of token3 is token2
	parentOfToken3, err7 := keeper.ParentOf(ctx, defaultTokenID3)
	require.NoError(t, err7)
	require.Equal(t, parentOfToken3.GetTokenID(), defaultTokenID2)

	// parent of token4 is token1
	parentOfToken4, err8 := keeper.ParentOf(ctx, defaultTokenID4)
	require.NoError(t, err8)
	require.Equal(t, parentOfToken4.GetTokenID(), defaultTokenID1)

	// children of token1 are token2, token4
	childrenOfToken1, err9 := keeper.ChildrenOf(ctx, defaultTokenID1)
	require.NoError(t, err9)
	require.Equal(t, len(childrenOfToken1), 2)
	require.True(t, (childrenOfToken1[0].GetTokenID() == defaultTokenID2 && childrenOfToken1[1].GetTokenID() == defaultTokenID4) || (childrenOfToken1[0].GetTokenID() == defaultTokenID4 && childrenOfToken1[1].GetTokenID() == defaultTokenID2))

	// child of token2 is token3
	childrenOfToken2, err10 := keeper.ChildrenOf(ctx, defaultTokenID2)
	require.NoError(t, err10)
	require.Equal(t, len(childrenOfToken2), 1)
	require.True(t, childrenOfToken2[0].GetTokenID() == defaultTokenID3)

	// child of token3 is empty
	childrenOfToken3, err11 := keeper.ChildrenOf(ctx, defaultTokenID3)
	require.NoError(t, err11)
	require.Equal(t, len(childrenOfToken3), 0)

	// child of token4 is empty
	childrenOfToken4, err12 := keeper.ChildrenOf(ctx, defaultTokenID4)
	require.NoError(t, err12)
	require.Equal(t, len(childrenOfToken4), 0)

	// query failure cases
	_, err := keeper.ParentOf(ctx, defaultTokenIDFT)
	require.EqualError(t, err, sdkerrors.Wrapf(types.ErrTokenNotNFT, "TokenID: %s", defaultTokenIDFT).Error())

	//
	// attach error cases
	//

	// attach non-root token : failure
	require.EqualError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID2), sdkerrors.Wrapf(types.ErrTokenAlreadyAChild, "TokenID: %s", defaultTokenID2).Error())

	// attach non-exist token : failure
	require.EqualError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID8), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenID8).Error())
	require.EqualError(t, keeper.Attach(ctx, addr1, defaultTokenID8, defaultTokenID1), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenID8).Error())

	// attach non-mine token : failure
	require.EqualError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID5), sdkerrors.Wrapf(types.ErrTokenNotOwnedBy, "TokenID: %s, Owner: %s", defaultTokenID5, addr1.String()).Error())
	require.EqualError(t, keeper.Attach(ctx, addr1, defaultTokenID5, defaultTokenID1), sdkerrors.Wrapf(types.ErrTokenNotOwnedBy, "TokenID: %s, Owner: %s", defaultTokenID5, addr1.String()).Error())

	// attach to itself : failure
	require.EqualError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID1), sdkerrors.Wrapf(types.ErrCannotAttachToItself, "TokenID: %s", defaultTokenID1).Error())

	// attach to a descendant : failure
	require.EqualError(t, keeper.Attach(ctx, addr1, defaultTokenID2, defaultTokenID1), sdkerrors.Wrapf(types.ErrCannotAttachToADescendant, "TokenID: %s, ToTokenID: %s", defaultTokenID1, defaultTokenID2).Error())
	require.EqualError(t, keeper.Attach(ctx, addr1, defaultTokenID3, defaultTokenID1), sdkerrors.Wrapf(types.ErrCannotAttachToADescendant, "TokenID: %s, ToTokenID: %s", defaultTokenID1, defaultTokenID3).Error())
	require.EqualError(t, keeper.Attach(ctx, addr1, defaultTokenID4, defaultTokenID1), sdkerrors.Wrapf(types.ErrCannotAttachToADescendant, "TokenID: %s, ToTokenID: %s", defaultTokenID1, defaultTokenID4).Error())

	//
	// detach error cases
	//

	// detach not a child : failure
	require.EqualError(t, keeper.Detach(ctx, addr1, defaultTokenID1), sdkerrors.Wrapf(types.ErrTokenNotAChild, "TokenID: %s", defaultTokenID1).Error())

	// detach non-mine token : failure
	require.EqualError(t, keeper.Detach(ctx, addr1, defaultTokenID5), sdkerrors.Wrapf(types.ErrTokenNotOwnedBy, "TokenID: %s, Owner: %s", defaultTokenID5, addr1.String()).Error())

	// detach non-exist token : failure
	require.EqualError(t, keeper.Detach(ctx, addr1, defaultTokenID8), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s, TokenID: %s", defaultContractID, defaultTokenID8).Error())

	//
	// detach success cases
	//

	// detach single child
	require.NoError(t, keeper.Detach(ctx, addr1, defaultTokenID4))

	// detach a child having child
	require.NoError(t, keeper.Detach(ctx, addr1, defaultTokenID2))

	// detach child
	require.NoError(t, keeper.Detach(ctx, addr1, defaultTokenID3))

	//
	// verify the relations
	//
	// parent of token2 is nil
	parentOfToken2, err6 = keeper.ParentOf(ctx, defaultTokenID2)
	require.NoError(t, err6)
	require.Nil(t, parentOfToken2)

	// parent of token3 is nil
	parentOfToken3, err7 = keeper.ParentOf(ctx, defaultTokenID3)
	require.NoError(t, err7)
	require.Nil(t, parentOfToken3)

	// parent of token4 is nil
	parentOfToken4, err8 = keeper.ParentOf(ctx, defaultTokenID4)
	require.NoError(t, err8)
	require.Nil(t, parentOfToken4)

	// children of token1 is empty
	childrenOfToken1, err1 = keeper.ChildrenOf(ctx, defaultTokenID1)
	require.NoError(t, err1)
	require.Equal(t, len(childrenOfToken1), 0)

	// owner of token3 is addr1
	token3, err13 := keeper.GetToken(ctx, defaultTokenID3)
	require.NoError(t, err13)

	require.Equal(t, (token3.(types.NFT)).GetOwner(), addr1)
}

func setupNFTs(t *testing.T, ctx sdk.Context) {
	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, "name", "{}", defaultImgURI), addr1))
	require.NoError(t, keeper.IssueNFT(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta), addr1))
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID1, defaultName, defaultMeta, addr1)))
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID2, defaultName, defaultMeta, addr1)))
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID3, defaultName, defaultMeta, addr1)))
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID4, defaultName, defaultMeta, addr1)))
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID5, defaultName, defaultMeta, addr1)))
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID6, defaultName, defaultMeta, addr1)))
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID7, defaultName, defaultMeta, addr1)))
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID8, defaultName, defaultMeta, addr1)))
	require.NoError(t, keeper.MintNFT(ctx, addr1, types.NewNFT(defaultContractID, defaultTokenID9, defaultName, defaultMeta, addr1)))
}

func TestGetDepthWidthTable(t *testing.T) {
	ctx := cacheKeeper()

	setupNFTs(t, ctx)
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID2))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID3))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID2, defaultTokenID4))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID2, defaultTokenID5))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID2, defaultTokenID6))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID3, defaultTokenID7))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID3, defaultTokenID8))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID7, defaultTokenID9))

	//                    +- token4
	//                    +- token5
	// token1 -+- token2 -+- token6
	//         |
	//         +- token3 -+- token7 --- token9
	//                    +- token8

	var table []int
	table = keeper.GetDepthWidthTable(ctx, defaultTokenID1)
	require.Equal(t, []int{1, 2, 5, 1}, table)

	table = keeper.GetDepthWidthTable(ctx, defaultTokenID2)
	require.Equal(t, []int{1, 3}, table)

	table = keeper.GetDepthWidthTable(ctx, defaultTokenID3)
	require.Equal(t, []int{1, 2, 1}, table)
}

func TestGetCurrentDepthFromRoot(t *testing.T) {
	ctx := cacheKeeper()

	setupNFTs(t, ctx)
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID2))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID3))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID2, defaultTokenID4))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID3, defaultTokenID5))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID3, defaultTokenID6))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID4, defaultTokenID7))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID5, defaultTokenID8))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID8, defaultTokenID9))

	// token1 -+- token2 --- token4 --- token7
	//         +- token3 -+- token5 --- token8 --- token9
	//                    +- token6

	// the depth of token1 should be 0
	require.Equal(t, 0, keeper.GetDepthFromRoot(ctx, defaultTokenID1))

	// the depth of token2, token3 should be 1
	require.Equal(t, 1, keeper.GetDepthFromRoot(ctx, defaultTokenID2))
	require.Equal(t, 1, keeper.GetDepthFromRoot(ctx, defaultTokenID3))

	// the depth of token4, token5, token6 should be
	require.Equal(t, 2, keeper.GetDepthFromRoot(ctx, defaultTokenID4))
	require.Equal(t, 2, keeper.GetDepthFromRoot(ctx, defaultTokenID5))
	require.Equal(t, 2, keeper.GetDepthFromRoot(ctx, defaultTokenID6))

	// the depth of token7, token8 should be 3
	require.Equal(t, 3, keeper.GetDepthFromRoot(ctx, defaultTokenID7))
	require.Equal(t, 3, keeper.GetDepthFromRoot(ctx, defaultTokenID8))

	// the depth of token9 should be 4
	require.Equal(t, 4, keeper.GetDepthFromRoot(ctx, defaultTokenID9))
}

// nolint:dupl
func TestCheckDepth(t *testing.T) {
	ctx := cacheKeeper()
	keeper.SetParams(ctx, types.NewParams(4, 4)) // Sets the max composable width/depth to 4

	setupNFTs(t, ctx)
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID2))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID2, defaultTokenID3))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID3, defaultTokenID4))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID5, defaultTokenID6))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID5, defaultTokenID7))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID6, defaultTokenID8))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID8, defaultTokenID9))

	// Given two composed tokens
	// token1 --- token2 --- token3 --- token4
	//
	// token5 -+- token6 --- token8 --- token9
	//         +- token7
	//
	// Sets the max composable width/depth to 4

	// if token5 is attached to token2 then,
	//
	// token1 --- token2 -+- token3 --- token4
	//                    +- token5 -+- token6 --- token8 --- token9
	//                               +- token7
	// deepest depth is 5 (path: token1-token2-token5-token6-token8-token9)
	err := keeper.Attach(ctx, addr1, defaultTokenID2, defaultTokenID5)
	require.Error(t, err)

	// if token5 is attached to token1 then,
	//
	// token1 -+- token2 --- token3 --- token4
	//         +- token5 -+- token6 --- token8 --- token9
	//                    +- token7
	// deepest depth is 4 (path: token1-token5-token6-token8-token9)
	err = keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID5)
	require.NoError(t, err)
}

// nolint:dupl
func TestCheckWidth(t *testing.T) {
	ctx := cacheKeeper()
	keeper.SetParams(ctx, types.NewParams(4, 4)) // Sets the max composable width/depth to 4

	setupNFTs(t, ctx)
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID2))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID2, defaultTokenID3))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID2, defaultTokenID4))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID2, defaultTokenID5))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID6, defaultTokenID7))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID6, defaultTokenID8))
	require.NoError(t, keeper.Attach(ctx, addr1, defaultTokenID7, defaultTokenID9))

	// Given two composed tokens
	//                    +- token3
	//                    +- token4
	// token1 -+- token2 -+- token5
	//
	// token6 -+- token7 --- token9
	//         +- token8
	//
	// Sets the max composable width/depth to 4

	// if token6 is attached to token1 then,
	//
	//                    +- token3
	//                    +- token4
	// token1 -+- token2 -+- token5
	//         |
	//         +- token6 -+- token7 --- token9
	//                    +- token8
	//
	// widest width is 5
	err := keeper.Attach(ctx, addr1, defaultTokenID1, defaultTokenID6)
	require.Error(t, err)

	// if token6 is attached to token2 then,
	//
	//                    +- token3
	//                    +- token4
	// token1 -+- token2 -+- token5
	//                    +- token6 -+- token7 --- token9
	//                               +- token8
	//
	// widest width is 4
	err = keeper.Attach(ctx, addr1, defaultTokenID2, defaultTokenID6)
	require.NoError(t, err)
}
