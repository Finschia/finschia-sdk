package testutil

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"time"

	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	"github.com/Finschia/finschia-sdk/simapp"
	simappparams "github.com/Finschia/finschia-sdk/simapp/params"
	clitestutil "github.com/Finschia/finschia-sdk/testutil/cli"
	sdk "github.com/Finschia/finschia-sdk/types"
	banktypes "github.com/Finschia/finschia-sdk/x/bank/types"
	dacli "github.com/Finschia/finschia-sdk/x/or/da/client/cli"
	"github.com/Finschia/finschia-sdk/x/or/da/testutil"
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
		"invalid CCBatch": {
			[]string{
				s.sequencer.GetAddress().String(),
				rollupName,
				"wrong batch",
			},
			false,
		},
	}

	for name, tc := range tcs {
		tc := tc

		s.Run(name, func() {
			cmd := dacli.CmdTxAppendCCBatch()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			if tc.valid {
				s.Require().EqualValues(0, res.Code, out.String())
				return
			} else {
				s.Require().NotEqualValues(0, res.Code, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestCmdTxEnqueue() {
	val := s.network.Validators[0]

	tcs := map[string]struct {
		args  []string
		valid bool
	}{
		"valid enqueue": {
			[]string{
				s.sequencer.GetAddress().String(),
				rollupName,
				string(s.genMockTxs(1)[0]),
				fmt.Sprintf("--%s=%d", dacli.FlagL2GasLimit, 500000),
			},
			true,
		},
		"invalid queue tx": {
			[]string{
				s.sequencer.GetAddress().String(),
				rollupName,
				"wrong tx",
				fmt.Sprintf("--%s=%d", dacli.FlagL2GasLimit, 500000),
			},
			false,
		},
	}

	for name, tc := range tcs {
		tc := tc

		s.Run(name, func() {
			cmd := dacli.CmdTxEnqueue()
			out, _ := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			if tc.valid {
				s.Require().EqualValues(0, res.Code, out.String())
				return
			} else {
				s.Require().NotEqualValues(0, res.Code, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) genMockTxs(numToGenerate int) [][]byte {
	msg := banktypes.NewMsgSend(s.sequencer.GetAddress(), testutil.AccAddress(), sdk.NewCoins(sdk.NewInt64Coin("fnsa", 100)))
	txCfg := simappparams.MakeTestEncodingConfig().TxConfig
	txs, err := simapp.GenSequenceOfTxs(txCfg, []sdk.Msg{msg}, []uint64{0}, []uint64{uint64(0)}, numToGenerate, secp256k1.GenPrivKey())
	s.Require().NoError(err)

	serializedTxs := make([][]byte, numToGenerate)
	for i := 0; i < numToGenerate; i++ {
		serializedTxs[i], err = txCfg.TxEncoder()(txs[i])
		s.Require().NoError(err)
	}

	return serializedTxs
}
