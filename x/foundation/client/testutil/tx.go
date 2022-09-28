package testutil

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/line/lbm-sdk/client/flags"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	"github.com/line/lbm-sdk/testutil/testdata"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/foundation/client/cli"
)

func (s *IntegrationTestSuite) TestNewProposalCmdUpdateFoundationParams() {
	val := s.network.Validators[0]

	commonFlags := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)))),
	}

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		expectedCode uint32
		respType     proto.Message
	}{
		{
			"with wrong # of args",
			commonFlags,
			true, 0, nil,
		},
		{
			"valid transaction",
			append([]string{
				fmt.Sprintf(`{"foundation_tax": "%s"}`, sdk.ZeroDec()),
			}, commonFlags...),
			false, 0, &sdk.TxResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewProposalCmdUpdateFoundationParams()
			flags.AddTxFlagsToCmd(cmd)

			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out)

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdFundTreasury() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)))),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				val.Address.String(),
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.OneInt())).String(),
			},
			true,
		},
		"extra args": {
			[]string{
				val.Address.String(),
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.OneInt())).String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				val.Address.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdFundTreasury()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
			s.Require().Zero(res.Code, out)
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdWithdrawFromTreasury() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)))),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.operator.String(),
				s.stranger.String(),
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.OneInt())).String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.operator.String(),
				s.stranger.String(),
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.OneInt())).String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.operator.String(),
				s.stranger.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdWithdrawFromTreasury()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
			s.Require().Zero(res.Code, out)
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdUpdateMembers() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)))),
	}

	updates := `[{"address":"%s"}]`
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.operator.String(),
				fmt.Sprintf(updates, s.comingMember),
			},
			true,
		},
		"extra args": {
			[]string{
				s.operator.String(),
				fmt.Sprintf(updates, s.comingMember),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.operator.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdUpdateMembers()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
			s.Require().Zero(res.Code, out)
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdUpdateDecisionPolicy() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)))),
	}

	doMarshal := func(policy foundation.DecisionPolicy) string {
		bz, err := val.ClientCtx.Codec.MarshalInterfaceJSON(policy)
		s.Require().NoError(err)
		return string(bz)
	}
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.operator.String(),
				doMarshal(&foundation.ThresholdDecisionPolicy{
					Threshold: sdk.NewDec(10),
					Windows: &foundation.DecisionPolicyWindows{
						VotingPeriod: time.Hour,
					},
				}),
			},
			true,
		},
		"extra args": {
			[]string{
				s.operator.String(),
				doMarshal(&foundation.ThresholdDecisionPolicy{
					Threshold: sdk.NewDec(10),
					Windows: &foundation.DecisionPolicyWindows{
						VotingPeriod: time.Hour,
					},
				}),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.operator.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdUpdateDecisionPolicy()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
			s.Require().Zero(res.Code, out)
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdSubmitProposal() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)))),
	}

	proposers := `["%s"]`
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				"test proposal",
				fmt.Sprintf(proposers, val.Address),
				s.msgToString(testdata.NewTestMsg()),
			},
			true,
		},
		"extra args": {
			[]string{
				"test proposal",
				fmt.Sprintf(proposers, val.Address),
				s.msgToString(testdata.NewTestMsg()),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				"test proposal",
				fmt.Sprintf(proposers, val.Address),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdSubmitProposal()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
			s.Require().Zero(res.Code, out)
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdWithdrawProposal() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)))),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		// TODO: make it work
		// "valid transaction": {
		// 	[]string{
		// 		"1",
		// 		val.Address.String(),
		// 	},
		// 	true,
		// },
		"extra args": {
			[]string{
				"1",
				val.Address.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				"1",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdWithdrawProposal()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
			s.Require().Zero(res.Code, out)
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdVote() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)))),
	}

	id := s.submitProposal(testdata.NewTestMsg(s.operator), false)
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				fmt.Sprint(id),
				val.Address.String(),
				"VOTE_OPTION_YES",
				"test vote",
			},
			true,
		},
		"extra args": {
			[]string{
				fmt.Sprint(id),
				val.Address.String(),
				"VOTE_OPTION_YES",
				"test vote",
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				fmt.Sprint(id),
				val.Address.String(),
				"VOTE_OPTION_YES",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdVote()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
			s.Require().Zero(res.Code, out)
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdExec() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)))),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				"1",
				val.Address.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				"1",
				val.Address.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				"1",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdExec()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
			s.Require().Zero(res.Code, out)
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdLeaveFoundation() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)))),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.leavingMember.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.leavingMember.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdLeaveFoundation()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
			s.Require().Zero(res.Code, out)
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdGrant() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)))),
	}

	doMarshal := func(authorization foundation.Authorization) string {
		bz, err := val.ClientCtx.Codec.MarshalInterfaceJSON(authorization)
		s.Require().NoError(err)
		return string(bz)
	}
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.operator.String(),
				s.comingMember.String(),
				doMarshal(&foundation.ReceiveFromTreasuryAuthorization{}),
			},
			true,
		},
		"extra args": {
			[]string{
				s.operator.String(),
				s.comingMember.String(),
				doMarshal(&foundation.ReceiveFromTreasuryAuthorization{}),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.operator.String(),
				s.comingMember.String(),
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
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
			s.Require().Zero(res.Code, out)
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdRevoke() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)))),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.operator.String(),
				s.leavingMember.String(),
				foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.operator.String(),
				s.leavingMember.String(),
				foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.operator.String(),
				s.leavingMember.String(),
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

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
			s.Require().Zero(res.Code, out)
		})
	}
}
