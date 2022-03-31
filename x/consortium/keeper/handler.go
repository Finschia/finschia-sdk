package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/consortium"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

func NewProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		if !k.GetEnabled(ctx) {
			return nil
		}

		switch c := content.(type) {
		case *consortium.UpdateConsortiumParamsProposal:
			return k.handleUpdateConsortiumParamsProposal(ctx, c)
		case *consortium.UpdateValidatorAuthsProposal:
			return k.handleUpdateValidatorAuthsProposal(ctx, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized consortium proposal content type: %T", c)
		}
	}
}
