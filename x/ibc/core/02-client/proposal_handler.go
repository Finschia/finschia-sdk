package client

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	govtypes "github.com/line/lbm-sdk/x/gov/types"

	"github.com/line/lbm-sdk/x/ibc/core/02-client/keeper"
	"github.com/line/lbm-sdk/x/ibc/core/02-client/types"
)

// NewClientProposalHandler defines the 02-client proposal handler
func NewClientProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.ClientUpdateProposal:
			return k.ClientUpdateProposal(ctx, c)
		case *types.UpgradeProposal:
			return k.HandleUpgradeProposal(ctx, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized ibc proposal content type: %T", c)
		}
	}
}
