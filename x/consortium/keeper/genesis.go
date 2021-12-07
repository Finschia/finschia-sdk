package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/consortium/types"

	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, sk types.StakingKeeper, data *types.GenesisState) error {
	if !data.Enabled {
		return nil
	}

	k.SetEnabled(ctx, data.Enabled)

	allowedValidators := []sdk.ValAddress{}
	if len(data.AllowedValidators) != 0 {
		// Import the allowed validators from the previous chain data.
		for _, addr := range data.AllowedValidators {
			if err := sdk.ValidateValAddress(addr); err != nil {
				return err
			}
			allowedValidators = append(allowedValidators, sdk.ValAddress(addr))
		}
	} else {
		// Allowed validators must exist if the module is enabled,
		// so it should be the very first block of the chain.
		// We gather the information from staking module.
		sk.IterateValidators(ctx, func(_ int64, addr stakingtypes.ValidatorI) (stop bool) {
			allowedValidators = append(allowedValidators, addr.GetOperator())
			return false
		})
	}

	for _, addr := range allowedValidators {
		k.SetAllowedValidator(ctx, addr, true)
	}

	return nil
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	enabled := k.GetEnabled(ctx)

	allowedValidators := []string{}
	k.IterateAllowedValidators(ctx, func(valAddr sdk.ValAddress) (stop bool) {
		allowedValidators = append(allowedValidators, valAddr.String())
		return false
	})

	return &types.GenesisState{
		Enabled:           enabled,
		AllowedValidators: allowedValidators,
	}
}
