package keeper_test

import (
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/staking/teststaking"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnbondingToUnbondedPanic(t *testing.T) {
	_, app, ctx := createTestInput()

	//create a validator
	addrDels := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(0))
	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)

	validator := teststaking.NewValidator(t, addrVals[0], PKs[0])
	app.StakingKeeper.SetValidatorByConsAddr(ctx, validator)

	// unbond the validator
	require.Panics(t, func() {
		app.StakingKeeper.UnbondingToUnbonded(ctx, validator)
	})
}
