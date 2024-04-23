package keeper_test

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/testutil"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (s *KeeperTestSuite) TestInitAndExportGenesis() {
	s.keeper.InitGenesis(s.sdkCtx, s.bankKeeper, *types.DefaultGenesis())
	got := s.keeper.ExportGenesis(s.sdkCtx)
	s.Require().NotNil(got)
	s.Require().Equal(types.DefaultSwapped(), got.Swapped)

	gotTotalSupply := s.keeper.GetTotalSupply(s.sdkCtx)
	totalNewCoinsSupply := types.DefaultConfig().SwapRate.MulInt(testutil.MockOldCoin.Amount)
	expectedCoin := sdk.NewDecCoinFromDec(types.DefaultConfig().NewCoinDenom, totalNewCoinsSupply)
	s.Require().Equal(expectedCoin, gotTotalSupply)
}
