package foundation

import (
	"fmt"

	cmtcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/gogoproto/proto"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/Finschia/finschia-sdk/x/foundation/client/cli"
)

func (s *E2ETestSuite) TestNewQueryCmdParams() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", cmtcli.OutputFlag),
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
					FoundationTax: math.LegacyMustNewDecFromStr("0.2"),
				},
			},
		},
		"wrong number of args": {
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

func (s *E2ETestSuite) TestNewQueryCmdTreasury() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", cmtcli.OutputFlag),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{},
			true,
		},
		"wrong number of args": {
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

func (s *E2ETestSuite) TestNewQueryCmdFoundationInfo() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", cmtcli.OutputFlag),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{},
			true,
		},
		"wrong number of args": {
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

func (s *E2ETestSuite) TestNewQueryCmdMember() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", cmtcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected *foundation.Member
	}{
		"valid query": {
			[]string{
				s.permanentMember.String(),
			},
			true,
			&foundation.Member{
				Address:  s.permanentMember.String(),
				Metadata: "permanent member",
			},
		},
		"wrong number of args": {
			[]string{
				s.permanentMember.String(),
				"extra",
			},
			false,
			nil,
		},
		"invalid member": {
			[]string{
				"",
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

func (s *E2ETestSuite) TestNewQueryCmdMembers() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", cmtcli.OutputFlag),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{},
			true,
		},
		"wrong number of args": {
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

func (s *E2ETestSuite) TestNewQueryCmdProposal() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", cmtcli.OutputFlag),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{
				fmt.Sprintf("%d", s.proposalID),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				fmt.Sprintf("%d", s.proposalID),
				"extra",
			},
			false,
		},
		"invalid id": {
			[]string{
				fmt.Sprintf("%d", -1),
			},
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

func (s *E2ETestSuite) TestNewQueryCmdProposals() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", cmtcli.OutputFlag),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{},
			true,
		},
		"wrong number of args": {
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

func (s *E2ETestSuite) TestNewQueryCmdVote() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", cmtcli.OutputFlag),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
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
		"invalid proposal id": {
			[]string{
				fmt.Sprintf("%d", -1),
				s.permanentMember.String(),
			},
			false,
		},
		"invalid voter": {
			[]string{
				fmt.Sprintf("%d", s.proposalID),
				"",
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

func (s *E2ETestSuite) TestNewQueryCmdVotes() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", cmtcli.OutputFlag),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{
				fmt.Sprintf("%d", s.proposalID),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				fmt.Sprintf("%d", s.proposalID),
				"extra",
			},
			false,
		},
		"invalid proposal id": {
			[]string{
				fmt.Sprintf("%d", -1),
			},
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

func (s *E2ETestSuite) TestNewQueryCmdTallyResult() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", cmtcli.OutputFlag),
	}

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{
				fmt.Sprintf("%d", s.proposalID),
			},
			true,
		},
		"wrong number of args": {
			[]string{
				fmt.Sprintf("%d", s.proposalID),
				"extra",
			},
			false,
		},
		"invalid proposal id": {
			[]string{
				fmt.Sprintf("%d", -1),
			},
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

func (s *E2ETestSuite) TestNewQueryCmdCensorships() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", cmtcli.OutputFlag),
	}

	testCases := map[string]struct {
		args     []string
		valid    bool
		expected int
	}{
		"valid query": {
			[]string{},
			true,
			1,
		},
		"wrong number of args": {
			[]string{
				"extra",
			},
			false,
			0,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewQueryCmdCensorships()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var actual foundation.QueryCensorshipsResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &actual), out.String())
			s.Require().Len(actual.Censorships, tc.expected)
		})
	}
}

func (s *E2ETestSuite) TestNewQueryCmdGrants() {
	val := s.network.Validators[0]
	commonArgs := []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", cmtcli.OutputFlag),
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
		"wrong number of args": {
			[]string{
				s.stranger.String(),
				foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
				"extra",
			},
			false,
			0,
		},
		"invalid grantee": {
			[]string{
				"",
				foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
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
