package testutil

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	ostcli "github.com/line/ostracon/libs/cli"

	"github.com/line/lbm-sdk/client/flags"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/foundation/client/cli"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

func (s *IntegrationTestSuite) TestNewQueryCmdParams() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected proto.Message
	}{
		"valid query": {
			[]string{},
			true,
			&foundation.QueryParamsResponse{
				Params: foundation.Params{
					FoundationTax: sdk.MustNewDecFromStr("0.2"),
					CensoredMsgTypeUrls: []string{
						sdk.MsgTypeURL((*stakingtypes.MsgCreateValidator)(nil)),
						sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
					},
				},
			},
		},
		"extra args": {
			[]string{
				"extra",
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdParams()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual foundation.QueryParamsResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, &actual)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdTreasury() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{},
			true,
		},
		"extra args": {
			[]string{
				"extra",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdTreasury()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual foundation.QueryTreasuryResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdFoundationInfo() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{},
			true,
		},
		"extra args": {
			[]string{
				"extra",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdFoundationInfo()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual foundation.QueryFoundationInfoResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdMember() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected *foundation.Member
	}{
		"valid query": {
			[]string{
				val.Address.String(),
			},
			true,
			&foundation.Member{
				Address:  val.Address.String(),
				Metadata: "genesis member",
			},
		},
		"extra args": {
			[]string{
				"extra",
			},
			false,
			nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdMember()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual foundation.QueryMemberResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, actual.Member)
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdMembers() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{},
			true,
		},
		"extra args": {
			[]string{
				"extra",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdMembers()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual foundation.QueryMembersResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdProposal() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{
				"1",
			},
			true,
		},
		"extra args": {
			[]string{
				"1",
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
			cmd := cli.NewQueryCmdProposal()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual foundation.QueryProposalResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdProposals() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{},
			true,
		},
		"extra args": {
			[]string{
				"extra",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdProposals()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual foundation.QueryProposalsResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdVote() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
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
			cmd := cli.NewQueryCmdVote()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual foundation.QueryVoteResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdVotes() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{
				"1",
			},
			true,
		},
		"extra args": {
			[]string{
				"1",
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
			cmd := cli.NewQueryCmdVotes()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual foundation.QueryVotesResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdTallyResult() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{
				"1",
			},
			true,
		},
		"extra args": {
			[]string{
				"1",
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
			cmd := cli.NewQueryCmdTallyResult()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual foundation.QueryTallyResultResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestNewQueryCmdGrants() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected int
	}{
		"valid query": {
			[]string{
				s.stranger.String(),
				foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
			},
			true,
			1,
		},
		"no msg type url": {
			[]string{
				s.stranger.String(),
			},
			true,
			1,
		},
		"extra args": {
			[]string{
				s.stranger.String(),
				foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
				"extra",
			},
			false,
			0,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdGrants()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual foundation.QueryGrantsResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Equal(tc.expected, len(actual.Authorizations))
		})
	}
}
