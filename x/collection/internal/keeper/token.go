package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

type TokenKeeper interface {
	GetToken(ctx sdk.Context, symbol, tokenID string) (types.Token, sdk.Error)
	SetToken(ctx sdk.Context, token types.Token) sdk.Error
	UpdateToken(ctx sdk.Context, token types.Token) sdk.Error
}

func (k Keeper) GetToken(ctx sdk.Context, symbol, tokenID string) (types.Token, sdk.Error) {
	c, err := k.GetCollection(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.GetToken(tokenID)
}

func (k Keeper) GetNFT(ctx sdk.Context, symbol, tokenID string) (types.NFT, sdk.Error) {
	token, err := k.GetToken(ctx, symbol, tokenID)
	if err != nil {
		return nil, err
	}
	nft, ok := token.(types.NFT)
	if !ok {
		return nil, types.ErrTokenNotCNFT(types.DefaultCodespace, token.GetDenom())
	}
	return nft, nil
}

func (k Keeper) GetFT(ctx sdk.Context, symbol, tokenID string) (types.FT, sdk.Error) {
	token, err := k.GetToken(ctx, symbol, tokenID)
	if err != nil {
		return nil, err
	}
	ft, ok := token.(types.FT)
	if !ok {
		return nil, types.ErrTokenNotCNFT(types.DefaultCodespace, token.GetDenom())
	}
	return ft, nil
}

func (k Keeper) SetToken(ctx sdk.Context, token types.Token) sdk.Error {
	c, err := k.GetCollection(ctx, token.GetSymbol())
	if err != nil {
		return err
	}
	if t, ok := token.(types.NFT); ok {
		tokenType := t.GetTokenType()
		if !k.hasTokenType(ctx, token.GetSymbol(), tokenType) {
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

func (k Keeper) UpdateToken(ctx sdk.Context, token types.Token) sdk.Error {
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

func (k Keeper) ModifyTokenURI(ctx sdk.Context, owner sdk.AccAddress, symbol, tokenID, tokenURI string) sdk.Error {
	if !types.ValidTokenURI(tokenURI) {
		return types.ErrInvalidTokenURILength(types.DefaultCodespace, tokenURI)
	}

	token, err := k.GetToken(ctx, symbol, tokenID)
	if err != nil {
		return err
	}
	tokenURIModifyPerm := types.NewModifyTokenURIPermission(token.GetDenom())
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
			sdk.NewAttribute(types.AttributeKeyTokenID, token.GetTokenID()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyTokenURI, token.GetTokenURI()),
		),
	})
	return nil
}

func (k Keeper) GetSupply(ctx sdk.Context, symbol, tokenID string) (supply sdk.Int, err sdk.Error) {
	if _, err = k.GetToken(ctx, symbol, tokenID); err != nil {
		return sdk.NewInt(0), err
	}
	return k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(symbol + tokenID), nil
}
