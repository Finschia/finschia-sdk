package cli_test

import (
	"fmt"
	"io"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/stretchr/testify/suite"

	abci "github.com/cometbft/cometbft/abci/types"
	rpcclientmock "github.com/cometbft/cometbft/rpc/client/mock"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testutilmod "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/gov"

	"github.com/Finschia/finschia-sdk/x/collection"
	"github.com/Finschia/finschia-sdk/x/collection/client/cli"
)

type CLITestSuite struct {
	suite.Suite

	kr          keyring.Keyring
	encCfg      testutilmod.TestEncodingConfig
	baseCtx     client.Context
	clientCtx   client.Context
	commonFlags []string

	vendor   testutil.TestAccount
	operator testutil.TestAccount
	customer testutil.TestAccount
	stranger testutil.TestAccount

	contractID string
	classID    string
	tokenIdx   string
}

func TestCLITestSuite(t *testing.T) {
	suite.Run(t, new(CLITestSuite))
}

func (s *CLITestSuite) SetupSuite() {
	s.encCfg = testutilmod.MakeTestEncodingConfig(gov.AppModuleBasic{})
	s.kr = keyring.NewInMemory(s.encCfg.Codec)
	s.baseCtx = client.Context{}.
		WithKeyring(s.kr).
		WithTxConfig(s.encCfg.TxConfig).
		WithCodec(s.encCfg.Codec).
		WithClient(clitestutil.MockCometRPC{Client: rpcclientmock.Client{}}).
		WithAccountRetriever(client.MockAccountRetriever{}).
		WithOutput(io.Discard).
		WithChainID("test-chain")

	ctxGen := func() client.Context {
		bz, _ := s.encCfg.Codec.Marshal(&sdk.TxResponse{})
		c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
			Value: bz,
		})
		return s.baseCtx.WithClient(c)
	}
	s.clientCtx = ctxGen()

	s.commonFlags = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(10))).String()),
		fmt.Sprintf("--%s=test-chain", flags.FlagChainID),
	}

	val := testutil.CreateKeyringAccounts(s.T(), s.kr, 4)
	s.vendor = val[0]
	s.operator = val[1]
	s.customer = val[2]
	s.stranger = val[3]

	s.contractID = "678c146a"
	s.classID = "10000001"
	s.tokenIdx = "00000001"
}

func (s *CLITestSuite) TestNewTxCmdSendNFT() {

	tokenID := collection.NewNFTID(s.classID, 1)
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.vendor.Address.String(),
				s.customer.Address.String(),
				tokenID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.stranger.Address.String(),
				s.customer.Address.String(),
				tokenID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.stranger.Address.String(),
				s.customer.Address.String(),
			},
			false,
		},
		"amount out of range": {
			[]string{
				s.contractID,
				s.stranger.Address.String(),
				s.customer.Address.String(),
				fmt.Sprintf("%s:1%0127d", s.classID, 0),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.stranger.Address.String(),
				s.customer.Address.String(),
				tokenID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdSendNFT()
			out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, append(tc.args, s.commonFlags...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			resp := new(sdk.TxResponse)
			s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), resp))
		})
	}
}

func (s *CLITestSuite) TestNewTxCmdOperatorSendNFT() {

	tokenID := collection.NewNFTID(s.classID, 1)
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.Address.String(),
				s.vendor.Address.String(),
				s.customer.Address.String(),
				tokenID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.Address.String(),
				s.vendor.Address.String(),
				s.customer.Address.String(),
				tokenID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.Address.String(),
				s.vendor.Address.String(),
				s.customer.Address.String(),
			},
			false,
		},
		"amount out of range": {
			[]string{
				s.contractID,
				s.operator.Address.String(),
				s.vendor.Address.String(),
				s.customer.Address.String(),
				fmt.Sprintf("%s:1%0127d", s.classID, 0),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.Address.String(),
				s.vendor.Address.String(),
				s.customer.Address.String(),
				tokenID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdOperatorSendNFT()
			out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, append(tc.args, s.commonFlags...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			resp := new(sdk.TxResponse)
			s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), resp))
		})
	}
}

func (s *CLITestSuite) TestNewTxCmdCreateContract() {
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.vendor.Address.String(),
				fmt.Sprintf("--%s=%s", cli.FlagName, "arctic fox"),
				fmt.Sprintf("--%s=%s", cli.FlagMeta, "nft metadata"),
				fmt.Sprintf("--%s=%s", cli.FlagBaseImgURI, "contract base img uri"),
			},
			true,
		},
		"extra args": {
			[]string{
				s.vendor.Address.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{},
			false,
		},
		"invalid creator": {
			[]string{
				"",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdCreateContract()
			out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, append(tc.args, s.commonFlags...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			resp := new(sdk.TxResponse)
			s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), resp))
		})
	}
}

