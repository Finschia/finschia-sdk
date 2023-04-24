package testutil

import (
	"fmt"

	"github.com/Finschia/finschia-sdk/client/flags"
	clitestutil "github.com/Finschia/finschia-sdk/testutil/cli"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/collection"
	"github.com/Finschia/finschia-sdk/x/collection/client/cli"
)

func (s *IntegrationTestSuite) TestNewTxCmdSendFT() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	amount := collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance))
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.stranger.String(),
				s.customer.String(),
				amount.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.stranger.String(),
				s.customer.String(),
				amount.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.stranger.String(),
				s.customer.String(),
			},
			false,
		},
		"amount out of range": {
			[]string{
				s.contractID,
				s.stranger.String(),
				s.customer.String(),
				fmt.Sprintf("%s:1%0127d", s.ftClassID, 0),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.stranger.String(),
				s.customer.String(),
				amount.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdSendFT()
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

func (s *IntegrationTestSuite) TestNewTxCmdOperatorSendFT() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	amount := collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance))
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.vendor.String(),
				amount.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.vendor.String(),
				amount.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.vendor.String(),
			},
			false,
		},
		"amount out of range": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.vendor.String(),
				fmt.Sprintf("%s:1%0127d", s.ftClassID, 0),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				s.customer.String(),
				s.vendor.String(),
				amount.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdOperatorSendFT()
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

func (s *IntegrationTestSuite) TestNewTxCmdSendNFT() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	tokenID := collection.NewNFTID(s.nftClassID, s.lenChain*3*3+1)
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.stranger.String(),
				s.customer.String(),
				tokenID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.stranger.String(),
				s.customer.String(),
				tokenID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.stranger.String(),
				s.customer.String(),
			},
			false,
		},
		"amount out of range": {
			[]string{
				s.contractID,
				s.stranger.String(),
				s.customer.String(),
				fmt.Sprintf("%s:1%0127d", s.ftClassID, 0),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.stranger.String(),
				s.customer.String(),
				tokenID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdSendNFT()
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

func (s *IntegrationTestSuite) TestNewTxCmdOperatorSendNFT() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	tokenID := collection.NewNFTID(s.nftClassID, 1)
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.vendor.String(),
				tokenID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.vendor.String(),
				tokenID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.vendor.String(),
			},
			false,
		},
		"amount out of range": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.vendor.String(),
				fmt.Sprintf("%s:1%0127d", s.ftClassID, 0),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				s.customer.String(),
				s.vendor.String(),
				tokenID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdOperatorSendNFT()
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

