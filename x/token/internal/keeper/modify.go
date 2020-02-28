package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) ModifyToken(ctx sdk.Context, owner sdk.AccAddress, contractID string,
	change linktype.Change) sdk.Error {
	token, err := k.GetToken(ctx, contractID)
	if err != nil {
		return err
	}

	tokenModifyPerm := types.NewModifyPermission(token.GetContractID())
	if !k.HasPermission(ctx, owner, tokenModifyPerm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, owner, tokenModifyPerm)
	}

	switch change.Field {
	case types.AttributeKeyName:
		token = token.SetName(change.Value)
	case types.AttributeKeyTokenURI:
		token = token.SetImageURI(change.Value)
	}

	err = k.UpdateToken(ctx, token)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeModifyToken,
			sdk.NewAttribute(types.AttributeKeyContractID, token.GetContractID()),
			sdk.NewAttribute(types.AttributeKeyModifiedField, change.Field),
			sdk.NewAttribute(types.AttributeKeyName, token.GetName()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyTokenURI, token.GetImageURI()),
		),
	})
	return nil
}
