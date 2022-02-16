package testutil

import (
	"fmt"

	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/client/flags"
	"github.com/line/lbm-sdk/crypto/hd"
	"github.com/line/lbm-sdk/crypto/keyring"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	"github.com/line/lbm-sdk/testutil/network"
	sdk "github.com/line/lbm-sdk/types"
	bankcli "github.com/line/lbm-sdk/x/bank/client/cli"
	"github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/token/client/cli"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	setupHeight int64

	vendor   sdk.AccAddress
	customer sdk.AccAddress

	classes []token.Token

	balance sdk.Int
}

var commonArgs = []string{
	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.network = network.New(s.T(), s.cfg)
	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	s.vendor = s.createAccount("vendor")
	s.customer = s.createAccount("customer")

	s.classes = []token.Token{
		{
			Id:       "678c146a",
			Name:     "test",
			Symbol:   "ZERO",
			Decimals: 8,
			Mintable: true,
		},
		{
			Id:       "9be17165",
			Name:     "test",
			Symbol:   "ONE",
			Decimals: 8,
			Mintable: true,
		},
	}

	s.balance = sdk.NewInt(1000)

	// vendor creates 2 tokens
	s.createClass(s.vendor, s.vendor, s.classes[1].Name, s.classes[1].Symbol, s.balance, s.classes[1].Mintable)
	s.createClass(s.vendor, s.customer, s.classes[0].Name, s.classes[0].Symbol, s.balance, s.classes[0].Mintable)

	// customer approves vendor to transfer its tokens of the both classes, so vendor can do transferFrom later.
	for _, class := range s.classes {
		s.approve(class.Id, s.customer, s.vendor)
	}

	s.setupHeight, err = s.network.LatestHeight()
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) createClass(owner, to sdk.AccAddress, name, symbol string, supply sdk.Int, mintable bool) {
	val := s.network.Validators[0]
	args := append([]string{
		owner.String(),
		to.String(),
		name,
		symbol,
		fmt.Sprintf("--%s=%v", cli.FlagMintable, mintable),
		fmt.Sprintf("--%s=%s", cli.FlagSupply, supply),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, owner),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdIssue(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())

	s.Require().NoError(s.network.WaitForNextBlock())
}

// creates an account and send some coins to it for the future transactions.
func (s *IntegrationTestSuite) createAccount(uid string) sdk.AccAddress {
	val := s.network.Validators[0]
	keyInfo, _, err := val.ClientCtx.Keyring.NewMnemonic(uid, keyring.English, sdk.FullFundraiserPath, hd.Secp256k1)
	s.Require().NoError(err)
	addr := keyInfo.GetAddress()

	fee := sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1000)))
	args := append([]string{
		val.Address.String(),
		addr.String(),
		fee.String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
	}, commonArgs...)
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, bankcli.NewSendTxCmd(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())

	s.Require().NoError(s.network.WaitForNextBlock())

	return addr
}

func (s *IntegrationTestSuite) approve(classID string, approver, proxy sdk.AccAddress) {
	val := s.network.Validators[0]
	args := append([]string{
		classID,
		approver.String(),
		proxy.String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, approver),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdApprove(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())

	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}
