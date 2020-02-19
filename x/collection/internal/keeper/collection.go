package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

type CollectionKeeper interface {
	CreateCollection(ctx sdk.Context, collection types.Collection, owner sdk.AccAddress) sdk.Error
	ExistCollection(ctx sdk.Context, symbol string) bool
	GetCollection(ctx sdk.Context, symbol string) (collection types.Collection, err sdk.Error)
	SetCollection(ctx sdk.Context, collection types.Collection) sdk.Error
	UpdateCollection(ctx sdk.Context, collection types.Collection) sdk.Error
	GetNFTCount(ctx sdk.Context, symbol, baseID string) (count sdk.Int, err sdk.Error)
	GetAllCollections(ctx sdk.Context) types.Collections
}

var _ CollectionKeeper = (*Keeper)(nil)

func (k Keeper) CreateCollection(ctx sdk.Context, collection types.Collection, owner sdk.AccAddress) sdk.Error {
	err := k.SetCollection(ctx, collection)
	if err != nil {
		return err
	}

	perm := types.NewIssuePermission(collection.GetSymbol())
	k.AddPermission(ctx, owner, perm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateCollection,
			sdk.NewAttribute(types.AttributeKeyName, collection.GetName()),
			sdk.NewAttribute(types.AttributeKeySymbol, collection.GetSymbol()),
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
	k.iterateCollections(ctx, "", appendCollection)
	return collections
}

func (k Keeper) iterateCollections(ctx sdk.Context, symbol string, process func(types.Collection) (stop bool)) {
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
		token.(types.Token).SetCollection(collection)
	}
	return collection
}
