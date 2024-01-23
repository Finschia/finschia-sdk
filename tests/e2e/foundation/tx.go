package foundation

import (
	"fmt"
	"time"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"

	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/Finschia/finschia-sdk/x/foundation/client/cli"
)

func (s *E2ETestSuite) TestNewTxCmdUpdateParams() {
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
				s.bytesToString(s.authority),
				fmt.Sprintf(`{"foundation_tax": "%s"}`, math.LegacyZeroDec()),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				s.bytesToString(s.authority),
				fmt.Sprintf(`{"foundation_tax": "%s"}`, math.LegacyZeroDec()),
				"extra",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdUpdateParams()
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

func (s *E2ETestSuite) TestNewTxCmdFundTreasury() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.bytesToString(val.Address)),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10)))),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.bytesToString(val.Address),
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.OneInt())).String(),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				s.bytesToString(val.Address),
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.OneInt())).String(),
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

			s.Require().NoError(s.network.WaitForNextBlock())
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdWithdrawFromTreasury() {
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
				s.bytesToString(s.authority),
				s.bytesToString(s.stranger),
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.OneInt())).String(),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				s.bytesToString(s.authority),
				s.bytesToString(s.stranger),
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.OneInt())).String(),
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

func (s *E2ETestSuite) TestNewTxCmdUpdateMembers() {
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
				s.bytesToString(s.authority),
				fmt.Sprintf(updates, s.bytesToString(s.comingMember)),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				s.bytesToString(s.authority),
				fmt.Sprintf(updates, s.bytesToString(s.comingMember)),
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

func (s *E2ETestSuite) TestNewTxCmdUpdateDecisionPolicy() {
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
				s.bytesToString(s.authority),
				doMarshal(&foundation.ThresholdDecisionPolicy{
					Threshold: math.LegacyNewDec(10),
					Windows: &foundation.DecisionPolicyWindows{
						VotingPeriod: time.Hour,
					},
				}),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				s.bytesToString(s.authority),
				doMarshal(&foundation.ThresholdDecisionPolicy{
					Threshold: math.LegacyNewDec(10),
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

func (s *E2ETestSuite) TestNewTxCmdSubmitProposal() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.bytesToString(s.permanentMember)),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10)))),
	}

	proposers := `["%s"]`
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				"test proposal",
				fmt.Sprintf(proposers, s.bytesToString(s.permanentMember)),
				s.msgToString(&foundation.MsgWithdrawFromTreasury{
					Authority: s.bytesToString(s.authority),
					To:        s.bytesToString(s.stranger),
					Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(123))),
				}),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				"test proposal",
				fmt.Sprintf(proposers, s.bytesToString(s.permanentMember)),
				s.msgToString(&foundation.MsgWithdrawFromTreasury{
					Authority: s.bytesToString(s.authority),
					To:        s.bytesToString(s.stranger),
					Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(123))),
				}),
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

			s.Require().NoError(s.network.WaitForNextBlock())
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdWithdrawProposal() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.bytesToString(s.permanentMember)),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10)))),
	}

	proposalID := 2
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				fmt.Sprint(proposalID),
				s.bytesToString(s.permanentMember),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				fmt.Sprint(proposalID),
				s.bytesToString(s.permanentMember),
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

			s.Require().NoError(s.network.WaitForNextBlock())
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdVote() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.bytesToString(s.permanentMember)),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10)))),
	}

	proposalID := 3
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				fmt.Sprint(proposalID),
				s.bytesToString(s.permanentMember),
				"VOTE_OPTION_YES",
				"test vote",
			},
			true,
		},
		"wrong number of args": {
			[]string{
				fmt.Sprint(proposalID),
				s.bytesToString(s.permanentMember),
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

			s.Require().NoError(s.network.WaitForNextBlock())
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdExec() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.bytesToString(s.permanentMember)),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10)))),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				fmt.Sprintf("%d", s.proposalID),
				s.bytesToString(s.permanentMember),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				fmt.Sprintf("%d", s.proposalID),
				s.bytesToString(s.permanentMember),
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

			s.Require().NoError(s.network.WaitForNextBlock())
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdUpdateCensorship() {
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
				s.bytesToString(s.authority),
				foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
				foundation.CensorshipAuthorityGovernance.String(),
			},
			true,
		},
		"valid abbreviation": {
			[]string{
				s.bytesToString(s.authority),
				foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
				"governance",
			},
			true,
		},
		"wrong number of args": {
			[]string{
				s.bytesToString(s.authority),
				foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
				foundation.CensorshipAuthorityGovernance.String(),
				"extra",
			},
			false,
		},
		"invalid new authority": {
			[]string{
				s.bytesToString(s.authority),
				foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
				"invalid-new-authority",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdUpdateCensorship()
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

func (s *E2ETestSuite) TestNewTxCmdLeaveFoundation() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.bytesToString(s.leavingMember)),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10)))),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.bytesToString(s.leavingMember),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				s.bytesToString(s.leavingMember),
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

			s.Require().NoError(s.network.WaitForNextBlock())
		})
	}
}

func (s *E2ETestSuite) TestNewTxCmdGrant() {
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
				s.bytesToString(s.authority),
				s.bytesToString(s.comingMember),
				doMarshal(&foundation.ReceiveFromTreasuryAuthorization{}),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				s.bytesToString(s.authority),
				s.bytesToString(s.comingMember),
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

func (s *E2ETestSuite) TestNewTxCmdRevoke() {
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
				s.bytesToString(s.authority),
				s.bytesToString(s.leavingMember),
				foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				s.bytesToString(s.authority),
				s.bytesToString(s.leavingMember),
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
