package keeper

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
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

func (k Keeper) SetTokenType(ctx sdk.Context, symbol, tokenType string) sdk.Error {
	collection, err := k.GetCollection(ctx, symbol)
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	if store.Has(types.TokenTypeKey(collection.GetSymbol(), tokenType)) {
		return types.ErrCollectionTokenTypeExist(types.DefaultCodespace, collection.GetSymbol(), tokenType)
	}
	store.Set(types.TokenTypeKey(collection.GetSymbol(), tokenType), k.cdc.MustMarshalBinaryBare(tokenType))
	return nil
}

func (k Keeper) HasTokenType(ctx sdk.Context, symbol, tokenType string) bool {
	collection, err := k.GetCollection(ctx, symbol)
	if err != nil {
		return false
	}
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.TokenTypeKey(collection.GetSymbol(), tokenType))
}

func (k Keeper) GetNextTokenTypeForCNFT(ctx sdk.Context, symbol string) (tokenType string, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStoreReversePrefixIterator(store, types.TokenTypeKey(symbol, ""))
	defer iter.Close()
	if !iter.Valid() {
		return types.SmallestNFTType, nil
	}
	k.cdc.MustUnmarshalBinaryBare(iter.Value(), &tokenType)
	tokenType = types.NextID(tokenType, "")
	if tokenType[0] == types.FungibleFlag[0] {
		return "", types.ErrCollectionTokenTypeFull(types.DefaultCodespace, symbol)
	}
	return tokenType, nil
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

func (k Keeper) IssueFTCollection(ctx sdk.Context, token types.CollectiveFT, amount sdk.Int, owner sdk.AccAddress) sdk.Error {
	err := k.setTokenToCollection(ctx, token)
	if err != nil {
		return err
	}

	err = k.mintTokens(ctx, sdk.NewCoins(sdk.NewCoin(token.GetDenom(), amount)), owner)
	if err != nil {
		return err
	}

	tokenURIModifyPerm := types.NewModifyTokenURIPermission(token.GetDenom())
	k.AddPermission(ctx, owner, tokenURIModifyPerm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueCFT,
			sdk.NewAttribute(types.AttributeKeyName, token.GetName()),
			sdk.NewAttribute(types.AttributeKeySymbol, token.GetSymbol()),
			sdk.NewAttribute(types.AttributeKeyTokenID, token.GetTokenID()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyMintable, strconv.FormatBool(token.GetMintable())),
			sdk.NewAttribute(types.AttributeKeyDecimals, token.GetDecimals().String()),
			sdk.NewAttribute(types.AttributeKeyTokenURI, token.GetTokenURI()),
		),
		sdk.NewEvent(
			types.EventTypeGrantPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
			sdk.NewAttribute(types.AttributeKeyResource, tokenURIModifyPerm.GetResource()),
			sdk.NewAttribute(types.AttributeKeyAction, tokenURIModifyPerm.GetAction()),
		),
	})
	if token.GetMintable() {
		mintPerm := types.NewMintPermission(token.GetDenom())
		k.AddPermission(ctx, owner, mintPerm)
		burnPerm := types.NewBurnPermission(token.GetDenom())
		k.AddPermission(ctx, owner, burnPerm)

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeGrantPermToken,
				sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
				sdk.NewAttribute(types.AttributeKeyResource, mintPerm.GetResource()),
				sdk.NewAttribute(types.AttributeKeyAction, mintPerm.GetAction()),
			),
			sdk.NewEvent(
				types.EventTypeGrantPermToken,
				sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
				sdk.NewAttribute(types.AttributeKeyResource, burnPerm.GetResource()),
				sdk.NewAttribute(types.AttributeKeyAction, burnPerm.GetAction()),
			),
		})
	}

	return nil
}

func (k Keeper) MintCollectionTokens(ctx sdk.Context, amount linktype.CoinWithTokenIDs, from, to sdk.AccAddress) sdk.Error {
	for _, coin := range amount {
		symbol, tokenID := coin.Symbol, coin.TokenID
		token, err := k.GetToken(ctx, symbol, tokenID)
		if err != nil {
			return err
		}
		if err := k.isMintable(ctx, token, from); err != nil {
			return err
		}
	}
	err := k.mintTokens(ctx, amount.ToCoins(), to)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintCFT,
			sdk.NewAttribute(types.AttributeKeyAmount, amount.ToCoins().String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
		),
	})
	return nil
}
func (k Keeper) BurnCollectionTokens(ctx sdk.Context, amount linktype.CoinWithTokenIDs, from sdk.AccAddress) sdk.Error {
	coins := amount.ToCoins()

	err := k.isBurnable(ctx, coins, from)
	if err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("%v has not enough coins for %v", from, amount))
	}

	err = k.burnTokens(ctx, coins, from)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnCFT,
			sdk.NewAttribute(types.AttributeKeyAmount, amount.ToCoins().String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
		),
	})
	return nil
}

