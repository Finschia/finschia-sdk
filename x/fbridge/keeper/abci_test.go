package keeper_test

import (
	"github.com/Finschia/finschia-sdk/simapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func (s *IntegrationTestSuite) TestBeginBlocker() {
	dummy := simapp.AddTestAddrs(s.app, s.ctx, 1, sdk.NewInt(1000000000))[0]
	_, err := s.app.FbridgeKeeper.RegisterRoleProposal(s.ctx, s.guardians[0], dummy, types.RoleGuardian)
	s.Require().NoError(err)

	bh := s.ctx.BlockHeader()
	bh.Time = s.ctx.BlockHeader().Time.AddDate(0, 0, 1)
	s.ctx = s.ctx.WithBlockHeader(bh)
	s.app.FbridgeKeeper.BeginBlocker(s.ctx)

	_, err = s.app.FbridgeKeeper.RegisterRoleProposal(s.ctx, s.guardians[0], dummy, types.RoleGuardian)
	s.Require().NoError(err)
}
