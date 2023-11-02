package keeper

import (
	gogotypes "github.com/cosmos/gogoproto/types"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/collection"
)

func legacyNFTNotFoundError(k Keeper, ctx sdk.Context, contractID string, tokenID string) error {
	if err2 := collection.ValidateLegacyNFTID(tokenID); err2 == nil /* "==" is intentional */ {
		return collection.ErrTokenNotExist.Wrap(tokenID)
	}

	if err2 := collection.ValidateFTID(tokenID); err2 != nil {
		return collection.ErrTokenNotExist.Wrap(tokenID)
	}
	classID := collection.SplitTokenID(tokenID)
	if _, err2 := k.GetTokenClass(ctx, contractID, classID); err2 != nil {
		return collection.ErrTokenNotExist.Wrap(tokenID)
	}

	return collection.ErrTokenNotNFT.Wrap(tokenID)
}

func (k Keeper) hasNFT(ctx sdk.Context, contractID string, tokenID string) error {
	store := ctx.KVStore(k.storeKey)
	key := nftKey(contractID, tokenID)
	if !store.Has(key) {
		return legacyNFTNotFoundError(k, ctx, contractID, tokenID)
	}
	return nil
}

func (k Keeper) GetNFT(ctx sdk.Context, contractID string, tokenID string) (*collection.NFT, error) {
	store := ctx.KVStore(k.storeKey)
	key := nftKey(contractID, tokenID)
	bz := store.Get(key)
	if bz == nil {
		return nil, legacyNFTNotFoundError(k, ctx, contractID, tokenID)
	}

	var token collection.NFT
	k.cdc.MustUnmarshal(bz, &token)

	return &token, nil
}

func (k Keeper) setNFT(ctx sdk.Context, contractID string, token collection.NFT) {
	store := ctx.KVStore(k.storeKey)
	key := nftKey(contractID, token.TokenId)

	bz, err := token.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(key, bz)
}

func (k Keeper) deleteNFT(ctx sdk.Context, contractID string, tokenID string) {
	store := ctx.KVStore(k.storeKey)
	key := nftKey(contractID, tokenID)
	store.Delete(key)
}

func (k Keeper) pruneNFT(ctx sdk.Context, contractID string, tokenID string) []string {
	burnt := []string{}
	for _, child := range k.GetChildren(ctx, contractID, tokenID) {
		k.deleteChild(ctx, contractID, tokenID, child)
		k.deleteParent(ctx, contractID, child)
		k.deleteNFT(ctx, contractID, child)
		burnt = append(burnt, child)

		pruned := k.pruneNFT(ctx, contractID, child)
		burnt = append(burnt, pruned...)
	}
	return burnt
}

func (k Keeper) Attach(ctx sdk.Context, contractID string, owner sdk.AccAddress, subject, target string) error {
	// validate subject
	if err := k.hasNFT(ctx, contractID, subject); err != nil {
		return err
	}

	if _, err := k.GetParent(ctx, contractID, subject); err == nil {
		return collection.ErrTokenAlreadyAChild.Wrap(subject)
	}

	if !k.GetBalance(ctx, contractID, owner, subject).IsPositive() {
		return collection.ErrTokenNotOwnedBy.Wrapf("%s is not owner of %s", owner, subject)
	}

	// validate target
	if err := k.hasNFT(ctx, contractID, target); err != nil {
		return err
	}

	root := k.GetRoot(ctx, contractID, target)
	if !owner.Equals(k.getOwner(ctx, contractID, root)) {
		return collection.ErrTokenNotOwnedBy.Wrapf("%s is not owner of %s", owner, target)
	}
	if root == subject {
		return collection.ErrCannotAttachToADescendant.Wrap("cycles not allowed")
	}

	if err := k.subtractCoins(ctx, contractID, owner, collection.NewCoins(collection.NewCoin(subject, sdk.OneInt()))); err != nil {
		panic(collection.ErrInsufficientToken.Wrapf("%s not owns %s", owner, subject))
	}

	// update relation
	k.setParent(ctx, contractID, subject, target)
	k.setChild(ctx, contractID, target, subject)

	// finally, check the invariant
	if err := k.validateDepthAndWidth(ctx, contractID, root); err != nil {
		return err
	}

	// legacy
	k.iterateDescendants(ctx, contractID, subject, func(descendantID string, _ int) (stop bool) {
		event := collection.EventRootChanged{
			ContractId: contractID,
			TokenId:    descendantID,
			From:       subject,
			To:         root,
		}
		if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
			panic(err)
		}
		return false
	})

	return nil
}

