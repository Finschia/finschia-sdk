package keeper

import (
	"github.com/line/link/x/collection/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenTypeKeeper interface {
	GetNextTokenType(ctx sdk.Context, symbol string) (tokenType string, err sdk.Error)
	GetTokenTypes(ctx sdk.Context, symbol string) (types.TokenTypes, sdk.Error)
	GetTokenType(ctx sdk.Context, symbol, tokenType string) (types.TokenType, sdk.Error)
	HasTokenType(ctx sdk.Context, symbol, tokenType string) bool
	SetTokenType(ctx sdk.Context, symbol string, token types.TokenType) sdk.Error
	UpdateTokenType(ctx sdk.Context, symbol string, token types.TokenType) sdk.Error
}

var _ TokenTypeKeeper = (*Keeper)(nil)

func (k Keeper) SetTokenType(ctx sdk.Context, symbol string, tokenType types.TokenType) sdk.Error {
	_, err := k.GetCollection(ctx, symbol)
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	if store.Has(types.TokenTypeKey(symbol, tokenType.GetTokenType())) {
		return types.ErrTokenTypeExist(types.DefaultCodespace, symbol, tokenType.GetTokenType())
	}
	store.Set(types.TokenTypeKey(symbol, tokenType.GetTokenType()), k.cdc.MustMarshalBinaryBare(tokenType))
	return nil
}

func (k Keeper) UpdateTokenType(ctx sdk.Context, symbol string, tokenType types.TokenType) sdk.Error {
	_, err := k.GetCollection(ctx, symbol)
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.TokenTypeKey(symbol, tokenType.GetTokenType())) {
		return types.ErrTokenTypeNotExist(types.DefaultCodespace, symbol, tokenType.GetTokenType())
	}
	store.Set(types.TokenTypeKey(symbol, tokenType.GetTokenType()), k.cdc.MustMarshalBinaryBare(tokenType))
	return nil
}

func (k Keeper) GetTokenType(ctx sdk.Context, symbol string, tokenTypeID string) (types.TokenType, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	tokenTypeKey := types.TokenTypeKey(symbol, tokenTypeID)
	bz := store.Get(tokenTypeKey)
	if bz == nil {
		return nil, types.ErrTokenTypeNotExist(types.DefaultCodespace, symbol, tokenTypeID)
	}
	tokenType := k.mustDecodeTokenType(bz)
	return tokenType, nil
}

func (k Keeper) GetTokenTypes(ctx sdk.Context, symbol string) (tokenTypes types.TokenTypes, err sdk.Error) {
	_, err = k.GetCollection(ctx, symbol)
	if err != nil {
		return nil, err
	}
	k.iterateTokenTypes(ctx, symbol, "", false, func(t types.TokenType) bool {
		tokenTypes = append(tokenTypes, t)
		return false
	})
	return tokenTypes, nil
}

func (k Keeper) HasTokenType(ctx sdk.Context, symbol, tokenType string) bool {
	store := ctx.KVStore(k.storeKey)
	tokenTypeKey := types.TokenTypeKey(symbol, tokenType)
	return store.Has(tokenTypeKey)
}

func (k Keeper) GetNextTokenType(ctx sdk.Context, symbol string) (tokenType string, err sdk.Error) {
	var lastTokenType types.TokenType
	k.iterateTokenTypes(ctx, symbol, "", true, func(t types.TokenType) bool {
		lastTokenType = t
		return true
	})

	if lastTokenType == nil {
		return types.SmallestNFTType, nil
	}
	tokenType = nextID(lastTokenType.GetTokenType(), "")
	if tokenType[0] == types.FungibleFlag[0] {
		return "", types.ErrTokenTypeFull(types.DefaultCodespace, symbol)
	}
	return tokenType, nil
}

func (k Keeper) iterateTokenTypes(ctx sdk.Context, symbol, prefix string, reverse bool, process func(types.TokenType) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	var iter sdk.Iterator
	if reverse {
		iter = sdk.KVStoreReversePrefixIterator(store, types.TokenTypeKey(symbol, prefix))
	} else {
		iter = sdk.KVStorePrefixIterator(store, types.TokenTypeKey(symbol, prefix))
	}
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		tokenType := k.mustDecodeTokenType(val)
		if process(tokenType) {
			return
		}
		iter.Next()
	}
}

func (k Keeper) mustDecodeTokenType(bz []byte) (tokenType types.TokenType) {
	err := k.cdc.UnmarshalBinaryBare(bz, &tokenType)
	if err != nil {
		panic(err)
	}
	return tokenType
}
