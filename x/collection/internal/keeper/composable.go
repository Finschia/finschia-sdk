package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

var ChildExists = []byte{1}

type ComposeKeeper interface {
	Attach(ctx sdk.Context, contractID string, from sdk.AccAddress, toTokenID string, tokenID string) sdk.Error
	AttachFrom(ctx sdk.Context, contractID string, proxy sdk.AccAddress, from sdk.AccAddress, toTokenID string, tokenID string) sdk.Error
	Detach(ctx sdk.Context, contractID string, from sdk.AccAddress, to sdk.AccAddress, tokenID string) sdk.Error
	DetachFrom(ctx sdk.Context, contractID string, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, tokenID string) sdk.Error
	RootOf(ctx sdk.Context, contractID string, tokenID string) (types.NFT, sdk.Error)
	ParentOf(ctx sdk.Context, contractID string, tokenID string) (types.NFT, sdk.Error)
	ChildrenOf(ctx sdk.Context, contractID string, tokenID string) (types.Tokens, sdk.Error)
}

func (k Keeper) Attach(ctx sdk.Context, contractID string, from sdk.AccAddress, toTokenID string, tokenID string) sdk.Error {
	if err := k.attach(ctx, contractID, from, toTokenID, tokenID); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAttachToken,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyToTokenID, toTokenID),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
		),
	})

	return nil
}

func (k Keeper) AttachFrom(ctx sdk.Context, contractID string, proxy sdk.AccAddress, from sdk.AccAddress, toTokenID string, tokenID string) sdk.Error {
	if !k.IsApproved(ctx, contractID, proxy, from) {
		return types.ErrCollectionNotApproved(types.DefaultCodespace, proxy.String(), from.String(), contractID)
	}

	if err := k.attach(ctx, contractID, from, toTokenID, tokenID); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAttachFrom,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyToTokenID, toTokenID),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
		),
	})

	return nil
}

func (k Keeper) attach(ctx sdk.Context, contractID string, from sdk.AccAddress, parentID string, childID string) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	if parentID == childID {
		return types.ErrCannotAttachToItself(types.DefaultCodespace, childID)
	}

	childToken, err := k.GetNFT(ctx, contractID, childID)
	if err != nil {
		return err
	}

	if !from.Equals(childToken.GetOwner()) {
		return types.ErrTokenNotOwnedBy(types.DefaultCodespace, childID, from)
	}

	toToken, err := k.GetNFT(ctx, contractID, parentID)
	if err != nil {
		return err
	}

	if !from.Equals(toToken.GetOwner()) {
		return types.ErrTokenNotOwnedBy(types.DefaultCodespace, parentID, from)
	}

	// verify token should be a root
	childToParentKey := types.TokenChildToParentKey(contractID, childID)
	if store.Has(childToParentKey) {
		return types.ErrTokenAlreadyAChild(types.DefaultCodespace, childID)
	}

	// verify no circulation(toToken must not be a descendant of token)
	rootOfToToken, err := k.RootOf(ctx, contractID, parentID)
	if err != nil {
		return err
	}
	parentToken, err := k.GetNFT(ctx, contractID, parentID)
	if err != nil {
		return err
	}
	if rootOfToToken.GetTokenID() == childID {
		return types.ErrCannotAttachToADescendant(types.DefaultCodespace, childID, parentID)
	}

	parentToChildKey := types.TokenParentToChildKey(contractID, parentID, childID)
	if store.Has(parentToChildKey) {
		panic("token is already a child of some other")
	}

	if !from.Equals(parentToken.GetOwner()) {
		if err := k.moveNFToken(ctx, contractID, from, parentToken.GetOwner(), childToken); err != nil {
			return err
		}
	}

	store.Set(childToParentKey, k.mustEncodeString(parentID))
	store.Set(parentToChildKey, k.mustEncodeString(childID))

	return nil
}

