package testutil

import (
	"fmt"

	ostcli "github.com/Finschia/ostracon/libs/cli"
	"github.com/gogo/protobuf/proto"

	"github.com/Finschia/finschia-sdk/client/flags"
	clitestutil "github.com/Finschia/finschia-sdk/testutil/cli"
	"github.com/Finschia/finschia-sdk/types/query"
	"github.com/Finschia/finschia-sdk/x/token"
	"github.com/Finschia/finschia-sdk/x/token/client/cli"
)

func (s *IntegrationTestSuite) TestNewQueryCmdBalance() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.classes[0].Id,
				s.customer.String(),
			},
			true,
			&token.QueryBalanceResponse{
				Amount: s.balance,
			},
		},
		"extra args": {
			[]string{
				s.classes[0].Id,
				s.customer.String(),
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.classes[0].Id,
			},
			false,
			nil,
		},
		"invalid address": {
			[]string{
				s.classes[0].Id,
				"invalid",
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
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
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
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.classes[0].Id,
			},
			true,
			&token.QueryContractResponse{
				Contract: s.classes[0],
			},
		},
		"extra args": {
			[]string{
				s.classes[0].Id,
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdContract()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.QueryContractResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdGranteeGrants() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.classes[0].Id,
				s.vendor.String(),
			},
			true,
			&token.QueryGranteeGrantsResponse{
				Grants: []token.Grant{
					{
						Grantee:    s.vendor.String(),
						Permission: token.PermissionModify,
					},
					{
						Grantee:    s.vendor.String(),
						Permission: token.PermissionMint,
					},
					{
						Grantee:    s.vendor.String(),
						Permission: token.PermissionBurn,
					},
				},
				Pagination: &query.PageResponse{
					Total: 3,
				},
			},
		},
		"extra args": {
			[]string{
				s.classes[0].Id,
				s.vendor.String(),
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.classes[0].Id,
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdGranteeGrants()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.QueryGranteeGrantsResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdIsOperatorFor() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.classes[0].Id,
				s.vendor.String(),
				s.customer.String(),
			},
			true,
			&token.QueryIsOperatorForResponse{
				Authorized: true,
			},
		},
		"extra args": {
			[]string{
				s.classes[0].Id,
				s.vendor.String(),
				s.customer.String(),
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.classes[0].Id,
				s.vendor.String(),
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdIsOperatorFor()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.QueryIsOperatorForResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdHoldersByOperator() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{
				s.classes[0].Id,
				s.vendor.String(),
			},
			true,
			&token.QueryHoldersByOperatorResponse{
				Holders:    []string{s.customer.String()},
				Pagination: &query.PageResponse{},
			},
		},
		"extra args": {
			[]string{
				s.classes[0].Id,
				s.vendor.String(),
				"extra",
			},
			false,
			nil,
		},
		"not enough args": {
			[]string{
				s.classes[0].Id,
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdHoldersByOperator()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.QueryHoldersByOperatorResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}
