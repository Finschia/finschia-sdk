package testutil

import (
	"fmt"

	"github.com/cosmos/gogoproto/proto"

	"github.com/Finschia/finschia-sdk/testutil/rest"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (s *IntegrationTestSuite) TestGRPCParams() {
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
					FoundationTax: sdk.MustNewDecFromStr("0.2"),
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := rest.GetRequest(tc.url)
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
