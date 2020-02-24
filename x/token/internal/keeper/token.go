package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) GetToken(ctx sdk.Context, symbol string) (types.Token, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TokenSymbolKey(symbol))
	if bz == nil {
		return nil, types.ErrTokenNotExist(types.DefaultCodespace, symbol)
	}
	return k.mustDecodeToken(bz), nil
}

func (k Keeper) SetToken(ctx sdk.Context, token types.Token) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenSymbolKey(token.GetSymbol())
	if store.Has(tokenKey) {
		return types.ErrTokenExist(types.DefaultCodespace, token.GetSymbol())
	}
	store.Set(tokenKey, k.cdc.MustMarshalBinaryBare(token))

	k.setSupply(ctx, types.DefaultSupply(token.GetSymbol()))

	return nil
}

func (k Keeper) UpdateToken(ctx sdk.Context, token types.Token) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenSymbolKey(token.GetSymbol())
	if !store.Has(tokenKey) {
		return types.ErrTokenNotExist(types.DefaultCodespace, token.GetSymbol())
	}
	store.Set(tokenKey, k.cdc.MustMarshalBinaryBare(token))
	return nil
}

func (k Keeper) ModifyTokenURI(ctx sdk.Context, owner sdk.AccAddress, symbol, tokenURI string) sdk.Error {
	if !types.ValidTokenURI(tokenURI) {
		return types.ErrInvalidTokenURILength(types.DefaultCodespace, tokenURI)
	}

	token, err := k.GetToken(ctx, symbol)
	if err != nil {
		return err
	}
	tokenURIModifyPerm := types.NewModifyTokenURIPermission(token.GetSymbol())
	if !k.HasPermission(ctx, owner, tokenURIModifyPerm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, owner, tokenURIModifyPerm)
	}
	token.SetTokenURI(tokenURI)

	err = k.UpdateToken(ctx, token)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeModifyTokenURI,
			sdk.NewAttribute(types.AttributeKeyName, token.GetName()),
			sdk.NewAttribute(types.AttributeKeySymbol, token.GetSymbol()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyTokenURI, token.GetTokenURI()),
		),
	})
	return nil
}

func (k Keeper) GetAllTokens(ctx sdk.Context) (tokens types.Tokens) {
	appendToken := func(token types.Token) (stop bool) {
		tokens = append(tokens, token)
		return false
	}
	k.iterateTokens(ctx, "", appendToken)
	return tokens
}

func (k Keeper) iterateTokens(ctx sdk.Context, prefix string, process func(types.Token) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.TokenSymbolKey(prefix))
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
