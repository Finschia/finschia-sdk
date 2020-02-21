package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

type IssueKeeper interface {
	IssueCFT(ctx sdk.Context, owner sdk.AccAddress, token types.FT, amount sdk.Int) sdk.Error
	IssueCNFT(ctx sdk.Context, owner sdk.AccAddress, symbol, tokenType string) sdk.Error
}

func (k Keeper) IssueCFT(ctx sdk.Context, owner sdk.AccAddress, token types.FT, amount sdk.Int) sdk.Error {
	if !types.ValidTokenURI(token.GetTokenURI()) {
		return types.ErrInvalidTokenURILength(types.DefaultCodespace, token.GetTokenURI())
	}
	err := k.SetToken(ctx, token)
	if err != nil {
		return err
	}

	err = k.mintTokens(ctx, sdk.NewCoins(sdk.NewCoin(token.GetDenom(), amount)), owner)
	if err != nil {
		return err
	}

	tokenURIModifyPerm := types.NewModifyTokenURIPermission(token.GetDenom())
	k.AddPermission(ctx, owner, tokenURIModifyPerm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueCFT,
			sdk.NewAttribute(types.AttributeKeyName, token.GetName()),
			sdk.NewAttribute(types.AttributeKeySymbol, token.GetSymbol()),
			sdk.NewAttribute(types.AttributeKeyTokenID, token.GetTokenID()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyMintable, strconv.FormatBool(token.GetMintable())),
			sdk.NewAttribute(types.AttributeKeyDecimals, token.GetDecimals().String()),
			sdk.NewAttribute(types.AttributeKeyTokenURI, token.GetTokenURI()),
		),
		sdk.NewEvent(
			types.EventTypeGrantPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
			sdk.NewAttribute(types.AttributeKeyResource, tokenURIModifyPerm.GetResource()),
			sdk.NewAttribute(types.AttributeKeyAction, tokenURIModifyPerm.GetAction()),
		),
	})
	if token.GetMintable() {
		mintPerm := types.NewMintPermission(token.GetDenom())
		k.AddPermission(ctx, owner, mintPerm)
		burnPerm := types.NewBurnPermission(token.GetDenom())
		k.AddPermission(ctx, owner, burnPerm)

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeGrantPermToken,
				sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
				sdk.NewAttribute(types.AttributeKeyResource, mintPerm.GetResource()),
				sdk.NewAttribute(types.AttributeKeyAction, mintPerm.GetAction()),
			),
			sdk.NewEvent(
				types.EventTypeGrantPermToken,
				sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
				sdk.NewAttribute(types.AttributeKeyResource, burnPerm.GetResource()),
				sdk.NewAttribute(types.AttributeKeyAction, burnPerm.GetAction()),
			),
		})
	}

	return nil
}

func (k Keeper) IssueCNFT(ctx sdk.Context, owner sdk.AccAddress, symbol string) sdk.Error {
	tokenType, err := k.getNextTokenType(ctx, symbol)
	if err != nil {
		return err
	}

	err = k.setTokenType(ctx, symbol, tokenType)
	if err != nil {
		return err
	}

	mintPerm := types.NewMintPermission(symbol + tokenType)
	k.AddPermission(ctx, owner, mintPerm)
	burnPerm := types.NewBurnPermission(symbol + tokenType)
	k.AddPermission(ctx, owner, burnPerm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueCNFT,
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyTokenType, tokenType),
		),
		sdk.NewEvent(
			types.EventTypeGrantPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
			sdk.NewAttribute(types.AttributeKeyResource, mintPerm.GetResource()),
			sdk.NewAttribute(types.AttributeKeyAction, mintPerm.GetAction()),
		),
		sdk.NewEvent(
			types.EventTypeGrantPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
			sdk.NewAttribute(types.AttributeKeyResource, burnPerm.GetResource()),
			sdk.NewAttribute(types.AttributeKeyAction, burnPerm.GetAction()),
		),
	})

	return nil
}

func (k Keeper) setTokenType(ctx sdk.Context, symbol, tokenType string) sdk.Error {
	collection, err := k.GetCollection(ctx, symbol)
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	if store.Has(types.TokenTypeKey(collection.GetSymbol(), tokenType)) {
		return types.ErrCollectionTokenTypeExist(types.DefaultCodespace, collection.GetSymbol(), tokenType)
	}
	store.Set(types.TokenTypeKey(collection.GetSymbol(), tokenType), k.cdc.MustMarshalBinaryBare(tokenType))
	return nil
}

func (k Keeper) hasTokenType(ctx sdk.Context, symbol, tokenType string) bool {
	collection, err := k.GetCollection(ctx, symbol)
	if err != nil {
		return false
	}
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.TokenTypeKey(collection.GetSymbol(), tokenType))
}

func (k Keeper) getNextTokenType(ctx sdk.Context, symbol string) (tokenType string, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStoreReversePrefixIterator(store, types.TokenTypeKey(symbol, ""))
	defer iter.Close()
	if !iter.Valid() {
		return types.SmallestNFTType, nil
	}
	k.cdc.MustUnmarshalBinaryBare(iter.Value(), &tokenType)
	tokenType = types.NextID(tokenType, "")
	if tokenType[0] == types.FungibleFlag[0] {
		return "", types.ErrCollectionTokenTypeFull(types.DefaultCodespace, symbol)
	}
	return tokenType, nil
}
