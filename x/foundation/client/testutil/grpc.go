package testutil

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/line/lbm-sdk/testutil"
	"github.com/line/lbm-sdk/testutil/rest"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
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
				Params: &foundation.Params{
					Enabled:       true,
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

func (s *IntegrationTestSuite) TestGRPCValidatorAuth() {
	val := s.network.Validators[0]

	apiTemplate := "%s/lbm/foundation/v1/validators/%s"
	testCases := []struct {
		name   string
		url    string
		expErr bool
	}{
		{
			"with an empty validator address",
			fmt.Sprintf(apiTemplate, val.APIAddress, ""),
			true,
		},
		{
			"with an invalid validator address",
			fmt.Sprintf(apiTemplate, val.APIAddress, "this-is-an-invalid-address"),
			true,
		},
		{
			"valid request",
			fmt.Sprintf(apiTemplate, val.APIAddress, val.ValAddress),
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)

			var auth foundation.QueryValidatorAuthResponse
			err = s.cfg.Codec.UnmarshalJSON(resp, &auth)

			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(auth.Auth)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGRPCValidatorAuths() {
	val := s.network.Validators[0]

	testCases := []struct {
		name         string
		url          string
		headers      map[string]string
		wantNumAuths int
		expErr       bool
	}{
		{
			"valid request",
			fmt.Sprintf("%s/lbm/foundation/v1/validators", val.APIAddress),
			map[string]string{},
			len(s.network.Validators),
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			s.Require().NoError(err)

			var auths foundation.QueryValidatorAuthsResponse
			err = s.cfg.Codec.UnmarshalJSON(resp, &auths)

			if tc.expErr {
				s.Require().Empty(auths.Auths)
			} else {
				s.Require().NoError(err)
				s.Require().Len(auths.Auths, tc.wantNumAuths)
			}
		})
	}
}
