package keeper

import (
	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/collection"
)

const (
	DepthLimit = 15
)

func (k Keeper) GetNFT(ctx sdk.Context, contractID string, tokenID string) (*collection.NFT, error) {
	store := ctx.KVStore(k.storeKey)
	key := nftKey(contractID, tokenID)
	bz := store.Get(key)
	if bz == nil {
		return nil, sdkerrors.ErrNotFound.Wrapf("nft not exists: %s", tokenID)
	}

	var token collection.NFT
	k.cdc.MustUnmarshal(bz, &token)

	return &token, nil
}

func (k Keeper) setNFT(ctx sdk.Context, contractID string, token collection.NFT) {
	store := ctx.KVStore(k.storeKey)
	key := nftKey(contractID, token.Id)

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

func (k Keeper) Attach(ctx sdk.Context, contractID string, owner sdk.AccAddress, subject, target string) error {
	// validate subject
	if !k.GetBalance(ctx, contractID, owner, subject).IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrapf("%s is not owner of %s", owner, subject)
	}

	// validate target
	if _, err := k.GetNFT(ctx, contractID, target); err != nil {
		return err
	}

	root, err := k.GetRoot(ctx, contractID, target)
	if err != nil {
		return err
	}

	if owner != k.getOwner(ctx, contractID, *root) {
		return sdkerrors.ErrInvalidRequest.Wrapf("%s is not owner of %s", owner, subject)
	}

	// update subject
	k.deleteOwner(ctx, contractID, subject)
	k.setParent(ctx, contractID, subject, target)

	// update target
	k.setChild(ctx, contractID, target, subject)

	return nil
}

func (k Keeper) Detach(ctx sdk.Context, contractID string, owner sdk.AccAddress, subject string) error {
	if _, err := k.GetNFT(ctx, contractID, subject); err != nil {
		return err
	}

	parent, err := k.GetParent(ctx, contractID, subject)
	if err != nil {
		return err
	}

	root, err := k.GetRoot(ctx, contractID, *parent)
	if err != nil {
		return err
	}

	if owner != k.getOwner(ctx, contractID, *root) {
		return sdkerrors.ErrInvalidRequest.Wrapf("%s is not owner of %s", owner, subject)
	}

	// update subject
	k.deleteParent(ctx, contractID, subject)
	k.setOwner(ctx, contractID, subject, owner)

	// update parent
	k.deleteChild(ctx, contractID, *parent, subject)

	return nil
}

func (k Keeper) getRootOwnerUnbounded(ctx sdk.Context, contractID string, tokenID string) sdk.AccAddress {
	rootID := k.getRootUnbounded(ctx, contractID, tokenID)
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
		return nil, sdkerrors.ErrNotFound.Wrapf("%s has no parent", tokenID)
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

//nolint:unused
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

func (k Keeper) getRootUnbounded(ctx sdk.Context, contractID string, tokenID string) string {
	id := tokenID
	for {
		parent, err := k.GetParent(ctx, contractID, id)
		if err != nil {
			return id
		}

		id = *parent
	}
}

func (k Keeper) GetRoot(ctx sdk.Context, contractID string, tokenID string) (*string, error) {
	id := tokenID
	for depth := 0; depth <= DepthLimit; depth++ {
		parent, err := k.GetParent(ctx, contractID, id)
		if err != nil {
			return &id, nil
		}

		id = *parent
	}

	return nil, sdkerrors.ErrInvalidRequest.Wrapf("depth of %s exceeds the limit: %d", tokenID, DepthLimit)
}

// legacy index
func (k Keeper) setLegacyToken(ctx sdk.Context, contractID string, tokenID string) {
	store := ctx.KVStore(k.storeKey)
	key := legacyTokenKey(contractID, tokenID)
	store.Set(key, []byte{})
}

func (k Keeper) deleteLegacyToken(ctx sdk.Context, contractID string, tokenID string) {
	store := ctx.KVStore(k.storeKey)
	key := legacyTokenKey(contractID, tokenID)
	store.Delete(key)
}

func (k Keeper) setLegacyTokenType(ctx sdk.Context, contractID string, tokenType string) {
	store := ctx.KVStore(k.storeKey)
	key := legacyTokenTypeKey(contractID, tokenType)
	store.Set(key, []byte{})
}
