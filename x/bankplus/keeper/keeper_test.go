package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	ocproto "github.com/line/ostracon/proto/ostracon/types"

	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	authkeeper "github.com/line/lbm-sdk/x/auth/keeper"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	"github.com/line/lbm-sdk/x/bank/types"
	bankpluskeeper "github.com/line/lbm-sdk/x/bankplus/keeper"
	minttypes "github.com/line/lbm-sdk/x/mint/types"
)

const (
	initialPower = int64(100)
	holder       = "holder"
	blocker      = "blocker"
)

var (
	holderAcc  = authtypes.NewEmptyModuleAccount(holder)
	blockedAcc = authtypes.NewEmptyModuleAccount(blocker)
	burnerAcc  = authtypes.NewEmptyModuleAccount(authtypes.Burner, authtypes.Burner)

	initTokens = sdk.TokensFromConsensusPower(initialPower, sdk.DefaultPowerReduction)
	initCoins  = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))
)

// nolint: interfacer
func getCoinsByName(ctx sdk.Context, bk bankpluskeeper.Keeper, ak types.AccountKeeper, moduleName string) sdk.Coins {
	moduleAddress := ak.GetModuleAddress(moduleName)
	macc := ak.GetAccount(ctx, moduleAddress)
	if macc == nil {
		return sdk.NewCoins()
	}

	return bk.GetAllBalances(ctx, macc.GetAddress())
}

type IntegrationTestSuite struct {
	suite.Suite
}

func (suite *IntegrationTestSuite) TestSupply_SendCoins() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{Height: 1})
	appCodec := app.AppCodec()

	// add module accounts to supply keeper
	maccPerms := simapp.GetMaccPerms()
	maccPerms[holder] = nil
	maccPerms[authtypes.Burner] = []string{authtypes.Burner}
	maccPerms[authtypes.Minter] = []string{authtypes.Minter}

	authKeeper := authkeeper.NewAccountKeeper(
		appCodec, app.GetKey(types.StoreKey), app.GetSubspace(types.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms,
	)
	keeper := bankpluskeeper.NewBaseKeeper(
		appCodec, app.GetKey(types.StoreKey), authKeeper,
		app.GetSubspace(types.ModuleName), make(map[string]bool), false,
	)

	baseAcc := authKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("baseAcc"))

	// set initial balances
	suite.
		Require().
		NoError(keeper.MintCoins(ctx, minttypes.ModuleName, initCoins))

	suite.
		Require().
		NoError(keeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, holderAcc.GetAddress(), initCoins))

	suite.Require().NoError(keeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	authKeeper.SetModuleAccount(ctx, holderAcc)
	authKeeper.SetModuleAccount(ctx, burnerAcc)
	authKeeper.SetAccount(ctx, baseAcc)

	suite.Require().Panics(func() {
		keeper.SendCoinsFromModuleToModule(ctx, "", holderAcc.GetName(), initCoins) // nolint:errcheck
	})

	suite.Require().Panics(func() {
		keeper.SendCoinsFromModuleToModule(ctx, authtypes.Burner, "", initCoins) // nolint:errcheck
	})

	suite.Require().Panics(func() {
		keeper.SendCoinsFromModuleToAccount(ctx, "", baseAcc.GetAddress(), initCoins) // nolint:errcheck
	})

	// not enough balance (100stake - 200stake)
	suite.Require().Error(
		keeper.SendCoinsFromModuleToAccount(ctx, holderAcc.GetName(), baseAcc.GetAddress(), initCoins.Add(initCoins...)),
	)

	suite.Require().NoError(
		keeper.SendCoinsFromModuleToModule(ctx, holderAcc.GetName(), authtypes.Burner, initCoins),
	)
	suite.Require().Equal(sdk.NewCoins().String(), getCoinsByName(ctx, keeper, authKeeper, holderAcc.GetName()).String())
	suite.Require().Equal(initCoins, getCoinsByName(ctx, keeper, authKeeper, authtypes.Burner))

	suite.Require().NoError(
		keeper.SendCoinsFromModuleToAccount(ctx, authtypes.Burner, baseAcc.GetAddress(), initCoins),
	)
	suite.Require().Equal(sdk.NewCoins().String(), getCoinsByName(ctx, keeper, authKeeper, authtypes.Burner).String())
	suite.Require().Equal(initCoins, keeper.GetAllBalances(ctx, baseAcc.GetAddress()))

	suite.Require().NoError(keeper.SendCoinsFromAccountToModule(ctx, baseAcc.GetAddress(), authtypes.Burner, initCoins))
	suite.Require().Equal(sdk.NewCoins().String(), keeper.GetAllBalances(ctx, baseAcc.GetAddress()).String())
	suite.Require().Equal(initCoins, getCoinsByName(ctx, keeper, authKeeper, authtypes.Burner))
}

