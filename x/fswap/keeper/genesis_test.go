package keeper_test

import (
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (s *KeeperTestSuite) TestInitAndExportGenesis() {
	ctx, _ := s.ctx.CacheContext()
	defaultGenesis := types.DefaultGenesis()
	err := s.keeper.InitGenesis(ctx, defaultGenesis)
	s.Require().NoError(err)

	exportGenesis := s.keeper.ExportGenesis(ctx)
	s.Require().Equal(defaultGenesis, exportGenesis)
	s.Require().Equal(defaultGenesis.GetFswapInit(), exportGenesis.GetFswapInit())
	s.Require().Equal(defaultGenesis.GetSwapped(), exportGenesis.GetSwapped())
}
