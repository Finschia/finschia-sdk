package keeper_test

import (
	"github.com/Finschia/finschia-sdk/x/collection"
)

// TODO: Add more test cases
func (s *KeeperTestSuite) TestImportExportGenesis() {
	// export
	genesis := s.keeper.ExportGenesis(s.ctx)

	// forge & import
	ctx, _ := s.ctx.CacheContext()
	amount := collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance))
	err := s.keeper.SendCoins(ctx, s.contractID, s.vendor, s.customer, amount)
	s.Require().NoError(err)

	err = s.keeper.SendCoins(ctx, s.contractID, s.customer, s.operator, amount)
	s.Require().NoError(err)

	_, err = s.keeper.BurnCoins(ctx, s.contractID, s.operator, amount)
	s.Require().NoError(err)

	s.keeper.Abandon(ctx, s.contractID, s.vendor, collection.PermissionMint)

	s.keeper.InitGenesis(ctx, genesis)

	// export again and compare
	newGenesis := s.keeper.ExportGenesis(s.ctx)
	s.Require().Equal(genesis, newGenesis)
}
