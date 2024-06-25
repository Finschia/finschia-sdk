package cli_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	rpcclientmock "github.com/tendermint/tendermint/rpc/client/mock"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/crypto/hd"
	"github.com/Finschia/finschia-sdk/crypto/keyring"
	"github.com/Finschia/finschia-sdk/tests/mocks"
	clitestutil "github.com/Finschia/finschia-sdk/testutil/cli"
	sdk "github.com/Finschia/finschia-sdk/types"
	testutilmod "github.com/Finschia/finschia-sdk/types/module/testutil"
	"github.com/Finschia/finschia-sdk/x/fbridge/client/cli"
	fbridgem "github.com/Finschia/finschia-sdk/x/fbridge/module"
)

type CLITestSuite struct {
	suite.Suite

	kr        keyring.Keyring
	encCfg    testutilmod.TestEncodingConfig
	baseCtx   client.Context
	clientCtx client.Context
	addrs     []sdk.AccAddress
}

func TestCLITestSuite(t *testing.T) {
	suite.Run(t, new(CLITestSuite))
}

func (s *CLITestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	ar := mocks.NewMockAccountRetriever(ctrl)

	s.encCfg = testutilmod.MakeTestEncodingConfig(fbridgem.AppModule{})
	s.kr = keyring.NewInMemory()

	s.baseCtx = client.Context{}.
		WithKeyring(s.kr).
		WithTxConfig(s.encCfg.TxConfig).
		WithCodec(s.encCfg.Codec).
		WithInterfaceRegistry(s.encCfg.InterfaceRegistry).
		WithLegacyAmino(s.encCfg.Amino).
		WithClient(clitestutil.MockTendermintRPC{Client: rpcclientmock.Client{}}).
		WithAccountRetriever(ar).
		WithOutput(io.Discard).
		WithChainID("test-chain")

	ctxGen := func() client.Context {
		bz, _ := s.encCfg.Codec.Marshal(&sdk.TxResponse{})
		c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
			Value: bz,
		})
		return s.baseCtx.WithClient(c)
	}

	s.clientCtx = ctxGen()
	s.addrs = make([]sdk.AccAddress, 0)
	for i := 0; i < 3; i++ {
		k, _, err := s.clientCtx.Keyring.NewMnemonic(fmt.Sprintf("TestAccount-%d", i), keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
		s.Require().NoError(err)
		pub := k.GetPubKey()
		newAddr := sdk.AccAddress(pub.Address())
		s.addrs = append(s.addrs, newAddr)

		ar.EXPECT().EnsureExists(gomock.Any(), newAddr).Return(nil).AnyTimes()
		ar.EXPECT().GetAccountNumberSequence(gomock.Any(), newAddr).Return(uint64(i), uint64(1), nil).AnyTimes()
	}
}

func cliArgs(args ...string) []string {
	return append(args, []string{
		fmt.Sprintf("--%s=json", FlagOutput),
		fmt.Sprintf("--%s=home", flags.FlagKeyringDir),
		fmt.Sprintf("--%s=mynote", flags.FlagNote),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
		fmt.Sprintf("--%s=1.2", flags.FlagGasAdjustment),
		fmt.Sprintf("--%s=false", flags.FlagUseLedger),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=false", flags.FlagDryRun),
		fmt.Sprintf("--%s=false", flags.FlagGenerateOnly),
		fmt.Sprintf("--%s=false", flags.FlagOffline),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=direct", flags.FlagSignMode),
		fmt.Sprintf("--%s=%d", flags.FlagTimeoutHeight, 0),
	}...,
	)
}

func (s *CLITestSuite) TestNewTxCmd() {
	cmdQuery := []string{
		"add-vote-for-role",
		"set-bridge-status",
		"suggest-role",
		"transfer",
	}

	cmd := cli.NewTxCmd()
	for i, c := range cmd.Commands() {
		s.Require().Equal(cmdQuery[i], c.Name())
	}
}

