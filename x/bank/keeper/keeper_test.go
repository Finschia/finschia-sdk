package keeper_test

import (
	"testing"
	"time"

	abci "github.com/line/ostracon/abci/types"
	ostproto "github.com/line/ostracon/proto/ostracon/types"
	osttime "github.com/line/ostracon/types/time"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/baseapp"
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	authkeeper "github.com/line/lbm-sdk/x/auth/keeper"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	vesting "github.com/line/lbm-sdk/x/auth/vesting/types"
	"github.com/line/lbm-sdk/x/bank/keeper"
	"github.com/line/lbm-sdk/x/bank/types"
	minttypes "github.com/line/lbm-sdk/x/mint/types"
)

const (
	fooDenom     = "foo"
	barDenom     = "bar"
	initialPower = int64(100)
	holder       = "holder"
	multiPerm    = "multiple permissions account"
	randomPerm   = "random permission"
)

var (
	holderAcc     = authtypes.NewEmptyModuleAccount(holder)
	burnerAcc     = authtypes.NewEmptyModuleAccount(authtypes.Burner, authtypes.Burner)
	minterAcc     = authtypes.NewEmptyModuleAccount(authtypes.Minter, authtypes.Minter)
	multiPermAcc  = authtypes.NewEmptyModuleAccount(multiPerm, authtypes.Burner, authtypes.Minter, authtypes.Staking)
	randomPermAcc = authtypes.NewEmptyModuleAccount(randomPerm, "random")

	initTokens = sdk.TokensFromConsensusPower(initialPower)
	initCoins  = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))
)

func newFooCoin(amt int64) sdk.Coin {
	return sdk.NewInt64Coin(fooDenom, amt)
}

func newBarCoin(amt int64) sdk.Coin {
	return sdk.NewInt64Coin(barDenom, amt)
}

// nolint: interfacer
func getCoinsByName(ctx sdk.Context, bk keeper.Keeper, ak types.AccountKeeper, moduleName string) sdk.Coins {
	moduleAddress := ak.GetModuleAddress(moduleName)
	macc := ak.GetAccount(ctx, moduleAddress)
	if macc == nil {
		return sdk.Coins(nil)
	}

	return bk.GetAllBalances(ctx, macc.GetAddress())
}

type IntegrationTestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	queryClient types.QueryClient
}

func (suite *IntegrationTestSuite) SetupTest() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ostproto.Header{})

	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	app.BankKeeper.SetParams(ctx, types.DefaultParams())

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.BankKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = app
	suite.ctx = ctx
	suite.queryClient = queryClient
}

func (suite *IntegrationTestSuite) TestSupply() {
	app, ctx := suite.app, suite.ctx

	initialPower := int64(100)
	initTokens := sdk.TokensFromConsensusPower(initialPower)

	totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))
	suite.NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, totalSupply))

	total := app.BankKeeper.GetTotalSupply(ctx)
	suite.Require().Equal(totalSupply, total)
}

func (suite *IntegrationTestSuite) TestSendCoinsFromModuleToAccount_Blacklist() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ostproto.Header{Height: 1})
	appCodec := app.AppCodec()

	// add module accounts to supply keeper
	maccPerms := simapp.GetMaccPerms()
	maccPerms[holder] = nil
	maccPerms[authtypes.Burner] = []string{authtypes.Burner}
	maccPerms[authtypes.Minter] = []string{authtypes.Minter}
	maccPerms[multiPerm] = []string{authtypes.Burner, authtypes.Minter, authtypes.Staking}
	maccPerms[randomPerm] = []string{"random"}

	addr1 := sdk.AccAddress([]byte("addr1_______________"))

	authKeeper := authkeeper.NewAccountKeeper(
		appCodec, app.GetKey(types.StoreKey), app.GetSubspace(types.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms,
	)
	keeper := keeper.NewBaseKeeper(
		appCodec, app.GetKey(types.StoreKey), authKeeper,
		app.GetSubspace(types.ModuleName), map[string]bool{addr1.String(): true},
	)

	baseAcc := authKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("baseAcc"))
	suite.Require().NoError(keeper.SetBalances(ctx, holderAcc.GetAddress(), initCoins))

	keeper.SetSupply(ctx, types.NewSupply(initCoins))
	authKeeper.SetModuleAccount(ctx, holderAcc)
	authKeeper.SetAccount(ctx, baseAcc)

	suite.Require().Error(keeper.SendCoinsFromModuleToAccount(ctx, holderAcc.GetName(), addr1, initCoins))
}

func (suite *IntegrationTestSuite) TestSupply_SendCoins() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ostproto.Header{Height: 1})
	appCodec := app.AppCodec()

	// add module accounts to supply keeper
	maccPerms := simapp.GetMaccPerms()
	maccPerms[holder] = nil
	maccPerms[authtypes.Burner] = []string{authtypes.Burner}
	maccPerms[authtypes.Minter] = []string{authtypes.Minter}
	maccPerms[multiPerm] = []string{authtypes.Burner, authtypes.Minter, authtypes.Staking}
	maccPerms[randomPerm] = []string{"random"}

	authKeeper := authkeeper.NewAccountKeeper(
		appCodec, app.GetKey(types.StoreKey), app.GetSubspace(types.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms,
	)
	keeper := keeper.NewBaseKeeper(
		appCodec, app.GetKey(types.StoreKey), authKeeper,
		app.GetSubspace(types.ModuleName), make(map[string]bool),
	)

	baseAcc := authKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("baseAcc"))

	// set initial balances
	suite.
		Require().
		NoError(keeper.MintCoins(ctx, minttypes.ModuleName, initCoins))

	suite.
		Require().
		NoError(keeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, holderAcc.GetAddress(), initCoins))

	authKeeper.SetModuleAccount(ctx, holderAcc)
	authKeeper.SetModuleAccount(ctx, burnerAcc)
	authKeeper.SetAccount(ctx, baseAcc)

	suite.Require().Panics(func() {
		_ = keeper.SendCoinsFromModuleToModule(ctx, "", holderAcc.GetName(), initCoins) // nolint:errcheck
	})

	suite.Require().Panics(func() {
		_ = keeper.SendCoinsFromModuleToModule(ctx, authtypes.Burner, "", initCoins) // nolint:errcheck
	})

	suite.Require().Panics(func() {
		_ = keeper.SendCoinsFromModuleToAccount(ctx, "", baseAcc.GetAddress(), initCoins) // nolint:errcheck
	})

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

