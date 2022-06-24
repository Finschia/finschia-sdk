package keeper

import (
	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/collection"
)

func (k Keeper) ValidateNFTID(ctx sdk.Context, contractID string, tokenID string) error {
	classID := collection.SplitTokenID(tokenID)
	class, err := k.GetTokenClass(ctx, contractID, classID)
	if err != nil {
		return err
	}

	if _, ok := class.(*collection.NFTClass); !ok {
		return sdkerrors.ErrInvalidType.Wrapf("invalid class: %s", classID)
	}

	return nil
}

func (k Keeper) GetNFT(ctx sdk.Context, contractID string, tokenID string) (*collection.NFT, error) {
	store := ctx.KVStore(k.storeKey)
	key := nftKey(contractID, tokenID)
	bz := store.Get(key)
	if bz == nil {
		return nil, sdkerrors.ErrNotFound.Wrapf("nft not exists: %s", tokenID)
	}

	var nft collection.NFT
	k.cdc.MustUnmarshal(bz, &nft)

	return &nft, nil
}

func (k Keeper) setNFT(ctx sdk.Context, contractID string, nft collection.NFT) {
	store := ctx.KVStore(k.storeKey)
	key := nftKey(contractID, nft.Id)

	bz, err := nft.Marshal()
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

func (k Keeper) GetRoot(ctx sdk.Context, contractID string, tokenID string) string {
	id := tokenID
	for {
		parent, err := k.GetParent(ctx, contractID, id)
		if err != nil {
			return id
		}

		id = *parent
	}
}

func (k Keeper) isRoot(ctx sdk.Context, contractID string, tokenID string) bool {
	_, err := k.GetParent(ctx, contractID, tokenID)
	return err != nil
}

// legacy
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

func (k Keeper) deleteLegacyTokenType(ctx sdk.Context, contractID string, tokenType string) {
	store := ctx.KVStore(k.storeKey)
	key := legacyTokenTypeKey(contractID, tokenType)
	store.Delete(key)
}
