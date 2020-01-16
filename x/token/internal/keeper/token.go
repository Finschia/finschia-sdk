package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) SetToken(ctx sdk.Context, token types.Token) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenDenomKey(token.GetDenom())
	if store.Has(tokenKey) {
		return types.ErrTokenExist(types.DefaultCodespace, token.GetDenom())
	}
	store.Set(tokenKey, k.cdc.MustMarshalBinaryBare(token))
	return nil
}

func (k Keeper) GetToken(ctx sdk.Context, denom string) (types.Token, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TokenDenomKey(denom))
	if bz == nil {
		return nil, types.ErrTokenNotExist(types.DefaultCodespace, denom)
	}
	return k.mustDecodeToken(bz), nil
}

func (k Keeper) ModifyToken(ctx sdk.Context, token types.Token) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	tsk := types.TokenDenomKey(token.GetDenom())
	if !store.Has(tsk) {
		return types.ErrTokenNotExist(types.DefaultCodespace, token.GetDenom())
	}
	store.Set(tsk, k.cdc.MustMarshalBinaryBare(token))
	return nil
}

func (k Keeper) GetAllTokens(ctx sdk.Context) types.Tokens {
	return k.GetPrefixedTokens(ctx, "")
}

func (k Keeper) GetPrefixedTokens(ctx sdk.Context, prefix string) (tokens types.Tokens) {
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

func (k Keeper) GetTokenWithSupply(ctx sdk.Context, denom string) (token types.Token, supply sdk.Int, err sdk.Error) {
	token, err = k.GetToken(ctx, denom)
	if err != nil {
		return token, supply, err
	}
	supply = k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(denom)

	return token, supply, nil
}

func (k Keeper) IterateTokens(ctx sdk.Context, denom string, process func(types.Token) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.TokenDenomKey(denom))
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

func (k Keeper) mustDecodeToken(tokenByte []byte) (token types.Token) {
	err := k.cdc.UnmarshalBinaryBare(tokenByte, &token)
	if err != nil {
		panic(err)
	}
	return token
}
