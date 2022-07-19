package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/foundation"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

func NewProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		if !k.GetEnabled(ctx) {
			return nil
		}

		switch c := content.(type) {
		case *foundation.UpdateFoundationParamsProposal:
			return k.handleUpdateFoundationParamsProposal(ctx, c)
		case *foundation.UpdateValidatorAuthsProposal:
			return k.handleUpdateValidatorAuthsProposal(ctx, c)
		default:
			return sdkerrors.ErrUnknownRequest.Wrapf("unrecognized foundation proposal content type: %T", c)
		}
	}
}
