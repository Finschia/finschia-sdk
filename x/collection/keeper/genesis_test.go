package keeper_test

import (
	"github.com/line/lbm-sdk/x/collection"
)

func (s *KeeperTestSuite) TestImportExportGenesis() {
	// export
	genesis := s.keeper.ExportGenesis(s.ctx)

	// forge
	amount := collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance))
	err := s.keeper.SendCoins(s.ctx, s.contractID, s.vendor, s.customer, amount)
	s.Require().NoError(err)

	err = s.keeper.SendCoins(s.ctx, s.contractID, s.customer, s.operator, amount)
	s.Require().NoError(err)

	_, err = s.keeper.BurnCoins(s.ctx, s.contractID, s.operator, amount)
	s.Require().NoError(err)

	s.keeper.Abandon(s.ctx, s.contractID, s.vendor, collection.PermissionMint)

	// restore
	s.keeper.InitGenesis(s.ctx, genesis)

	// export again and compare
	newGenesis := s.keeper.ExportGenesis(s.ctx)
	s.Require().Equal(genesis, newGenesis)
}
