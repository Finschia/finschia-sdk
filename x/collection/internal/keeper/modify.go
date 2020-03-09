package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
	"github.com/line/link/x/collection/internal/types"
)

func (k Keeper) Modify(ctx sdk.Context, owner sdk.AccAddress, contractID, tokenType, tokenIndex string,
	changes linktype.Changes) sdk.Error {
	if tokenType != "" {
		if tokenIndex != "" {
			return k.modifyToken(ctx, owner, contractID, tokenType+tokenIndex, changes)
		}
		return k.modifyTokenType(ctx, owner, contractID, tokenType, changes)
	}
	if tokenIndex == "" {
		return k.modifyCollection(ctx, owner, contractID, changes)
	}
	return types.ErrTokenIndexWithoutType(types.DefaultCodespace)
}

//nolint:dupl
func (k Keeper) modifyCollection(ctx sdk.Context, owner sdk.AccAddress, contractID string,
	changes linktype.Changes) sdk.Error {
	collection, err := k.GetCollection(ctx, contractID)
	if err != nil {
		return err
	}
	modifyPerm := types.NewModifyPermission(contractID)
	if !k.HasPermission(ctx, owner, modifyPerm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, owner, modifyPerm)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeModifyCollection,
			sdk.NewAttribute(types.AttributeKeyContractID, collection.GetContractID()),
		),
	})

	for _, change := range changes {
		switch change.Field {
		case types.AttributeKeyName:
			collection.SetName(change.Value)
		case types.AttributeKeyMeta:
			collection.SetMeta(change.Value)
		case types.AttributeKeyBaseImgURI:
			collection.SetBaseImgURI(change.Value)
		default:
			return types.ErrInvalidChangesField(types.DefaultCodespace, change.Field)
		}

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeModifyCollection,
				sdk.NewAttribute(change.Field, change.Value),
			),
		})
	}
	err = k.UpdateCollection(ctx, collection)
	if err != nil {
		return err
	}
	return nil
}

//nolint:dupl
func (k Keeper) modifyTokenType(ctx sdk.Context, owner sdk.AccAddress, contractID, tokenTypeID string,
	changes linktype.Changes) sdk.Error {
	tokenType, err := k.GetTokenType(ctx, contractID, tokenTypeID)
	if err != nil {
		return err
	}
	modifyPerm := types.NewModifyPermission(contractID)
	if !k.HasPermission(ctx, owner, modifyPerm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, owner, modifyPerm)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeModifyTokenType,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyTokenType, tokenType.GetTokenType()),
		),
	})

	for _, change := range changes {
		switch change.Field {
		case types.AttributeKeyName:
			tokenType.SetName(change.Value)
		case types.AttributeKeyMeta:
			tokenType.SetMeta(change.Value)
		default:
			return types.ErrInvalidChangesField(types.DefaultCodespace, change.Field)
		}

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeModifyTokenType,
				sdk.NewAttribute(change.Field, change.Value),
			),
		})
	}
	err = k.UpdateTokenType(ctx, tokenType)
	if err != nil {
		return err
	}
	return nil
}

//nolint:dupl
func (k Keeper) modifyToken(ctx sdk.Context, owner sdk.AccAddress, contractID, tokenID string,
	changes linktype.Changes) sdk.Error {
	token, err := k.GetToken(ctx, contractID, tokenID)
	if err != nil {
		return err
	}
	modifyPerm := types.NewModifyPermission(contractID)
	if !k.HasPermission(ctx, owner, modifyPerm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, owner, modifyPerm)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeModifyToken,
			sdk.NewAttribute(types.AttributeKeyContractID, token.GetContractID()),
			sdk.NewAttribute(types.AttributeKeyTokenID, token.GetTokenID()),
		),
	})

	for _, change := range changes {
		switch change.Field {
		case types.AttributeKeyName:
			token.SetName(change.Value)
		case types.AttributeKeyMeta:
			token.SetMeta(change.Value)
		default:
			return types.ErrInvalidChangesField(types.DefaultCodespace, change.Field)
		}
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeModifyToken,
				sdk.NewAttribute(change.Field, change.Value),
			),
		})
	}
	err = k.UpdateToken(ctx, token)
	if err != nil {
		return err
	}
	return nil
}
