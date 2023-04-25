package keeper_test

import (
	"github.com/Finschia/finschia-sdk/x/staking/keeper"
	"github.com/Finschia/finschia-sdk/x/staking/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnbondingToUnbondedPanic(t *testing.T) {
	app, ctx, _, _, validators := initValidators(t, 100, 2, []int64{0, 100})

	for i, validator := range validators {
		validators[i] = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, false)
	}

	assert.Equal(t, validators[0].Status, types.Unbonded)
	assert.Equal(t, validators[1].Status, types.Bonded)

	// unbond validator which is in unbonded status
	require.Panics(t, func() {
		app.StakingKeeper.UnbondingToUnbonded(ctx, validators[0])
	})

	// unbond validator which is in bonded status
	require.Panics(t, func() {
		app.StakingKeeper.UnbondingToUnbonded(ctx, validators[1])
	})
}