func (s *CLITestSuite) TestNewTransferTxCmd() {
	cmd := cli.NewTransferTxCmd()
	s.Require().NotNil(cmd)

	tcs := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			name: "valid request",
			args: cliArgs(
				s.addrs[1].String(),
				"10stake",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addrs[0]),
			),
			expectErr:    false,
			respType:     &sdk.TxResponse{},
			expectedCode: 0,
		},
		{
			name: "invalid from address",
			args: cliArgs(
				s.addrs[0].String(),
				"10stake",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, "link1..."),
			),
			expectErr: true,
		},
		{
			name: "invalid decimal coin",
			args: cliArgs(
				s.addrs[1].String(),
				fmt.Sprintf("10%s", strings.Repeat("a", 300)),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addrs[0]),
			),
			expectErr: true,
		},
		{
			name: "more than one coin",
			args: cliArgs(
				s.addrs[1].String(),
				"10stake,20cony",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addrs[0]),
			),
			expectErr: true,
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
				tsResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, tsResp.Code, out.String())
			}
		})
	}
}

func (s *CLITestSuite) TestNewSuggestRoleTxCmd() {
	cmd := cli.NewSuggestRoleTxCmd()
	s.Require().NotNil(cmd)

	tcs := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			name: "invalid from address",
			args: cliArgs(
				s.addrs[1].String(),
				"guardian",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, "link1..."),
			),
			expectErr: true,
		},
		{
			name: "invalid role",
			args: cliArgs(
				s.addrs[1].String(),
				"random",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addrs[0].String()),
			),
			expectErr: true,
		},
		{
			name: "valid request",
			args: cliArgs(
				s.addrs[1].String(),
				"guardian",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addrs[0].String()),
			),
			expectErr:    false,
			respType:     &sdk.TxResponse{},
			expectedCode: 0,
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
				tsResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, tsResp.Code, out.String())
			}
		})
	}
}

func (s *CLITestSuite) TestNewAddVoteForRoleTxCmd() {
	cmd := cli.NewAddVoteForRoleTxCmd()
	s.Require().NotNil(cmd)

	tcs := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			name: "invalid from address",
			args: cliArgs(
				"1",
				"yes",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, "link1..."),
			),
			expectErr: true,
		},
		{
			name: "invalid proposal ID",
			args: cliArgs(
				"0xf",
				"yes",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addrs[0].String()),
			),
			expectErr: true,
		},
		{
			name: "invalid vote option",
			args: cliArgs(
				"1",
				"n/a",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addrs[0].String()),
			),
			expectErr: true,
		},
		{
			name: "valid request - yes",
			args: cliArgs(
				"1",
				"yes",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addrs[0].String()),
			),
			expectErr:    false,
			respType:     &sdk.TxResponse{},
			expectedCode: 0,
		},
		{
			name: "valid request - no",
			args: cliArgs(
				"1",
				"no",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addrs[0].String()),
			),
			expectErr:    false,
			respType:     &sdk.TxResponse{},
			expectedCode: 0,
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
				tsResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, tsResp.Code, out.String())
			}
		})
	}
}

func (s *CLITestSuite) TestNewSetBridgeStatusTxCmd() {
	cmd := cli.NewSetBridgeStatusTxCmd()
	s.Require().NotNil(cmd)
	tcs := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			name: "invalid from address",
			args: cliArgs(
				"halt",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, "link1..."),
			),
			expectErr: true,
		},
		{
			name: "invalid brdige status",
			args: cliArgs(
				"wrongstatus",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addrs[0].String()),
			),
			expectErr: true,
		},
		{
			name: "valid request - halt",
			args: cliArgs(
				"halt",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addrs[0].String()),
			),
			expectErr:    false,
			respType:     &sdk.TxResponse{},
			expectedCode: 0,
		},
		{
			name: "valid request - resume",
			args: cliArgs(
				"resume",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addrs[0].String()),
			),
			expectErr:    false,
			respType:     &sdk.TxResponse{},
			expectedCode: 0,
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
				tsResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, tsResp.Code, out.String())
			}
		})
	}
}