func (k Keeper) IssueCNFT(ctx sdk.Context, symbol, tokenType string, owner sdk.AccAddress) sdk.Error {
	err := k.SetTokenType(ctx, symbol, tokenType)
	if err != nil {
		return err
	}

	mintPerm := types.NewMintPermission(symbol + tokenType)
	k.AddPermission(ctx, owner, mintPerm)
	burnPerm := types.NewBurnPermission(symbol + tokenType)
	k.AddPermission(ctx, owner, burnPerm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueCNFT,
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyTokenType, tokenType),
		),
		sdk.NewEvent(
			types.EventTypeGrantPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
			sdk.NewAttribute(types.AttributeKeyResource, mintPerm.GetResource()),
			sdk.NewAttribute(types.AttributeKeyAction, mintPerm.GetAction()),
		),
		sdk.NewEvent(
			types.EventTypeGrantPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
			sdk.NewAttribute(types.AttributeKeyResource, burnPerm.GetResource()),
			sdk.NewAttribute(types.AttributeKeyAction, burnPerm.GetAction()),
		),
	})

	return nil
}

func (k Keeper) MintCollectionNFT(ctx sdk.Context, token types.CollectiveNFT, from sdk.AccAddress) sdk.Error {
	if !k.HasTokenType(ctx, token.GetSymbol(), token.GetTokenType()) {
		return types.ErrCollectionTokenTypeNotExist(types.DefaultCodespace, token.GetSymbol(), token.GetTokenType())
	}

	perm := types.NewMintPermission(token.GetSymbol() + token.GetTokenType())
	if !k.HasPermission(ctx, from, perm) {
		return types.ErrTokenPermission(types.DefaultCodespace, from, perm)
	}

	err := k.setTokenToCollection(ctx, token)
	if err != nil {
		return err
	}

	err = k.mintTokens(ctx, sdk.NewCoins(sdk.NewCoin(token.GetDenom(), sdk.NewInt(1))), token.GetOwner())
	if err != nil {
		return err
	}

	tokenURIModifyPerm := types.NewModifyTokenURIPermission(token.GetDenom())
	k.AddPermission(ctx, token.GetOwner(), tokenURIModifyPerm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintCNFT,
			sdk.NewAttribute(types.AttributeKeyName, token.GetName()),
			sdk.NewAttribute(types.AttributeKeySymbol, token.GetSymbol()),
			sdk.NewAttribute(types.AttributeKeyTokenID, token.GetTokenID()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, token.GetOwner().String()),
			sdk.NewAttribute(types.AttributeKeyTokenURI, token.GetTokenURI()),
		),
		sdk.NewEvent(
			types.EventTypeGrantPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, token.GetOwner().String()),
			sdk.NewAttribute(types.AttributeKeyResource, tokenURIModifyPerm.GetResource()),
			sdk.NewAttribute(types.AttributeKeyAction, tokenURIModifyPerm.GetAction()),
		),
	})

	return nil
}

func (k Keeper) BurnCollectionNFT(ctx sdk.Context, symbol, tokenID string, from sdk.AccAddress) sdk.Error {
	token, err := k.GetCNFT(ctx, symbol, tokenID)
	if err != nil {
		return err
	}

	perm := types.NewBurnPermission(symbol + tokenID[:types.TokenTypeLength])
	if !k.HasPermission(ctx, from, perm) {
		return types.ErrTokenPermission(types.DefaultCodespace, from, perm)
	}

	if !token.GetOwner().Equals(from) {
		return types.ErrTokenPermission(types.DefaultCodespace, from, perm)
	}

	err = k.burnCollectionNFT(ctx, token, from)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnCNFT,
			sdk.NewAttribute(types.AttributeKeySymbol, token.GetSymbol()),
			sdk.NewAttribute(types.AttributeKeyTokenID, token.GetTokenID()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
		),
	})
	return nil
}
func (k Keeper) burnCollectionNFT(ctx sdk.Context, token types.CollectiveNFT, from sdk.AccAddress) sdk.Error {
	children, err := k.ChildrenOf(ctx, token.GetSymbol(), token.GetTokenID())
	if err != nil {
		return err
	}

	for _, child := range children {
		err = k.burnCollectionNFT(ctx, child.(types.CollectiveNFT), from)
		if err != nil {
			return err
		}
	}

	parent, err := k.ParentOf(ctx, token.GetSymbol(), token.GetTokenID())
	if err != nil {
		return err
	}
	if parent != nil {
		err = k.detach(ctx, from, from, token.GetSymbol(), token.GetTokenID())
		if err != nil {
			return err
		}
	}
	collection, err := k.GetCollection(ctx, token.GetSymbol())
	if err != nil {
		return err
	}
	collection, err = collection.DeleteToken(token)
	if err != nil {
		return err
	}
	err = k.UpdateCollection(ctx, collection)
	if err != nil {
		return err
	}
	err = k.burnTokens(ctx, sdk.NewCoins(sdk.NewCoin(token.GetDenom(), sdk.NewInt(1))), token.GetOwner())
	if err != nil {
		return err
	}
	return nil
}