func (suite *IntegrationTestSuite) TestSupply_MintCoins() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ostproto.Header{Height: 1})
	appCodec := app.AppCodec()

	// add module accounts to supply keeper
	maccPerms := simapp.GetMaccPerms()
	maccPerms[holder] = nil
	maccPerms[authtypes.Burner] = []string{authtypes.Burner}
	maccPerms[authtypes.Minter] = []string{authtypes.Minter}
	maccPerms[multiPerm] = []string{authtypes.Burner, authtypes.Minter, authtypes.Staking}
	maccPerms[randomPerm] = []string{"random"}

	authKeeper := authkeeper.NewAccountKeeper(
		appCodec, app.GetKey(types.StoreKey), app.GetSubspace(types.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms,
	)
	keeper := keeper.NewBaseKeeper(
		appCodec, app.GetKey(types.StoreKey), authKeeper,
		app.GetSubspace(types.ModuleName), make(map[string]bool),
	)

	authKeeper.SetModuleAccount(ctx, burnerAcc)
	authKeeper.SetModuleAccount(ctx, minterAcc)
	authKeeper.SetModuleAccount(ctx, multiPermAcc)
	authKeeper.SetModuleAccount(ctx, randomPermAcc)

	initialSupply := keeper.GetTotalSupply(ctx)

	suite.Require().Panics(func() { keeper.MintCoins(ctx, "", initCoins) }, "no module account")                // nolint:errcheck
	suite.Require().Panics(func() { keeper.MintCoins(ctx, authtypes.Burner, initCoins) }, "invalid permission") // nolint:errcheck

	err := keeper.MintCoins(ctx, authtypes.Minter, sdk.Coins{sdk.Coin{Denom: "denom", Amount: sdk.NewInt(-10)}})
	suite.Require().Error(err, "insufficient coins")

	suite.Require().Panics(func() { keeper.MintCoins(ctx, randomPerm, initCoins) }) // nolint:errcheck

	err = keeper.MintCoins(ctx, authtypes.Minter, initCoins)
	suite.Require().NoError(err)

	suite.Require().Equal(initCoins, getCoinsByName(ctx, keeper, authKeeper, authtypes.Minter))
	suite.Require().Equal(initialSupply.Add(initCoins...), keeper.GetTotalSupply(ctx))

	// test same functionality on module account with multiple permissions
	initialSupply = keeper.GetTotalSupply(ctx)

	err = keeper.MintCoins(ctx, multiPermAcc.GetName(), initCoins)
	suite.Require().NoError(err)

	suite.Require().Equal(initCoins, getCoinsByName(ctx, keeper, authKeeper, multiPermAcc.GetName()))
	suite.Require().Equal(initialSupply.Add(initCoins...), keeper.GetTotalSupply(ctx))
	suite.Require().Panics(func() { keeper.MintCoins(ctx, authtypes.Burner, initCoins) }) // nolint:errcheck
}

func (suite *IntegrationTestSuite) TestSupply_BurnCoins() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ostproto.Header{Height: 1})
	appCodec := simapp.MakeTestEncodingConfig().Marshaler

	// add module accounts to supply keeper
	maccPerms := simapp.GetMaccPerms()
	maccPerms[holder] = nil
	maccPerms[authtypes.Burner] = []string{authtypes.Burner}
	maccPerms[authtypes.Minter] = []string{authtypes.Minter}
	maccPerms[multiPerm] = []string{authtypes.Burner, authtypes.Minter, authtypes.Staking}
	maccPerms[randomPerm] = []string{"random"}

	authKeeper := authkeeper.NewAccountKeeper(
		appCodec, app.GetKey(types.StoreKey), app.GetSubspace(types.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms,
	)
	keeper := keeper.NewBaseKeeper(
		appCodec, app.GetKey(types.StoreKey), authKeeper,
		app.GetSubspace(types.ModuleName), make(map[string]bool),
	)

	// set burnerAcc balance
	authKeeper.SetModuleAccount(ctx, burnerAcc)
	suite.
		Require().
		NoError(keeper.MintCoins(ctx, authtypes.Minter, initCoins))
	suite.
		Require().
		NoError(keeper.SendCoinsFromModuleToAccount(ctx, authtypes.Minter, burnerAcc.GetAddress(), initCoins))

	// inflate supply
	suite.
		Require().
		NoError(keeper.MintCoins(ctx, authtypes.Minter, initCoins))
	supplyAfterInflation := keeper.GetTotalSupply(ctx)

	suite.Require().Panics(func() { keeper.BurnCoins(ctx, "", initCoins) }, "no module account")                    // nolint:errcheck
	suite.Require().Panics(func() { keeper.BurnCoins(ctx, authtypes.Minter, initCoins) }, "invalid permission")     // nolint:errcheck
	suite.Require().Panics(func() { keeper.BurnCoins(ctx, randomPerm, supplyAfterInflation) }, "random permission") // nolint:errcheck
	err := keeper.BurnCoins(ctx, authtypes.Burner, supplyAfterInflation)
	suite.Require().Error(err, "insufficient coins")

	err = keeper.BurnCoins(ctx, authtypes.Burner, initCoins)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.NewCoins().String(), getCoinsByName(ctx, keeper, authKeeper, authtypes.Burner).String())
	suite.Require().Equal(supplyAfterInflation.Sub(initCoins), keeper.GetTotalSupply(ctx))

	// test same functionality on module account with multiple permissions
	suite.
		Require().
		NoError(keeper.MintCoins(ctx, authtypes.Minter, initCoins))
	supplyAfterInflation = keeper.GetTotalSupply(ctx)

	suite.Require().NoError(keeper.SendCoins(ctx, authtypes.NewModuleAddress(authtypes.Minter), multiPermAcc.GetAddress(), initCoins))
	authKeeper.SetModuleAccount(ctx, multiPermAcc)

	err = keeper.BurnCoins(ctx, multiPermAcc.GetName(), initCoins)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.NewCoins().String(), getCoinsByName(ctx, keeper, authKeeper, multiPermAcc.GetName()).String())
	suite.Require().Equal(supplyAfterInflation.Sub(initCoins), keeper.GetTotalSupply(ctx))
}

