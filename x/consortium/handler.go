package consortium

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/consortium/keeper"
	"github.com/line/lbm-sdk/x/consortium/types"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

// NewHandler creates an sdk.Handler for all the consortium type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	// msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		if !k.GetEnabled(ctx) {
			return nil, nil
		}

		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func NewProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		if !k.GetEnabled(ctx) {
			return nil
		}

		switch c := content.(type) {
		case *types.DisableConsortiumProposal:
			return HandleDisableConsortiumProposal(ctx, k, c)
		case *types.EditAllowedValidatorsProposal:
			return HandleEditAllowedValidatorsProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized staking proposal content type: %T", c)
		}
	}
}
