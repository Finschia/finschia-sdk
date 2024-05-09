package fswap

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/fswap/keeper"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
	govtypes "github.com/Finschia/finschia-sdk/x/gov/types"
)

// NewSwapHandler creates a governance handler to manage new proposal types.
// It enables Swap to propose a swap init
func NewSwapHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.MakeSwapProposal:
			return handleMakeSwapProposal(ctx, k, c)

		default:
			return sdkerrors.ErrUnknownRequest.Wrapf("unrecognized sawp proposal content type: %T", c)
		}
	}
}

func handleMakeSwapProposal(ctx sdk.Context, k keeper.Keeper, p *types.MakeSwapProposal) error {
	return k.MakeSwap(ctx, p.Swap, p.ToDenomMetadata)
}
