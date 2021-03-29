package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
)

type CollectionKeeper interface {
	CreateCollection(ctx sdk.Context, collection types.Collection, owner sdk.AccAddress) error
	ExistCollection(ctx sdk.Context) bool
	GetCollection(ctx sdk.Context) (collection types.Collection, err error)
	SetCollection(ctx sdk.Context, collection types.Collection) error
	UpdateCollection(ctx sdk.Context, collection types.Collection) error
}

var _ CollectionKeeper = (*Keeper)(nil)

func (k Keeper) CreateCollection(ctx sdk.Context, collection types.Collection, owner sdk.AccAddress) error {
	err := k.SetCollection(ctx, collection)
	if err != nil {
		return err
	}
	k.SetSupply(ctx, types.DefaultSupply(collection.GetContractID()))

	perms := types.NewPermissions(
		types.NewIssuePermission(),
		types.NewMintPermission(),
		types.NewBurnPermission(),
		types.NewModifyPermission(),
	)
	for _, perm := range perms {
		k.AddPermission(ctx, owner, perm)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateCollection,
			sdk.NewAttribute(types.AttributeKeyContractID, collection.GetContractID()),
			sdk.NewAttribute(types.AttributeKeyName, collection.GetName()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
		),
		sdk.NewEvent(
			types.EventTypeGrantPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
			sdk.NewAttribute(types.AttributeKeyContractID, collection.GetContractID()),
		),
	})
	for _, perm := range perms {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeGrantPermToken,
				sdk.NewAttribute(types.AttributeKeyPerm, perm.String()),
			),
		})
	}

	return nil
}

func (k Keeper) ExistCollection(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.CollectionKey(k.getContractID(ctx)))
}

func (k Keeper) GetCollection(ctx sdk.Context) (collection types.Collection, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CollectionKey(k.getContractID(ctx)))
	if bz == nil {
		return collection, sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", k.getContractID(ctx))
	}

	collection = k.mustDecodeCollection(bz)
	return collection, nil
}

func (k Keeper) SetCollection(ctx sdk.Context, collection types.Collection) error {
	store := ctx.KVStore(k.storeKey)
	if store.Has(types.CollectionKey(collection.GetContractID())) {
		return sdkerrors.Wrapf(types.ErrCollectionExist, "ContractID: %s", collection.GetContractID())
	}

	store.Set(types.CollectionKey(collection.GetContractID()), k.cdc.MustMarshalBinaryBare(collection))
	k.setNextTokenTypeFT(ctx, types.ReservedEmpty)
	k.setNextTokenTypeNFT(ctx, types.ReservedEmptyNFT)
	return nil
}

func (k Keeper) UpdateCollection(ctx sdk.Context, collection types.Collection) error {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.CollectionKey(collection.GetContractID())) {
		return sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", collection.GetContractID())
	}

	store.Set(types.CollectionKey(collection.GetContractID()), k.cdc.MustMarshalBinaryBare(collection))
	return nil
}

func (k Keeper) GetAllCollections(ctx sdk.Context) types.Collections {
	var collections types.Collections
	appendCollection := func(collection types.Collection) (stop bool) {
		collections = append(collections, collection)
		return false
	}
	k.iterateCollections(ctx, appendCollection)
	return collections
}

func (k Keeper) iterateCollections(ctx sdk.Context, process func(types.Collection) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.CollectionKey(""))
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
