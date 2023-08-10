package testutil

import (
	"fmt"
	"time"

	"github.com/Finschia/finschia-sdk/client/flags"
	clitestutil "github.com/Finschia/finschia-sdk/testutil/cli"
	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"
	txtypes "github.com/Finschia/finschia-sdk/types/tx"
	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/Finschia/finschia-sdk/x/foundation/client/cli"
)

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
		"wrong number of args": {
			[]string{
				val.Address.String(),
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.OneInt())).String(),
				"extra",
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
		fmt.Sprintf("--%s", flags.FlagGenerateOnly),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.authority.String(),
				s.stranger.String(),
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.OneInt())).String(),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				s.authority.String(),
				s.stranger.String(),
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.OneInt())).String(),
				"extra",
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

			var res txtypes.Tx
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdUpdateMembers() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s", flags.FlagGenerateOnly),
	}

	updates := `[{"address":"%s"}]`
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.authority.String(),
				fmt.Sprintf(updates, s.comingMember),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				s.authority.String(),
				fmt.Sprintf(updates, s.comingMember),
				"extra",
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

			var res txtypes.Tx
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdUpdateDecisionPolicy() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s", flags.FlagGenerateOnly),
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
				s.authority.String(),
				doMarshal(&foundation.ThresholdDecisionPolicy{
					Threshold: sdk.NewDec(10),
					Windows: &foundation.DecisionPolicyWindows{
						VotingPeriod: time.Hour,
					},
				}),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				s.authority.String(),
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

			var res txtypes.Tx
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdSubmitProposal() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.permanentMember),
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
				fmt.Sprintf(proposers, s.permanentMember),
				s.msgToString(testdata.NewTestMsg()),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				"test proposal",
				fmt.Sprintf(proposers, s.permanentMember),
				s.msgToString(testdata.NewTestMsg()),
				"extra",
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
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.permanentMember),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)))),
	}

	id := s.submitProposal(testdata.NewTestMsg(s.authority), false)
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				fmt.Sprint(id),
				s.permanentMember.String(),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				fmt.Sprint(id),
				s.permanentMember.String(),
				"extra",
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
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.permanentMember),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)))),
	}

	id := s.submitProposal(testdata.NewTestMsg(s.authority), false)
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				fmt.Sprint(id),
				s.permanentMember.String(),
				"VOTE_OPTION_YES",
				"test vote",
			},
			true,
		},
		"wrong number of args": {
			[]string{
				fmt.Sprint(id),
				s.permanentMember.String(),
				"VOTE_OPTION_YES",
				"test vote",
				"extra",
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
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.permanentMember),
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
				fmt.Sprintf("%d", s.proposalID),
				s.permanentMember.String(),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				fmt.Sprintf("%d", s.proposalID),
				s.permanentMember.String(),
				"extra",
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
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.leavingMember),
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
		"wrong number of args": {
			[]string{
				s.leavingMember.String(),
				"extra",
			},
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
		fmt.Sprintf("--%s", flags.FlagGenerateOnly),
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
				s.authority.String(),
				s.comingMember.String(),
				doMarshal(&foundation.ReceiveFromTreasuryAuthorization{}),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				s.authority.String(),
				s.comingMember.String(),
				doMarshal(&foundation.ReceiveFromTreasuryAuthorization{}),
				"extra",
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

			var res txtypes.Tx
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
		})
	}
}

func (s *IntegrationTestSuite) TestNewTxCmdRevoke() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s", flags.FlagGenerateOnly),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.authority.String(),
				s.leavingMember.String(),
				foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				s.authority.String(),
				s.leavingMember.String(),
				foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
				"extra",
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

			var res txtypes.Tx
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out)
		})
	}
}
