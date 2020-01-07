package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) SetToken(ctx sdk.Context, token types.Token) sdk.Error {

	store := ctx.KVStore(k.storeKey)
	if store.Has(types.TokenSymbolKey(token.Symbol)) {
		return types.ErrTokenExist(types.DefaultCodespace, token.Symbol)
	}
	store.Set(types.TokenSymbolKey(token.Symbol), k.cdc.MustMarshalBinaryBare(token))
	return nil
}

func (k Keeper) GetToken(ctx sdk.Context, symbol string) (token types.Token, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TokenSymbolKey(symbol))
	if bz == nil {
		return token, types.ErrTokenNotExist(types.DefaultCodespace, symbol)
	}

	token = k.mustDecodeToken(bz)
	return token, nil
}

func (k Keeper) GetAllTokens(ctx sdk.Context) types.Tokens {
	return k.GetPrefixedTokens(ctx, "")
}

func (k Keeper) GetPrefixedTokens(ctx sdk.Context, prefix string) types.Tokens {
	var tokens types.Tokens
	appendToken := func(token types.Token) (stop bool) {
		tokens = append(tokens, token)
		return false
	}
	k.IterateTokens(ctx, prefix, appendToken)
	return tokens
}

func (k Keeper) GetSupply(ctx sdk.Context, symbol string) (supply sdk.Int, err sdk.Error) {
	if _, err = k.GetToken(ctx, symbol); err != nil {
		return sdk.NewInt(0), err
	}
	return k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(symbol), nil
}

func (k Keeper) GetTokenWithSupply(ctx sdk.Context, symbol string) (token types.Token, supply sdk.Int, err sdk.Error) {
	token, err = k.GetToken(ctx, symbol)
	if err != nil {
		return token, supply, err
	}
	supply = k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(symbol)

	return token, supply, nil
}

func (k Keeper) IterateTokens(ctx sdk.Context, denom string, process func(types.Token) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.TokenSymbolKey(denom))
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

func (k Keeper) mustDecodeToken(tokenByte []byte) types.Token {
	var token types.Token
	err := k.cdc.UnmarshalBinaryBare(tokenByte, &token)
	if err != nil {
		panic(err)
	}
	return token
}