func (s *CLITestSuite) TestNewTxCmdIssueNFT() {
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.vendor.Address.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.vendor.Address.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.vendor.Address.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdIssueNFT()
			out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, append(tc.args, s.commonFlags...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			resp := new(sdk.TxResponse)
			s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), resp))
		})
	}
}

func (s *CLITestSuite) TestNewTxCmdMintNFT() {
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.vendor.Address.String(),
				s.customer.Address.String(),
				s.classID,
				fmt.Sprintf("--%s=%s", cli.FlagName, "arctic fox"),
				fmt.Sprintf("--%s=%s", cli.FlagMeta, "nft metadata"),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.vendor.Address.String(),
				s.customer.Address.String(),
				s.classID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.vendor.Address.String(),
				s.customer.Address.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.vendor.Address.String(),
				s.customer.Address.String(),
				s.classID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdMintNFT()
			out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, append(tc.args, s.commonFlags...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			resp := new(sdk.TxResponse)
			s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), resp))
		})
	}
}

func (s *CLITestSuite) TestNewTxCmdBurnNFT() {
	tokenID := collection.NewNFTID(s.classID, 2)
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.customer.Address.String(),
				tokenID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.customer.Address.String(),
				tokenID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.customer.Address.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.customer.Address.String(),
				tokenID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdBurnNFT()
			out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, append(tc.args, s.commonFlags...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			resp := new(sdk.TxResponse)
			s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), resp))
		})
	}
}

func (s *CLITestSuite) TestNewTxCmdOperatorOperatorBurnNFT() {
	tokenID := collection.NewNFTID(s.classID, 1)
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.Address.String(),
				s.customer.Address.String(),
				tokenID,
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.Address.String(),
				s.customer.Address.String(),
				tokenID,
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.Address.String(),
				s.customer.Address.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.Address.String(),
				s.customer.Address.String(),
				tokenID,
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdOperatorBurnNFT()
			out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, append(tc.args, s.commonFlags...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			resp := new(sdk.TxResponse)
			s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), resp))
		})
	}
}

func (s *CLITestSuite) TestNewTxCmdModify() {
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.Address.String(),
				s.classID,
				s.tokenIdx,
				collection.AttributeKeyName.String(),
				"tibetian fox",
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.Address.String(),
				s.classID,
				s.tokenIdx,
				collection.AttributeKeyName.String(),
				"tibetian fox",
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.Address.String(),
				s.classID,
				s.tokenIdx,
				collection.AttributeKeyName.String(),
			},
			false,
		},
		"invalid contract id": {
			[]string{
				"",
				s.operator.Address.String(),
				s.classID,
				s.tokenIdx,
				collection.AttributeKeyName.String(),
				"tibetian fox",
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdModify()
			out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, append(tc.args, s.commonFlags...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			resp := new(sdk.TxResponse)
			s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), resp))
		})
	}
}

func (s *CLITestSuite) TestNewTxCmdGrantPermission() {
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.operator.Address.String(),
				s.vendor.Address.String(),
				collection.LegacyPermissionMint.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.operator.Address.String(),
				s.vendor.Address.String(),
				collection.LegacyPermissionMint.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.operator.Address.String(),
				s.vendor.Address.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdGrantPermission()
			out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, append(tc.args, s.commonFlags...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			resp := new(sdk.TxResponse)
			s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), resp))
		})
	}
}

func (s *CLITestSuite) TestNewTxCmdRevokePermission() {
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.vendor.Address.String(),
				collection.LegacyPermissionModify.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.vendor.Address.String(),
				collection.LegacyPermissionModify.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.vendor.Address.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdRevokePermission()
			out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, append(tc.args, s.commonFlags...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), &res))
		})
	}
}

func (s *CLITestSuite) TestNewTxCmdAuthorizeOperator() {
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.customer.Address.String(),
				s.operator.Address.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.customer.Address.String(),
				s.operator.Address.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.customer.Address.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdAuthorizeOperator()
			out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, append(tc.args, s.commonFlags...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			resp := new(sdk.TxResponse)
			s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), resp))
		})
	}
}

func (s *CLITestSuite) TestNewTxCmdRevokeOperator() {
	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid transaction": {
			[]string{
				s.contractID,
				s.customer.Address.String(),
				s.operator.Address.String(),
			},
			true,
		},
		"extra args": {
			[]string{
				s.contractID,
				s.customer.Address.String(),
				s.operator.Address.String(),
				"extra",
			},
			false,
		},
		"not enough args": {
			[]string{
				s.contractID,
				s.customer.Address.String(),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.NewTxCmdRevokeOperator()
			out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, append(tc.args, s.commonFlags...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			resp := new(sdk.TxResponse)
			s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), resp))
		})
	}
}
