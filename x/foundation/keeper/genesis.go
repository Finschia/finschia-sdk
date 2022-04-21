package keeper

import (
	sdk "github.com/line/lbm-sdk/types"

	"github.com/line/lbm-sdk/x/foundation"

	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, sk foundation.StakingKeeper, data *foundation.GenesisState) error {
	k.SetParams(ctx, data.Params)

	validatorAuths := data.ValidatorAuths
	if k.GetEnabled(ctx) && len(validatorAuths) == 0 {
		// Allowed validators must exist if the module is enabled,
		// so it should be the very first block of the chain.
		// We gather the information from staking module.
		sk.IterateValidators(ctx, func(_ int64, addr stakingtypes.ValidatorI) (stop bool) {
			auth := foundation.ValidatorAuth{
				OperatorAddress: addr.GetOperator().String(),
				CreationAllowed: true,
			}
			validatorAuths = append(validatorAuths, auth)
			return false
		})
	}

	for _, auth := range validatorAuths {
		if err := k.SetValidatorAuth(ctx, auth); err != nil {
			return err
		}
	}

	// TODO: check null
	k.setFoundationInfo(ctx, foundation.FoundationInfo{
		Version: 1,
		Operator: string(k.GetAdmin(ctx)),
	})

	return nil
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *foundation.GenesisState {
	return &foundation.GenesisState{
		Params:         k.GetParams(ctx),
		ValidatorAuths: k.GetValidatorAuths(ctx),
	}
}
