package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/collection/internal/types"
)

var ChildExists = []byte{1}

type ComposeKeeper interface {
	Attach(ctx sdk.Context, from sdk.AccAddress, toTokenID string, tokenID string) error
	AttachFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, toTokenID string, tokenID string) error
	Detach(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, tokenID string) error
	DetachFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, tokenID string) error
	RootOf(ctx sdk.Context, tokenID string) (types.NFT, error)
	ParentOf(ctx sdk.Context, tokenID string) (types.NFT, error)
	ChildrenOf(ctx sdk.Context, tokenID string) (types.Tokens, error)
}

func (k Keeper) eventRootChanged(ctx sdk.Context, tokenID string) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeOperationRootChanged,
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
		),
	})

	children := k.getChildren(ctx, tokenID)
	for _, child := range children {
		k.eventRootChanged(ctx, child.GetTokenID())
	}
}

func (k Keeper) Attach(ctx sdk.Context, from sdk.AccAddress, toTokenID string, tokenID string) error {
	if err := k.attach(ctx, from, toTokenID, tokenID); err != nil {
		return err
	}

	newRoot, err := k.RootOf(ctx, toTokenID)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAttachToken,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyToTokenID, toTokenID),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
			sdk.NewAttribute(types.AttributeKeyOldRoot, tokenID),
			sdk.NewAttribute(types.AttributeKeyNewRoot, newRoot.GetTokenID()),
		),
	})

	k.eventRootChanged(ctx, tokenID)

	return nil
}

func (k Keeper) AttachFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, toTokenID string, tokenID string) error {
	if !k.IsApproved(ctx, proxy, from) {
		return sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", proxy.String(), from.String(), k.getContractID(ctx))
	}

	if err := k.attach(ctx, from, toTokenID, tokenID); err != nil {
		return err
	}

	newRoot, err := k.RootOf(ctx, toTokenID)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAttachFrom,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyToTokenID, toTokenID),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
			sdk.NewAttribute(types.AttributeKeyOldRoot, tokenID),
			sdk.NewAttribute(types.AttributeKeyNewRoot, newRoot.GetTokenID()),
		),
	})

	k.eventRootChanged(ctx, tokenID)

	return nil
}

func (k Keeper) attach(ctx sdk.Context, from sdk.AccAddress, parentID string, childID string) error {
	store := ctx.KVStore(k.storeKey)

	if parentID == childID {
		return sdkerrors.Wrapf(types.ErrCannotAttachToItself, "TokenID: %s", childID)
	}

	childToken, err := k.GetNFT(ctx, childID)
	if err != nil {
		return err
	}

	if !from.Equals(childToken.GetOwner()) {
		return sdkerrors.Wrapf(types.ErrTokenNotOwnedBy, "TokenID: %s, Owner: %s", childID, from.String())
	}

	toToken, err := k.GetNFT(ctx, parentID)
	if err != nil {
		return err
	}

	if !from.Equals(toToken.GetOwner()) {
		return sdkerrors.Wrapf(types.ErrTokenNotOwnedBy, "TokenID: %s, Owner: %s", parentID, from.String())
	}

	// verify token should be a root
	childToParentKey := types.TokenChildToParentKey(k.getContractID(ctx), childID)
	if store.Has(childToParentKey) {
		return sdkerrors.Wrapf(types.ErrTokenAlreadyAChild, "TokenID: %s", childID)
	}

	// verify no circulation(toToken must not be a descendant of token)
	rootOfToToken, err := k.RootOf(ctx, parentID)
	if err != nil {
		return err
	}
	parentToken, err := k.GetNFT(ctx, parentID)
	if err != nil {
		return err
	}
	if rootOfToToken.GetTokenID() == childID {
		return sdkerrors.Wrapf(types.ErrCannotAttachToADescendant, "TokenID: %s, ToTokenID: %s", childID, parentID)
	}

	if err := k.checkDepthAndWidth(ctx, rootOfToToken.GetTokenID(), parentID, childID); err != nil {
		return err
	}

	parentToChildKey := types.TokenParentToChildKey(k.getContractID(ctx), parentID, childID)
	if store.Has(parentToChildKey) {
		panic("token is already a child of some other")
	}

	if !from.Equals(parentToken.GetOwner()) {
		if err := k.moveNFToken(ctx, from, parentToken.GetOwner(), childToken); err != nil {
			return err
		}
	}

	store.Set(childToParentKey, k.mustEncodeString(parentID))
	store.Set(parentToChildKey, k.mustEncodeString(childID))

	return nil
}

