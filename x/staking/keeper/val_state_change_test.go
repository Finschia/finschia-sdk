package keeper_test

import (
	"github.com/line/lbm-sdk/x/staking/keeper"
	"github.com/line/lbm-sdk/x/staking/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnbondingToUnbondedPanic(t *testing.T) {
	app, ctx, _, _, validators := initValidators(t, 100, 2, []int64{0, 100})

	validators[0] = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validators[0], false)
	validators[1] = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validators[1], false)
	applyValidatorSetUpdates(t, ctx, app.StakingKeeper, 1)

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
