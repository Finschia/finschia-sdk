package cli_test

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	clitestutil "github.com/Finschia/finschia-sdk/testutil/cli"
	"github.com/Finschia/finschia-sdk/x/fbridge/client/cli"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

const FlagOutput = "output"

func (s *CLITestSuite) TestNewQueryCmd() {
	cmdQuery := []string{
		"member",
		"members",
		"params",
		"proposal",
		"proposals",
		"sending-next-seq",
		"seq-to-blocknums",
		"status",
		"vote",
		"votes",
	}

	cmd := cli.NewQueryCmd()
	for i, c := range cmd.Commands() {
		s.Require().Equal(cmdQuery[i], c.Name())
	}
}

func (s *CLITestSuite) TestQueryParams() {
	cmd := cli.NewQueryParamsCmd()
	s.Require().NotNil(cmd)
	cmd.SetOut(io.Discard)

	tcs := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
	}{
		{
			"json output",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryParamsResponse{
					Params: types.DefaultParams(),
				})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{fmt.Sprintf("--%s=json", FlagOutput)},
			&types.QueryParamsResponse{},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			var outBuf bytes.Buffer
			clientCtx := tc.ctxGen().WithOutput(&outBuf)
			cmd.SetArgs(tc.args)
			res, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)
			s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(res.Bytes(), tc.expectResult))
		})
	}
}

func (s *CLITestSuite) TestQueryNextSeqSend() {
	cmd := cli.NewQueryNextSeqSendCmd()
	s.Require().NotNil(cmd)
	cmd.SetOut(io.Discard)

	tcs := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
	}{
		{
			"json output",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryNextSeqSendResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{fmt.Sprintf("--%s=json", FlagOutput)},
			&types.QueryNextSeqSendResponse{},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cmd.SetArgs(tc.args)
			res, err := clitestutil.ExecTestCLICmd(tc.ctxGen(), cmd, tc.args)
			s.Require().NoError(err)
			s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(res.Bytes(), tc.expectResult))
		})
	}
}

func (s *CLITestSuite) TestQuerySeqToBlocknumsCmd() {
	cmd := cli.NewQuerySeqToBlocknumsCmd()
	s.Require().NotNil(cmd)
	cmd.SetOut(io.Discard)

	tcs := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
		expectErr    bool
	}{
		{
			"json output",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QuerySeqToBlocknumsResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				fmt.Sprintf("--sequences=1"),
				fmt.Sprintf("--%s=json", FlagOutput),
			},
			&types.QuerySeqToBlocknumsResponse{},
			false,
		},
		{
			"invalid seq",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QuerySeqToBlocknumsResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{fmt.Sprintf("--sequences=1.3"), fmt.Sprintf("--%s=json", FlagOutput)},
			&types.QuerySeqToBlocknumsResponse{},
			true,
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cmd.SetArgs(tc.args)
			res, err := clitestutil.ExecTestCLICmd(tc.ctxGen(), cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Error(s.encCfg.Codec.UnmarshalJSON(res.Bytes(), tc.expectResult))
			} else {
				s.Require().NoError(err)
				s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(res.Bytes(), tc.expectResult))
			}
		})
	}
}

func (s *CLITestSuite) TestQueryMembersCmd() {
	cmd := cli.NewQueryMembersCmd()
	s.Require().NotNil(cmd)
	cmd.SetOut(io.Discard)

	tcs := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
	}{
		{
			"json output",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryMembersResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{"guardian", fmt.Sprintf("--%s=json", FlagOutput)},
			&types.QueryMembersResponse{},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cmd.SetArgs(tc.args)
			res, err := clitestutil.ExecTestCLICmd(tc.ctxGen(), cmd, tc.args)
			s.Require().NoError(err)
			s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(res.Bytes(), tc.expectResult))
		})
	}
}

func (s *CLITestSuite) TestQueryMemberCmd() {
	cmd := cli.NewQueryMemberCmd()
	s.Require().NotNil(cmd)
	cmd.SetOut(io.Discard)

	tcs := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
	}{
		{
			"json output",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryMemberResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				s.addrs[0].String(),
				fmt.Sprintf("--%s=json", FlagOutput),
			},
			&types.QueryMemberResponse{},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cmd.SetArgs(tc.args)
			res, err := clitestutil.ExecTestCLICmd(tc.ctxGen(), cmd, tc.args)
			s.Require().NoError(err)
			s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(res.Bytes(), tc.expectResult))
		})
	}
}

func (s *CLITestSuite) TestNewQueryProposalsCmd() {
	cmd := cli.NewQueryProposalsCmd()
	s.Require().NotNil(cmd)
	cmd.SetOut(io.Discard)

	tcs := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
	}{
		{
			"json output",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryProposalsResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{fmt.Sprintf("--%s=json", FlagOutput)},
			&types.QueryProposalsResponse{},
		},
		{
			"pagination",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryProposalsResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				fmt.Sprintf("--%s=100", flags.FlagLimit),
				fmt.Sprintf("--%s=20", flags.FlagOffset),
				fmt.Sprintf("--%s=true", flags.FlagCountTotal),
				fmt.Sprintf("--%s=false", flags.FlagReverse),
				fmt.Sprintf("--%s=json", FlagOutput),
			},
			&types.QueryProposalsResponse{},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cmd.SetArgs(tc.args)
			res, err := clitestutil.ExecTestCLICmd(tc.ctxGen(), cmd, tc.args)
			s.Require().NoError(err)
			s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(res.Bytes(), tc.expectResult))
		})
	}
}

