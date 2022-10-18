package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/foundation/keeper"
)

func (s *KeeperTestSuite) TestBeginBlocker() {
	ctx, _ := s.ctx.CacheContext()

	s.keeper.SetParams(ctx, foundation.Params{
		FoundationTax: sdk.MustNewDecFromStr("0.5"),
		CensoredMsgTypeUrls: []string{
			sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
		},
	})

	before := s.keeper.GetTreasury(ctx)
	s.Require().Equal(1, len(before))
	s.Require().Equal(sdk.NewDecFromInt(s.balance), before[0].Amount)

	// collect
	keeper.BeginBlocker(ctx, s.keeper)

	after := s.keeper.GetTreasury(ctx)
	s.Require().Equal(1, len(after))
	// s.balance + s.balance * 0.5
	s.Require().Equal(sdk.NewDecFromInt(s.balance.Add(s.balance.Quo(sdk.NewInt(2)))), after[0].Amount)
}

func (s *KeeperTestSuite) TestEndBlocker() {
	ctx, _ := s.ctx.CacheContext()

	// check preconditions
	for name, tc := range map[string]struct {
		id     uint64
		status foundation.ProposalStatus
	}{
		"active proposal": {
			s.activeProposal,
			foundation.PROPOSAL_STATUS_SUBMITTED,
		},
		"voted proposal": {
			s.votedProposal,
			foundation.PROPOSAL_STATUS_SUBMITTED,
		},
		"withdrawn proposal": {
			s.withdrawnProposal,
			foundation.PROPOSAL_STATUS_WITHDRAWN,
		},
		"invalid proposal": {
			s.invalidProposal,
			foundation.PROPOSAL_STATUS_SUBMITTED,
		},
	} {
		proposal, err := s.keeper.GetProposal(ctx, tc.id)
		s.Require().NoError(err, name)
		s.Require().NotNil(proposal, name)
		s.Require().Equal(tc.status, proposal.Status, name)
	}

	// voting periods end
	votingPeriod := s.keeper.GetFoundationInfo(ctx).GetDecisionPolicy().GetVotingPeriod()
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(votingPeriod))
	keeper.EndBlocker(ctx, s.keeper)

	for name, tc := range map[string]struct {
		id     uint64
		status foundation.ProposalStatus
	}{
		"active proposal": {
			s.activeProposal,
			foundation.PROPOSAL_STATUS_ACCEPTED,
		},
		"voted proposal": {
			s.votedProposal,
			foundation.PROPOSAL_STATUS_REJECTED,
		},
		"withdrawn proposal": {
			s.withdrawnProposal,
			foundation.PROPOSAL_STATUS_WITHDRAWN,
		},
		"invalid proposal": {
			s.invalidProposal,
			foundation.PROPOSAL_STATUS_ACCEPTED,
		},
	} {
		proposal, err := s.keeper.GetProposal(ctx, tc.id)
		s.Require().NoError(err, name)
		s.Require().NotNil(proposal, name)
		s.Require().Equal(tc.status, proposal.Status, name)
	}

	// proposals expire
	maxExecutionPeriod := foundation.DefaultConfig().MaxExecutionPeriod
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(maxExecutionPeriod))
	keeper.EndBlocker(ctx, s.keeper)

	// all proposals must be pruned
	s.Require().Empty(s.keeper.GetProposals(ctx))
}
