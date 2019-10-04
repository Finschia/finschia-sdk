package token

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

type Keeper struct {
	bankKeeper bank.Keeper

	supplyKeeper supply.Keeper

	storeKey sdk.StoreKey

	cdc *codec.Codec
}

func NewKeeper(cdc *codec.Codec, bankKeeper bank.Keeper, supplyKeeper supply.Keeper, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		bankKeeper:   bankKeeper,
		supplyKeeper: supplyKeeper,
		storeKey:     storeKey,
		cdc:          cdc,
	}
}

func (k Keeper) GetModuleAddress() sdk.AccAddress {
	return k.supplyKeeper.GetModuleAddress(ModuleName)
}

func (k Keeper) SetToken(ctx sdk.Context, token Token) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	if store.Has([]byte(token.Symbol)) {
		return ErrTokenExist(DefaultCodespace)
	}
	store.Set(TokenSymbolKey(token.Symbol), k.cdc.MustMarshalBinaryBare(token))
	return nil
}

func (k Keeper) mustDecodeToken(tokenByte []byte) Token {
	var token Token
	err := k.cdc.UnmarshalBinaryBare(tokenByte, &token)
	if err != nil {
		panic(err)
	}
	return token
}

func (k Keeper) GetToken(ctx sdk.Context, symbol string) (token Token, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(TokenSymbolKey(symbol))
	if bz == nil {
		return token, ErrTokenNotExist(DefaultCodespace)
	}

	token = k.mustDecodeToken(bz)
	return token, nil
}

func (k Keeper) GetAllTokens(ctx sdk.Context) []Token {
	var tokens []Token
	appendToken := func(token Token) (stop bool) {
		tokens = append(tokens, token)
		return false
	}
	k.IterateTokens(ctx, appendToken)
	return tokens
}

func (k Keeper) IterateTokens(ctx sdk.Context, process func(Token) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, TokenSymbolKeyPrefix)
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

func (k Keeper) MintToken(ctx sdk.Context, amount sdk.Coin, to sdk.AccAddress) sdk.Error {
	err := k.supplyKeeper.MintCoins(ctx, ModuleName, sdk.NewCoins(amount))
	if err != nil {
		return err
	}

	err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, ModuleName, to, sdk.NewCoins(amount))
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) BurnToken(ctx sdk.Context, amount sdk.Coin, from sdk.AccAddress) sdk.Error {

	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, from, ModuleName, sdk.NewCoins(amount))
	if err != nil {
		return err
	}

	err = k.supplyKeeper.BurnCoins(ctx, ModuleName, sdk.NewCoins(amount))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetSupply(ctx sdk.Context, symbol string) (amount sdk.Int, err sdk.Error) {
	_, err = k.GetToken(ctx, symbol)
	if err != nil {
		return sdk.NewInt(0), err
	}
	return k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(symbol), nil
}
