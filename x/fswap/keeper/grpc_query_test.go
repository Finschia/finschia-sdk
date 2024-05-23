package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Finschia/finschia-sdk/baseapp"
	"github.com/Finschia/finschia-sdk/simapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	bank "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/fswap/keeper"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func TestFSwapQueryTestSuite(t *testing.T) {
	suite.Run(t, &FSwapQueryTestSuite{})
}

type FSwapQueryTestSuite struct {
	suite.Suite

	app             *simapp.SimApp
	ctx             sdk.Context
	queryClient     types.QueryClient
	keeper          keeper.Keeper
	swap            types.Swap
	toDenomMetadata bank.Metadata
	fromDenom       string
	toDenom         string
}

func (s *FSwapQueryTestSuite) SetupTest() {
	s.app = simapp.Setup(false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, s.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQueryServer(s.app.FswapKeeper))
	s.queryClient = types.NewQueryClient(queryHelper)
	s.keeper = s.app.FswapKeeper

	pebSwapRateForCony, err := sdk.NewDecFromStr("148079656000000")
	s.Require().NoError(err)
	swapCap := sdk.NewInt(1000)
	s.Require().NoError(err)
	s.fromDenom = "cony"
	s.toDenom = "kei"
	s.swap = types.Swap{
		FromDenom:           s.fromDenom,
		ToDenom:             s.toDenom,
		AmountCapForToDenom: swapCap,
		SwapRate:            pebSwapRateForCony,
	}
	s.toDenomMetadata = bank.Metadata{
		Description: "This is metadata for to-coin",
		DenomUnits: []*bank.DenomUnit{
			{Denom: s.toDenom, Exponent: 0},
		},
		Base:    s.toDenom,
		Display: s.toDenom,
		Name:    "DUMMY",
		Symbol:  "DUM",
	}
	err = s.toDenomMetadata.Validate()
	s.Require().NoError(err)

	fromDenom := bank.Metadata{
		Description: "This is metadata for from-coin",
		DenomUnits: []*bank.DenomUnit{
			{Denom: s.fromDenom, Exponent: 0},
		},
		Base:    s.fromDenom,
		Display: s.fromDenom,
		Name:    "FROM",
		Symbol:  "FROM",
	}
	err = fromDenom.Validate()
	s.Require().NoError(err)

	s.app.BankKeeper.SetDenomMetaData(s.ctx, fromDenom)
	err = s.keeper.SetSwap(s.ctx, s.swap, s.toDenomMetadata)
	s.Require().NoError(err)
}

func (s *FSwapQueryTestSuite) TestQuerySwapRequest() {
	tests := []struct {
		name             string
		FromDenom        string
		ToDenom          string
		wantErr          bool
		expectedGrpcCode codes.Code
	}{
		{
			name:             "valid",
			FromDenom:        s.fromDenom,
			ToDenom:          s.toDenom,
			wantErr:          false,
			expectedGrpcCode: codes.OK,
		},
		{
			name:             "invalid: empty fromDenom",
			FromDenom:        "",
			ToDenom:          s.toDenom,
			wantErr:          true,
			expectedGrpcCode: codes.InvalidArgument,
		},
		{
			name:             "invalid: empty toDenom",
			FromDenom:        s.fromDenom,
			ToDenom:          "",
			wantErr:          true,
			expectedGrpcCode: codes.InvalidArgument,
		},
		{
			name:             "invalid: the same fromDenom and toDenom",
			FromDenom:        s.fromDenom,
			ToDenom:          s.fromDenom,
			wantErr:          true,
			expectedGrpcCode: codes.InvalidArgument,
		},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			m := &types.QuerySwapRequest{
				FromDenom: tc.FromDenom,
				ToDenom:   tc.ToDenom,
			}

			_, err := s.queryClient.Swap(s.ctx.Context(), m)
			actualGrpcCode := status.Code(err)
			s.Require().Equal(tc.expectedGrpcCode, actualGrpcCode, actualGrpcCode.String())
		})
	}
}
