package internal_test

import (
	"time"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (s *KeeperTestSuite) TestVote() {
	// no such a vote
	_, err := s.impl.GetVote(s.ctx, s.nextProposal, s.members[0])
	s.Require().Error(err)

	testCases := map[string]struct {
		proposalID uint64
		voter      sdk.AccAddress
		option     foundation.VoteOption
		metadata   string
		after      time.Duration
		valid      bool
	}{
		"vote yes": {
			proposalID: s.activeProposal,
			voter:      s.members[0],
			option:     foundation.VOTE_OPTION_YES,
			valid:      true,
		},
		"vote no": {
			proposalID: s.activeProposal,
			voter:      s.members[0],
			option:     foundation.VOTE_OPTION_NO,
			valid:      true,
		},
		"already voted": {
			proposalID: s.votedProposal,
			voter:      s.members[0],
			option:     foundation.VOTE_OPTION_YES,
		},
		"no such a proposal": {
			proposalID: s.nextProposal,
			voter:      s.members[0],
			option:     foundation.VOTE_OPTION_YES,
		},
		"inactive proposal": {
			proposalID: s.withdrawnProposal,
			voter:      s.members[0],
			option:     foundation.VOTE_OPTION_YES,
		},
		"long metadata": {
			proposalID: s.activeProposal,
			voter:      s.members[0],
			option:     foundation.VOTE_OPTION_YES,
			metadata:   string(make([]rune, 256)),
		},
		"voting too late": {
			proposalID: s.activeProposal,
			voter:      s.members[0],
			option:     foundation.VOTE_OPTION_YES,
			after:      s.impl.GetFoundationInfo(s.ctx).GetDecisionPolicy().GetVotingPeriod(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			ctx = ctx.WithBlockTime(ctx.BlockTime().Add(tc.after))

			vote := foundation.Vote{
				ProposalId: tc.proposalID,
				Voter:      tc.voter.String(),
				Option:     tc.option,
				Metadata:   tc.metadata,
			}
			err := s.impl.Vote(ctx, vote)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			_, err = s.impl.GetVote(ctx, vote.ProposalId, sdk.MustAccAddressFromBech32(vote.Voter))
			s.Require().NoError(err)
		})
	}
}