func (suite *IntegrationTestSuite) TestSendCoinsNewAccount() {
	app, ctx := suite.app, suite.ctx
	balances := sdk.NewCoins(newFooCoin(100), newBarCoin(50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	app.AccountKeeper.SetAccount(ctx, acc1)
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr1, balances))

	acc1Balances := app.BankKeeper.GetAllBalances(ctx, addr1)
	suite.Require().Equal(balances, acc1Balances)

	addr2 := sdk.AccAddress([]byte("addr2_______________"))

	suite.Require().Nil(app.AccountKeeper.GetAccount(ctx, addr2))
	app.BankKeeper.GetAllBalances(ctx, addr2)
	suite.Require().Empty(app.BankKeeper.GetAllBalances(ctx, addr2))

	sendAmt := sdk.NewCoins(newFooCoin(50), newBarCoin(25))
	suite.Require().NoError(app.BankKeeper.SendCoins(ctx, addr1, addr2, sendAmt))

	acc2Balances := app.BankKeeper.GetAllBalances(ctx, addr2)
	suite.Require().Equal(sendAmt, acc2Balances)
	suite.Require().NotNil(app.AccountKeeper.GetAccount(ctx, addr2))
}

func (suite *IntegrationTestSuite) TestInputOutputNewAccount() {
	app, ctx := suite.app, suite.ctx

	balances := sdk.NewCoins(newFooCoin(100), newBarCoin(50))
	addr1 := sdk.BytesToAccAddress([]byte("addr1_______________"))
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	app.AccountKeeper.SetAccount(ctx, acc1)
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr1, balances))

	acc1Balances := app.BankKeeper.GetAllBalances(ctx, addr1)
	suite.Require().Equal(balances, acc1Balances)

	addr2 := sdk.BytesToAccAddress([]byte("addr2_______________"))

	suite.Require().Nil(app.AccountKeeper.GetAccount(ctx, addr2))
	suite.Require().Empty(app.BankKeeper.GetAllBalances(ctx, addr2))

	inputs := []types.Input{
		{Address: addr1.String(), Coins: sdk.NewCoins(newFooCoin(30), newBarCoin(10))},
	}
	outputs := []types.Output{
		{Address: addr2.String(), Coins: sdk.NewCoins(newFooCoin(30), newBarCoin(10))},
	}

	suite.Require().NoError(app.BankKeeper.InputOutputCoins(ctx, inputs, outputs))

	expected := sdk.NewCoins(newFooCoin(30), newBarCoin(10))
	acc2Balances := app.BankKeeper.GetAllBalances(ctx, addr2)
	suite.Require().Equal(expected, acc2Balances)
	suite.Require().NotNil(app.AccountKeeper.GetAccount(ctx, addr2))
}

func (suite *IntegrationTestSuite) TestInputOutputCoins() {
	app, ctx := suite.app, suite.ctx
	balances := sdk.NewCoins(newFooCoin(90), newBarCoin(30))

	addr1 := sdk.BytesToAccAddress([]byte("addr1_______________"))
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	app.AccountKeeper.SetAccount(ctx, acc1)

	addr2 := sdk.BytesToAccAddress([]byte("addr2_______________"))
	acc2 := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)
	app.AccountKeeper.SetAccount(ctx, acc2)

	addr3 := sdk.BytesToAccAddress([]byte("addr3_______________"))
	acc3 := app.AccountKeeper.NewAccountWithAddress(ctx, addr3)
	app.AccountKeeper.SetAccount(ctx, acc3)

	inputs := []types.Input{
		{Address: addr1.String(), Coins: sdk.NewCoins(newFooCoin(30), newBarCoin(10))},
		{Address: addr1.String(), Coins: sdk.NewCoins(newFooCoin(30), newBarCoin(10))},
	}
	outputs := []types.Output{
		{Address: addr2.String(), Coins: sdk.NewCoins(newFooCoin(30), newBarCoin(10))},
		{Address: addr3.String(), Coins: sdk.NewCoins(newFooCoin(30), newBarCoin(10))},
	}

	suite.Require().Error(app.BankKeeper.InputOutputCoins(ctx, inputs, []types.Output{}))
	suite.Require().Error(app.BankKeeper.InputOutputCoins(ctx, inputs, outputs))

	suite.Require().NoError(simapp.FundAccount(app, ctx, addr1, balances))

	insufficientInputs := []types.Input{
		{Address: addr1.String(), Coins: sdk.NewCoins(newFooCoin(300), newBarCoin(100))},
		{Address: addr1.String(), Coins: sdk.NewCoins(newFooCoin(300), newBarCoin(100))},
	}
	insufficientOutputs := []types.Output{
		{Address: addr2.String(), Coins: sdk.NewCoins(newFooCoin(300), newBarCoin(100))},
		{Address: addr3.String(), Coins: sdk.NewCoins(newFooCoin(300), newBarCoin(100))},
	}
	suite.Require().Error(app.BankKeeper.InputOutputCoins(ctx, insufficientInputs, insufficientOutputs))
	suite.Require().NoError(app.BankKeeper.InputOutputCoins(ctx, inputs, outputs))

	acc1Balances := app.BankKeeper.GetAllBalances(ctx, addr1)
	expected := sdk.NewCoins(newFooCoin(30), newBarCoin(10))
	suite.Require().Equal(expected, acc1Balances)

	acc2Balances := app.BankKeeper.GetAllBalances(ctx, addr2)
	suite.Require().Equal(expected, acc2Balances)

	acc3Balances := app.BankKeeper.GetAllBalances(ctx, addr3)
	suite.Require().Equal(expected, acc3Balances)
}

func (suite *IntegrationTestSuite) TestSendCoins() {
	app, ctx := suite.app, suite.ctx
	balances := sdk.NewCoins(newFooCoin(100), newBarCoin(50))

	addr1 := sdk.AccAddress("addr1_______________")
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	app.AccountKeeper.SetAccount(ctx, acc1)

	addr2 := sdk.AccAddress("addr2_______________")
	acc2 := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)
	app.AccountKeeper.SetAccount(ctx, acc2)
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr2, balances))

	sendAmt := sdk.NewCoins(newFooCoin(50), newBarCoin(25))
	suite.Require().Error(app.BankKeeper.SendCoins(ctx, addr1, addr2, sendAmt))

	suite.Require().NoError(simapp.FundAccount(app, ctx, addr1, balances))
	suite.Require().NoError(app.BankKeeper.SendCoins(ctx, addr1, addr2, sendAmt))

	acc1Balances := app.BankKeeper.GetAllBalances(ctx, addr1)
	expected := sdk.NewCoins(newFooCoin(50), newBarCoin(25))
	suite.Require().Equal(expected, acc1Balances)

	acc2Balances := app.BankKeeper.GetAllBalances(ctx, addr2)
	expected = sdk.NewCoins(newFooCoin(150), newBarCoin(75))
	suite.Require().Equal(expected, acc2Balances)
}

