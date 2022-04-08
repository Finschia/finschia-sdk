package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

// handleUpdateFoundationParamsProposal is a handler for update foundation params proposal
func (k Keeper) handleUpdateFoundationParamsProposal(ctx sdk.Context, p *foundation.UpdateFoundationParamsProposal) error {
	params := p.Params
	k.SetParams(ctx, params)

	if !params.Enabled {
		k.Cleanup(ctx)
	}

	return ctx.EventManager().EmitTypedEvent(&foundation.EventUpdateFoundationParams{
		Params: params,
	})
}

// handleUpdateValidatorAuthsProposal is a handler for update validator auths proposal
func (k Keeper) handleUpdateValidatorAuthsProposal(ctx sdk.Context, p *foundation.UpdateValidatorAuthsProposal) error {
	for _, auth := range p.Auths {
		if err := k.SetValidatorAuth(ctx, auth); err != nil {
			return err
		}
	}

	return ctx.EventManager().EmitTypedEvent(&foundation.EventUpdateValidatorAuths{
		Auths: p.Auths,
	})
}
