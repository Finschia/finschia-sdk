package testutil

import (
	"fmt"
	"strings"

	"github.com/line/lbm-sdk/client/flags"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	"github.com/line/lbm-sdk/x/token/client/cli"
	ostcli "github.com/line/ostracon/libs/cli"
)

// func (s *IntegrationTestSuite) TestNewQueryCmdBalance() {
// 	val := s.network.Validators[0]

// 	testCases := []struct {
// 		name           string
// 		args           []string
// 		expectErr      bool
// 		expectedOutput string
// 	}{
// 		{
// 			"json output",
// 			[]string{
// 				fmt.Sprintf("--%s=1", flags.FlagHeight),
// 				fmt.Sprintf("--%s=json", ostcli.OutputFlag),
// 				val.ValAddress.String(),
// 			},
// 			false,
// 			fmt.Sprintf(`{"auth":{"operator_address":"%s","creation_allowed":true}}`,
// 				val.ValAddress.String(),
// 			),
// 		},
// 		{
// 			"text output",
// 			[]string{
// 				fmt.Sprintf("--%s=1", flags.FlagHeight),
// 				fmt.Sprintf("--%s=text", ostcli.OutputFlag),
// 				val.ValAddress.String(),
// 			},
// 			false,
// 			fmt.Sprintf(`auth:
//   creation_allowed: true
//   operator_address: %s`,
// 				val.ValAddress.String(),
// 			),
// 		},
// 		{
// 			"with no args",
// 			[]string{
// 				fmt.Sprintf("--%s=1", flags.FlagHeight),
// 			},
// 			true,
// 			"",
// 		},
// 		{
// 			"with an invalid address",
// 			[]string{
// 				fmt.Sprintf("--%s=1", flags.FlagHeight),
// 				"this-is-an-invalid-address",
// 			},
// 			true,
// 			"",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc

// 		s.Run(tc.name, func() {
// 			cmd := cli.NewQueryCmdBalance()
// 			clientCtx := val.ClientCtx

// 			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
// 			if tc.expectErr {
// 				s.Require().Error(err)
// 			} else {
// 				s.Require().NoError(err)
// 				s.Require().Equal(tc.expectedOutput, strings.TrimSpace(out.String()))
// 			}
// 		})
// 	}
// }

// func (s *IntegrationTestSuite) TestNewQueryCmdBalance() {
// 	val := s.network.Validators[0]

// 	testCases := []struct {
// 		name           string
// 		args           []string
// 		expectErr      bool
// 		expectedOutput string
// 	}{
// 		{
// 			"json output",
// 			[]string{
// 				fmt.Sprintf("--%s=1", flags.FlagHeight),
// 				fmt.Sprintf("--%s=json", ostcli.OutputFlag),
// 				val.ValAddress.String(),
// 			},
// 			false,
// 			fmt.Sprintf(`{"auth":{"operator_address":"%s","creation_allowed":true}}`,
// 				val.ValAddress.String(),
// 			),
// 		},
// 		{
// 			"text output",
// 			[]string{
// 				fmt.Sprintf("--%s=1", flags.FlagHeight),
// 				fmt.Sprintf("--%s=text", ostcli.OutputFlag),
// 				val.ValAddress.String(),
// 			},
// 			false,
// 			fmt.Sprintf(`auth:
//   creation_allowed: true
//   operator_address: %s`,
// 				val.ValAddress.String(),
// 			),
// 		},
// 		{
// 			"with no args",
// 			[]string{
// 				fmt.Sprintf("--%s=1", flags.FlagHeight),
// 			},
// 			true,
// 			"",
// 		},
// 		{
// 			"with an invalid address",
// 			[]string{
// 				fmt.Sprintf("--%s=1", flags.FlagHeight),
// 				"this-is-an-invalid-address",
// 			},
// 			true,
// 			"",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc

// 		s.Run(tc.name, func() {
// 			cmd := cli.NewQueryCmdBalance()
// 			clientCtx := val.ClientCtx

// 			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
// 			if tc.expectErr {
// 				s.Require().Error(err)
// 			} else {
// 				s.Require().NoError(err)
// 				s.Require().Equal(tc.expectedOutput, strings.TrimSpace(out.String()))
// 			}
// 		})
// 	}
// }

func (s *IntegrationTestSuite) TestNewQueryCmdTokens() {
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
			fmt.Sprintf(`{"pagination":{"next_key":null,"total":"0"},"tokens":[{"operator_address":"%s","creation_allowed":true}]}`,
				val.ValAddress.String(),
			),
		},
		{
			"text output",
			[]string{
				fmt.Sprintf("--%s=1", flags.FlagHeight),
				fmt.Sprintf("--%s=text", ostcli.OutputFlag),
			},
			fmt.Sprintf(`pagination:
  next_key: null
  total: "0"
tokens:
- creation_allowed: true
  operator_address: %s
`,
				val.ValAddress.String(),
			),
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewQueryCmdTokens()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)
			s.Require().Equal(tc.expectedOutput, strings.TrimSpace(out.String()))
		})
	}
}
