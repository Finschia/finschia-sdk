package internal_test

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/Finschia/finschia-sdk/x/foundation/keeper/internal"
)

func (s *KeeperTestSuite) TestBeginBlocker() {
	for name, tc := range map[string]struct {
		taxRatio math.LegacyDec
		valid    bool
	}{
		"valid ratio": {
			taxRatio: math.LegacyOneDec(),
			valid:    true,
		},
		"ratio > 1": {
			taxRatio: math.LegacyMustNewDecFromStr("1.00000001"),
		},
	} {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			// collect
			testing := func() {
				s.impl.SetParams(ctx, foundation.Params{
					FoundationTax: tc.taxRatio,
				})
				internal.BeginBlocker(ctx, s.impl)
			}
			if !tc.valid {
				s.Require().Panics(testing)
				return
			}
			s.Require().NotPanics(testing)

			if s.deterministic {
				expectedEvents := sdk.Events{} // TODO
				s.Require().Equal(expectedEvents, ctx.EventManager().Events())
			}
		})
	}
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
		s.Run(name, func() {
			proposal, err := s.impl.GetProposal(ctx, tc.id)
			s.Require().NoError(err)
			s.Require().NotNil(proposal)
			s.Require().Equal(tc.status, proposal.Status)
		})
	}

	// voting periods end
	votingPeriod := s.impl.GetFoundationInfo(ctx).GetDecisionPolicy().GetVotingPeriod()
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(votingPeriod))
	internal.EndBlocker(ctx, s.impl)

	for name, tc := range map[string]struct {
		id      uint64
		removed bool
		status  foundation.ProposalStatus
	}{
		"active proposal": {
			id:     s.activeProposal,
			status: foundation.PROPOSAL_STATUS_ACCEPTED,
		},
		"voted proposal": {
			id:     s.votedProposal,
			status: foundation.PROPOSAL_STATUS_REJECTED,
		},
		"withdrawn proposal": {
			id:      s.withdrawnProposal,
			removed: true,
		},
		"invalid proposal": {
			id:     s.invalidProposal,
			status: foundation.PROPOSAL_STATUS_ACCEPTED,
		},
	} {
		s.Run(name, func() {
			proposal, err := s.impl.GetProposal(ctx, tc.id)
			if tc.removed {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(proposal)
			s.Require().Equal(tc.status, proposal.Status)
		})
	}

	// proposals expire
	maxExecutionPeriod := foundation.DefaultConfig().MaxExecutionPeriod
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(maxExecutionPeriod))
	internal.EndBlocker(ctx, s.impl)

	// all proposals must be pruned
	s.Require().Empty(s.impl.GetProposals(ctx))

	if s.deterministic {
		expectedEvents := sdk.Events{}
		s.Require().Equal(expectedEvents, ctx.EventManager().Events())
	}
}