func (k Keeper) Detach(ctx sdk.Context, contractID string, from sdk.AccAddress, tokenID string) sdk.Error {
	parentTokenID, err := k.detach(ctx, contractID, from, tokenID)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDetachToken,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyFromTokenID, parentTokenID),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
		),
	})
	return nil
}

//nolint:dupl
func (k Keeper) DetachFrom(ctx sdk.Context, contractID string, proxy sdk.AccAddress, from sdk.AccAddress, tokenID string) sdk.Error {
	if !k.IsApproved(ctx, contractID, proxy, from) {
		return types.ErrCollectionNotApproved(types.DefaultCodespace, proxy.String(), from.String(), contractID)
	}

	parentTokenID, err := k.detach(ctx, contractID, from, tokenID)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDetachFrom,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyFromTokenID, parentTokenID),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
		),
	})

	return nil
}

func (k Keeper) detach(ctx sdk.Context, contractID string, from sdk.AccAddress, childID string) (string, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	childToken, err := k.GetNFT(ctx, contractID, childID)
	if err != nil {
		return "", err
	}

	if !from.Equals(childToken.GetOwner()) {
		return "", types.ErrTokenNotOwnedBy(types.DefaultCodespace, childID, from)
	}

	childToParentKey := types.TokenChildToParentKey(contractID, childID)
	if !store.Has(childToParentKey) {
		return "", types.ErrTokenNotAChild(types.DefaultCodespace, childID)
	}

	bz := store.Get(childToParentKey)
	parentID := k.mustDecodeString(bz)

	_, err = k.GetNFT(ctx, contractID, parentID)
	if err != nil {
		return "", err
	}

	parentToChildKey := types.TokenParentToChildKey(contractID, parentID, childID)
	if !store.Has(parentToChildKey) {
		panic("token is not a child of some other")
	}

	store.Delete(childToParentKey)
	store.Delete(parentToChildKey)

	return parentID, nil
}

func (k Keeper) RootOf(ctx sdk.Context, contractID string, tokenID string) (types.NFT, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	token, err := k.GetNFT(ctx, contractID, tokenID)
	if err != nil {
		return nil, err
	}

	for childToParentKey := types.TokenChildToParentKey(contractID, token.GetTokenID()); store.Has(childToParentKey); {
		bz := store.Get(childToParentKey)
		parentTokenID := k.mustDecodeString(bz)

		token, err = k.GetNFT(ctx, contractID, parentTokenID)
		if err != nil {
			return nil, err
		}
		childToParentKey = types.TokenChildToParentKey(contractID, token.GetTokenID())
	}

	return token, nil
}

func (k Keeper) ParentOf(ctx sdk.Context, contractID string, tokenID string) (types.NFT, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	token, err := k.GetNFT(ctx, contractID, tokenID)
	if err != nil {
		return nil, err
	}
	childToParentKey := types.TokenChildToParentKey(contractID, token.GetTokenID())
	if store.Has(childToParentKey) {
		bz := store.Get(childToParentKey)
		parentTokenID := k.mustDecodeString(bz)

		return k.GetNFT(ctx, contractID, parentTokenID)
	}
	return nil, nil
}

func (k Keeper) ChildrenOf(ctx sdk.Context, contractID string, tokenID string) (types.Tokens, sdk.Error) {
	_, err := k.GetNFT(ctx, contractID, tokenID)
	if err != nil {
		return nil, err
	}
	tokens := k.getChildren(ctx, contractID, tokenID)

	return tokens, nil
}

func (k Keeper) getChildren(ctx sdk.Context, contractID, parentID string) (tokens types.Tokens) {
	getToken := func(tokenID string) bool {
		token, err := k.GetNFT(ctx, contractID, tokenID)
		if err != nil {
			panic(err)
		}
		tokens = append(tokens, token)
		return false
	}

	k.iterateChildren(ctx, contractID, parentID, getToken)

	return tokens
}

func (k Keeper) iterateChildren(ctx sdk.Context, contractID, parentID string, process func(string) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.TokenParentToChildSubKey(contractID, parentID))
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
