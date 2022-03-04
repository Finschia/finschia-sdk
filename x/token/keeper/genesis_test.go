package keeper_test

import (
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestImportExportGenesis() {
	// export
	genesis := s.keeper.ExportGenesis(s.ctx)

	// forge
	amount := token.FT{ClassId: s.classID, Amount: s.balance}
	err := s.keeper.Burn(s.ctx, s.vendor, []token.FT{amount})
	s.Require().NoError(err)
	err = s.keeper.Mint(s.ctx, s.vendor, s.customer, []token.FT{amount})
	s.Require().NoError(err)
	err = s.keeper.Revoke(s.ctx, s.vendor, s.classID, token.ActionMint)
	s.Require().NoError(err)

	// restore
	s.keeper.InitGenesis(s.ctx, genesis)

	// export again and compare
	newGenesis := s.keeper.ExportGenesis(s.ctx)
	s.Require().Equal(genesis, newGenesis)

	// nil class state
	s.keeper.InitGenesis(s.ctx, &token.GenesisState{})
}
