package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/token/internal/types"
)

func (k Keeper) ModifyToken(ctx sdk.Context, owner sdk.AccAddress, changes types.Changes) error {
	token, err := k.GetToken(ctx)
	if err != nil {
		return err
	}

	tokenModifyPerm := types.NewModifyPermission()
	if !k.HasPermission(ctx, owner, tokenModifyPerm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", owner.String(), tokenModifyPerm.String())
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
			token.SetName(change.Value)
		case types.AttributeKeyMeta:
			token.SetMeta(change.Value)
		case types.AttributeKeyImageURI:
			token.SetImageURI(change.Value)
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
