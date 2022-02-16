package testutil

import (
	"fmt"
	"testing"

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

	mintableClass    token.Token
	notMintableClass token.Token

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

	if testing.Short() {
		s.T().Skip("skipping test in unit-tests mode.")
	}

	s.network = network.New(s.T(), s.cfg)
	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	s.vendor = s.createAccount("vendor")
	s.customer = s.createAccount("customer")

	s.mintableClass = token.Token{
		Id:       "9be17165",
		Name:     "Mintable",
		Symbol:   "OK",
		Decimals: 8,
		Mintable: true,
	}
	s.notMintableClass = token.Token{
		Id:       "678c146a",
		Name:     "NotMintable",
		Symbol:   "NO",
		Decimals: 8,
		Mintable: false,
	}

	s.balance = sdk.NewInt(1000000)

	// vendor creates 2 tokens: mintable and not mintable one. And mint to customer.
	s.createClass(s.vendor, s.customer, s.mintableClass.Name, s.mintableClass.Symbol, s.balance, s.mintableClass.Mintable)
	s.createClass(s.vendor, s.customer, s.notMintableClass.Name, s.notMintableClass.Symbol, s.balance, s.notMintableClass.Mintable)

	// customer approves vendor to transfer its tokens of the both classes, so vendor can do transferFrom later.
	s.approve(s.mintableClass.Id, s.customer, s.vendor)
	s.approve(s.notMintableClass.Id, s.customer, s.vendor)

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

	_, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdIssue(), args)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())
}

// creates an account and send some coins to it for the future transactions.
func (s *IntegrationTestSuite) createAccount(uid string) sdk.AccAddress {
	val := s.network.Validators[0]
	keyInfo, _, err := val.ClientCtx.Keyring.NewMnemonic(uid, keyring.English, sdk.FullFundraiserPath, hd.Secp256k1)
	s.Require().NoError(err)
	addr := keyInfo.GetAddress()

	fee := sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1000000)))
	args := append([]string{
		val.Address.String(),
		addr.String(),
		fee.String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
	}, commonArgs...)
	_, err = clitestutil.ExecTestCLICmd(val.ClientCtx, bankcli.NewSendTxCmd(), args)
	s.Require().NoError(err)
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

	_, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewTxCmdApprove(), args)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}
