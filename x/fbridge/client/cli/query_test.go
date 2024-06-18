package cli_test

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Finschia/finschia-sdk/client"
	clitestutil "github.com/Finschia/finschia-sdk/testutil/cli"
	"github.com/Finschia/finschia-sdk/x/fbridge/client/cli"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

const FlagOutput = "output"

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
		expectErr    bool
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
			false,
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cmd.SetArgs(tc.args)
			_, err := clitestutil.ExecTestCLICmd(tc.ctxGen(), cmd, tc.args)
			s.Require().NoError(err)
		})
	}
}
