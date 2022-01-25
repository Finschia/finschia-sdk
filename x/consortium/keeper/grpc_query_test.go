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
	"github.com/line/lbm-sdk/x/consortium/types"
)

type ConsortiumTestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	queryClient types.QueryClient
}

func (suite *ConsortiumTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, ocproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.ConsortiumKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func (suite *ConsortiumTestSuite) TestQueryParams() {
	var (
		req         *types.QueryParamsRequest
		expResponse types.QueryParamsResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"with enabled",
			func() {
				params := &types.Params{Enabled: true}
				suite.app.ConsortiumKeeper.SetParams(suite.ctx, params)

				req = &types.QueryParamsRequest{}
				expResponse = types.QueryParamsResponse{Params: params}
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

func (suite *ConsortiumTestSuite) TestQueryValidatorAuth() {
	var (
		req         *types.QueryValidatorAuthRequest
		expResponse types.QueryValidatorAuthResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"with non-existent auth",
			func() {
				req = &types.QueryValidatorAuthRequest{ValidatorAddress: valAddr.String()}
				expResponse = types.QueryValidatorAuthResponse{}
			},
			false,
		},
		{
			"with existing auth",
			func() {
				auth := &types.ValidatorAuth{
					OperatorAddress: valAddr.String(),
					CreationAllowed: true,
				}
				suite.app.ConsortiumKeeper.SetValidatorAuth(suite.ctx, auth)

				req = &types.QueryValidatorAuthRequest{ValidatorAddress: valAddr.String()}
				expResponse = types.QueryValidatorAuthResponse{Auth: auth}
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

func (suite *ConsortiumTestSuite) TestQueryValidatorAuths() {
	var	req *types.QueryValidatorAuthsRequest
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
				req = &types.QueryValidatorAuthsRequest{}
			},
			true,
			0,
			false,
		},
		{
			"with empty auths",
			func() {
				req = &types.QueryValidatorAuthsRequest{
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
				auth := &types.ValidatorAuth{
					OperatorAddress: valAddr.String(),
					CreationAllowed: true,
				}
				suite.app.ConsortiumKeeper.SetValidatorAuth(suite.ctx, auth)

				req = &types.QueryValidatorAuthsRequest{
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

func TestConsortiumTestSuite(t *testing.T) {
	suite.Run(t, new(ConsortiumTestSuite))
}