func (suite *IntegrationTestSuite) TestValidateBalance() {
	app, ctx := suite.app, suite.ctx
	now := osttime.Now()
	ctx = ctx.WithBlockHeader(ostproto.Header{Time: now})
	endTime := now.Add(24 * time.Hour)

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))

	suite.Require().Error(app.BankKeeper.ValidateBalance(ctx, addr1))

	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	app.AccountKeeper.SetAccount(ctx, acc)

	balances := sdk.NewCoins(newFooCoin(100))
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr1, balances))
	suite.Require().NoError(app.BankKeeper.ValidateBalance(ctx, addr1))

	bacc := authtypes.NewBaseAccountWithAddress(addr2)
	vacc := vesting.NewContinuousVestingAccount(bacc, balances.Add(balances...), now.Unix(), endTime.Unix())

	app.AccountKeeper.SetAccount(ctx, vacc)
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr2, balances))
	suite.Require().Error(app.BankKeeper.ValidateBalance(ctx, addr2))
}

func (suite *IntegrationTestSuite) TestSendEnabled() {
	app, ctx := suite.app, suite.ctx
	enabled := true
	params := types.DefaultParams()
	suite.Require().Equal(enabled, params.DefaultSendEnabled)

	app.BankKeeper.SetParams(ctx, params)

	bondCoin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())
	fooCoin := sdk.NewCoin("foocoin", sdk.OneInt())
	barCoin := sdk.NewCoin("barcoin", sdk.OneInt())

	// assert with default (all denom) send enabled both Bar and Bond Denom are enabled
	suite.Require().Equal(enabled, app.BankKeeper.SendEnabledCoin(ctx, barCoin))
	suite.Require().Equal(enabled, app.BankKeeper.SendEnabledCoin(ctx, bondCoin))

	// Both coins should be send enabled.
	err := app.BankKeeper.SendEnabledCoins(ctx, fooCoin, bondCoin)
	suite.Require().NoError(err)

	// Set default send_enabled to !enabled, add a foodenom that overrides default as enabled
	params.DefaultSendEnabled = !enabled
	params = params.SetSendEnabledParam(fooCoin.Denom, enabled)
	app.BankKeeper.SetParams(ctx, params)

	// Expect our specific override to be enabled, others to be !enabled.
	suite.Require().Equal(enabled, app.BankKeeper.SendEnabledCoin(ctx, fooCoin))
	suite.Require().Equal(!enabled, app.BankKeeper.SendEnabledCoin(ctx, barCoin))
	suite.Require().Equal(!enabled, app.BankKeeper.SendEnabledCoin(ctx, bondCoin))

	// Foo coin should be send enabled.
	err = app.BankKeeper.SendEnabledCoins(ctx, fooCoin)
	suite.Require().NoError(err)

	// Expect an error when one coin is not send enabled.
	err = app.BankKeeper.SendEnabledCoins(ctx, fooCoin, bondCoin)
	suite.Require().Error(err)

	// Expect an error when all coins are not send enabled.
	err = app.BankKeeper.SendEnabledCoins(ctx, bondCoin, barCoin)
	suite.Require().Error(err)
}

func (suite *IntegrationTestSuite) TestHasBalance() {
	app, ctx := suite.app, suite.ctx
	addr := sdk.AccAddress([]byte("addr1_______________"))

	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)

	balances := sdk.NewCoins(newFooCoin(100))
	suite.Require().False(app.BankKeeper.HasBalance(ctx, addr, newFooCoin(99)))

	suite.Require().NoError(simapp.FundAccount(app, ctx, addr, balances))
	suite.Require().False(app.BankKeeper.HasBalance(ctx, addr, newFooCoin(101)))
	suite.Require().True(app.BankKeeper.HasBalance(ctx, addr, newFooCoin(100)))
	suite.Require().True(app.BankKeeper.HasBalance(ctx, addr, newFooCoin(1)))
}

func (suite *IntegrationTestSuite) TestMsgSendEvents() {
	app, ctx := suite.app, suite.ctx
	addr := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)

	app.AccountKeeper.SetAccount(ctx, acc)
	newCoins := sdk.NewCoins(sdk.NewInt64Coin(fooDenom, 50))

	suite.Require().Error(app.BankKeeper.SendCoins(ctx, addr, addr2, newCoins))

	events := ctx.EventManager().ABCIEvents()
	suite.Require().Equal(2, len(events))

	event1 := sdk.Event{
		Type:       types.EventTypeTransfer,
		Attributes: []abci.EventAttribute{},
	}
	event1.Attributes = append(
		event1.Attributes,
		abci.EventAttribute{Key: []byte(types.AttributeKeyRecipient), Value: []byte(addr2.String())},
	)
	event1.Attributes = append(
		event1.Attributes,
		abci.EventAttribute{Key: []byte(types.AttributeKeySender), Value: []byte(addr.String())},
	)
	event1.Attributes = append(
		event1.Attributes,
		abci.EventAttribute{Key: []byte(sdk.AttributeKeyAmount), Value: []byte(newCoins.String())},
	)

	event2 := sdk.Event{
		Type:       sdk.EventTypeMessage,
		Attributes: []abci.EventAttribute{},
	}
	event2.Attributes = append(
		event2.Attributes,
		abci.EventAttribute{Key: []byte(types.AttributeKeySender), Value: []byte(addr.String())},
	)

	suite.Require().Equal(abci.Event(event1), events[0])
	suite.Require().Equal(abci.Event(event2), events[1])

	suite.Require().NoError(simapp.FundAccount(app, ctx, addr, sdk.NewCoins(sdk.NewInt64Coin(fooDenom, 50))))
	newCoins = sdk.NewCoins(sdk.NewInt64Coin(fooDenom, 50))

	suite.Require().NoError(app.BankKeeper.SendCoins(ctx, addr, addr2, newCoins))

	// events are shifted due to the funding account events
	events = ctx.EventManager().ABCIEvents()
	suite.Require().Equal(12, len(events))
	suite.Require().Equal(abci.Event(event1), events[8])
	suite.Require().Equal(abci.Event(event2), events[9])
}

