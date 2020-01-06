package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) OccupySymbol(ctx sdk.Context, symbol string, owner sdk.AccAddress) sdk.Error {
	if k.ExistCollection(ctx, symbol) {
		return types.ErrCollectionExist(types.DefaultCodespace, symbol)
	}

	err := k.SetCollection(ctx, types.NewCollection(symbol))
	if err != nil {
		return err
	}

	perm := types.NewIssuePermission(symbol)
	k.AddPermission(ctx, owner, perm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeOccupySymbol,
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
		),
		sdk.NewEvent(
			types.EventTypeGrantPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
			sdk.NewAttribute(types.AttributeKeyResource, perm.GetResource()),
			sdk.NewAttribute(types.AttributeKeyAction, perm.GetAction()),
		),
	})

	return nil
}

func (k Keeper) ExistCollection(ctx sdk.Context, denom string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.CollectionKey(denom))
}

func (k Keeper) GetCollection(ctx sdk.Context, denom string) (collection types.Collection, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CollectionKey(denom))
	if bz == nil {
		return collection, types.ErrCollectionNotExist(types.DefaultCodespace, denom)
	}

	collection = k.mustDecodeCollection(bz)
	return collection, nil
}

func (k Keeper) GetAllCollections(ctx sdk.Context) types.Collections {
	var collections types.Collections
	appendCollection := func(collection types.Collection) (stop bool) {
		collections = append(collections, collection)
		return false
	}
	k.IterateCollections(ctx, "", appendCollection)
	return collections
}

func (k Keeper) SetCollection(ctx sdk.Context, collection types.Collection) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	if store.Has(types.CollectionKey(collection.Symbol)) {
		return types.ErrCollectionExist(types.DefaultCodespace, collection.Symbol)
	}

	store.Set(types.CollectionKey(collection.Symbol), k.cdc.MustMarshalBinaryBare(collection))
	return nil
}

func (k Keeper) IterateCollections(ctx sdk.Context, denom string, process func(types.Collection) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.CollectionKey(denom))
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		collection := k.mustDecodeCollection(val)
		if process(collection) {
			return
		}
		iter.Next()
	}
}

func (k Keeper) mustDecodeCollection(collectionByte []byte) types.Collection {
	var collection types.Collection
	err := k.cdc.UnmarshalBinaryBare(collectionByte, &collection)
	if err != nil {
		panic(err)
	}
	return collection
}
