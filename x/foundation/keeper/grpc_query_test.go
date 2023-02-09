package keeper_test

import (
	gocontext "context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/line/lbm-sdk/baseapp"
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
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
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})

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
			"with foundation tax",
			func() {
				params := foundation.Params{
					FoundationTax: sdk.OneDec(),
				}
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

func TestFoundationTestSuite(t *testing.T) {
	suite.Run(t, new(FoundationTestSuite))
}
