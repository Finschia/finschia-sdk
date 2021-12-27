package testutil

import (
	"fmt"
	"strings"
	
	"github.com/line/lbm-sdk/x/consortium/client/cli"
	"github.com/line/lbm-sdk/client/flags"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	ostcli "github.com/line/ostracon/libs/cli"
)

func (s *IntegrationTestSuite) TestNewQueryCmdParams() {
	val := s.network.Validators[0]

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{
			"json output",
			[]string{
				fmt.Sprintf("--%s=1", flags.FlagHeight),
				fmt.Sprintf("--%s=json", ostcli.OutputFlag),
			},
			`{"params":{"enabled":true}}`,
		},
		{
			"text output",
			[]string{
				fmt.Sprintf("--%s=1", flags.FlagHeight),
				fmt.Sprintf("--%s=text", ostcli.OutputFlag),
			},
			`params:
  enabled: true`,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewQueryCmdParams()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)
			s.Require().Equal(tc.expectedOutput, strings.TrimSpace(out.String()))
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdValidatorAuth() {
	val := s.network.Validators[0]

	testCases := []struct {
		name           string
		args           []string
		expectErr      bool
		expectedOutput string
	}{
		{
			"json output",
			[]string{
				fmt.Sprintf("--%s=1", flags.FlagHeight),
				fmt.Sprintf("--%s=json", ostcli.OutputFlag),
				string(val.ValAddress),
			},
			false,
			fmt.Sprintf(`{"auth":{"operator_address":"%s","creation_allowed":true}}`,
				string(val.ValAddress),
			),
		},
		{
			"text output",
			[]string{
				fmt.Sprintf("--%s=1", flags.FlagHeight),
				fmt.Sprintf("--%s=text", ostcli.OutputFlag),
				string(val.ValAddress),
			},
			false,
			fmt.Sprintf(`auth:
  creation_allowed: true
  operator_address: %s`,
				string(val.ValAddress),
			),
		},
		{
			"with no args",
			[]string{
				fmt.Sprintf("--%s=1", flags.FlagHeight),
			},
			true,
			"",
		},
		{
			"with an invalid address",
			[]string{
				fmt.Sprintf("--%s=1", flags.FlagHeight),
				fmt.Sprintf("this-is-an-invalid-address"),
			},
			true,
			"",
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewQueryCmdValidatorAuth()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedOutput, strings.TrimSpace(out.String()))
			}
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdValidatorAuths() {
	val := s.network.Validators[0]

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{
			"json output",
			[]string{
				fmt.Sprintf("--%s=1", flags.FlagHeight),
				fmt.Sprintf("--%s=json", ostcli.OutputFlag),
			},
			fmt.Sprintf(`{"auths":[{"operator_address":"%s","creation_allowed":true}],"pagination":{"next_key":null,"total":"0"}}`,
				string(val.ValAddress),
			),
		},
		{
			"text output",
			[]string{
				fmt.Sprintf("--%s=1", flags.FlagHeight),
				fmt.Sprintf("--%s=text", ostcli.OutputFlag),
			},
			fmt.Sprintf(`auths:
- creation_allowed: true
  operator_address: %s
pagination:
  next_key: null
  total: "0"`,
				string(val.ValAddress),
			),
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewQueryCmdValidatorAuths()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)
			s.Require().Equal(tc.expectedOutput, strings.TrimSpace(out.String()))
		})
	}
}
