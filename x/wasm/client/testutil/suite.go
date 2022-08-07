package testutil

import (
	"fmt"
	"github.com/line/lbm-sdk/client/flags"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	"github.com/line/lbm-sdk/testutil/network"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/wasm/client/cli"
	"github.com/line/lbm-sdk/x/wasm/keeper"
	ostcli "github.com/line/ostracon/libs/cli"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"testing"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	setupHeight int64

	codeId          string
	contractAddress string

	// for hackatom contract
	verifier    string
	beneficiary sdk.AccAddress
}

var commonArgs = []string{
	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))).String()),
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("Setting up integration test suite.")

	if testing.Short() {
		s.T().Skip("skipping test in unit-tests mode.")
	}

	s.network = network.New(s.T(), s.cfg)
	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	// deploy contract
	s.codeId = s.deployContract()
	fmt.Printf("codeId: %s\n", s.codeId)

	s.verifier = s.network.Validators[0].Address.String()
	s.beneficiary = keeper.RandomAccountAddress(s.T())
	params := fmt.Sprintf("{\"verifier\": \"%s\", \"beneficiary\": \"%s\"}", s.verifier, s.beneficiary)
	s.contractAddress = s.instantiate(s.codeId, params)
	fmt.Printf("contractAddress: %s\n", s.contractAddress)

	s.setupHeight, err = s.network.LatestHeight()
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) queryCommonArgs() []string {
	return []string{
		fmt.Sprintf("--%s=%d", flags.FlagHeight, s.setupHeight),
		fmt.Sprintf("--%s=json", ostcli.OutputFlag),
	}
}

func (s *IntegrationTestSuite) deployContract() string {
	val := s.network.Validators[0]

	wasmPath := "../../keeper/testdata/hackatom.wasm"
	_, err := ioutil.ReadFile(wasmPath)
	s.Require().NoError(err)

	args := append([]string{
		wasmPath,
		fmt.Sprintf("--%s=%v", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=%v", flags.FlagGas, 1500000),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.StoreCodeCmd(), args)
	fmt.Printf("err: %v\n", err)
	s.Require().NoError(err)
	fmt.Printf("out: %v\n", out)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())

	// parse codeId
	for _, v := range res.Events {
		if v.Type == "store_code" {
			return string(v.Attributes[0].Value)
		}
	}

	return ""
}

func (s *IntegrationTestSuite) instantiate(codeId, params string) string {
	val := s.network.Validators[0]
	owner := val.Address.String()

	args := append([]string{
		codeId,
		params,
		fmt.Sprintf("--label=%v", "TestContract"),
		fmt.Sprintf("--admin=%v", owner),
		fmt.Sprintf("--%s=%v", flags.FlagFrom, val.Address.String()),
	}, commonArgs...)

	//fmt.Printf("instantiate args: %v\n", args)
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.InstantiateContractCmd(), args)
	s.Require().NoError(err)
	//fmt.Printf("instantiate out: %v\n", out)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())

	// parse contractAddress
	for _, v := range res.Events {
		if v.Type == "instantiate" {
			return string(v.Attributes[0].Value)
		}
	}

	return ""
}
