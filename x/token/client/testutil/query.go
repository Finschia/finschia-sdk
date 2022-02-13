package testutil

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/line/lbm-sdk/client/flags"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	"github.com/line/lbm-sdk/types/query"
	"github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/token/client/cli"
	ostcli "github.com/line/ostracon/libs/cli"
)

func (s *IntegrationTestSuite) TestNewQueryCmdBalance() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args           []string
		valid bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.mintableClass.Id,
				s.vendor.Address.String(),
			},
			true,
			&token.QueryBalanceResponse{
				Amount: s.balance,
			},
		},
		"extra args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.Address.String(),
				"extra",
			},
			false,
			nil,
		},
		"no such an id": {
			[]string{
				"invalid",
				s.vendor.Address.String(),
			},
			true,
			&token.QueryBalanceResponse{
				Amount: sdk.ZeroInt(),
			},
		},
		"no such an address": {
			[]string{
				s.mintableClass.Id,
				"invalid",
			},
			true,
			&token.QueryBalanceResponse{
				Amount: sdk.ZeroInt(),
			},
		},
		"no arguments": {
			[]string{},
			false,
			nil,
		},
		"empty address": {
			[]string{
				s.mintableClass.Id,
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdBalance()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.QueryBalanceResponse
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdToken() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args           []string
		valid bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.mintableClass.Id,
			},
			true,
			&token.QueryTokenResponse{
				Token: s.mintableClass,
			},
		},
		"extra args": {
			[]string{
				s.mintableClass.Id,
				"extra",
			},
			false,
			nil,
		},
		"no such an id": {
			[]string{
				"invalid",
			},
			false,
			nil,
		},
		"no class id": {
			[]string{},
			false,
			&token.QueryTokenResponse{
				Token: s.mintableClass,
			},
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdToken()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.QueryTokenResponse
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdTokens() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args           []string
		valid bool
		expected proto.Message
	}{
		"query all": {
			[]string{},
			true,
			&token.QueryTokensResponse{
				Tokens: []token.Token{s.notMintableClass, s.mintableClass},
				Pagination: &query.PageResponse{},
			},
		},
		"extra args": {
			[]string{
				"extra",
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdTokens()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.QueryTokensResponse
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}
