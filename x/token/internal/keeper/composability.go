package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

var ChildExists = []byte{1}

func (k Keeper) Attach(ctx sdk.Context, from sdk.AccAddress, symbol string, toTokenID string, tokenID string) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	if toTokenID == tokenID {
		return types.ErrCannotAttachToItself(types.DefaultCodespace, symbol+tokenID)
	}

	token, err := k.getIDNFT(ctx, symbol+tokenID)
	if err != nil {
		return err
	}

	if !from.Equals(token.Owner) {
		return types.ErrTokenNotOwnedBy(types.DefaultCodespace, token.GetDenom(), from)
	}

	// verify token should be a root
	childToParentKey := types.TokenChildToParentKey(token)
	if store.Has(childToParentKey) {
		return types.ErrTokenAlreadyAChild(types.DefaultCodespace, token.GetDenom())
	}

	// verify no circulation(toToken must not be a descendant of token)
	rootOfToToken, err := k.RootOf(ctx, symbol, toTokenID)
	if err != nil {
		return err
	}
	toToken, err := k.getIDNFT(ctx, symbol+toTokenID)
	if err != nil {
		return err
	}
	if rootOfToToken != nil && rootOfToToken.GetTokenID() == tokenID {
		return types.ErrCannotAttachToADescendant(types.DefaultCodespace, token.GetDenom(), toToken.GetDenom())
	}

	parentToChildKey := types.TokenParentToChildKey(toToken)
	childrenStore := prefix.NewStore(store, parentToChildKey)
	parentToChildSubKey := types.TokenParentToChildSubKey(token)
	if childrenStore.Has(parentToChildSubKey) {
		panic("token is already a child of some other")
	}

	if !from.Equals(toToken.Owner) {
		if err := k.moveIDNFToken(ctx, store, from, toToken.Owner, token); err != nil {
			return err
		}
	}

	store.Set(childToParentKey, k.encodeToken(toToken))
	childrenStore.Set(parentToChildSubKey, ChildExists)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAttachToken,
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyToTokenID, toTokenID),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
		),
	})

	return nil
}

func (k Keeper) Detach(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	token, err := k.getIDNFT(ctx, symbol+tokenID)
	if err != nil {
		return err
	}

	if !from.Equals(token.Owner) {
		return types.ErrTokenNotOwnedBy(types.DefaultCodespace, token.GetDenom(), from)
	}

	childToParentKey := types.TokenChildToParentKey(token)
	if !store.Has(childToParentKey) {
		return types.ErrTokenNotAChild(types.DefaultCodespace, token.GetDenom())
	}

	bz := store.Get(childToParentKey)
	parentToken, err := k.decodeToken(ctx, bz)
	if err != nil {
		return err
	}

	parentToChildKey := types.TokenParentToChildKey(parentToken)
	childrenStore := prefix.NewStore(store, parentToChildKey)
	parentToChildSubKey := types.TokenParentToChildSubKey(token)
	if !childrenStore.Has(parentToChildSubKey) {
		panic("token is not a child of some other")
	}

	if !from.Equals(to) {
		if err := k.moveIDNFToken(ctx, store, from, to, token); err != nil {
			return err
		}
	}
	store.Delete(childToParentKey)
	childrenStore.Delete(parentToChildSubKey)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDetachToken,
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
		),
	})

	return nil
}

func (k Keeper) RootOf(ctx sdk.Context, symbol string, tokenID string) (*types.BaseIDNFT, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	token, err := k.getIDNFT(ctx, symbol+tokenID)
	if err != nil {
		return nil, err
	}
	myself := token

	for childToParentKey := types.TokenChildToParentKey(token); store.Has(childToParentKey); {
		bz := store.Get(childToParentKey)
		token, err = k.decodeToken(ctx, bz)
		if err != nil {
			return nil, err
		}
		childToParentKey = types.TokenChildToParentKey(token)
	}

	if token == myself {
		return nil, nil
	}
	return token, nil
}

func (k Keeper) ParentOf(ctx sdk.Context, symbol string, tokenID string) (*types.BaseIDNFT, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	token, err := k.getIDNFT(ctx, symbol+tokenID)
	if err != nil {
		return nil, err
	}
	childToParentKey := types.TokenChildToParentKey(token)
	if store.Has(childToParentKey) {
		bz := store.Get(childToParentKey)
		return k.decodeToken(ctx, bz)
	}
	return nil, nil
}

func (k Keeper) ChildrenOf(ctx sdk.Context, symbol string, tokenID string) (types.Tokens, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	token, err := k.getIDNFT(ctx, symbol+tokenID)
	if err != nil {
		return nil, err
	}

	result := make([]types.Token, 0)
	parentToChildKey := types.TokenParentToChildKey(token)
	iter := sdk.KVStorePrefixIterator(store, parentToChildKey)
	defer iter.Close()
	for {
		if !iter.Valid() {
			break
		}
		tokenDenom := types.ParentToChildSubKeyToToken(parentToChildKey, iter.Key())
		childToken, err := k.getIDNFT(ctx, tokenDenom)
		if err != nil {
			panic("child token does not exist")
		}
		result = append(result, childToken)
		iter.Next()
	}
	return result, nil
}

func (k Keeper) getIDNFT(ctx sdk.Context, tokenDenom string) (*types.BaseIDNFT, sdk.Error) {
	token, err := k.GetToken(ctx, tokenDenom)
	if err != nil {
		return nil, err
	}
	tokenNFT, ok := token.(*types.BaseIDNFT)
	if !ok {
		return nil, types.ErrTokenNotIDNF(types.DefaultCodespace, token.GetDenom())
	}
	return tokenNFT, nil
}

func (k Keeper) encodeToken(token *types.BaseIDNFT) []byte {
	return k.cdc.MustMarshalBinaryBare(token.GetDenom())
}

func (k Keeper) decodeToken(ctx sdk.Context, tokenByte []byte) (*types.BaseIDNFT, sdk.Error) {
	var tokenDenom string
	k.cdc.MustUnmarshalBinaryBare(tokenByte, &tokenDenom)
	return k.getIDNFT(ctx, tokenDenom)
}

func (k Keeper) decodeOwner(tokenByte []byte) (owner sdk.AccAddress) {
	err := k.cdc.UnmarshalBinaryBare(tokenByte, &owner)
	if err != nil {
		panic(err)
	}
	return owner
}
