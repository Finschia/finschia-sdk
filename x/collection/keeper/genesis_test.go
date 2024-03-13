package keeper_test

import (
	"github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/collection"
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

func (s *KeeperTestSuite) TestShouldPanicWhenZeroNextTokenIdInGenesis() {
	dontCare := "12345678"
	state := &collection.GenesisState{
		Params: collection.Params{},
		NextTokenIds: []collection.ContractNextTokenIDs{
			{ContractId: dontCare, TokenIds: []collection.NextTokenID{
				{ClassId: dontCare, Id: types.NewUint(0)},
			}},
		},
	}

	s.Require().Panics(func() { s.keeper.InitGenesis(s.ctx, state) })
}
