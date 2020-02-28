package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) ModifyToken(ctx sdk.Context, owner sdk.AccAddress, contractID string,
	changes linktype.Changes) sdk.Error {
	token, err := k.GetToken(ctx, contractID)
	if err != nil {
		return err
	}

	tokenModifyPerm := types.NewModifyPermission(token.GetContractID())
	if !k.HasPermission(ctx, owner, tokenModifyPerm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, owner, tokenModifyPerm)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeModifyToken,
			sdk.NewAttribute(types.AttributeKeyContractID, token.GetContractID()),
		),
	})

	for _, change := range changes {
		switch change.Field {
		case types.AttributeKeyName:
			token = token.SetName(change.Value)
		case types.AttributeKeyTokenURI:
			token = token.SetImageURI(change.Value)
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