func (k Keeper) Detach(ctx sdk.Context, contractID string, owner sdk.AccAddress, subject string) error {
	if err := k.hasNFT(ctx, contractID, subject); err != nil {
		return err
	}

	parent, err := k.GetParent(ctx, contractID, subject)
	if err != nil {
		return collection.ErrTokenNotAChild.Wrap(err.Error())
	}

	if !owner.Equals(k.GetRootOwner(ctx, contractID, subject)) {
		return collection.ErrTokenNotOwnedBy.Wrapf("%s is not owner of %s", owner, subject)
	}

	k.addCoins(ctx, contractID, owner, collection.NewCoins(collection.NewCoin(subject, sdk.OneInt())))

	// update relation
	k.deleteParent(ctx, contractID, subject)
	k.deleteChild(ctx, contractID, *parent, subject)

	// legacy
	root := k.GetRoot(ctx, contractID, *parent)
	k.iterateDescendants(ctx, contractID, subject, func(descendantID string, _ int) (stop bool) {
		event := collection.EventRootChanged{
			ContractId: contractID,
			TokenId:    descendantID,
			From:       root,
			To:         subject,
		}
		if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
			panic(err)
		}
		return false
	})

	return nil
}

func (k Keeper) iterateAncestors(ctx sdk.Context, contractID string, tokenID string, fn func(tokenID string) error) error {
	var err error
	for id := &tokenID; err == nil; id, err = k.GetParent(ctx, contractID, *id) {
		if fnErr := fn(*id); fnErr != nil {
			return fnErr
		}
	}

	return nil
}

func (k Keeper) GetRootOwner(ctx sdk.Context, contractID string, tokenID string) sdk.AccAddress {
	rootID := k.GetRoot(ctx, contractID, tokenID)
	return k.getOwner(ctx, contractID, rootID)
}

func (k Keeper) getOwner(ctx sdk.Context, contractID string, tokenID string) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	key := ownerKey(contractID, tokenID)
	bz := store.Get(key)
	if bz == nil {
		panic("owner must exist")
	}

	var owner sdk.AccAddress
	if err := owner.Unmarshal(bz); err != nil {
		panic(err)
	}
	return owner
}

func (k Keeper) setOwner(ctx sdk.Context, contractID string, tokenID string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := ownerKey(contractID, tokenID)

	bz, err := owner.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(key, bz)
}

func (k Keeper) deleteOwner(ctx sdk.Context, contractID string, tokenID string) {
	store := ctx.KVStore(k.storeKey)
	key := ownerKey(contractID, tokenID)
	store.Delete(key)
}

func (k Keeper) GetParent(ctx sdk.Context, contractID string, tokenID string) (*string, error) {
	store := ctx.KVStore(k.storeKey)
	key := parentKey(contractID, tokenID)
	bz := store.Get(key)
	if bz == nil {
		return nil, collection.ErrTokenNotAChild.Wrapf("%s has no parent", tokenID)
	}

	var parent gogotypes.StringValue
	k.cdc.MustUnmarshal(bz, &parent)
	return &parent.Value, nil
}

func (k Keeper) setParent(ctx sdk.Context, contractID string, tokenID, parentID string) {
	store := ctx.KVStore(k.storeKey)
	key := parentKey(contractID, tokenID)

	val := &gogotypes.StringValue{Value: parentID}
	bz := k.cdc.MustMarshal(val)
	store.Set(key, bz)
}

