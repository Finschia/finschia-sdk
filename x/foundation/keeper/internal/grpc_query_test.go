package internal_test

import (
	gocontext "context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/Finschia/finschia-sdk/x/foundation/keeper"
	"github.com/Finschia/finschia-sdk/x/foundation/keeper/internal"
)

type FoundationTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	queryClient foundation.QueryClient

	impl internal.Keeper
}

func (s *FoundationTestSuite) SetupTest() {
	var encCfg moduletestutil.TestEncodingConfig
	var k keeper.Keeper
	s.impl, k, _, _, encCfg, _, _, s.ctx = setupFoundationKeeper(s.T(), nil, nil)

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, encCfg.InterfaceRegistry)
	foundation.RegisterQueryServer(queryHelper, keeper.NewQueryServer(k))
	s.queryClient = foundation.NewQueryClient(queryHelper)
}

func (s *FoundationTestSuite) TestQueryParams() {
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
					FoundationTax: math.LegacyOneDec(),
				}
				s.impl.SetParams(s.ctx, params)

				req = &foundation.QueryParamsRequest{}
				expResponse = foundation.QueryParamsResponse{Params: params}
			},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			s.SetupTest() // reset

			tc.malleate()

			res, err := s.queryClient.Params(gocontext.Background(), req)

			if tc.expPass {
				s.Require().NoError(err)
				s.Require().NotNil(res)
				s.Require().Equal(&expResponse, res)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func TestFoundationTestSuite(t *testing.T) {
	suite.Run(t, new(FoundationTestSuite))
}
