package testutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	"github.com/line/lbm-sdk/testutil/network"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/token/client/cli"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	setupHeight int64

	vendor *network.Validator
	operator *network.Validator
	customer *network.Validator

	mintableClass token.Token
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

	s.vendor = s.network.Validators[0]
	s.operator = s.network.Validators[1]
	s.customer = s.network.Validators[2]

	s.mintableClass = token.Token{
		Id: "9be17165",
		Name: "Mintable",
		Symbol: "OK",
		Decimals: 8,
		Mintable: true,
	}
	s.notMintableClass = token.Token{
		Id: "678c146a",
		Name: "NotMintable",
		Symbol: "NO",
		Decimals: 8,
		Mintable: false,
	}

	s.balance = sdk.NewInt(1000000)

	err = createClass(s.vendor.ClientCtx, s.vendor.Address,
		s.mintableClass.Name, s.mintableClass.Symbol, s.balance, s.mintableClass.Mintable)
	s.Require().NoError(err)
	err = createClass(s.vendor.ClientCtx, s.vendor.Address,
		s.notMintableClass.Name, s.notMintableClass.Symbol, s.balance, s.notMintableClass.Mintable)
	s.Require().NoError(err)

	s.Require().NoError(s.network.WaitForNextBlock())
	s.setupHeight, err = s.network.LatestHeight()
	s.Require().NoError(err)
}

func createClass(clientCtx client.Context, owner sdk.AccAddress, name, symbol string, supply sdk.Int, mintable bool) error {
	args := append([]string{
		owner.String(),
		owner.String(),
		name,
		symbol,
	}, commonArgs...)

	args = append(args, fmt.Sprintf("--%s=%v", cli.FlagMintable, mintable))
	args = append(args, fmt.Sprintf("--%s=%s", cli.FlagSupply, supply))
	args = append(args, fmt.Sprintf("--%s=%s", flags.FlagFrom, owner))

	_, err := clitestutil.ExecTestCLICmd(clientCtx, cli.NewTxCmdIssue(), args)
	return err
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}
