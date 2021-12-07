package consortium

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/consortium/types"
	"github.com/line/lbm-sdk/x/consortium/keeper"
)

// HandleDisableConsortiumProposal is a handler for disable consortium proposal
func HandleDisableConsortiumProposal(ctx sdk.Context, k keeper.Keeper, p *types.DisableConsortiumProposal) error {
	enabled := false
	if !enabled {
		k.SetEnabled(ctx, enabled)
		return nil
	}

	return nil
}

// HandleEditAllowedValidatorsProposal is a handler for edit allowed validators proposal
func HandleEditAllowedValidatorsProposal(ctx sdk.Context, k keeper.Keeper, p *types.EditAllowedValidatorsProposal) error {
	for _, addr := range p.AddingAddresses {
		err := sdk.ValidateValAddress(addr)
		if err != nil {
			return err
		}
		k.SetAllowedValidator(ctx, sdk.ValAddress(addr), true)
	}

	for _, addr := range p.RemovingAddresses {
		err := sdk.ValidateValAddress(addr)
		if err != nil {
			return err
		}
		k.SetAllowedValidator(ctx, sdk.ValAddress(addr), false)
	}

	return nil
}
