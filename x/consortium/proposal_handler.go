package consortium

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/consortium/keeper"
	"github.com/line/lbm-sdk/x/consortium/types"
)

// HandleUpdateConsortiumParamsProposal is a handler for update consortium params proposal
func HandleUpdateConsortiumParamsProposal(ctx sdk.Context, k keeper.Keeper, p *types.UpdateConsortiumParamsProposal) error {
	params := p.Params
	k.SetParams(ctx, params)

	if !params.Enabled {
		k.Cleanup(ctx)
	}

	return nil
}

// HandleUpdateValidatorAuthsProposal is a handler for update validator auths proposal
func HandleUpdateValidatorAuthsProposal(ctx sdk.Context, k keeper.Keeper, p *types.UpdateValidatorAuthsProposal) error {
	for _, auth := range p.Auths {
		if err := k.SetValidatorAuth(ctx, auth); err != nil {
			return err
		}
	}

	return nil
}
