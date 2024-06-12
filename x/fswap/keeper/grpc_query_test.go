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
	"github.com/Finschia/finschia-sdk/types/query"
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

	keiSwapRateForCony, err := sdk.NewDecFromStr("148079656000000")
	s.Require().NoError(err)
	swapCap := sdk.NewInt(1000)
	s.Require().NoError(err)
	s.fromDenom = "cony"
	s.toDenom = "kei"
	s.swap = types.Swap{
		FromDenom:           s.fromDenom,
		ToDenom:             s.toDenom,
		AmountCapForToDenom: swapCap,
		SwapRate:            keiSwapRateForCony,
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
		request          *types.QuerySwapRequest
		expectedResponse *types.QuerySwapResponse
		expectedGrpcCode codes.Code
	}{
		{
			name: "valid",
			request: &types.QuerySwapRequest{
				FromDenom: s.fromDenom,
				ToDenom:   s.toDenom,
			},
			expectedResponse: &types.QuerySwapResponse{
				Swap: types.Swap{
					FromDenom:           s.swap.FromDenom,
					ToDenom:             s.swap.ToDenom,
					AmountCapForToDenom: s.swap.AmountCapForToDenom,
					SwapRate:            s.swap.SwapRate,
				},
			},
			expectedGrpcCode: codes.OK,
		},
		{
			name: "invalid: empty fromDenom",
			request: &types.QuerySwapRequest{
				FromDenom: "",
				ToDenom:   s.toDenom,
			},
			expectedGrpcCode: codes.InvalidArgument,
		},
		{
			name: "invalid: empty toDenom",
			request: &types.QuerySwapRequest{
				FromDenom: s.fromDenom,
				ToDenom:   "",
			},
			expectedGrpcCode: codes.InvalidArgument,
		},
		{
			name: "invalid: the same fromDenom and toDenom",
			request: &types.QuerySwapRequest{
				FromDenom: s.fromDenom,
				ToDenom:   s.fromDenom,
			},
			expectedGrpcCode: codes.InvalidArgument,
		},
		{
			name: "invalid: unregistered swap",
			request: &types.QuerySwapRequest{
				FromDenom: s.toDenom,
				ToDenom:   s.fromDenom,
			},
			expectedGrpcCode: codes.NotFound,
		},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			response, err := s.queryClient.Swap(s.ctx.Context(), tc.request)
			s.Require().Equal(tc.expectedResponse, response)
			actualGrpcCode := status.Code(err)
			s.Require().Equal(tc.expectedGrpcCode, actualGrpcCode, actualGrpcCode.String())
		})
	}
}

func (s *FSwapQueryTestSuite) TestQuerySwappedRequest() {
	tests := []struct {
		name             string
		request          *types.QuerySwappedRequest
		expectedResponse *types.QuerySwappedResponse
		expectedGrpcCode codes.Code
	}{
		{
			name: "valid",
			request: &types.QuerySwappedRequest{
				FromDenom: s.fromDenom,
				ToDenom:   s.toDenom,
			},
			expectedResponse: &types.QuerySwappedResponse{
				FromCoinAmount: sdk.NewCoin(s.fromDenom, sdk.ZeroInt()),
				ToCoinAmount:   sdk.NewCoin(s.toDenom, sdk.ZeroInt()),
			},
			expectedGrpcCode: codes.OK,
		},
		{
			name: "invalid: empty fromDenom",
			request: &types.QuerySwappedRequest{
				FromDenom: "",
				ToDenom:   s.toDenom,
			},
			expectedResponse: nil,
			expectedGrpcCode: codes.Unknown,
		},
		{
			name: "invalid: empty toDenom",
			request: &types.QuerySwappedRequest{
				FromDenom: s.fromDenom,
				ToDenom:   "",
			},
			expectedResponse: nil,
			expectedGrpcCode: codes.Unknown,
		},
		{
			name: "invalid: unregistered swap",
			request: &types.QuerySwappedRequest{
				FromDenom: s.toDenom,
				ToDenom:   s.fromDenom,
			},
			expectedResponse: nil,
			expectedGrpcCode: codes.Unknown,
		},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			response, err := s.queryClient.Swapped(s.ctx.Context(), tc.request)
			s.Require().Equal(tc.expectedResponse, response)
			actualGrpcCode := status.Code(err)
			s.Require().Equal(tc.expectedGrpcCode, actualGrpcCode, actualGrpcCode.String())
		})
	}
}

