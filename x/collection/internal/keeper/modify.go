package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
)

func (k Keeper) Modify(ctx sdk.Context, owner sdk.AccAddress, tokenType, tokenIndex string,
	changes types.Changes) error {
	if tokenType != "" {
		if tokenIndex != "" {
			return k.modifyToken(ctx, owner, tokenType+tokenIndex, changes)
		}
		return k.modifyTokenType(ctx, owner, tokenType, changes)
	}
	if tokenIndex == "" {
		return k.modifyCollection(ctx, owner, changes)
	}
	return types.ErrTokenIndexWithoutType
}

// nolint:dupl
func (k Keeper) modifyCollection(ctx sdk.Context, owner sdk.AccAddress, changes types.Changes) error {
	collection, err := k.GetCollection(ctx)
	if err != nil {
		return err
	}
	modifyPerm := types.NewModifyPermission()
	if !k.HasPermission(ctx, owner, modifyPerm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", owner.String(), modifyPerm.String())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeModifyCollection,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
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
			return sdkerrors.Wrapf(types.ErrInvalidChangesField, "Field: %s", change.Field)
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

// nolint:dupl
func (k Keeper) modifyTokenType(ctx sdk.Context, owner sdk.AccAddress, tokenTypeID string,
	changes types.Changes) error {
	tokenType, err := k.GetTokenType(ctx, tokenTypeID)
	if err != nil {
		return err
	}
	modifyPerm := types.NewModifyPermission()
	if !k.HasPermission(ctx, owner, modifyPerm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", owner.String(), modifyPerm.String())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeModifyTokenType,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
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
			return sdkerrors.Wrapf(types.ErrInvalidChangesField, "Field: %s", change.Field)
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

// nolint:dupl
func (k Keeper) modifyToken(ctx sdk.Context, owner sdk.AccAddress, tokenID string,
	changes types.Changes) error {
	token, err := k.GetToken(ctx, tokenID)
	if err != nil {
		return err
	}
	modifyPerm := types.NewModifyPermission()
	if !k.HasPermission(ctx, owner, modifyPerm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", owner.String(), modifyPerm.String())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeModifyToken,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
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
			return sdkerrors.Wrapf(types.ErrInvalidChangesField, "Field: %s", change.Field)
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