func (k Keeper) checkDepthAndWidth(ctx sdk.Context, rootID, parentID, childID string) error {
	rootTable := k.GetDepthWidthTable(ctx, rootID)
	childTable := k.GetDepthWidthTable(ctx, childID)

	parentDepth := k.GetDepthFromRoot(ctx, parentID)

	// root: token1 - token2 - token3 => depth of token3 is 2
	// child: token4 - token5 => length 2
	// attach result: token1 - token2 - token3 - token4 - token5 => depth 4
	// [depth of token3] + [len([token4,token5])] should be result depth
	resultDepth := uint64(parentDepth + len(childTable))
	if resultDepth > k.GetParams(ctx).MaxComposableDepth {
		return sdkerrors.Wrapf(types.ErrCompositionTooDeep, "Depth: %d", resultDepth)
	}

	//  root table: [1, 2, 3, 4, 5, 6, 7, 8]
	// child table:             [1, 3, 5, 7, 9]
	// if the child is attached after depth 3,
	// then the merged width table should be [1, 2, 3, 4, 5+1, 6+3, 7+5, 8+7, 9]
	maxComposableWidth := k.GetParams(ctx).MaxComposableWidth
	for curDepth, idx := parentDepth+1, 0; curDepth < len(rootTable) && idx < len(childTable); {
		totalWidth := uint64(rootTable[curDepth] + childTable[idx])
		if totalWidth > maxComposableWidth {
			return sdkerrors.Wrapf(types.ErrCompositionTooWide, "Width: %d (at depth %d)", totalWidth, curDepth)
		}
		curDepth++
		idx++
	}

	return nil
}

// Gets the depth-width(count) table (array)
//
// lv0(1)     lv1(2)     lv2(3)     lv3(2)     lv4(1)
// token1 -+- token2 --- token4 --- token7
//         +- token3 -+- token5 --- token8 --- token9
//                    +- token6
//
// then returns [1, 2, 3, 2, 1]
// and len([1, 2, 3, 2, 1]) - 1 represents the depth of token1's children
func (k Keeper) GetDepthWidthTable(ctx sdk.Context, tokenID string) []int {
	table := make([]int, 1)
	table[0] = 1
	k.fillDepthWidthTable(ctx, tokenID, &table, 1)
	return table
}

func (k Keeper) fillDepthWidthTable(ctx sdk.Context, tokenID string, table *[]int, index int) {
	count := 0

	k.iterateChildren(ctx, tokenID, func(tokenID string) bool {
		k.fillDepthWidthTable(ctx, tokenID, table, index+1)
		count++
		return false
	})

	// if count = 0, the current is leaf, so doesn't need to insert a row
	if count > 0 {
		for len(*table) <= index {
			*table = append(*table, 0)
		}
		// fills the children count of current depth
		(*table)[index] += count
	}
}

func (k Keeper) GetDepthFromRoot(ctx sdk.Context, tokenID string) int {
	store := ctx.KVStore(k.storeKey)

	depth := 0
	for nextID := tokenID; ; depth++ {
		childToParentKey := types.TokenChildToParentKey(k.getContractID(ctx), nextID)
		bz := store.Get(childToParentKey)
		if bz == nil {
			break
		}
		nextID = k.mustDecodeString(bz)
	}
	return depth
}

func (k Keeper) Detach(ctx sdk.Context, from sdk.AccAddress, tokenID string) error {
	oldRoot, err := k.RootOf(ctx, tokenID)
	if err != nil {
		return err
	}

	parentTokenID, err := k.detach(ctx, from, tokenID)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDetachToken,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyFromTokenID, parentTokenID),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
			sdk.NewAttribute(types.AttributeKeyOldRoot, oldRoot.GetTokenID()),
			sdk.NewAttribute(types.AttributeKeyNewRoot, tokenID),
		),
	})

	k.eventRootChanged(ctx, tokenID)

	return nil
}

