package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
	"github.com/line/link/x/collection/internal/types"
)

func (k Keeper) Modify(ctx sdk.Context, owner sdk.AccAddress, symbol, tokenType, tokenIndex string,
	change linktype.Change) sdk.Error {
	if tokenType != "" {
		if tokenIndex != "" {
			return k.modifyToken(ctx, owner, symbol, tokenType+tokenIndex, change)
		}
		return k.modifyTokenType(ctx, owner, symbol, tokenType, change)
	}
	if tokenIndex == "" {
		return k.modifyCollection(ctx, owner, symbol, change)
	}
	return types.ErrTokenIndexWithoutType(types.DefaultCodespace)
}

//nolint:dupl
func (k Keeper) modifyCollection(ctx sdk.Context, owner sdk.AccAddress, symbol string,
	change linktype.Change) sdk.Error {
	collection, err := k.GetCollection(ctx, symbol)
	if err != nil {
		return err
	}
	modifyPerm := types.NewModifyPermission(symbol)
	if !k.HasPermission(ctx, owner, modifyPerm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, owner, modifyPerm)
	}

	switch change.Field {
	case types.AttributeKeyName:
		collection = collection.SetName(change.Value)
	case types.AttributeKeyBaseImgURI:
		collection = collection.SetBaseImgURI(change.Value)
	default:
		return types.ErrInvalidChangesField(types.DefaultCodespace, change.Field)
	}

	err = k.UpdateCollection(ctx, collection)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeModifyCollection,
			sdk.NewAttribute(types.AttributeKeyModifiedField, change.Field),
			sdk.NewAttribute(types.AttributeKeyName, collection.GetName()),
			sdk.NewAttribute(types.AttributeKeySymbol, collection.GetSymbol()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyBaseImgURI, collection.GetBaseImgURI()),
		),
	})
	return nil
}

//nolint:dupl
func (k Keeper) modifyTokenType(ctx sdk.Context, owner sdk.AccAddress, symbol, tokenTypeID string,
	change linktype.Change) sdk.Error {
	tokenType, err := k.GetTokenType(ctx, symbol, tokenTypeID)
	if err != nil {
		return err
	}
	modifyPerm := types.NewModifyPermission(symbol)
	if !k.HasPermission(ctx, owner, modifyPerm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, owner, modifyPerm)
	}

	switch change.Field {
	case types.AttributeKeyName:
		tokenType = tokenType.SetName(change.Value)
	default:
		return types.ErrInvalidChangesField(types.DefaultCodespace, change.Field)
	}

	err = k.UpdateTokenType(ctx, symbol, tokenType)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeModifyTokenType,
			sdk.NewAttribute(types.AttributeKeyModifiedField, change.Field),
			sdk.NewAttribute(types.AttributeKeyName, tokenType.GetName()),
			sdk.NewAttribute(types.AttributeKeySymbol, tokenType.GetSymbol()),
			sdk.NewAttribute(types.AttributeKeyTokenType, tokenType.GetTokenType()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
		),
	})
	return nil
}

//nolint:dupl
func (k Keeper) modifyToken(ctx sdk.Context, owner sdk.AccAddress, symbol, tokenID string,
	change linktype.Change) sdk.Error {
	token, err := k.GetToken(ctx, symbol, tokenID)
	if err != nil {
		return err
	}
	modifyPerm := types.NewModifyPermission(symbol)
	if !k.HasPermission(ctx, owner, modifyPerm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, owner, modifyPerm)
	}

	switch change.Field {
	case types.AttributeKeyName:
		token = token.SetName(change.Value)
	default:
		return types.ErrInvalidChangesField(types.DefaultCodespace, change.Field)
	}

	err = k.UpdateToken(ctx, symbol, token)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeModifyToken,
			sdk.NewAttribute(types.AttributeKeyModifiedField, change.Field),
			sdk.NewAttribute(types.AttributeKeyName, token.GetName()),
			sdk.NewAttribute(types.AttributeKeySymbol, token.GetSymbol()),
			sdk.NewAttribute(types.AttributeKeyTokenID, token.GetTokenID()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
		),
	})
	return nil
}
