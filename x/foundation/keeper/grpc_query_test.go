package keeper_test

import (
	gocontext "context"
	"fmt"
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/baseapp"
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/query"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/foundation/keeper"
)

type FoundationTestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	queryClient foundation.QueryClient
}

func (suite *FoundationTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, ocproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	foundation.RegisterQueryServer(queryHelper, keeper.NewQueryServer(suite.app.FoundationKeeper))
	suite.queryClient = foundation.NewQueryClient(queryHelper)
}

func (suite *FoundationTestSuite) TestQueryParams() {
	var (
		req         *foundation.QueryParamsRequest
		expResponse foundation.QueryParamsResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"with enabled",
			func() {
				params := &foundation.Params{Enabled: true}
				suite.app.FoundationKeeper.SetParams(suite.ctx, params)

				req = &foundation.QueryParamsRequest{}
				expResponse = foundation.QueryParamsResponse{Params: params}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()

			res, err := suite.queryClient.Params(gocontext.Background(), req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(&expResponse, res)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *FoundationTestSuite) TestQueryValidatorAuth() {
	var (
		req         *foundation.QueryValidatorAuthRequest
		expResponse foundation.QueryValidatorAuthResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"with non-existent auth",
			func() {
				req = &foundation.QueryValidatorAuthRequest{ValidatorAddress: valAddr.String()}
				expResponse = foundation.QueryValidatorAuthResponse{}
			},
			false,
		},
		{
			"with existing auth",
			func() {
				auth := foundation.ValidatorAuth{
					OperatorAddress: valAddr.String(),
					CreationAllowed: true,
				}
				suite.app.FoundationKeeper.SetValidatorAuth(suite.ctx, auth)

				req = &foundation.QueryValidatorAuthRequest{ValidatorAddress: valAddr.String()}
				expResponse = foundation.QueryValidatorAuthResponse{Auth: &auth}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()

			res, err := suite.queryClient.ValidatorAuth(gocontext.Background(), req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(&expResponse, res)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *FoundationTestSuite) TestQueryValidatorAuths() {
	var req *foundation.QueryValidatorAuthsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		numVals  int
		hasNext  bool
	}{
		{
			"empty request",
			func() {
				req = &foundation.QueryValidatorAuthsRequest{}
			},
			true,
			0,
			false,
		},
		{
			"with empty auths",
			func() {
				req = &foundation.QueryValidatorAuthsRequest{
					Pagination: &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
			0,
			false,
		},
		{
			"with non-empty auths",
			func() {
				auth := foundation.ValidatorAuth{
					OperatorAddress: valAddr.String(),
					CreationAllowed: true,
				}
				suite.app.FoundationKeeper.SetValidatorAuth(suite.ctx, auth)

				req = &foundation.QueryValidatorAuthsRequest{
					Pagination: &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
			1,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()

			res, err := suite.queryClient.ValidatorAuths(gocontext.Background(), req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(tc.numVals, len(res.Auths))

				if tc.hasNext {
					suite.Require().NotNil(res.Pagination.NextKey)
				} else {
					suite.Require().Nil(res.Pagination.NextKey)
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func TestFoundationTestSuite(t *testing.T) {
	suite.Run(t, new(FoundationTestSuite))
}
