package keeper_test

import (
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestImportExportGenesis() {
	// export
	genesis := s.keeper.ExportGenesis(s.ctx)

	// forge
	err := s.keeper.Burn(s.ctx, s.classID, s.vendor, s.balance)
	s.Require().NoError(err)
	err = s.keeper.Mint(s.ctx, s.classID, s.vendor, s.customer, s.balance)
	s.Require().NoError(err)
	err = s.keeper.Abandon(s.ctx, s.classID, s.vendor, token.Permission_Mint)
	s.Require().NoError(err)

	// restore
	s.keeper.InitGenesis(s.ctx, genesis)

	// export again and compare
	newGenesis := s.keeper.ExportGenesis(s.ctx)
	s.Require().Equal(genesis, newGenesis)

	// nil class state
	s.keeper.InitGenesis(s.ctx, &token.GenesisState{})
}
