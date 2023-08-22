package testutil

import (
	"bytes"
	"compress/zlib"
	"fmt"
	clitestutil "github.com/Finschia/finschia-sdk/testutil/cli"
	sdk "github.com/Finschia/finschia-sdk/types"
	"time"

	dacli "github.com/Finschia/finschia-sdk/x/or/da/client/cli"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

func (s *IntegrationTestSuite) TestCmdTxAppendCCBatch() {
	val := s.network.Validators[0]

	tcs := map[string]struct {
		args  []string
		valid bool
	}{
		"valid CCBatch": {
			[]string{
				s.sequencer.GetAddress().String(),
				rollupName,
				func() string {
					batch := types.CCBatch{
						ShouldStartAtFrame: 1,
						Frames: []*types.CCBatchFrame{
							{
								Header: &types.CCBatchHeader{
									ParentHash: []byte("parent_hash"),
									Timestamp:  time.Now().UTC(),
									L2Height:   2,
									L1Height:   100,
								},
								Elements: []*types.CCBatchElement{
									{
										QueueIndex: 1,
									},
									{
										Txraw: []byte("Sequencer Tx"),
									},
								},
							},
						},
					}
					serializedBatch := s.cfg.Codec.MustMarshal(&batch)
					var b bytes.Buffer
					w := zlib.NewWriter(&b)
					_, err := w.Write(serializedBatch)
					s.Require().NoError(err)
					err = w.Close()
					return string(b.Bytes())
				}(),
				fmt.Sprintf("--%s=%d", dacli.FlagCompression, 1),
			},
			true,
		},
	}

	for name, tc := range tcs {
		tc := tc

		s.Run(name, func() {
			cmd := dacli.CmdTxAppendCCBatch()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			s.Require().EqualValues(0, res.Code, out.String())
		})
	}
}
