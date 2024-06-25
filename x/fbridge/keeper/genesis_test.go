package keeper_test

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func (s *IntegrationTestSuite) TestExportImportGenesis() {
	goctx := sdk.WrapSDKContext(s.ctx)
	const expProposalID uint64 = 5

	_, err := s.msgServer.Transfer(goctx, &types.MsgTransfer{
		Sender:   s.guardians[0].String(),
		Receiver: s.ethAddr,
		Amount:   sdk.NewInt(100),
	})
	s.Require().NoError(err)

	_, err = s.msgServer.SuggestRole(goctx, &types.MsgSuggestRole{
		From:   s.guardians[0].String(),
		Target: s.operator.String(),
		Role:   types.RoleJudge,
	})
	s.Require().NoError(err)

	_, err = s.msgServer.AddVoteForRole(goctx, &types.MsgAddVoteForRole{
		From:       s.guardians[0].String(),
		ProposalId: expProposalID,
		Option:     types.OptionYes,
	})
	s.Require().NoError(err)

	gen := s.app.FbridgeKeeper.ExportGenesis(s.ctx)
	gen.SendingState.SeqToBlocknum[0].Blocknum = 1
	err = types.ValidateGenesis(*gen)
	s.Require().NoError(err)

	err = s.app.FbridgeKeeper.InitGenesis(s.ctx, gen)
	s.Require().NoError(err)
}
