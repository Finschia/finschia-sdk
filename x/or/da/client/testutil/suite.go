package testutil

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/Finschia/finschia-sdk/codec"
	codectypes "github.com/Finschia/finschia-sdk/codec/types"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/crypto/hd"
	"github.com/Finschia/finschia-sdk/crypto/keyring"
	clitestutil "github.com/Finschia/finschia-sdk/testutil/cli"
	"github.com/Finschia/finschia-sdk/testutil/network"
	sdk "github.com/Finschia/finschia-sdk/types"
	bankcli "github.com/Finschia/finschia-sdk/x/bank/client/cli"
	dacli "github.com/Finschia/finschia-sdk/x/or/da/client/cli"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
	rollupcli "github.com/Finschia/finschia-sdk/x/or/rollup/client/cli"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg         network.Config
	network     *network.Network
	setupHeight int64

	sequencer keyring.Info
}

var commonArgs = []string{
	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
}

const rollupName = "test-rollup"

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	var gs types.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(s.cfg.GenesisState[types.ModuleName], &gs))
	gsBz, err := s.cfg.Codec.MarshalJSON(&gs)
	s.Require().NoError(err)
	s.cfg.GenesisState[types.ModuleName] = gsBz

	s.network = network.New(s.T(), s.cfg)
	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	s.sequencer = s.createAccount("sequencer")

	s.createRollup()
	s.registerSequencer()
	s.enqueue(s.sequencer.GetAddress())

	batch := types.CCBatch{
		ShouldStartAtFrame: 0,
		Frames: []*types.CCBatchFrame{
			{
				Header: &types.CCBatchHeader{
					ParentHash: []byte("parent_hash"),
					Timestamp:  time.Now().UTC(),
					L2Height:   1,
					L1Height:   100,
				},
				Elements: []*types.CCBatchElement{
					{
						Txraw: []byte("Transactions"),
					},
				},
			},
		},
	}
	serializedBatch := s.cfg.Codec.MustMarshal(&batch)
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err = w.Write(serializedBatch)
	s.Require().NoError(err)
	err = w.Close()
	s.Require().NoError(err)
	s.appendCCBatch(s.sequencer.GetAddress(), b.Bytes())

	s.setupHeight, err = s.network.LatestHeight()
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) createAccount(uid string) keyring.Info {
	val := s.network.Validators[0]
	keyInfo, _, err := val.ClientCtx.Keyring.NewMnemonic(uid, keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	s.Require().NoError(err)
	addr := keyInfo.GetAddress()

	fee := sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1000000)))
	args := append([]string{
		val.Address.String(),
		addr.String(),
		fee.String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
	}, commonArgs...)
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, bankcli.NewSendTxCmd(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())

	return keyInfo
}

func (s *IntegrationTestSuite) createRollup() {
	val := s.network.Validators[0]
	sequencer := s.sequencer.GetAddress()

	args := append([]string{
		rollupName,
		"100",
		fmt.Sprintf(`{"addresses": ["%s"]}`, sequencer.String()),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, sequencer.String()),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, rollupcli.NewCreateRollupCmd(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())
}

func (s *IntegrationTestSuite) registerSequencer() {
	val := s.network.Validators[0]
	apk, err := codectypes.NewAnyWithValue(s.sequencer.GetPubKey())
	s.Require().NoError(err)
	bz, err := codec.ProtoMarshalJSON(apk, nil)
	s.Require().NoError(err)

	args := append([]string{
		rollupName,
		string(bz),
		sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1000)).String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.sequencer.GetAddress().String()),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, rollupcli.NewRegisterSequencerCmd(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())
}

func (s *IntegrationTestSuite) enqueue(sequencer sdk.Address) {
	val := s.network.Validators[0]
	args := append([]string{
		sequencer.String(),
		rollupName,
		string(s.genMockTxs(1)[0]),
		fmt.Sprintf("--%s=%d", dacli.FlagL2GasLimit, 300000),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, dacli.CmdTxEnqueue(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())
}

func (s *IntegrationTestSuite) appendCCBatch(sequencer sdk.Address, batch []byte) {
	val := s.network.Validators[0]
	args := append([]string{
		sequencer.String(),
		rollupName,
		string(batch),
		fmt.Sprintf("--%s=%d", dacli.FlagCompression, 1),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, dacli.CmdTxAppendCCBatch(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())
}