func (suite *IntegrationTestSuite) TestMsgMultiSendEvents() {
	app, ctx := suite.app, suite.ctx

	app.BankKeeper.SetParams(ctx, types.DefaultParams())

	addr := sdk.BytesToAccAddress([]byte("addr1_______________"))
	addr2 := sdk.BytesToAccAddress([]byte("addr2_______________"))
	addr3 := sdk.BytesToAccAddress([]byte("addr3_______________"))
	addr4 := sdk.BytesToAccAddress([]byte("addr4_______________"))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	acc2 := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)

	app.AccountKeeper.SetAccount(ctx, acc)
	app.AccountKeeper.SetAccount(ctx, acc2)

	newCoins := sdk.NewCoins(sdk.NewInt64Coin(fooDenom, 50))
	newCoins2 := sdk.NewCoins(sdk.NewInt64Coin(barDenom, 100))
	inputs := []types.Input{
		{Address: addr.String(), Coins: newCoins},
		{Address: addr2.String(), Coins: newCoins2},
	}
	outputs := []types.Output{
		{Address: addr3.String(), Coins: newCoins},
		{Address: addr4.String(), Coins: newCoins2},
	}

	suite.Require().Error(app.BankKeeper.InputOutputCoins(ctx, inputs, outputs))

	events := ctx.EventManager().ABCIEvents()
	suite.Require().Equal(0, len(events))

	// Set addr's coins but not addr2's coins
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr, sdk.NewCoins(sdk.NewInt64Coin(fooDenom, 50))))
	suite.Require().Error(app.BankKeeper.InputOutputCoins(ctx, inputs, outputs))

	events = ctx.EventManager().ABCIEvents()
	suite.Require().Equal(8, len(events)) // 7 events because account funding causes extra minting + coin_spent + coin_recv events

	event1 := sdk.Event{
		Type:       sdk.EventTypeMessage,
		Attributes: []abci.EventAttribute{},
	}
	event1.Attributes = append(
		event1.Attributes,
		abci.EventAttribute{Key: []byte(types.AttributeKeySender), Value: []byte(addr.String())},
	)
	suite.Require().Equal(abci.Event(event1), events[7])

	// Set addr's coins and addr2's coins
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr, sdk.NewCoins(sdk.NewInt64Coin(fooDenom, 50))))
	newCoins = sdk.NewCoins(sdk.NewInt64Coin(fooDenom, 50))

	suite.Require().NoError(simapp.FundAccount(app, ctx, addr2, sdk.NewCoins(sdk.NewInt64Coin(barDenom, 100))))
	newCoins2 = sdk.NewCoins(sdk.NewInt64Coin(barDenom, 100))

	suite.Require().NoError(app.BankKeeper.InputOutputCoins(ctx, inputs, outputs))

	events = ctx.EventManager().ABCIEvents()
	suite.Require().Equal(28, len(events)) // 25 due to account funding + coin_spent + coin_recv events

	event2 := sdk.Event{
		Type:       sdk.EventTypeMessage,
		Attributes: []abci.EventAttribute{},
	}
	event2.Attributes = append(
		event2.Attributes,
		abci.EventAttribute{Key: []byte(types.AttributeKeySender), Value: []byte(addr2.String())},
	)
	event3 := sdk.Event{
		Type:       types.EventTypeTransfer,
		Attributes: []abci.EventAttribute{},
	}
	event3.Attributes = append(
		event3.Attributes,
		abci.EventAttribute{Key: []byte(types.AttributeKeyRecipient), Value: []byte(addr3.String())},
	)
	event3.Attributes = append(
		event3.Attributes,
		abci.EventAttribute{Key: []byte(sdk.AttributeKeyAmount), Value: []byte(newCoins.String())})
	event4 := sdk.Event{
		Type:       types.EventTypeTransfer,
		Attributes: []abci.EventAttribute{},
	}
	event4.Attributes = append(
		event4.Attributes,
		abci.EventAttribute{Key: []byte(types.AttributeKeyRecipient), Value: []byte(addr4.String())},
	)
	event4.Attributes = append(
		event4.Attributes,
		abci.EventAttribute{Key: []byte(sdk.AttributeKeyAmount), Value: []byte(newCoins2.String())},
	)
	// events are shifted due to the funding account events
	suite.Require().Equal(abci.Event(event1), events[21])
	suite.Require().Equal(abci.Event(event2), events[23])
	suite.Require().Equal(abci.Event(event3), events[25])
	suite.Require().Equal(abci.Event(event4), events[27])
}

func (suite *IntegrationTestSuite) TestSpendableCoins() {
	app, ctx := suite.app, suite.ctx
	now := osttime.Now()
	ctx = ctx.WithBlockHeader(ostproto.Header{Time: now})
	endTime := now.Add(24 * time.Hour)

	origCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	delCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))
	addrModule := sdk.AccAddress([]byte("moduleAcc___________"))

	macc := app.AccountKeeper.NewAccountWithAddress(ctx, addrModule)
	bacc := authtypes.NewBaseAccountWithAddress(addr1)
	vacc := vesting.NewContinuousVestingAccount(bacc, origCoins, ctx.BlockHeader().Time.Unix(), endTime.Unix())
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)

	app.AccountKeeper.SetAccount(ctx, macc)
	app.AccountKeeper.SetAccount(ctx, vacc)
	app.AccountKeeper.SetAccount(ctx, acc)
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr1, origCoins))
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr2, origCoins))

	suite.Require().Equal(origCoins, app.BankKeeper.SpendableCoins(ctx, addr2))

	ctx = ctx.WithBlockTime(now.Add(12 * time.Hour))
	suite.Require().NoError(app.BankKeeper.DelegateCoins(ctx, addr2, addrModule, delCoins))
	suite.Require().Equal(origCoins.Sub(delCoins), app.BankKeeper.SpendableCoins(ctx, addr1))
}

func (suite *IntegrationTestSuite) TestVestingAccountSend() {
	app, ctx := suite.app, suite.ctx
	now := osttime.Now()
	ctx = ctx.WithBlockHeader(ostproto.Header{Time: now})
	endTime := now.Add(24 * time.Hour)

	origCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))

	bacc := authtypes.NewBaseAccountWithAddress(addr1)
	vacc := vesting.NewContinuousVestingAccount(bacc, origCoins, now.Unix(), endTime.Unix())

	app.AccountKeeper.SetAccount(ctx, vacc)
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr1, origCoins))

	// require that no coins be sendable at the beginning of the vesting schedule
	suite.Require().Error(app.BankKeeper.SendCoins(ctx, addr1, addr2, sendCoins))

	// receive some coins
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr1, sendCoins))
	// require that all vested coins are spendable plus any received
	ctx = ctx.WithBlockTime(now.Add(12 * time.Hour))
	suite.Require().NoError(app.BankKeeper.SendCoins(ctx, addr1, addr2, sendCoins))
	suite.Require().Equal(origCoins, app.BankKeeper.GetAllBalances(ctx, addr1))
}