func (k Keeper) deleteParent(ctx sdk.Context, contractID string, tokenID string) {
	store := ctx.KVStore(k.storeKey)
	key := parentKey(contractID, tokenID)
	store.Delete(key)
}

func (k Keeper) GetChildren(ctx sdk.Context, contractID string, tokenID string) []string {
	var children []string
	k.iterateChildren(ctx, contractID, tokenID, func(childID string) (stop bool) {
		children = append(children, childID)
		return false
	})
	return children
}

func (k Keeper) iterateChildren(ctx sdk.Context, contractID string, tokenID string, fn func(childID string) (stop bool)) {
	k.iterateChildrenImpl(ctx, childKeyPrefixByTokenID(contractID, tokenID), func(_ string, _ string, childID string) (stop bool) {
		return fn(childID)
	})
}

func (k Keeper) iterateDescendants(ctx sdk.Context, contractID string, tokenID string, fn func(descendantID string, depth int) (stop bool)) {
	k.iterateDescendantsImpl(ctx, contractID, tokenID, 1, fn)
}

func (k Keeper) iterateDescendantsImpl(ctx sdk.Context, contractID string, tokenID string, depth int, fn func(descendantID string, depth int) (stop bool)) {
	k.iterateChildren(ctx, contractID, tokenID, func(childID string) (stop bool) {
		if stop := fn(childID, depth); stop {
			return true
		}

		k.iterateDescendantsImpl(ctx, contractID, childID, depth+1, fn)
		return false
	})
}

func (k Keeper) setChild(ctx sdk.Context, contractID string, tokenID, childID string) {
	store := ctx.KVStore(k.storeKey)
	key := childKey(contractID, tokenID, childID)
	store.Set(key, []byte{})
}

func (k Keeper) deleteChild(ctx sdk.Context, contractID string, tokenID, childID string) {
	store := ctx.KVStore(k.storeKey)
	key := childKey(contractID, tokenID, childID)
	store.Delete(key)
}

func (k Keeper) GetRoot(ctx sdk.Context, contractID string, tokenID string) string {
	id := tokenID
	k.iterateAncestors(ctx, contractID, tokenID, func(tokenID string) error {
		id = tokenID
		return nil
	})

	return id
}

// Deprecated
func (k Keeper) setLegacyToken(ctx sdk.Context, contractID string, tokenID string) {
	store := ctx.KVStore(k.storeKey)
	key := legacyTokenKey(contractID, tokenID)
	store.Set(key, []byte{})
}

// Deprecated
func (k Keeper) deleteLegacyToken(ctx sdk.Context, contractID string, tokenID string) {
	store := ctx.KVStore(k.storeKey)
	key := legacyTokenKey(contractID, tokenID)
	store.Delete(key)
}

// Deprecated
func (k Keeper) setLegacyTokenType(ctx sdk.Context, contractID string, tokenType string) {
	store := ctx.KVStore(k.storeKey)
	key := legacyTokenTypeKey(contractID, tokenType)
	store.Set(key, []byte{})
}

// Deprecated
func (k Keeper) validateDepthAndWidth(ctx sdk.Context, contractID string, tokenID string) error {
	widths := map[int]int{0: 1}
	k.iterateDescendants(ctx, contractID, tokenID, func(descendantID string, depth int) (stop bool) {
		widths[depth]++
		return false
	})

	params := k.GetParams(ctx)

	depth := len(widths)
	if legacyDepth := depth - 1; legacyDepth > int(params.DepthLimit) {
		return collection.ErrCompositionTooDeep.Wrapf("resulting depth exceeds its limit: %d", params.DepthLimit)
	}

	for _, width := range widths {
		if width > int(params.WidthLimit) {
			return collection.ErrCompositionTooWide.Wrapf("resulting width exceeds its limit: %d", params.WidthLimit)
		}
	}

	return nil
}
