package consortium

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/consortium/types"
	"github.com/line/lbm-sdk/x/consortium/keeper"

	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, sk types.StakingKeeper, data *types.GenesisState) error {
	keeper.SetParams(ctx, data.Params)

	validatorAuths := data.ValidatorAuths
	if keeper.GetEnabled(ctx) && len(validatorAuths) == 0 {
		// Allowed validators must exist if the module is enabled,
		// so it should be the very first block of the chain.
		// We gather the information from staking module.
		sk.IterateValidators(ctx, func(_ int64, addr stakingtypes.ValidatorI) (stop bool) {
			auth := &types.ValidatorAuth{
				OperatorAddress: addr.GetOperator().String(),
				CreationAllowed: true,
			}
			validatorAuths = append(validatorAuths, auth)
			return false
		})
	}

	for _, auth := range validatorAuths {
		keeper.SetValidatorAuth(ctx, auth)
	}

	return nil
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:         keeper.GetParams(ctx),
		ValidatorAuths: keeper.GetValidatorAuths(ctx),
	}
}
