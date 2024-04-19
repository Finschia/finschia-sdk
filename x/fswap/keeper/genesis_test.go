package keeper_test

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/testutil"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (s *KeeperTestSuite) TestInitAndExportGenesis() {
	err := s.keeper.InitGenesis(s.sdkCtx, s.bankKeeper, *types.DefaultGenesis())
	s.Require().NoError(err)
	got := s.keeper.ExportGenesis(s.sdkCtx)
	s.Require().NotNil(got)
	s.Require().Equal(types.DefaultParams(), got.Params)
	s.Require().Equal(types.DefaultSwapped(), got.Swapped)

	gotTotalSupply := s.keeper.GetTotalSupply(s.sdkCtx)
	totalNewCoinsSupply := types.DefaultSwapRate.MulInt(testutil.MockOldCoin.Amount)
	expectedCoin := sdk.NewDecCoinFromDec(types.DefaultNewCoinDenom, totalNewCoinsSupply)
	s.Require().Equal(expectedCoin, gotTotalSupply)
}
