package keeper

import (
	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

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
	store := ctx.KVStore(k.storeKey)
	prefix := childKeyPrefixByTokenID(contractID, tokenID)
	iter := sdk.KVStorePrefixIterator(store, prefix)

	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		_, _, childID := splitChildKey(iter.Key())
		if fn(childID) {
			break
		}
	}
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
