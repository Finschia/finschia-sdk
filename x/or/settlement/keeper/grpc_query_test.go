package keeper_test

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/or/settlement/types"
)

func (s *KeeperTestSuite) TestChallenge() {
	ctx := sdk.WrapSDKContext(s.ctx)

	testCases := map[string]struct {
		req         *types.QueryChallengeRequest
		expectedRes *types.QueryChallengeResponse
		err         error
	}{
		"valid request": {
			req:         &types.QueryChallengeRequest{ChallengeId: "challenge_1"},
			expectedRes: &types.QueryChallengeResponse{Challenge: &types.Challenge{}},
		},
		"nil request": {
			err: status.Error(codes.InvalidArgument, "empty request"),
		},
		"non-existent challenge id": {
			req: &types.QueryChallengeRequest{ChallengeId: "challenge_2"},
			err: status.Error(codes.NotFound, types.ErrChallengeNotExist.Wrapf("no challenge for %s", "challenge_2").Error()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			if tc.expectedRes != nil {
				s.app.SettlementKeeper.SetChallenge(s.ctx, tc.req.ChallengeId, *tc.expectedRes.Challenge)
			}
			res, err := s.app.SettlementKeeper.Challenge(ctx, tc.req)
			if tc.err != nil {
				s.Require().EqualError(tc.err, err.Error())
				return
			}
			s.Require().NoError(err)
			s.Require().Equal(tc.expectedRes, res)
		})
	}
}