func (suite *IntegrationTestSuite) TestPeriodicVestingAccountSend() {
	app, ctx := suite.app, suite.ctx
	now := osttime.Now()
	ctx = ctx.WithBlockHeader(ostproto.Header{Time: now})
	origCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))
	periods := vesting.Periods{
		vesting.Period{Length: int64(12 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin("stake", 50)}},
		vesting.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin("stake", 25)}},
		vesting.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin("stake", 25)}},
	}

	bacc := authtypes.NewBaseAccountWithAddress(addr1)
	vacc := vesting.NewPeriodicVestingAccount(bacc, origCoins, ctx.BlockHeader().Time.Unix(), periods)

	app.AccountKeeper.SetAccount(ctx, vacc)
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr1, origCoins))

	// require that no coins be sendable at the beginning of the vesting schedule
	suite.Require().Error(app.BankKeeper.SendCoins(ctx, addr1, addr2, sendCoins))

	// receive some coins
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr1, sendCoins))

	// require that all vested coins are spendable plus any received
	ctx = ctx.WithBlockTime(now.Add(12 * time.Hour))
	suite.Require().NoError(app.BankKeeper.SendCoins(ctx, addr1, addr2, sendCoins))
	suite.Require().Equal(origCoins, app.BankKeeper.GetAllBalances(ctx, addr1))
}

func (suite *IntegrationTestSuite) TestVestingAccountReceive() {
	app, ctx := suite.app, suite.ctx
	now := osttime.Now()
	ctx = ctx.WithBlockHeader(ostproto.Header{Time: now})
	endTime := now.Add(24 * time.Hour)

	origCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))

	bacc := authtypes.NewBaseAccountWithAddress(addr1)
	vacc := vesting.NewContinuousVestingAccount(bacc, origCoins, ctx.BlockHeader().Time.Unix(), endTime.Unix())
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)

	app.AccountKeeper.SetAccount(ctx, vacc)
	app.AccountKeeper.SetAccount(ctx, acc)
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr1, origCoins))
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr2, origCoins))

	// send some coins to the vesting account
	suite.Require().NoError(app.BankKeeper.SendCoins(ctx, addr2, addr1, sendCoins))

	// require the coins are spendable
	vacc = app.AccountKeeper.GetAccount(ctx, addr1).(*vesting.ContinuousVestingAccount)
	balances := app.BankKeeper.GetAllBalances(ctx, addr1)
	suite.Require().Equal(origCoins.Add(sendCoins...), balances)
	suite.Require().Equal(balances.Sub(vacc.LockedCoins(now)), sendCoins)

	// require coins are spendable plus any that have vested
	suite.Require().Equal(balances.Sub(vacc.LockedCoins(now.Add(12*time.Hour))), origCoins)
}

func (suite *IntegrationTestSuite) TestPeriodicVestingAccountReceive() {
	app, ctx := suite.app, suite.ctx
	now := osttime.Now()
	ctx = ctx.WithBlockHeader(ostproto.Header{Time: now})

	origCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))

	bacc := authtypes.NewBaseAccountWithAddress(addr1)
	periods := vesting.Periods{
		vesting.Period{Length: int64(12 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin("stake", 50)}},
		vesting.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin("stake", 25)}},
		vesting.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin("stake", 25)}},
	}

	vacc := vesting.NewPeriodicVestingAccount(bacc, origCoins, ctx.BlockHeader().Time.Unix(), periods)
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)

	app.AccountKeeper.SetAccount(ctx, vacc)
	app.AccountKeeper.SetAccount(ctx, acc)
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr1, origCoins))
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr2, origCoins))

	// send some coins to the vesting account
	suite.Require().NoError(app.BankKeeper.SendCoins(ctx, addr2, addr1, sendCoins))

	// require the coins are spendable
	vacc = app.AccountKeeper.GetAccount(ctx, addr1).(*vesting.PeriodicVestingAccount)
	balances := app.BankKeeper.GetAllBalances(ctx, addr1)
	suite.Require().Equal(origCoins.Add(sendCoins...), balances)
	suite.Require().Equal(balances.Sub(vacc.LockedCoins(now)), sendCoins)

	// require coins are spendable plus any that have vested
	suite.Require().Equal(balances.Sub(vacc.LockedCoins(now.Add(12*time.Hour))), origCoins)
}

func (suite *IntegrationTestSuite) TestDelegateCoins() {
	app, ctx := suite.app, suite.ctx
	now := osttime.Now()
	ctx = ctx.WithBlockHeader(ostproto.Header{Time: now})
	endTime := now.Add(24 * time.Hour)

	origCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	delCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))
	addrModule := sdk.AccAddress([]byte("moduleAcc___________"))

	macc := app.AccountKeeper.NewAccountWithAddress(ctx, addrModule) // we don't need to define an actual module account bc we just need the address for testing
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)
	bacc := authtypes.NewBaseAccountWithAddress(addr1)
	vacc := vesting.NewContinuousVestingAccount(bacc, origCoins, ctx.BlockHeader().Time.Unix(), endTime.Unix())

	app.AccountKeeper.SetAccount(ctx, vacc)
	app.AccountKeeper.SetAccount(ctx, acc)
	app.AccountKeeper.SetAccount(ctx, macc)
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr1, origCoins))
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr2, origCoins))

	ctx = ctx.WithBlockTime(now.Add(12 * time.Hour))

	// require the ability for a non-vesting account to delegate
	suite.Require().NoError(app.BankKeeper.DelegateCoins(ctx, addr2, addrModule, delCoins))
	suite.Require().Equal(origCoins.Sub(delCoins), app.BankKeeper.GetAllBalances(ctx, addr2))
	suite.Require().Equal(delCoins, app.BankKeeper.GetAllBalances(ctx, addrModule))

	// require the ability for a vesting account to delegate
	suite.Require().NoError(app.BankKeeper.DelegateCoins(ctx, addr1, addrModule, delCoins))
	suite.Require().Equal(delCoins, app.BankKeeper.GetAllBalances(ctx, addr1))
}

