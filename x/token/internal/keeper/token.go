package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) setToken(ctx sdk.Context, token types.Token) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenDenomKey(token.GetSymbol())
	if store.Has(tokenKey) {
		return types.ErrTokenExist(types.DefaultCodespace, token.GetSymbol())
	}
	store.Set(tokenKey, k.cdc.MustMarshalBinaryBare(token))
	return nil
}

func (k Keeper) setTokenToCollection(ctx sdk.Context, token types.CollectiveToken) sdk.Error {
	c, err := k.GetCollection(ctx, token.GetSymbol())
	if err != nil {
		return err
	}
	if t, ok := token.(types.CollectiveNFT); ok {
		tokenType := t.GetTokenType()
		if !k.HasTokenType(ctx, token.GetSymbol(), tokenType) {
			return types.ErrCollectionTokenTypeNotExist(types.DefaultCodespace, token.GetSymbol(), tokenType)
		}
		if t.GetTokenIndex() == types.ReservedEmpty {
			return types.ErrCollectionTokenIndexFull(types.DefaultCodespace, token.GetSymbol(), tokenType)
		}
	}
	c, err = c.AddToken(token)
	if err != nil {
		return err
	}
	err = k.UpdateCollection(ctx, c)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) GetToken(ctx sdk.Context, symbol, tokenID string) (types.Token, sdk.Error) {
	if len(tokenID) == 0 {
		return k.getToken(ctx, symbol)
	}
	return k.getTokenFromCollection(ctx, symbol, tokenID)
}

func (k Keeper) getToken(ctx sdk.Context, denom string) (types.Token, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TokenDenomKey(denom))
	if bz == nil {
		return nil, types.ErrTokenNotExist(types.DefaultCodespace, denom)
	}
	return k.mustDecodeToken(bz), nil
}
func (k Keeper) getTokenFromCollection(ctx sdk.Context, symbol, tokenID string) (types.Token, sdk.Error) {
	c, err := k.GetCollection(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.GetToken(tokenID)
}

func (k Keeper) ModifyToken(ctx sdk.Context, token types.Token) sdk.Error {
	if len(token.GetTokenID()) == 0 {
		return k.modifyToken(ctx, token)
	}
	return k.modifyTokenToCollection(ctx, token)
}

func (k Keeper) modifyToken(ctx sdk.Context, token types.Token) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenDenomKey(token.GetSymbol())
	if !store.Has(tokenKey) {
		return types.ErrTokenNotExist(types.DefaultCodespace, token.GetSymbol())
	}
	store.Set(tokenKey, k.cdc.MustMarshalBinaryBare(token))
	return nil
}

func (k Keeper) modifyTokenToCollection(ctx sdk.Context, token types.Token) sdk.Error {
	c, err := k.GetCollection(ctx, token.GetSymbol())
	if err != nil {
		return err
	}
	c, err = c.UpdateToken(token)
	if err != nil {
		return err
	}
	err = k.UpdateCollection(ctx, c)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) GetAllTokens(ctx sdk.Context) (tokens types.Tokens) {
	appendToken := func(token types.Token) (stop bool) {
		tokens = append(tokens, token)
		return false
	}
	k.IterateTokens(ctx, "", appendToken)
	return tokens
}

func (k Keeper) IterateTokens(ctx sdk.Context, prefix string, process func(types.Token) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.TokenDenomKey(prefix))
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		token := k.mustDecodeToken(val)
		if process(token) {
			return
		}
		iter.Next()
	}
}

func (k Keeper) GetSupply(ctx sdk.Context, symbol, tokenID string) (supply sdk.Int, err sdk.Error) {
	if _, err = k.GetToken(ctx, symbol, tokenID); err != nil {
		return sdk.NewInt(0), err
	}
	return k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(symbol + tokenID), nil
}

func (k Keeper) mustDecodeToken(tokenByte []byte) (token types.Token) {
	err := k.cdc.UnmarshalBinaryBare(tokenByte, &token)
	if err != nil {
		panic(err)
	}
	return token
}
