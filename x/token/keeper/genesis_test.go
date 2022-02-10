package keeper_test

import (
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestImportExportGenesis() {
	// export
	genesis := s.keeper.ExportGenesis(s.ctx)

	// forge
	err := s.keeper.Mint(s.ctx, s.vendor, s.customer, []token.FT{{ClassId: s.mintableClass, Amount: s.balance}})
	s.Require().NoError(err)
	err = s.keeper.Revoke(s.ctx, s.vendor, s.mintableClass, "mint")
	s.Require().NoError(err)

	// restore
	s.keeper.InitGenesis(s.ctx, genesis)

	// export again and compare
	newGenesis := s.keeper.ExportGenesis(s.ctx)
	s.Require().Equal(genesis, newGenesis)
}
