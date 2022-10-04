package keeper_test

func (s *KeeperTestSuite) TestExec() {
	testCases := map[string]struct {
		proposalID uint64
		valid      bool
	}{
		"valid exec": {
			proposalID: s.votedProposal,
			valid:      true,
		},
		"not enough votes": {
			proposalID: s.activeProposal,
			valid:      true,
		},
		"invalid msg in proposal": {
			proposalID: s.invalidProposal,
			valid:      true,
		},
		"withdrawn proposal": {
			proposalID: s.withdrawnProposal,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.Exec(ctx, tc.proposalID)
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}
