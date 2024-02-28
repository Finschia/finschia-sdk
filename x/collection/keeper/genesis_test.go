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
	_, err := s.keeper.MintNFT(ctx, s.contractID, s.customer, []collection.MintNFTParam{{TokenType: s.nftClassID}})
	s.Require().NoError(err)

	s.keeper.Abandon(ctx, s.contractID, s.vendor, collection.PermissionMint)

	s.keeper.InitGenesis(ctx, genesis)

	// export again and compare
	newGenesis := s.keeper.ExportGenesis(s.ctx)
	s.Require().Equal(genesis, newGenesis)
}
