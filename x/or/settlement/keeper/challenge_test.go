package keeper_test

import (
	"github.com/Finschia/finschia-sdk/x/or/settlement/types"
)

func (s *KeeperTestSuite) TestGetChallenge() {
	testCases := map[string]struct {
		challengeID string
		expected    *types.Challenge
		err         error
	}{
		"valid request": {
			challengeID: "challenge_1",
			expected:    &types.Challenge{},
		},
		"not found (not existing challenge id)": {
			challengeID: "challenge_2",
			err:         types.ErrChallengeNotExist,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			if tc.expected != nil {
				s.app.SettlementKeeper.SetChallenge(s.ctx, tc.challengeID, *tc.expected)
			}
			challenge, err := s.app.SettlementKeeper.GetChallenge(s.ctx, tc.challengeID)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			s.Require().Equal(tc.expected, challenge)
		})
	}
}