func (s *CLITestSuite) TestNewQueryProposalCmd() {
	cmd := cli.NewQueryProposalCmd()
	s.Require().NotNil(cmd)
	cmd.SetOut(io.Discard)

	tcs := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
		expectErr    bool
	}{
		{
			"json output",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryProposalResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				"1",
				fmt.Sprintf("--%s=json", FlagOutput),
			},
			&types.QueryProposalResponse{},
			false,
		},
		{
			"no proposal ID",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryProposalResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				fmt.Sprintf("--%s=json", FlagOutput),
			},
			&types.QueryProposalResponse{},
			true,
		},
		{
			"wrong proposal ID",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryProposalResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				"one",
				fmt.Sprintf("--%s=json", FlagOutput),
			},
			&types.QueryProposalResponse{},
			true,
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cmd.SetArgs(tc.args)
			res, err := clitestutil.ExecTestCLICmd(tc.ctxGen(), cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Error(s.encCfg.Codec.UnmarshalJSON(res.Bytes(), tc.expectResult))
			} else {
				s.Require().NoError(err)
				s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(res.Bytes(), tc.expectResult))
			}
		})
	}
}

func (s *CLITestSuite) TestNewQueryVotesCmd() {
	cmd := cli.NewQueryVotesCmd()
	s.Require().NotNil(cmd)
	cmd.SetOut(io.Discard)

	tcs := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
		expectErr    bool
	}{
		{
			"json output",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryVotesResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				"1",
				fmt.Sprintf("--%s=json", FlagOutput),
			},
			&types.QueryVotesResponse{},
			false,
		},
		{
			"no proposal ID",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryVotesResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				fmt.Sprintf("--%s=json", FlagOutput),
			},
			&types.QueryVotesResponse{},
			true,
		},
		{
			"wrong proposal ID",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryVotesResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				"one",
				fmt.Sprintf("--%s=json", FlagOutput),
			},
			&types.QueryVotesResponse{},
			true,
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cmd.SetArgs(tc.args)
			res, err := clitestutil.ExecTestCLICmd(tc.ctxGen(), cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Error(s.encCfg.Codec.UnmarshalJSON(res.Bytes(), tc.expectResult))
			} else {
				s.Require().NoError(err)
				s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(res.Bytes(), tc.expectResult))
			}
		})
	}
}

func (s *CLITestSuite) TestNewQueryVoteCmd() {
	cmd := cli.NewQueryVoteCmd()
	s.Require().NotNil(cmd)
	cmd.SetOut(io.Discard)

	tcs := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
		expectErr    bool
	}{
		{
			"json output",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryVoteResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				"1",
				s.addrs[0].String(),
				fmt.Sprintf("--%s=json", FlagOutput),
			},
			&types.QueryVoteResponse{},
			false,
		},
		{
			"no voter",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryVoteResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				"1",
				fmt.Sprintf("--%s=json", FlagOutput),
			},
			&types.QueryVoteResponse{},
			true,
		},
		{
			"no proposal ID",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryVoteResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				fmt.Sprintf("--%s=json", FlagOutput),
			},
			&types.QueryVoteResponse{},
			true,
		},
		{
			"wrong proposal ID",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryVoteResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				"one",
				s.addrs[0].String(),
				fmt.Sprintf("--%s=json", FlagOutput),
			},
			&types.QueryVoteResponse{},
			true,
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cmd.SetArgs(tc.args)
			res, err := clitestutil.ExecTestCLICmd(tc.ctxGen(), cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Error(s.encCfg.Codec.UnmarshalJSON(res.Bytes(), tc.expectResult))
			} else {
				s.Require().NoError(err)
				s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(res.Bytes(), tc.expectResult))
			}
		})
	}
}

func (s *CLITestSuite) TestQueryBridgeStatusCmd() {
	cmd := cli.NewQueryBridgeStatusCmd()
	s.Require().NotNil(cmd)
	cmd.SetOut(io.Discard)

	tcs := []struct {
		name         string
		ctxGen       func() client.Context
		args         []string
		expectResult proto.Message
	}{
		{
			"json output",
			func() client.Context {
				bz, _ := s.encCfg.Codec.Marshal(&types.QueryBridgeStatusResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return s.baseCtx.WithClient(c)
			},
			[]string{
				fmt.Sprintf("--%s=json", FlagOutput),
			},
			&types.QueryBridgeStatusResponse{},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cmd.SetArgs(tc.args)
			res, err := clitestutil.ExecTestCLICmd(tc.ctxGen(), cmd, tc.args)
			s.Require().NoError(err)
			s.Require().NoError(s.encCfg.Codec.UnmarshalJSON(res.Bytes(), tc.expectResult))
		})
	}
}
