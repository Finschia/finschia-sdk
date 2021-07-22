package keeper_test

import (
	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/x/bank/types"
)

func (suite *IntegrationTestSuite) TestExportGenesis() {
	app, ctx := suite.app, suite.ctx

	expectedMetadata := suite.getTestMetadata()
	expectedBalances := suite.getTestBalances()
	for i := range []int{1, 2} {
		app.BankKeeper.SetDenomMetaData(ctx, expectedMetadata[i])
		err1 := sdk.ValidateAccAddress(expectedBalances[i].Address)
		if err1 != nil {
			panic(err1)
		}
		err := app.BankKeeper.SetBalances(ctx, sdk.AccAddress(expectedBalances[i].Address), expectedBalances[i].Coins)
		suite.Require().NoError(err)
	}

	totalSupply := types.NewSupply(sdk.NewCoins(sdk.NewInt64Coin("test", 400000000)))
	app.BankKeeper.SetSupply(ctx, totalSupply)
	app.BankKeeper.SetParams(ctx, types.DefaultParams())

	exportGenesis := app.BankKeeper.ExportGenesis(ctx)

	suite.Require().Len(exportGenesis.Params.SendEnabled, 0)
	suite.Require().Equal(types.DefaultParams().DefaultSendEnabled, exportGenesis.Params.DefaultSendEnabled)
	suite.Require().Equal(totalSupply.GetTotal(), exportGenesis.Supply)
	suite.Require().Equal(expectedBalances, exportGenesis.Balances)
	suite.Require().Equal(expectedMetadata, exportGenesis.DenomMetadata)
}

func (suite *IntegrationTestSuite) getTestBalances() []types.Balance {
	addr2 := sdk.AccAddress("link15klks9yty6rwvnqk47q4cg92r38qj4gsvxdfj8")
	addr1 := sdk.AccAddress("link14uxsrqf2cakphyw9ywwy9a9fv7yjspku4lkrny")
	return []types.Balance{
		{Address: addr1.String(), Coins: sdk.Coins{sdk.NewInt64Coin("testcoin3", 10)}},
		{Address: addr2.String(), Coins: sdk.Coins{sdk.NewInt64Coin("testcoin1", 32), sdk.NewInt64Coin("testcoin2", 34)}},
	}

}

func (suite *IntegrationTestSuite) TestInitGenesis() {
	m := types.Metadata{Description: sdk.DefaultBondDenom, Base: sdk.DefaultBondDenom, Display: sdk.DefaultBondDenom}
	g := types.DefaultGenesisState()
	g.DenomMetadata = []types.Metadata{m}
	bk := suite.app.BankKeeper
	bk.InitGenesis(suite.ctx, g)

	m2 := bk.GetDenomMetaData(suite.ctx, m.Base)
	suite.Require().Equal(m, m2)
}
