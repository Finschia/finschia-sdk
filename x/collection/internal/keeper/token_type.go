package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/collection/internal/types"
)

type TokenTypeKeeper interface {
	GetNextTokenType(ctx sdk.Context) (tokenType string, err error)
	GetTokenTypes(ctx sdk.Context) (types.TokenTypes, error)
	GetTokenType(ctx sdk.Context, tokenType string) (types.TokenType, error)
	HasTokenType(ctx sdk.Context, tokenType string) bool
	SetTokenType(ctx sdk.Context, token types.TokenType) error
	UpdateTokenType(ctx sdk.Context, token types.TokenType) error
}

var _ TokenTypeKeeper = (*Keeper)(nil)

func (k Keeper) SetTokenType(ctx sdk.Context, tokenType types.TokenType) error {
	_, err := k.GetCollection(ctx)
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	if store.Has(types.TokenTypeKey(k.getContractID(ctx), tokenType.GetTokenType())) {
		return sdkerrors.Wrapf(types.ErrTokenTypeExist, "ContractID: %s, TokenType: %s", k.getContractID(ctx), tokenType.GetTokenType())
	}
	store.Set(types.TokenTypeKey(k.getContractID(ctx), tokenType.GetTokenType()), k.cdc.MustMarshalBinaryBare(tokenType))
	k.setNextTokenTypeNFT(ctx, tokenType.GetTokenType())
	k.setNextTokenIndexNFT(ctx, tokenType.GetTokenType(), types.ReservedEmpty)
	return nil
}

func (k Keeper) UpdateTokenType(ctx sdk.Context, tokenType types.TokenType) error {
	_, err := k.GetCollection(ctx)
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.TokenTypeKey(k.getContractID(ctx), tokenType.GetTokenType())) {
		return sdkerrors.Wrapf(types.ErrTokenTypeNotExist, "ContractID: %s, TokenType: %s", k.getContractID(ctx), tokenType.GetTokenType())
	}
	store.Set(types.TokenTypeKey(k.getContractID(ctx), tokenType.GetTokenType()), k.cdc.MustMarshalBinaryBare(tokenType))
	return nil
}

func (k Keeper) GetTokenType(ctx sdk.Context, tokenTypeID string) (types.TokenType, error) {
	store := ctx.KVStore(k.storeKey)
	tokenTypeKey := types.TokenTypeKey(k.getContractID(ctx), tokenTypeID)
	bz := store.Get(tokenTypeKey)
	if bz == nil {
		return nil, sdkerrors.Wrapf(types.ErrTokenTypeNotExist, "ContractID: %s, TokenType: %s", k.getContractID(ctx), tokenTypeID)
	}
	tokenType := k.mustDecodeTokenType(bz)
	return tokenType, nil
}

func (k Keeper) GetTokenTypes(ctx sdk.Context) (tokenTypes types.TokenTypes, err error) {
	_, err = k.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	k.iterateTokenTypes(ctx, "", false, func(t types.TokenType) bool {
		tokenTypes = append(tokenTypes, t)
		return false
	})
	return tokenTypes, nil
}

func (k Keeper) HasTokenType(ctx sdk.Context, tokenType string) bool {
	store := ctx.KVStore(k.storeKey)
	tokenTypeKey := types.TokenTypeKey(k.getContractID(ctx), tokenType)
	return store.Has(tokenTypeKey)
}

func (k Keeper) GetNextTokenType(ctx sdk.Context) (tokenType string, err error) {
	if !k.ExistCollection(ctx) {
		return "", sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", k.getContractID(ctx))
	}
	tokenType, err = k.getNextTokenTypeNFT(ctx)
	if err != nil {
		return "", err
	}
	if tokenType[0] == types.FungibleFlag[0] {
		return "", sdkerrors.Wrapf(types.ErrTokenTypeFull, "ContractID: %s", k.getContractID(ctx))
	}
	return tokenType, nil
}

func (k Keeper) iterateTokenTypes(ctx sdk.Context, prefix string, reverse bool, process func(types.TokenType) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	var iter sdk.Iterator
	if reverse {
		iter = sdk.KVStoreReversePrefixIterator(store, types.TokenTypeKey(k.getContractID(ctx), prefix))
	} else {
		iter = sdk.KVStorePrefixIterator(store, types.TokenTypeKey(k.getContractID(ctx), prefix))
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
