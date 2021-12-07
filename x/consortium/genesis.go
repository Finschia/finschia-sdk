package consortium

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/consortium/keeper"
	"github.com/line/lbm-sdk/x/consortium/types"

	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, sk types.StakingKeeper, data *types.GenesisState) {
	if !data.Enabled {
		return
	}

	k.SetEnabled(ctx, data.Enabled)

	allowedValidators := []sdk.ValAddress{}
	if len(data.AllowedValidators) != 0 {
		// Import the allowed validators from the previous chain data.
		for _, addr := range data.AllowedValidators {
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
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
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
