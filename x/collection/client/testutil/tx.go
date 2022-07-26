package testutil

import (
	"fmt"

	"github.com/line/lbm-sdk/client/flags"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
	"github.com/line/lbm-sdk/x/collection/client/cli"
)

func (s *IntegrationTestSuite) TestNewTxCmdSend() {
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
				s.vendor.String(),
				s.customer.String(),
				amount.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.vendor.String(),
				s.customer.String(),
				amount.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.vendor.String(),
				s.customer.String(),
			},
			false,
		},
		"amount out of range": {
			[]string{
				s.contractID,
				s.vendor.String(),
				s.customer.String(),
				fmt.Sprintf("%s:1%0127d", s.ftClassID, 0),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.vendor.String(),
				s.customer.String(),
				amount.String(),
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

func (s *IntegrationTestSuite) TestNewTxCmdOperatorSend() {
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
			cmd := cli.NewTxCmdOperatorSend()
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

func (s *IntegrationTestSuite) TestNewTxCmdCreateFTClass() {
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
			cmd := cli.NewTxCmdCreateFTClass()
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

func (s *IntegrationTestSuite) TestNewTxCmdCreateNFTClass() {
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
			cmd := cli.NewTxCmdCreateNFTClass()
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

func (s *IntegrationTestSuite) TestNewTxCmdBurn() {
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

func (s *IntegrationTestSuite) TestNewTxCmdOperatorBurn() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	amount := collection.NewCoins(collection.NewNFTCoin(s.nftClassID, s.lenChain*2+1))
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				amount.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
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
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				s.customer.String(),
				amount.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdOperatorBurn()
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

func (s *IntegrationTestSuite) TestNewTxCmdModifyContract() {
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
				collection.AttributeKeyName.String(),
				"fox",
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				collection.AttributeKeyName.String(),
				"fox",
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
				collection.AttributeKeyName.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				collection.AttributeKeyName.String(),
				"fox",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdModifyContract()
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

func (s *IntegrationTestSuite) TestNewTxCmdModifyTokenClass() {
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
				s.ftClassID,
				collection.AttributeKeyName.String(),
				"tibetian fox",
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.ftClassID,
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
				s.ftClassID,
				collection.AttributeKeyName.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				s.ftClassID,
				collection.AttributeKeyName.String(),
				"tibetian fox",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdModifyTokenClass()
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

func (s *IntegrationTestSuite) TestNewTxCmdModifyNFT() {
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
				tokenID,
				collection.AttributeKeyName.String(),
				"fennec fox 1",
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				tokenID,
				collection.AttributeKeyName.String(),
				"fennec fox 1",
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.String(),
				tokenID,
				collection.AttributeKeyName.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.String(),
				tokenID,
				collection.AttributeKeyName.String(),
				"fennec fox 1",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdModifyNFT()
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

	subjectID := collection.NewNFTID(s.nftClassID, s.lenChain*4+1)
	targetID := collection.NewNFTID(s.nftClassID, s.lenChain*3+1)
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

	subjectID := collection.NewNFTID(s.nftClassID, s.lenChain*5)
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.vendor.String(),
				subjectID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.vendor.String(),
				subjectID,
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
		"invalid contract id": {
			[]string{
				"",
				s.vendor.String(),
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

	subjectID := collection.NewNFTID(s.nftClassID, s.lenChain+1)
	targetID := collection.NewNFTID(s.nftClassID, 1)
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

	subjectID := collection.NewNFTID(s.nftClassID, s.lenChain*2)
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
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				subjectID,
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

func (s *IntegrationTestSuite) TestNewTxCmdGrant() {
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
				collection.PermissionMint.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.String(),
				s.customer.String(),
				collection.PermissionMint.String(),
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
			cmd := cli.NewTxCmdGrant()
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

func (s *IntegrationTestSuite) TestNewTxCmdAbandon() {
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
				collection.PermissionModify.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.vendor.String(),
				collection.PermissionModify.String(),
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
			cmd := cli.NewTxCmdAbandon()
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
