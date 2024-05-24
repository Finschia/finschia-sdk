package testutil

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/Finschia/finschia-sdk/testutil"
	"github.com/Finschia/finschia-sdk/testutil/rest"
	sdk "github.com/Finschia/finschia-sdk/types"
	grpctypes "github.com/Finschia/finschia-sdk/types/grpc"
	"github.com/Finschia/finschia-sdk/types/query"
	fswaptypes "github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (s *IntegrationTestSuite) TestGRPCQuerySwap() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name        string
		url         string
		expectedErr bool
		expected    proto.Message
	}{
		{
			"test query swap with valid query string",
			fmt.Sprintf("%s/lbm/fswap/v1/swap?from_denom=stake&to_denom=dummy", baseURL),
			false,
			&fswaptypes.QuerySwapResponse{
				Swap: s.dummySwap,
			},
		},
		{
			"test query swap with not existed swap pairs",
			fmt.Sprintf("%s/lbm/fswap/v1/swap?from_denom=fake&to_denom=dummy", baseURL),
			true,
			&fswaptypes.QuerySwapResponse{},
		},
		{
			"test query swap with nil to_denom",
			fmt.Sprintf("%s/lbm/fswap/v1/swap?from_denom=stake", baseURL),
			true,
			&fswaptypes.QuerySwapResponse{},
		},
		{
			"test query swap with nil from_denom",
			fmt.Sprintf("%s/lbm/fswap/v1/swap?to_denom=dummy", baseURL),
			true,
			&fswaptypes.QuerySwapResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)
			var valRes fswaptypes.QuerySwapResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &valRes)

			if tc.expectedErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expected.String(), valRes.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGRPCQuerySwaps() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name        string
		url         string
		expectedErr bool
		expected    proto.Message
	}{
		{
			"test query swaps",
			fmt.Sprintf("%s/lbm/fswap/v1/swaps", baseURL),
			false,
			&fswaptypes.QuerySwapsResponse{
				Swaps:      []fswaptypes.Swap{s.dummySwap},
				Pagination: &query.PageResponse{Total: 1},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)
			var valRes fswaptypes.QuerySwapsResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &valRes)

			if tc.expectedErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expected.String(), valRes.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGRPCQuerySwapped() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name        string
		url         string
		expectedErr bool
		expected    proto.Message
	}{
		{
			"test query swapped with valid query string",
			fmt.Sprintf("%s/lbm/fswap/v1/swapped?from_denom=stake&to_denom=dummy", baseURL),
			false,
			&fswaptypes.QuerySwappedResponse{
				FromCoinAmount: sdk.NewCoin("stake", sdk.ZeroInt()),
				ToCoinAmount:   sdk.NewCoin("dummy", sdk.ZeroInt()),
			},
		},
		{
			"test query swapped with not existed swap pairs",
			fmt.Sprintf("%s/lbm/fswap/v1/swapped?from_denom=fake&to_denom=dummy", baseURL),
			true,
			&fswaptypes.QuerySwappedResponse{},
		},
		{
			"test query swapped with nil to_denom",
			fmt.Sprintf("%s/lbm/fswap/v1/swapped?from_denom=stake", baseURL),
			true,
			&fswaptypes.QuerySwappedResponse{},
		},
		{
			"test query swapped with nil from_denom",
			fmt.Sprintf("%s/lbm/fswap/v1/swapped?to_denom=dummy", baseURL),
			true,
			&fswaptypes.QuerySwappedResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			})
			s.Require().NoError(err)
			var valRes fswaptypes.QuerySwappedResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &valRes)

			if tc.expectedErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expected.String(), valRes.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGRPCQueryTotalSwappableAmount() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name        string
		url         string
		expectedErr bool
		expected    proto.Message
	}{
		{
			"test query total_swappable_to_coin_amount with valid query string",
			fmt.Sprintf("%s/lbm/fswap/v1/total_swappable_to_coin_amount?from_denom=stake&to_denom=dummy", baseURL),
			false,
			&fswaptypes.QueryTotalSwappableToCoinAmountResponse{
				SwappableAmount: sdk.NewCoin("dummy", s.dummySwap.AmountCapForToDenom),
			},
		},
		{
			"test query total_swappable_to_coin_amount with not existed swap pairs",
			fmt.Sprintf("%s/lbm/fswap/v1/total_swappable_to_coin_amount?from_denom=fake&to_denom=dummy", baseURL),
			true,
			&fswaptypes.QueryTotalSwappableToCoinAmountResponse{},
		},
		{
			"test query total_swappable_to_coin_amount with nil to_denom",
			fmt.Sprintf("%s/lbm/fswap/v1/total_swappable_to_coin_amount?from_denom=stake", baseURL),
			true,
			&fswaptypes.QueryTotalSwappableToCoinAmountResponse{},
		},
		{
			"test query total_swappable_to_coin_amount with nil from_denom",
			fmt.Sprintf("%s/lbm/fswap/v1/total_swappable_to_coin_amount?to_denom=dummy", baseURL),
			true,
			&fswaptypes.QueryTotalSwappableToCoinAmountResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			})
			s.Require().NoError(err)
			var valRes fswaptypes.QueryTotalSwappableToCoinAmountResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &valRes)

			if tc.expectedErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expected.String(), valRes.String())
			}
		})
	}
}
