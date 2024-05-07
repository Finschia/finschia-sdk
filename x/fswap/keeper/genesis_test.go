package keeper_test

import (
	"fmt"

	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (s *KeeperTestSuite) TestInitAndExportGenesis() {
	ctx, _ := s.ctx.CacheContext()
	defaultGenesis := types.DefaultGenesis()
	err := s.keeper.InitGenesis(ctx, defaultGenesis)
	s.Require().NoError(err)

	exportGenesis := s.keeper.ExportGenesis(ctx)
	fmt.Println(len(exportGenesis.GetSwaps()))
	s.Require().Equal(defaultGenesis, exportGenesis)
	s.Require().Equal(defaultGenesis.GetSwaps(), exportGenesis.GetSwaps())
	s.Require().Equal(defaultGenesis.GetSwapStats(), exportGenesis.GetSwapStats())
	s.Require().Equal(defaultGenesis.GetSwappeds(), exportGenesis.GetSwappeds())
}
