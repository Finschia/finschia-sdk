package testutil

import (
	"fmt"
	"strings"

	ostcli "github.com/line/ostracon/libs/cli"

	"github.com/line/lbm-sdk/client/flags"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	"github.com/line/lbm-sdk/x/foundation/client/cli"
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
			`{"params":{"enabled":true,"foundation_tax":"0.000000000000000000"}}`,
		},
		{
			"text output",
			[]string{
				fmt.Sprintf("--%s=1", flags.FlagHeight),
				fmt.Sprintf("--%s=text", ostcli.OutputFlag),
			},
			`params:
  enabled: true
  foundation_tax: "0.000000000000000000"`,
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
				val.ValAddress.String(),
			},
			false,
			fmt.Sprintf(`{"auth":{"operator_address":"%s","creation_allowed":true}}`,
				val.ValAddress.String(),
			),
		},
		{
			"text output",
			[]string{
				fmt.Sprintf("--%s=1", flags.FlagHeight),
				fmt.Sprintf("--%s=text", ostcli.OutputFlag),
				val.ValAddress.String(),
			},
			false,
			fmt.Sprintf(`auth:
  creation_allowed: true
  operator_address: %s`,
				val.ValAddress.String(),
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
				"this-is-an-invalid-address",
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

	jsonAuth := `{"operator_address":"%s","creation_allowed":true}`
	textAuth := `- creation_allowed: true
  operator_address: %s
`
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
			fmt.Sprintf(`{"auths":[%s,%s],"pagination":{"next_key":null,"total":"0"}}`,
				fmt.Sprintf(jsonAuth, s.network.Validators[0].ValAddress),
				fmt.Sprintf(jsonAuth, s.network.Validators[1].ValAddress),
			),
		},
		{
			"text output",
			[]string{
				fmt.Sprintf("--%s=1", flags.FlagHeight),
				fmt.Sprintf("--%s=text", ostcli.OutputFlag),
			},
			fmt.Sprintf(`auths:
%s%s
pagination:
  next_key: null
  total: "0"`,
				fmt.Sprintf(textAuth, s.network.Validators[0].ValAddress),
				fmt.Sprintf(textAuth, s.network.Validators[1].ValAddress),
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
