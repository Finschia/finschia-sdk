package fswap

import (
	"fmt"

	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/fswap/keeper"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
	govtypes "github.com/Finschia/finschia-sdk/x/gov/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	// this line is used by starport scaffolding # handler/msgServer

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		// ctx = ctx.WithEventManager(sdk.NewEventManager())

		// switch msg := msg.(type) {
		// this line is used by starport scaffolding # 1
		// default:
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		// }
	}
}

// NewSwapInitHandler creates a governance handler to manage new proposal types.
// It enables SwapInit to propose an fswap init
func NewSwapInitHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.SwapInitProposal:
			return handleSwapInit(ctx, k, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized sawp proposal content type: %T", c)
		}
	}
}

func handleSwapInit(ctx sdk.Context, k keeper.Keeper, p *types.SwapInitProposal) error {
	return k.SwapInit(ctx, p.SwapInit)
}