func (suite *IntegrationTestSuite) TestInactiveAddrOfSendCoins() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{Height: 1})
	appCodec := app.AppCodec()

	// add module accounts to supply keeper
	maccPerms := simapp.GetMaccPerms()
	maccPerms[holder] = nil
	maccPerms[authtypes.Burner] = []string{authtypes.Burner}
	maccPerms[authtypes.Minter] = []string{authtypes.Minter}

	authKeeper := authkeeper.NewAccountKeeper(
		appCodec, app.GetKey(types.StoreKey), app.GetSubspace(types.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms,
	)

	keeper := bankpluskeeper.NewBaseKeeper(
		appCodec, app.GetKey(types.StoreKey), authKeeper,
		app.GetSubspace(types.ModuleName), make(map[string]bool), false,
	)

	// set initial balances
	suite.
		Require().
		NoError(keeper.MintCoins(ctx, minttypes.ModuleName, initCoins))

	suite.
		Require().
		NoError(keeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, holderAcc.GetAddress(), initCoins))

	suite.Require().NoError(keeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	suite.Require().Equal(initCoins, keeper.GetAllBalances(ctx, holderAcc.GetAddress()))

	suite.Require().False(keeper.IsInactiveAddr(blockedAcc.GetAddress()))

	// add blocked address
	keeper.AddToInactiveAddr(ctx, blockedAcc.GetAddress())
	suite.Require().True(keeper.IsInactiveAddr(blockedAcc.GetAddress()))

	err := keeper.SendCoins(ctx, holderAcc.GetAddress(), blockedAcc.GetAddress(), initCoins)
	suite.Require().Contains(err.Error(), "is not allowed to receive funds")
	suite.Require().Equal(initCoins, keeper.GetAllBalances(ctx, holderAcc.GetAddress()))

	// delete blocked address
	keeper.DeleteFromInactiveAddr(ctx, blockedAcc.GetAddress())
	suite.Require().False(keeper.IsInactiveAddr(blockedAcc.GetAddress()))

	suite.Require().NoError(keeper.SendCoins(ctx, holderAcc.GetAddress(), blockedAcc.GetAddress(), initCoins))
	suite.Require().Equal(sdk.NewCoins().String(), keeper.GetAllBalances(ctx, holderAcc.GetAddress()).String())
}

func (suite *IntegrationTestSuite) TestInitializeBankPlus() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{Height: 1})
	appCodec := app.AppCodec()

	// add module accounts to supply keeper
	maccPerms := simapp.GetMaccPerms()
	maccPerms[holder] = nil
	maccPerms[authtypes.Burner] = []string{authtypes.Burner}
	maccPerms[authtypes.Minter] = []string{authtypes.Minter}

	authKeeper := authkeeper.NewAccountKeeper(
		appCodec, app.GetKey(types.StoreKey), app.GetSubspace(types.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms,
	)

	{
		keeper := bankpluskeeper.NewBaseKeeper(
			appCodec, app.GetKey(types.StoreKey), authKeeper,
			app.GetSubspace(types.ModuleName), make(map[string]bool), false,
		)

		// add blocked address
		keeper.AddToInactiveAddr(ctx, blockedAcc.GetAddress())
		suite.Require().True(keeper.IsInactiveAddr(blockedAcc.GetAddress()))
	}

	{
		keeper := bankpluskeeper.NewBaseKeeper(
			appCodec, app.GetKey(types.StoreKey), authKeeper,
			app.GetSubspace(types.ModuleName), make(map[string]bool), false,
		)
		keeper.InitializeBankPlus(ctx)
		suite.Require().True(keeper.IsInactiveAddr(blockedAcc.GetAddress()))
	}
}

