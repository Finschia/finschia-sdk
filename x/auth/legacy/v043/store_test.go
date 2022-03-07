package v043_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	ocproto "github.com/line/ostracon/proto/ostracon/types"

	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	authkeeper "github.com/line/lbm-sdk/x/auth/keeper"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	"github.com/line/lbm-sdk/x/auth/vesting/exported"
	"github.com/line/lbm-sdk/x/auth/vesting/types"
	"github.com/line/lbm-sdk/x/staking"
	stakingkeeper "github.com/line/lbm-sdk/x/staking/keeper"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

func TestMigrateVestingAccounts(t *testing.T) {
	testCases := []struct {
		name        string
		prepareFunc func(app *simapp.SimApp, ctx sdk.Context, validator stakingtypes.Validator, delegatorAddr sdk.AccAddress)
		garbageFunc func(ctx sdk.Context, vesting exported.VestingAccount, app *simapp.SimApp) error
		tokenAmount int64
		expVested   int64
		expFree     int64
		blockTime   int64
	}{
		{
			"delayed vesting has vested, multiple delegations less than the total account balance",
			func(app *simapp.SimApp, ctx sdk.Context, validator stakingtypes.Validator, delegatorAddr sdk.AccAddress) {

				baseAccount := authtypes.NewBaseAccountWithAddress(delegatorAddr)
				vestedCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), sdk.NewInt(200)))
				delayedAccount := types.NewDelayedVestingAccount(baseAccount, vestedCoins, ctx.BlockTime().Unix())

				ctx = ctx.WithBlockTime(ctx.BlockTime().AddDate(1, 0, 0))

				app.AccountKeeper.SetAccount(ctx, delayedAccount)

				_, err := app.StakingKeeper.Delegate(ctx, delegatorAddr, sdk.NewInt(100), stakingtypes.Unbonded, validator, true)
				require.NoError(t, err)
				_, err = app.StakingKeeper.Delegate(ctx, delegatorAddr, sdk.NewInt(100), stakingtypes.Unbonded, validator, true)
				require.NoError(t, err)
				_, err = app.StakingKeeper.Delegate(ctx, delegatorAddr, sdk.NewInt(100), stakingtypes.Unbonded, validator, true)
				require.NoError(t, err)
			},
			cleartTrackingFields,
			300,
			0,
			300,
			0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.BaseApp.NewContext(false, ocproto.Header{
				Time: time.Now(),
			})

			addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(tc.tokenAmount))
			delegatorAddr := addrs[0]

			_, valAddr := createValidator(t, ctx, app, tc.tokenAmount*2)
			validator, found := app.StakingKeeper.GetValidator(ctx, valAddr)
			require.True(t, found)

			tc.prepareFunc(app, ctx, validator, delegatorAddr)

			if tc.blockTime != 0 {
				ctx = ctx.WithBlockTime(time.Unix(tc.blockTime, 0))
			}

			// We introduce the bug
			savedAccount := app.AccountKeeper.GetAccount(ctx, delegatorAddr)
			vestingAccount, ok := savedAccount.(exported.VestingAccount)
			require.True(t, ok)
			require.NoError(t, tc.garbageFunc(ctx, vestingAccount, app))

			m := authkeeper.NewMigrator(app.AccountKeeper, app.GRPCQueryRouter())
			require.NoError(t, m.MigrateV040ToV043(ctx))

			var expVested sdk.Coins
			var expFree sdk.Coins

			if tc.expVested != 0 {
				expVested = sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), sdk.NewInt(tc.expVested)))
			}

			if tc.expFree != 0 {
				expFree = sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), sdk.NewInt(tc.expFree)))
			}

			trackingCorrected(
				ctx,
				t,
				app.AccountKeeper,
				savedAccount.GetAddress(),
				expVested,
				expFree,
			)
		})
	}

}

func trackingCorrected(ctx sdk.Context, t *testing.T, ak authkeeper.AccountKeeper, addr sdk.AccAddress, expDelVesting sdk.Coins, expDelFree sdk.Coins) {
	t.Helper()
	baseAccount := ak.GetAccount(ctx, addr)
	vDA, ok := baseAccount.(exported.VestingAccount)
	require.True(t, ok)

	vestedOk := expDelVesting.IsEqual(vDA.GetDelegatedVesting())
	freeOk := expDelFree.IsEqual(vDA.GetDelegatedFree())
	require.True(t, vestedOk, vDA.GetDelegatedVesting().String())
	require.True(t, freeOk, vDA.GetDelegatedFree().String())
}

func cleartTrackingFields(ctx sdk.Context, vesting exported.VestingAccount, app *simapp.SimApp) error {
	switch t := vesting.(type) {
	case *types.DelayedVestingAccount:
		t.DelegatedFree = nil
		t.DelegatedVesting = nil
		app.AccountKeeper.SetAccount(ctx, t)
	case *types.ContinuousVestingAccount:
		t.DelegatedFree = nil
		t.DelegatedVesting = nil
		app.AccountKeeper.SetAccount(ctx, t)
	case *types.PeriodicVestingAccount:
		t.DelegatedFree = nil
		t.DelegatedVesting = nil
		app.AccountKeeper.SetAccount(ctx, t)
	default:
		return fmt.Errorf("expected vesting account, found %t", t)
	}

	return nil
}

func dirtyTrackingFields(ctx sdk.Context, vesting exported.VestingAccount, app *simapp.SimApp) error {
	dirt := sdk.NewCoins(sdk.NewInt64Coin("stake", 42))

	switch t := vesting.(type) {
	case *types.DelayedVestingAccount:
		t.DelegatedFree = dirt
		t.DelegatedVesting = dirt
		app.AccountKeeper.SetAccount(ctx, t)
	case *types.ContinuousVestingAccount:
		t.DelegatedFree = dirt
		t.DelegatedVesting = dirt
		app.AccountKeeper.SetAccount(ctx, t)
	case *types.PeriodicVestingAccount:
		t.DelegatedFree = dirt
		t.DelegatedVesting = dirt
		app.AccountKeeper.SetAccount(ctx, t)
	default:
		return fmt.Errorf("expected vesting account, found %t", t)
	}

	return nil
}

func createValidator(t *testing.T, ctx sdk.Context, app *simapp.SimApp, powers int64) (sdk.AccAddress, sdk.ValAddress) {
	valTokens := sdk.TokensFromConsensusPower(powers, sdk.DefaultPowerReduction)
	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, valTokens)
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
	pks := simapp.CreateTestPubKeys(1)
	cdc := simapp.MakeTestEncodingConfig().Marshaler

	app.StakingKeeper = stakingkeeper.NewKeeper(
		cdc,
		app.GetKey(stakingtypes.StoreKey),
		app.AccountKeeper,
		app.BankKeeper,
		app.GetSubspace(stakingtypes.ModuleName),
	)

	val1, err := stakingtypes.NewValidator(valAddrs[0], pks[0], stakingtypes.Description{})
	require.NoError(t, err)

	app.StakingKeeper.SetValidator(ctx, val1)
	require.NoError(t, app.StakingKeeper.SetValidatorByConsAddr(ctx, val1))
	app.StakingKeeper.SetNewValidatorByPowerIndex(ctx, val1)

	_, err = app.StakingKeeper.Delegate(ctx, addrs[0], valTokens, stakingtypes.Unbonded, val1, true)
	require.NoError(t, err)

	_ = staking.EndBlocker(ctx, app.StakingKeeper)

	return addrs[0], valAddrs[0]
}