func (suite *IntegrationTestSuite) TestDelegateCoins_Invalid() {
	app, ctx := suite.app, suite.ctx

	origCoins := sdk.NewCoins(newFooCoin(100))
	delCoins := sdk.NewCoins(newFooCoin(50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addrModule := sdk.AccAddress([]byte("moduleAcc___________"))
	macc := app.AccountKeeper.NewAccountWithAddress(ctx, addrModule) // we don't need to define an actual module account bc we just need the address for testing
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)

	suite.Require().Error(app.BankKeeper.DelegateCoins(ctx, addr1, addrModule, delCoins))
	invalidCoins := sdk.Coins{sdk.Coin{Denom: "fooDenom", Amount: sdk.NewInt(-50)}}
	suite.Require().Error(app.BankKeeper.DelegateCoins(ctx, addr1, addrModule, invalidCoins))

	app.AccountKeeper.SetAccount(ctx, macc)
	suite.Require().Error(app.BankKeeper.DelegateCoins(ctx, addr1, addrModule, delCoins))
	app.AccountKeeper.SetAccount(ctx, acc)
	suite.Require().Error(app.BankKeeper.DelegateCoins(ctx, addr1, addrModule, origCoins.Add(origCoins...)))
}

func (suite *IntegrationTestSuite) TestUndelegateCoins() {
	app, ctx := suite.app, suite.ctx
	now := osttime.Now()
	ctx = ctx.WithBlockHeader(ostproto.Header{Time: now})
	endTime := now.Add(24 * time.Hour)

	origCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	delCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))
	addrModule := sdk.AccAddress([]byte("moduleAcc___________"))

	bacc := authtypes.NewBaseAccountWithAddress(addr1)
	macc := app.AccountKeeper.NewAccountWithAddress(ctx, addrModule) // we don't need to define an actual module account bc we just need the address for testing

	vacc := vesting.NewContinuousVestingAccount(bacc, origCoins, ctx.BlockHeader().Time.Unix(), endTime.Unix())
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)

	app.AccountKeeper.SetAccount(ctx, vacc)
	app.AccountKeeper.SetAccount(ctx, acc)
	app.AccountKeeper.SetAccount(ctx, macc)
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr1, origCoins))
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr2, origCoins))

	ctx = ctx.WithBlockTime(now.Add(12 * time.Hour))

	// require the ability for a non-vesting account to delegate
	err := app.BankKeeper.DelegateCoins(ctx, addr2, addrModule, delCoins)
	suite.Require().NoError(err)

	suite.Require().Equal(origCoins.Sub(delCoins), app.BankKeeper.GetAllBalances(ctx, addr2))
	suite.Require().Equal(delCoins, app.BankKeeper.GetAllBalances(ctx, addrModule))

	// require the ability for a non-vesting account to undelegate
	suite.Require().NoError(app.BankKeeper.UndelegateCoins(ctx, addrModule, addr2, delCoins))

	suite.Require().Equal(origCoins, app.BankKeeper.GetAllBalances(ctx, addr2))
	suite.Require().True(app.BankKeeper.GetAllBalances(ctx, addrModule).Empty())

	// require the ability for a vesting account to delegate
	suite.Require().NoError(app.BankKeeper.DelegateCoins(ctx, addr1, addrModule, delCoins))

	suite.Require().Equal(origCoins.Sub(delCoins), app.BankKeeper.GetAllBalances(ctx, addr1))
	suite.Require().Equal(delCoins, app.BankKeeper.GetAllBalances(ctx, addrModule))

	// require the ability for a vesting account to undelegate
	suite.Require().NoError(app.BankKeeper.UndelegateCoins(ctx, addrModule, addr1, delCoins))

	suite.Require().Equal(origCoins, app.BankKeeper.GetAllBalances(ctx, addr1))
	suite.Require().True(app.BankKeeper.GetAllBalances(ctx, addrModule).Empty())
}