// nolint:dupl
func (k Keeper) DetachFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, tokenID string) error {
	oldRoot, err := k.RootOf(ctx, tokenID)
	if err != nil {
		return err
	}

	if !k.IsApproved(ctx, proxy, from) {
		return sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", proxy.String(), from.String(), k.getContractID(ctx))
	}

	parentTokenID, err := k.detach(ctx, from, tokenID)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDetachFrom,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyFromTokenID, parentTokenID),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
			sdk.NewAttribute(types.AttributeKeyOldRoot, oldRoot.GetTokenID()),
			sdk.NewAttribute(types.AttributeKeyNewRoot, tokenID),
		),
	})

	k.eventRootChanged(ctx, tokenID)

	return nil
}

func (k Keeper) detach(ctx sdk.Context, from sdk.AccAddress, childID string) (string, error) {
	store := ctx.KVStore(k.storeKey)

	childToken, err := k.GetNFT(ctx, childID)
	if err != nil {
		return "", err
	}

	if !from.Equals(childToken.GetOwner()) {
		return "", sdkerrors.Wrapf(types.ErrTokenNotOwnedBy, "TokenID: %s, Owner: %s", childID, from.String())
	}

	childToParentKey := types.TokenChildToParentKey(k.getContractID(ctx), childID)
	if !store.Has(childToParentKey) {
		return "", sdkerrors.Wrapf(types.ErrTokenNotAChild, "TokenID: %s", childID)
	}

	bz := store.Get(childToParentKey)
	parentID := k.mustDecodeString(bz)

	_, err = k.GetNFT(ctx, parentID)
	if err != nil {
		return "", err
	}

	parentToChildKey := types.TokenParentToChildKey(k.getContractID(ctx), parentID, childID)
	if !store.Has(parentToChildKey) {
		panic("token is not a child of some other")
	}

	store.Delete(childToParentKey)
	store.Delete(parentToChildKey)

	return parentID, nil
}

func (k Keeper) RootOf(ctx sdk.Context, tokenID string) (types.NFT, error) {
	store := ctx.KVStore(k.storeKey)

	token, err := k.GetNFT(ctx, tokenID)
	if err != nil {
		return nil, err
	}

	for childToParentKey := types.TokenChildToParentKey(k.getContractID(ctx), token.GetTokenID()); store.Has(childToParentKey); {
		bz := store.Get(childToParentKey)
		parentTokenID := k.mustDecodeString(bz)

		token, err = k.GetNFT(ctx, parentTokenID)
		if err != nil {
			return nil, err
		}
		childToParentKey = types.TokenChildToParentKey(k.getContractID(ctx), token.GetTokenID())
	}

	return token, nil
}

func (k Keeper) ParentOf(ctx sdk.Context, tokenID string) (types.NFT, error) {
	store := ctx.KVStore(k.storeKey)

	token, err := k.GetNFT(ctx, tokenID)
	if err != nil {
		return nil, err
	}
	childToParentKey := types.TokenChildToParentKey(k.getContractID(ctx), token.GetTokenID())
	if store.Has(childToParentKey) {
		bz := store.Get(childToParentKey)
		parentTokenID := k.mustDecodeString(bz)

		return k.GetNFT(ctx, parentTokenID)
	}
	return nil, nil
}

func (k Keeper) ChildrenOf(ctx sdk.Context, tokenID string) (types.Tokens, error) {
	_, err := k.GetNFT(ctx, tokenID)
	if err != nil {
		return nil, err
	}
	tokens := k.getChildren(ctx, tokenID)

	return tokens, nil
}

func (k Keeper) getChildren(ctx sdk.Context, parentID string) (tokens types.Tokens) {
	getToken := func(tokenID string) bool {
		token, err := k.GetNFT(ctx, tokenID)
		if err != nil {
			panic(err)
		}
		tokens = append(tokens, token)
		return false
	}

	k.iterateChildren(ctx, parentID, getToken)

	return tokens
}

func (k Keeper) iterateChildren(ctx sdk.Context, parentID string, process func(string) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.TokenParentToChildSubKey(k.getContractID(ctx), parentID))
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		tokenID := k.mustDecodeString(val)
		if process(tokenID) {
			return
		}
		iter.Next()
	}
}