func (s *IntegrationTestSuite) TestNewTxCmdCreateContract() {
	val := s.network.Validators[0]
	commonArgs := []string{
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
			},
			true,
		},
		"extra args": {
			[]string{
				s.vendor.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{},
			false,
		},
		"invalid creator": {
			[]string{
				"",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdCreateContract()
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

func (s *IntegrationTestSuite) TestNewTxCmdIssueFT() {
	val := s.network.Validators[0]
	commonArgs := []string{
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
				s.contractID,
				s.operator.String(),
				fmt.Sprintf("--%s=%s", cli.FlagName, "tibetian fox"),
				fmt.Sprintf("--%s=%s", cli.FlagTo, s.operator),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				"extra",
				fmt.Sprintf("--%s=%s", cli.FlagName, "tibetian fox"),
				fmt.Sprintf("--%s=%s", cli.FlagTo, s.operator),
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				fmt.Sprintf("--%s=%s", cli.FlagName, "tibetian fox"),
				fmt.Sprintf("--%s=%s", cli.FlagTo, s.operator),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdIssueFT()
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

func (s *IntegrationTestSuite) TestNewTxCmdIssueNFT() {
	val := s.network.Validators[0]
	commonArgs := []string{
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
				s.contractID,
				s.operator.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdIssueNFT()
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

func (s *IntegrationTestSuite) TestNewTxCmdMintFT() {
	val := s.network.Validators[0]
	commonArgs := []string{
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
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.ftClassID,
				s.balance.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.ftClassID,
				s.balance.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.ftClassID,
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				s.customer.String(),
				s.ftClassID,
				s.balance.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdMintFT()
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

func (s *IntegrationTestSuite) TestNewTxCmdMintNFT() {
	val := s.network.Validators[0]
	commonArgs := []string{
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
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.nftClassID,
				fmt.Sprintf("--%s=%s", cli.FlagName, "arctic fox"),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				s.nftClassID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				s.customer.String(),
				s.nftClassID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdMintNFT()
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

func (s *IntegrationTestSuite) TestNewTxCmdBurnFT() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	amount := collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance))
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
				amount.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				amount.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				amount.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdBurnFT()
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

func (s *IntegrationTestSuite) TestNewTxCmdOperatorBurnFT() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	amount := collection.NewCoins(collection.NewFTCoin(s.ftClassID, s.balance))
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.vendor.String(),
				amount.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.vendor.String(),
				amount.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.vendor.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				s.vendor.String(),
				amount.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdOperatorBurnFT()
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

func (s *IntegrationTestSuite) TestNewTxCmdBurnNFT() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	tokenID := collection.NewNFTID(s.nftClassID, s.lenChain*3+1)
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
				tokenID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				tokenID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				tokenID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdBurnNFT()
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

func (s *IntegrationTestSuite) TestNewTxCmdOperatorOperatorBurnNFT() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	tokenID := collection.NewNFTID(s.nftClassID, s.lenChain*3*2+1)
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.vendor.String(),
				tokenID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.vendor.String(),
				tokenID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.vendor.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				s.vendor.String(),
				tokenID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdOperatorBurnNFT()
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
				s.contractID,
				s.operator.String(),
				s.nftClassID,
				"",
				collection.AttributeKeyName.String(),
				"tibetian fox",
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.nftClassID,
				"",
				collection.AttributeKeyName.String(),
				"tibetian fox",
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.nftClassID,
				"",
				collection.AttributeKeyName.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				s.nftClassID,
				"",
				collection.AttributeKeyName.String(),
				"tibetian fox",
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

func (s *IntegrationTestSuite) TestNewTxCmdAttach() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	subjectID := collection.NewNFTID(s.nftClassID, s.lenChain*(3*2+1)+1)
	targetID := collection.NewNFTID(s.nftClassID, s.lenChain*(3*2+2)+1)
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.vendor.String(),
				subjectID,
				targetID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.vendor.String(),
				subjectID,
				targetID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.vendor.String(),
				subjectID,
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.vendor.String(),
				subjectID,
				targetID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdAttach()
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

func (s *IntegrationTestSuite) TestNewTxCmdDetach() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	subjectID := collection.NewNFTID(s.nftClassID, s.lenChain*(3*3+1)+2)
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.stranger.String(),
				subjectID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.stranger.String(),
				subjectID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.stranger.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.stranger.String(),
				subjectID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdDetach()
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

func (s *IntegrationTestSuite) TestNewTxCmdOperatorAttach() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	subjectID := collection.NewNFTID(s.nftClassID, s.lenChain*2+1)
	targetID := collection.NewNFTID(s.nftClassID, s.lenChain+1)
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				subjectID,
				targetID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				subjectID,
				targetID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				subjectID,
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				s.customer.String(),
				subjectID,
				targetID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdOperatorAttach()
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

func (s *IntegrationTestSuite) TestNewTxCmdOperatorDetach() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	subjectID := collection.NewNFTID(s.nftClassID, s.lenChain*3*3)
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.vendor.String(),
				subjectID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.vendor.String(),
				subjectID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.vendor.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				s.vendor.String(),
				subjectID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdOperatorDetach()
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
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				collection.LegacyPermissionMint.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				collection.LegacyPermissionMint.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
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
				s.contractID,
				s.vendor.String(),
				collection.LegacyPermissionModify.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.vendor.String(),
				collection.LegacyPermissionModify.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
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

func (s *IntegrationTestSuite) TestNewTxCmdAuthorizeOperator() {
	val := s.network.Validators[0]
	commonArgs := []string{
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
				s.contractID,
				s.vendor.String(),
				s.customer.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.vendor.String(),
				s.customer.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.vendor.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdAuthorizeOperator()
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
				s.contractID,
				s.operator.String(),
				s.vendor.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.vendor.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
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
