package keeper_test

import (
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (s *KeeperTestSuite) TestInitAndExportGenesis() {
	s.keeper.InitGenesis(s.sdkCtx, *types.DefaultGenesis())
	got := s.keeper.ExportGenesis(s.sdkCtx)
	s.Require().NotNil(got)
	s.Require().Equal(types.DefaultParams(), got.Params)
	s.Require().Equal(types.DefaultSwapped(), got.Swapped)
}
