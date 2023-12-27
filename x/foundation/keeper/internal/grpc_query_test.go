package internal_test

import (
	gocontext "context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/baseapp"
	"github.com/Finschia/finschia-sdk/simapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/Finschia/finschia-sdk/x/foundation/keeper"
	"github.com/Finschia/finschia-sdk/x/foundation/keeper/internal"
)

type FoundationTestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	queryClient foundation.QueryClient

	impl internal.Keeper
}

func (s *FoundationTestSuite) SetupTest() {
	s.app = simapp.Setup(false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, s.app.InterfaceRegistry())
	foundation.RegisterQueryServer(queryHelper, keeper.NewQueryServer(s.app.FoundationKeeper))
	s.queryClient = foundation.NewQueryClient(queryHelper)

	s.impl = internal.NewKeeper(
		s.app.AppCodec(),
		s.app.GetKey(foundation.ModuleName),
		s.app.MsgServiceRouter(),
		s.app.AccountKeeper,
		s.app.BankKeeper,
		authtypes.FeeCollectorName,
		foundation.DefaultConfig(),
		foundation.DefaultAuthority().String(),
		s.app.GetSubspace(foundation.ModuleName),
	)
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
					FoundationTax: sdk.OneDec(),
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