func (s *FSwapQueryTestSuite) TestQueryTotalSwappableToCoinAmountRequest() {
	tests := []struct {
		name             string
		request          *types.QueryTotalSwappableToCoinAmountRequest
		expectedResponse *types.QueryTotalSwappableToCoinAmountResponse
		expectedGrpcCode codes.Code
	}{
		{
			name: "valid",
			request: &types.QueryTotalSwappableToCoinAmountRequest{
				FromDenom: s.fromDenom,
				ToDenom:   s.toDenom,
			},
			expectedResponse: &types.QueryTotalSwappableToCoinAmountResponse{
				SwappableAmount: sdk.NewCoin(s.toDenom, s.swap.AmountCapForToDenom),
			},
			expectedGrpcCode: codes.OK,
		},
		{
			name: "invalid: empty fromDenom",
			request: &types.QueryTotalSwappableToCoinAmountRequest{
				FromDenom: "",
				ToDenom:   s.toDenom,
			},
			expectedResponse: nil,
			expectedGrpcCode: codes.Unknown,
		},
		{
			name: "invalid: empty toDenom",
			request: &types.QueryTotalSwappableToCoinAmountRequest{
				FromDenom: s.fromDenom,
				ToDenom:   "",
			},
			expectedResponse: nil,
			expectedGrpcCode: codes.Unknown,
		},
		{
			name: "invalid: unregistered swap",
			request: &types.QueryTotalSwappableToCoinAmountRequest{
				FromDenom: s.toDenom,
				ToDenom:   s.fromDenom,
			},
			expectedResponse: nil,
			expectedGrpcCode: codes.Unknown,
		},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			response, err := s.queryClient.TotalSwappableToCoinAmount(s.ctx.Context(), tc.request)
			s.Require().Equal(tc.expectedResponse, response)
			actualGrpcCode := status.Code(err)
			s.Require().Equal(tc.expectedGrpcCode, actualGrpcCode, actualGrpcCode.String())
		})
	}
}

func (s *FSwapQueryTestSuite) TestQuerySwapsRequest() {
	tests := []struct {
		name             string
		request          *types.QuerySwapsRequest
		expectedResponse *types.QuerySwapsResponse
		expectedGrpcCode codes.Code
	}{
		{
			name: "valid",
			request: &types.QuerySwapsRequest{
				Pagination: nil,
			},
			expectedResponse: &types.QuerySwapsResponse{
				Swaps: []types.Swap{
					{
						FromDenom:           s.swap.FromDenom,
						ToDenom:             s.swap.ToDenom,
						AmountCapForToDenom: s.swap.AmountCapForToDenom,
						SwapRate:            s.swap.SwapRate,
					},
				},
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   1,
				},
			},
			expectedGrpcCode: codes.OK,
		},
		{
			name: "invalid request",
			request: &types.QuerySwapsRequest{
				Pagination: &query.PageRequest{
					Key:        []byte("invalid-key"),
					Offset:     1,
					Limit:      0,
					CountTotal: false,
					Reverse:    false,
				},
			},
			expectedResponse: nil,
			expectedGrpcCode: codes.Unknown,
		},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			response, err := s.queryClient.Swaps(s.ctx.Context(), tc.request)
			s.Require().Equal(tc.expectedResponse, response)
			actualGrpcCode := status.Code(err)
			s.Require().Equal(tc.expectedGrpcCode, actualGrpcCode, actualGrpcCode.String())
		})
	}
}