func (suite *IntegrationTestSuite) TestUndelegateCoins_Invalid() {
	app, ctx := suite.app, suite.ctx

	origCoins := sdk.NewCoins(newFooCoin(100))
	delCoins := sdk.NewCoins(newFooCoin(50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addrModule := sdk.AccAddress([]byte("moduleAcc___________"))
	macc := app.AccountKeeper.NewAccountWithAddress(ctx, addrModule) // we don't need to define an actual module account bc we just need the address for testing
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)

	suite.Require().Error(app.BankKeeper.UndelegateCoins(ctx, addrModule, addr1, delCoins))

	app.AccountKeeper.SetAccount(ctx, macc)
	suite.Require().NoError(simapp.FundAccount(app, ctx, addr1, origCoins))

	suite.Require().Error(app.BankKeeper.UndelegateCoins(ctx, addrModule, addr1, delCoins))
	app.AccountKeeper.SetAccount(ctx, acc)

	suite.Require().Error(app.BankKeeper.UndelegateCoins(ctx, addrModule, addr1, delCoins))
}

func (suite *IntegrationTestSuite) TestSetDenomMetaData() {
	app, ctx := suite.app, suite.ctx

	metadata := suite.getTestMetadata()

	for i := range []int{1, 2} {
		app.BankKeeper.SetDenomMetaData(ctx, metadata[i])
	}

	actualMetadata, found := app.BankKeeper.GetDenomMetaData(ctx, metadata[1].Base)
	suite.Require().True(found)
	suite.Require().Equal(metadata[1].GetBase(), actualMetadata.GetBase())
	suite.Require().Equal(metadata[1].GetDisplay(), actualMetadata.GetDisplay())
	suite.Require().Equal(metadata[1].GetDescription(), actualMetadata.GetDescription())
	suite.Require().Equal(metadata[1].GetDenomUnits()[1].GetDenom(), actualMetadata.GetDenomUnits()[1].GetDenom())
	suite.Require().Equal(metadata[1].GetDenomUnits()[1].GetExponent(), actualMetadata.GetDenomUnits()[1].GetExponent())
	suite.Require().Equal(metadata[1].GetDenomUnits()[1].GetAliases(), actualMetadata.GetDenomUnits()[1].GetAliases())
}

func (suite *IntegrationTestSuite) TestIterateAllDenomMetaData() {
	app, ctx := suite.app, suite.ctx

	expectedMetadata := suite.getTestMetadata()
	// set metadata
	for i := range []int{1, 2} {
		app.BankKeeper.SetDenomMetaData(ctx, expectedMetadata[i])
	}
	// retrieve metadata
	actualMetadata := make([]types.Metadata, 0)
	app.BankKeeper.IterateAllDenomMetaData(ctx, func(metadata types.Metadata) bool {
		actualMetadata = append(actualMetadata, metadata)
		return false
	})
	// execute checks
	for i := range []int{1, 2} {
		suite.Require().Equal(expectedMetadata[i].GetBase(), actualMetadata[i].GetBase())
		suite.Require().Equal(expectedMetadata[i].GetDisplay(), actualMetadata[i].GetDisplay())
		suite.Require().Equal(expectedMetadata[i].GetDescription(), actualMetadata[i].GetDescription())
		suite.Require().Equal(expectedMetadata[i].GetDenomUnits()[1].GetDenom(), actualMetadata[i].GetDenomUnits()[1].GetDenom())
		suite.Require().Equal(expectedMetadata[i].GetDenomUnits()[1].GetExponent(), actualMetadata[i].GetDenomUnits()[1].GetExponent())
		suite.Require().Equal(expectedMetadata[i].GetDenomUnits()[1].GetAliases(), actualMetadata[i].GetDenomUnits()[1].GetAliases())
	}
}

func (suite *IntegrationTestSuite) TestBalanceTrackingEvents() {
	// replace account keeper and bank keeper otherwise the account keeper won't be aware of the
	// existence of the new module account because GetModuleAccount checks for the existence via
	// permissions map and not via state... weird
	maccPerms := simapp.GetMaccPerms()
	maccPerms[multiPerm] = []string{authtypes.Burner, authtypes.Minter, authtypes.Staking}

	suite.app.AccountKeeper = authkeeper.NewAccountKeeper(
		suite.app.AppCodec(), suite.app.GetKey(authtypes.StoreKey), suite.app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms,
	)

	suite.app.BankKeeper = keeper.NewBaseKeeper(suite.app.AppCodec(), suite.app.GetKey(types.StoreKey),
		suite.app.AccountKeeper, suite.app.GetSubspace(types.ModuleName), nil)

	// set account with multiple permissions
	suite.app.AccountKeeper.SetModuleAccount(suite.ctx, multiPermAcc)
	// mint coins
	suite.Require().NoError(
		suite.app.BankKeeper.MintCoins(
			suite.ctx,
			multiPermAcc.Name,
			sdk.NewCoins(sdk.NewCoin("utxo", sdk.NewInt(100000)))),
	)
	// send coins to address
	addr1 := sdk.BytesToAccAddress([]byte("addr1_______________"))
	suite.Require().NoError(
		suite.app.BankKeeper.SendCoinsFromModuleToAccount(
			suite.ctx,
			multiPermAcc.Name,
			addr1,
			sdk.NewCoins(sdk.NewCoin("utxo", sdk.NewInt(50000))),
		),
	)

	// burn coins from module account
	suite.Require().NoError(
		suite.app.BankKeeper.BurnCoins(
			suite.ctx,
			multiPermAcc.Name,
			sdk.NewCoins(sdk.NewInt64Coin("utxo", 1000)),
		),
	)

	// process balances and supply from events
	supply := sdk.NewCoins()

	balances := make(map[string]sdk.Coins)

	for _, e := range suite.ctx.EventManager().ABCIEvents() {
		switch e.Type {
		case types.EventTypeCoinBurn:
			burnedCoins, err := sdk.ParseCoinsNormalized((string)(e.Attributes[1].Value))
			suite.Require().NoError(err)
			supply = supply.Sub(burnedCoins)

		case types.EventTypeCoinMint:
			mintedCoins, err := sdk.ParseCoinsNormalized((string)(e.Attributes[1].Value))
			suite.Require().NoError(err)
			supply = supply.Add(mintedCoins...)

		case types.EventTypeCoinSpent:
			coinsSpent, err := sdk.ParseCoinsNormalized((string)(e.Attributes[1].Value))
			suite.Require().NoError(err)
			err = sdk.ValidateAccAddress(string(e.Attributes[0].Value))
			suite.Require().NoError(err)
			spender := sdk.AccAddress(e.Attributes[0].Value)
			balances[spender.String()] = balances[spender.String()].Sub(coinsSpent)

		case types.EventTypeCoinReceived:
			coinsRecv, err := sdk.ParseCoinsNormalized((string)(e.Attributes[1].Value))
			suite.Require().NoError(err)
			err = sdk.ValidateAccAddress(string(e.Attributes[0].Value))
			suite.Require().NoError(err)
			receiver := sdk.AccAddress(e.Attributes[0].Value)
			balances[receiver.String()] = balances[receiver.String()].Add(coinsRecv...)
		}
	}

	// check balance and supply tracking
	savedSupply := suite.app.BankKeeper.GetSupply(suite.ctx, "utxo")
	utxoSupply := savedSupply
	suite.Require().Equal(utxoSupply.Amount, supply.AmountOf("utxo"))
	// iterate accounts and check balances
	suite.app.BankKeeper.IterateAllBalances(suite.ctx, func(address sdk.AccAddress, coin sdk.Coin) (stop bool) {
		// if it's not utxo coin then skip
		if coin.Denom != "utxo" {
			return false
		}

		balance, exists := balances[address.String()]
		suite.Require().True(exists)

		expectedUtxo := sdk.NewCoin("utxo", balance.AmountOf(coin.Denom))
		suite.Require().Equal(expectedUtxo.String(), coin.String())
		return false
	})
}

func (suite *IntegrationTestSuite) getTestMetadata() []types.Metadata {
	return []types.Metadata{{
		Name:        "Cosmos Hub Atom",
		Symbol:      "ATOM",
		Description: "The native staking token of the Cosmos Hub.",
		DenomUnits: []*types.DenomUnit{
			{"uatom", uint32(0), []string{"microatom"}},
			{"matom", uint32(3), []string{"milliatom"}},
			{"atom", uint32(6), nil},
		},
		Base:    "uatom",
		Display: "atom",
	},
		{
			Name:        "Token",
			Symbol:      "TOKEN",
			Description: "The native staking token of the Token Hub.",
			DenomUnits: []*types.DenomUnit{
				{"1token", uint32(5), []string{"decitoken"}},
				{"2token", uint32(4), []string{"centitoken"}},
				{"3token", uint32(7), []string{"dekatoken"}},
			},
			Base:    "utoken",
			Display: "token",
		},
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
