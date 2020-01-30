package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) ExistCollection(ctx sdk.Context, symbol string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.CollectionKey(symbol))
}

func (k Keeper) GetCollection(ctx sdk.Context, symbol string) (collection types.Collection, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CollectionKey(symbol))
	if bz == nil {
		return collection, types.ErrCollectionNotExist(types.DefaultCodespace, symbol)
	}

	collection = k.mustDecodeCollection(bz)
	return collection, nil
}

func (k Keeper) GetNFTCount(ctx sdk.Context, symbol, baseID string) (count sdk.Int, err sdk.Error) {
	collection, err := k.GetCollection(ctx, symbol)
	if err != nil {
		return count, err
	}
	tokens := collection.GetNFTokens()
	tokens = tokens.GetTokens(baseID)
	count = sdk.NewInt(int64(tokens.Len()))
	return count, nil
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
	if store.Has(types.CollectionKey(collection.GetSymbol())) {
		return types.ErrCollectionExist(types.DefaultCodespace, collection.GetSymbol())
	}

	store.Set(types.CollectionKey(collection.GetSymbol()), k.cdc.MustMarshalBinaryBare(collection))
	return nil
}

func (k Keeper) UpdateCollection(ctx sdk.Context, collection types.Collection) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.CollectionKey(collection.GetSymbol())) {
		return types.ErrCollectionNotExist(types.DefaultCodespace, collection.GetSymbol())
	}

	store.Set(types.CollectionKey(collection.GetSymbol()), k.cdc.MustMarshalBinaryBare(collection))
	return nil
}

func (k Keeper) IterateCollections(ctx sdk.Context, symbol string, process func(types.Collection) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.CollectionKey(symbol))
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
	//XXX:
	for _, token := range collection.GetAllTokens() {
		token.(types.CollectiveToken).SetCollection(collection)
	}
	return collection
}
