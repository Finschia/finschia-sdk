package testutil

import (
	"fmt"

	"github.com/line/lbm-sdk/client/flags"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/token/client/cli"
)

func (s *IntegrationTestSuite) TestNewTxCmdTransfer() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.customer),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.mintableClass.Id,
				s.customer.String(),
				s.vendor.String(),
				s.balance.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.mintableClass.Id,
				s.customer.String(),
				s.vendor.String(),
				s.balance.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.mintableClass.Id,
				s.customer.String(),
				s.vendor.String(),
			},
			false,
		},
		"amount out of range": {
			[]string{
				s.mintableClass.Id,
				s.customer.String(),
				s.vendor.String(),
				"10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdTransfer()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.MsgTransferResponse
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdTransferFrom() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.vendor),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.customer.String(),
				s.vendor.String(),
				s.balance.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.customer.String(),
				s.vendor.String(),
				s.balance.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.customer.String(),
				s.vendor.String(),
			},
			false,
		},
		"amount out of range": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.customer.String(),
				s.vendor.String(),
				"10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdTransferFrom()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.MsgTransferFromResponse
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdApprove() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.vendor),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.customer.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.customer.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdApprove()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.MsgApproveResponse
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdIssue() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.vendor),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.vendor.String(),
				s.vendor.String(),
				"Test class",
				"TT",
				fmt.Sprintf("--%s=%s", cli.FlagImageURI, "URI"),
				fmt.Sprintf("--%s=%s", cli.FlagMeta, "META"),
				fmt.Sprintf("--%s=%d", cli.FlagDecimals, 8),
				fmt.Sprintf("--%s=%s", cli.FlagSupply, "10000000000"),
				fmt.Sprintf("--%s=%v", cli.FlagMintable, true),
			},
			true,
		},
		"extra args": {
			[]string{
				s.vendor.String(),
				s.vendor.String(),
				"Test class",
				"TT",
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.vendor.String(),
				s.vendor.String(),
				"Test class",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdIssue()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.MsgIssueResponse
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdGrant() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.vendor),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.customer.String(),
				"mint",
			},
			true,
		},
		"extra args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.customer.String(),
				"mint",
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.customer.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdGrant()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.MsgGrantResponse
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdRevoke() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.vendor),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				"mint",
			},
			true,
		},
		"extra args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				"mint",
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdRevoke()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.MsgRevokeResponse
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdMint() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.vendor),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.customer.String(),
				s.balance.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.customer.String(),
				s.balance.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.customer.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdMint()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.MsgMintResponse
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdBurn() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.vendor),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.balance.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.balance.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdBurn()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.MsgBurnResponse
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdBurnFrom() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.vendor),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.customer.String(),
				s.balance.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.customer.String(),
				s.balance.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				s.customer.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdBurnFrom()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.MsgBurnFromResponse
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdModify() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.vendor),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				"name",
				"cool token",
			},
			true,
		},
		"extra args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				"name",
				"cool token",
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.mintableClass.Id,
				s.vendor.String(),
				"name",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdModify()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual token.MsgModifyResponse
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}
