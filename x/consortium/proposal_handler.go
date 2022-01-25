package consortium

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/consortium/keeper"
	"github.com/line/lbm-sdk/x/consortium/types"
)

// handleUpdateConsortiumParamsProposal is a handler for update consortium params proposal
func handleUpdateConsortiumParamsProposal(ctx sdk.Context, k keeper.Keeper, p *types.UpdateConsortiumParamsProposal) error {
	params := p.Params
	k.SetParams(ctx, params)

	if !params.Enabled {
		k.Cleanup(ctx)
	}

	return ctx.EventManager().EmitTypedEvent(&types.EventUpdateConsortiumParams{
		Params: params,
	})
}

// handleUpdateValidatorAuthsProposal is a handler for update validator auths proposal
func handleUpdateValidatorAuthsProposal(ctx sdk.Context, k keeper.Keeper, p *types.UpdateValidatorAuthsProposal) error {
	for _, auth := range p.Auths {
		if err := k.SetValidatorAuth(ctx, auth); err != nil {
			return err
		}
	}

	return ctx.EventManager().EmitTypedEvent(&types.EventUpdateValidatorAuths{
		Auths: p.Auths,
	})
}
