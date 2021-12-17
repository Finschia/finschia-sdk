package consortium

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/consortium/keeper"
	"github.com/line/lbm-sdk/x/consortium/types"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

func NewProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		if !k.GetEnabled(ctx) {
			return nil
		}

		switch c := content.(type) {
		case *types.UpdateConsortiumParamsProposal:
			return HandleUpdateConsortiumParamsProposal(ctx, k, c)
		case *types.UpdateValidatorAuthsProposal:
			return HandleUpdateValidatorAuthsProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized consortium proposal content type: %T", c)
		}
	}
}