func (suite *IntegrationTestSuite) TestSendCoinsFromModuleToAccount_Blacklist() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{Height: 1})
	appCodec := app.AppCodec()

	// add module accounts to supply keeper
	maccPerms := simapp.GetMaccPerms()
	maccPerms[holder] = nil
	maccPerms[authtypes.Burner] = []string{authtypes.Burner}
	maccPerms[authtypes.Minter] = []string{authtypes.Minter}

	addr1 := sdk.AccAddress([]byte("addr1_______________"))

	authKeeper := authkeeper.NewAccountKeeper(
		appCodec, app.GetKey(types.StoreKey), app.GetSubspace(types.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms,
	)
	keeper := bankpluskeeper.NewBaseKeeper(
		appCodec, app.GetKey(types.StoreKey), authKeeper,
		app.GetSubspace(types.ModuleName), map[string]bool{addr1.String(): true}, false)

	suite.Require().NoError(keeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	suite.Require().Error(keeper.SendCoinsFromModuleToAccount(
		ctx, minttypes.ModuleName, addr1, initCoins,
	))
}

func (suite *IntegrationTestSuite) TestInputOutputCoins() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{Height: 1})
	appCodec := app.AppCodec()

	// add module accounts to supply keeper
	maccPerms := simapp.GetMaccPerms()
	maccPerms[holder] = nil
	maccPerms[authtypes.Burner] = []string{authtypes.Burner}
	maccPerms[authtypes.Minter] = []string{authtypes.Minter}

	authKeeper := authkeeper.NewAccountKeeper(
		appCodec, app.GetKey(types.StoreKey), app.GetSubspace(types.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms,
	)
	keeper := bankpluskeeper.NewBaseKeeper(
		appCodec, app.GetKey(types.StoreKey), authKeeper,
		app.GetSubspace(types.ModuleName), make(map[string]bool), false,
	)

	baseAcc := authKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("baseAcc"))
	authKeeper.SetModuleAccount(ctx, holderAcc)
	authKeeper.SetModuleAccount(ctx, burnerAcc)
	authKeeper.SetAccount(ctx, baseAcc)

	// set initial balances
	suite.
		Require().
		NoError(keeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	suite.
		Require().
		NoError(keeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, baseAcc.GetAddress(), initCoins))
	suite.
		Require().
		NoError(keeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	suite.
		Require().
		NoError(keeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, holderAcc.GetName(), initCoins))

	input := []types.Input{types.NewInput(baseAcc.GetAddress(), initCoins), types.NewInput(holderAcc.GetAddress(), initCoins)}
	output := []types.Output{types.NewOutput(burnerAcc.GetAddress(), initCoins), types.NewOutput(burnerAcc.GetAddress(), initCoins)}

	targetKeeper := func(isDeact bool) bankpluskeeper.BaseKeeper {
		return bankpluskeeper.NewBaseKeeper(
			appCodec, app.GetKey(types.StoreKey), authKeeper,
			app.GetSubspace(types.ModuleName), make(map[string]bool), isDeact,
		)
	}
	tcs := map[string]struct {
		deactMultiSend bool
	}{
		"MultiSend Off": {
			true,
		},
		"MultiSend On": {
			false,
		},
	}

	for name, tc := range tcs {
		tc := tc
		suite.T().Run(name, func(t *testing.T) {
			if tc.deactMultiSend {
				suite.Panics(func() {
					_ = targetKeeper(tc.deactMultiSend).InputOutputCoins(ctx, input, output)
				})
			} else {
				err := targetKeeper(tc.deactMultiSend).InputOutputCoins(ctx, input, output)
				suite.Assert().NoError(err)
			}
		})
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
