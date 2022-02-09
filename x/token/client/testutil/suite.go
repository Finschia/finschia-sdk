package testutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	"github.com/line/lbm-sdk/testutil"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	"github.com/line/lbm-sdk/testutil/network"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token/client/cli"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	vendor *network.Validator
	operator *network.Validator
	customer *network.Validator

	mintableClass string
	notMintableClass string

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

	// genesisState := s.cfg.GenesisState

	// var tokenData token.GenesisState
	// s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[token.ModuleName], &tokenData))

	// classes := []token.Token{
	// 	// mintable class
	// 	token.Token{
	// 		Id: s.mintableClass,
	// 		Name: "mintable",
	// 		Symbol: "OK",
	// 		Mintable: true,
	// 	},
	// 	// not mintable class
	// 	token.Token{
	// 		Id: s.notMintableClass,
	// 		Name: "notmintable",
	// 		Symbol: "NO",
	// 		Mintable: false,
	// 	},
	// }
	// tokenData.Classes = classes

	// var balances []token.Balance
	// for _, addr := range []sdk.AccAddress{s.vendor, s.operator, s.customer} {
	// 	balances = append(balances, token.Balance{
	// 		Address: addr.String(),
	// 		Tokens: []token.FT{
	// 			token.FT{
	// 				ClassId: s.mintableClass,
	// 				Amount: s.balance,
	// 			},
	// 			token.FT{
	// 				ClassId: s.notMintableClass,
	// 				Amount: s.balance,
	// 			},
	// 		},
	// 	})
	// }
	// tokenData.Balances = balances

	// var grants []token.Grant
	// for _, grantee := range []sdk.AccAddress{s.vendor, s.operator} {
	// 	for _, class := range []string{s.mintableClass, s.notMintableClass} {
	// 		for _, action := range []string{"mint", "burn", "modify"} {
	// 			grants = append(grants, token.Grant{
	// 				Grantee: grantee.String(),
	// 				ClassId: class,
	// 				Action: action,
	// 			})
	// 		}
	// 	}
	// }
	// tokenData.Grants = grants

	// var approves []token.Approve
	// for _, class := range []string{s.mintableClass, s.notMintableClass} {
	// 	approves = append(approves, token.Approve{
	// 		Approver: s.customer.String(),
	// 		Proxy: s.operator.String(),
	// 		ClassId: class,
	// 	})
	// }
	// tokenData.Approves = approves
	
	// tokenDataBz, err := s.cfg.Codec.MarshalJSON(&tokenData)
	// s.Require().NoError(err)
	// genesisState[token.ModuleName] = tokenDataBz
	// s.cfg.GenesisState = genesisState

	s.network = network.New(s.T(), s.cfg)
	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	s.vendor = s.network.Validators[0]
	s.operator = s.network.Validators[1]
	s.customer = s.network.Validators[2]

	s.mintableClass = "foodbabe"
	s.notMintableClass = "fee1dead"

	s.balance = sdk.NewInt(1000000)

	createClass(s.vendor.ClientCtx, s.vendor.Address, s.operator.Address, "Mintable", "OK", nil)
	createClass(s.vendor.ClientCtx, s.vendor.Address, s.vendor.Address, "NotMintable", "NO", &s.balance)
	s.Require().NoError(s.network.WaitForNextBlock())
}

func createClass(clientCtx client.Context, owner, to sdk.AccAddress, name, symbol string, mint *sdk.Int) (testutil.BufferWriter, error) {
	args := append([]string{
		owner.String(),
		to.String(),
		name,
		symbol,
	}, commonArgs...)

	supply := sdk.ZeroInt()
	if mint == nil {
		args = append(args, fmt.Sprintf("--%s", cli.FlagMintable))
	} else {
		supply = *mint
	}
	args = append(args, fmt.Sprintf("--%s=%s", cli.FlagSupply, supply))

	return clitestutil.ExecTestCLICmd(clientCtx, cli.NewTxCmdIssue(), args)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}
