package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token/internal/types"
)

func (k Keeper) GetToken(ctx sdk.Context) (types.Token, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TokenKey(k.getContractID(ctx)))
	if bz == nil {
		return nil, sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s", k.getContractID(ctx))
	}
	return k.mustDecodeToken(bz), nil
}

func (k Keeper) SetToken(ctx sdk.Context, token types.Token) error {
	if k.getContractID(ctx) != token.GetContractID() {
		panic("cannot set token with different contract id")
	}
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenKey(k.getContractID(ctx))
	if store.Has(tokenKey) {
		return sdkerrors.Wrapf(types.ErrTokenExist, "ContractID: %s", k.getContractID(ctx))
	}
	store.Set(tokenKey, k.cdc.MustMarshalBinaryBare(token))

	k.setSupply(ctx, types.DefaultSupply(token.GetContractID()))

	return nil
}

func (k Keeper) UpdateToken(ctx sdk.Context, token types.Token) error {
	if k.getContractID(ctx) != token.GetContractID() {
		panic("cannot update token with different contract id")
	}
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.TokenKey(k.getContractID(ctx))
	if !store.Has(tokenKey) {
		return sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s", k.getContractID(ctx))
	}
	store.Set(tokenKey, k.cdc.MustMarshalBinaryBare(token))
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
	iter := sdk.KVStorePrefixIterator(store, types.TokenKey(prefix))
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
