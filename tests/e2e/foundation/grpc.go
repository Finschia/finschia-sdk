package foundation

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/testutil"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (s *E2ETestSuite) TestGRPCParams() {
	val := s.network.Validators[0]

	testCases := []struct {
		name       string
		url        string
		expErr     bool
		respType   proto.Message
		expectResp proto.Message
	}{
		{
			"valid request",
			fmt.Sprintf("%s/lbm/foundation/v1/params", val.APIAddress),
			false,
			&foundation.QueryParamsResponse{},
			&foundation.QueryParamsResponse{
				Params: foundation.Params{
					FoundationTax: math.LegacyMustNewDecFromStr("0.2"),
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequest(tc.url)
			s.Require().NoError(err)

			err = s.cfg.Codec.UnmarshalJSON(resp, tc.respType)

			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectResp.String(), tc.respType.String())
			}
		})
	}
}
