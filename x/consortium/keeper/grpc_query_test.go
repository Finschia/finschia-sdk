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
	"github.com/line/lbm-sdk/x/consortium"
	"github.com/line/lbm-sdk/x/consortium/keeper"
)

type ConsortiumTestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	queryClient consortium.QueryClient
}

func (suite *ConsortiumTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, ocproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	consortium.RegisterQueryServer(queryHelper, keeper.NewQueryServer(suite.app.ConsortiumKeeper))
	suite.queryClient = consortium.NewQueryClient(queryHelper)
}

func (suite *ConsortiumTestSuite) TestQueryParams() {
	var (
		req         *consortium.QueryParamsRequest
		expResponse consortium.QueryParamsResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"with enabled",
			func() {
				params := &consortium.Params{Enabled: true}
				suite.app.ConsortiumKeeper.SetParams(suite.ctx, params)

				req = &consortium.QueryParamsRequest{}
				expResponse = consortium.QueryParamsResponse{Params: params}
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
		req         *consortium.QueryValidatorAuthRequest
		expResponse consortium.QueryValidatorAuthResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"with non-existent auth",
			func() {
				req = &consortium.QueryValidatorAuthRequest{ValidatorAddress: valAddr.String()}
				expResponse = consortium.QueryValidatorAuthResponse{}
			},
			false,
		},
		{
			"with existing auth",
			func() {
				auth := &consortium.ValidatorAuth{
					OperatorAddress: valAddr.String(),
					CreationAllowed: true,
				}
				suite.app.ConsortiumKeeper.SetValidatorAuth(suite.ctx, auth)

				req = &consortium.QueryValidatorAuthRequest{ValidatorAddress: valAddr.String()}
				expResponse = consortium.QueryValidatorAuthResponse{Auth: auth}
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
	var req *consortium.QueryValidatorAuthsRequest
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
				req = &consortium.QueryValidatorAuthsRequest{}
			},
			true,
			0,
			false,
		},
		{
			"with empty auths",
			func() {
				req = &consortium.QueryValidatorAuthsRequest{
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
				auth := &consortium.ValidatorAuth{
					OperatorAddress: valAddr.String(),
					CreationAllowed: true,
				}
				suite.app.ConsortiumKeeper.SetValidatorAuth(suite.ctx, auth)

				req = &consortium.QueryValidatorAuthsRequest{
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
