package testutil

import (
	"fmt"

	"github.com/line/lbm-sdk/client/flags"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/token/client/cli"
)

func (s *IntegrationTestSuite) TestNewTxCmdSend() {
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
				s.classes[0].Id,
				s.customer.String(),
				s.vendor.String(),
				"1",
			},
			true,
		},
		"extra args": {
			[]string{
				s.classes[0].Id,
				s.customer.String(),
				s.vendor.String(),
				"1",
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.classes[0].Id,
				s.customer.String(),
				s.vendor.String(),
			},
			false,
		},
		"amount out of range": {
			[]string{
				s.classes[0].Id,
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
			cmd := cli.NewTxCmdSend()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			s.Require().EqualValues(0, res.Code, out.String())
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
				s.classes[0].Id,
				s.vendor.String(),
				s.customer.String(),
				s.vendor.String(),
				"1",
			},
			true,
		},
		"extra args": {
			[]string{
				s.classes[0].Id,
				s.vendor.String(),
				s.customer.String(),
				s.vendor.String(),
				"1",
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.classes[0].Id,
				s.vendor.String(),
				s.customer.String(),
				s.vendor.String(),
			},
			false,
		},
		"amount out of range": {
			[]string{
				s.classes[0].Id,
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

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			s.Require().EqualValues(0, res.Code, out.String())
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
				s.classes[0].Id,
				s.vendor.String(),
				s.customer.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.classes[0].Id,
				s.vendor.String(),
				s.customer.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.classes[0].Id,
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

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			s.Require().EqualValues(0, res.Code, out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdRevokeOperator() {
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
				s.classes[1].Id,
				s.customer.String(),
				s.vendor.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.classes[1].Id,
				s.customer.String(),
				s.vendor.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.classes[1].Id,
				s.customer.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdRevokeOperator()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			s.Require().EqualValues(0, res.Code, out.String())
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

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			s.Require().EqualValues(0, res.Code, out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdGrantPermission() {
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
				s.classes[1].Id,
				s.vendor.String(),
				s.customer.String(),
				token.LegacyPermissionMint.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.classes[1].Id,
				s.vendor.String(),
				s.customer.String(),
				token.LegacyPermissionMint.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.classes[1].Id,
				s.vendor.String(),
				s.customer.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdGrantPermission()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			s.Require().EqualValues(0, res.Code, out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdRevokePermission() {
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
				s.classes[1].Id,
				s.vendor.String(),
				token.LegacyPermissionModify.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.classes[1].Id,
				s.vendor.String(),
				token.LegacyPermissionModify.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.classes[1].Id,
				s.vendor.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdRevokePermission()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			s.Require().EqualValues(0, res.Code, out.String())
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
				s.classes[0].Id,
				s.vendor.String(),
				s.customer.String(),
				"1",
			},
			true,
		},
		"extra args": {
			[]string{
				s.classes[0].Id,
				s.vendor.String(),
				s.customer.String(),
				"1",
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.classes[0].Id,
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

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			s.Require().EqualValues(0, res.Code, out.String())
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
				s.classes[1].Id,
				s.vendor.String(),
				"1",
			},
			true,
		},
		"extra args": {
			[]string{
				s.classes[1].Id,
				s.vendor.String(),
				"1",
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.classes[1].Id,
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

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			s.Require().EqualValues(0, res.Code, out.String())
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
				s.classes[0].Id,
				s.vendor.String(),
				s.customer.String(),
				"1",
			},
			true,
		},
		"extra args": {
			[]string{
				s.classes[0].Id,
				s.vendor.String(),
				s.customer.String(),
				"1",
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.classes[0].Id,
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

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			s.Require().EqualValues(0, res.Code, out.String())
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
				s.classes[0].Id,
				s.vendor.String(),
				token.AttributeKeyName.String(),
				"cool token",
			},
			true,
		},
		"extra args": {
			[]string{
				s.classes[0].Id,
				s.vendor.String(),
				token.AttributeKeyName.String(),
				"cool token",
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.classes[0].Id,
				s.vendor.String(),
				token.AttributeKeyName.String(),
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

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			s.Require().EqualValues(0, res.Code, out.String())
		})
	}
}
