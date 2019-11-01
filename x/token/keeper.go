package token

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/link-chain/link/x/token/types"
)

type Keeper struct {
	supplyKeeper types.SupplyKeeper

	iamKeeper types.IamKeeper

	storeKey sdk.StoreKey

	cdc *codec.Codec
}

func NewKeeper(cdc *codec.Codec, supplyKeeper types.SupplyKeeper, iamKeeper types.IamKeeper, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		supplyKeeper: supplyKeeper,
		iamKeeper:    iamKeeper.WithPrefix(ModuleName),
		storeKey:     storeKey,
		cdc:          cdc,
	}
}

func (k Keeper) GetModuleAddress() sdk.AccAddress {
	return k.supplyKeeper.GetModuleAddress(ModuleName)
}

func (k Keeper) SetToken(ctx sdk.Context, token Token) sdk.Error {

	if token.Mintable {
		k.iamKeeper.GrantPermission(ctx, token.Owner, NewMintPermission(token.Symbol))
		k.iamKeeper.GrantPermission(ctx, token.Owner, NewBurnPermission(token.Symbol))
	}

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

func (k Keeper) GetAllTokens(ctx sdk.Context) Tokens {
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

func (k Keeper) AddPermission(ctx sdk.Context, addr sdk.AccAddress, perm PermissionI) {
	k.iamKeeper.GrantPermission(ctx, addr, perm)
}

func (k Keeper) RevokePermission(ctx sdk.Context, addr sdk.AccAddress, perm PermissionI) sdk.Error {
	if !k.HasPermission(ctx, addr, perm) {
		return ErrTokenPermission(DefaultCodespace)
	}
	k.iamKeeper.RevokePermission(ctx, addr, perm)
	return nil
}

func (k Keeper) HasPermission(ctx sdk.Context, addr sdk.AccAddress, perm PermissionI) bool {
	return k.iamKeeper.HasPermission(ctx, addr, perm)
}

func (k Keeper) InheritPermission(ctx sdk.Context, parent, child sdk.AccAddress) {
	k.iamKeeper.InheritPermission(ctx, parent, child)
}

func (k Keeper) GrantPermission(ctx sdk.Context, from, to sdk.AccAddress, perm PermissionI) sdk.Error {
	if !k.HasPermission(ctx, from, perm) {
		return ErrTokenPermission(DefaultCodespace)
	}
	k.AddPermission(ctx, to, perm)
	return nil
}

func (k Keeper) MintTokenWithPermission(ctx sdk.Context, amount sdk.Coin, to sdk.AccAddress) sdk.Error {
	if !k.HasPermission(ctx, to, NewMintPermission(amount.Denom)) {
		return ErrTokenPermissionMint(DefaultCodespace)
	}
	return k.mintToken(ctx, amount, to)
}
func (k Keeper) MintTokenWithOutPermission(ctx sdk.Context, amount sdk.Coin, to sdk.AccAddress) sdk.Error {
	return k.mintToken(ctx, amount, to)
}

func (k Keeper) mintToken(ctx sdk.Context, amount sdk.Coin, to sdk.AccAddress) sdk.Error {
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

func (k Keeper) BurnTokenWithPermission(ctx sdk.Context, amount sdk.Coin, from sdk.AccAddress) sdk.Error {
	if !k.HasPermission(ctx, from, NewBurnPermission(amount.Denom)) {
		return ErrTokenPermissionBurn(DefaultCodespace)
	}
	return k.burnToken(ctx, amount, from)

}

func (k Keeper) BurnTokenWithOutPermission(ctx sdk.Context, amount sdk.Coin, from sdk.AccAddress) sdk.Error {
	return k.burnToken(ctx, amount, from)
}

func (k Keeper) burnToken(ctx sdk.Context, amount sdk.Coin, from sdk.AccAddress) sdk.Error {
	if !k.HasPermission(ctx, from, NewBurnPermission(amount.Denom)) {
		return ErrTokenPermissionBurn(DefaultCodespace)
	}

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
