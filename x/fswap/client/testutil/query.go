package testutil

import (
	"github.com/gogo/protobuf/proto"

	clitestutil "github.com/Finschia/finschia-sdk/testutil/cli"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/query"
	"github.com/Finschia/finschia-sdk/x/fswap/client/cli"
	fswaptypes "github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (s *IntegrationTestSuite) TestCmdQuerySwapped() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	// avoid printing as yaml from CLI command
	clientCtx.OutputFormat = jsonOutputFormat

	fromDenom := s.cfg.BondDenom
	toDenom := s.toDenom.Base

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
		expected  proto.Message
	}{
		{
			"valid query",
			[]string{fromDenom, toDenom},
			false,
			&fswaptypes.QuerySwappedResponse{
				FromCoinAmount: sdk.NewCoin(fromDenom, sdk.ZeroInt()),
				ToCoinAmount:   sdk.NewCoin(toDenom, sdk.ZeroInt()),
			},
		},
		{
			"wrong number of args",
			[]string{fromDenom, toDenom, "extra"},
			true,
			nil,
		},
		{
			"invalid fromDenom",
			[]string{"", toDenom},
			true,
			nil,
		},
		{
			"invalid toDenom",
			[]string{fromDenom, ""},
			true,
			nil,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.CmdQuerySwapped()
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				var actual fswaptypes.QuerySwappedResponse
				s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual))
				s.Require().Equal(tc.expected, &actual)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestCmdQueryTotalSwappableAmount() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	// avoid printing as yaml from CLI command
	clientCtx.OutputFormat = jsonOutputFormat

	fromDenom := s.cfg.BondDenom
	toDenom := s.toDenom.Base

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
		expected  proto.Message
	}{
		{
			"valid query",
			[]string{fromDenom, toDenom},
			false,
			&fswaptypes.QueryTotalSwappableToCoinAmountResponse{
				SwappableAmount: sdk.NewCoin(toDenom, s.dummySwap.AmountCapForToDenom),
			},
		},
		{
			"wrong number of args",
			[]string{fromDenom, toDenom, "extra"},
			true,
			nil,
		},
		{
			"invalid fromDenom",
			[]string{"", toDenom},
			true,
			nil,
		},
		{
			"invalid toDenom",
			[]string{fromDenom, ""},
			true,
			nil,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.CmdQueryTotalSwappableAmount()
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				var actual fswaptypes.QueryTotalSwappableToCoinAmountResponse
				s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual))
				s.Require().Equal(tc.expected, &actual)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestCmdQuerySwap() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	// avoid printing as yaml from CLI command
	clientCtx.OutputFormat = jsonOutputFormat

	fromDenom := s.cfg.BondDenom
	toDenom := s.toDenom.Base

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
		expected  proto.Message
	}{
		{
			"valid query",
			[]string{fromDenom, toDenom},
			false,
			&fswaptypes.QuerySwapResponse{
				Swap: s.dummySwap,
			},
		},
		{
			"wrong number of args",
			[]string{fromDenom, toDenom, "extra"},
			true,
			nil,
		},
		{
			"invalid fromDenom",
			[]string{"", toDenom},
			true,
			nil,
		},
		{
			"invalid toDenom",
			[]string{fromDenom, ""},
			true,
			nil,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.CmdQuerySwap()
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				var actual fswaptypes.QuerySwapResponse
				s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual))
				s.Require().Equal(tc.expected, &actual)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestCmdQuerySwaps() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	// avoid printing as yaml from CLI command
	clientCtx.OutputFormat = jsonOutputFormat

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
		expected  proto.Message
	}{
		{
			"valid query (default pagination)",
			[]string{},
			false,
			&fswaptypes.QuerySwapsResponse{
				Swaps:      []fswaptypes.Swap{s.dummySwap},
				Pagination: &query.PageResponse{Total: 1},
			},
		},
		{
			"invalid query",
			[]string{"extra"},
			true,
			nil,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.CmdQuerySwaps()
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				var actual fswaptypes.QuerySwapsResponse
				s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual))
				s.Require().Equal(tc.expected, &actual)
			}
		})
	}
}
