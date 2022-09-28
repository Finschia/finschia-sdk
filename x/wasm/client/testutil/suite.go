package testutil

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/line/lbm-sdk/client/flags"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	"github.com/line/lbm-sdk/testutil/network"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/wasm/client/cli"
	"github.com/line/lbm-sdk/x/wasm/keeper"
	"github.com/line/lbm-sdk/x/wasm/types"
	ostcli "github.com/line/ostracon/libs/cli"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	setupHeight int64

	codeID                  string
	contractAddress         string
	nonExistValidAddress    string
	inactiveContractAddress string

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

	s.inactiveContractAddress = "link14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sgf2vn8"

	// add inactive contract to genesis
	var wasmData types.GenesisState
	genesisState := s.cfg.GenesisState
	genesisData, err := ioutil.ReadFile("./testdata/wasm_genesis.json")
	s.Require().NoError(err)
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisData, &wasmData))
	wasmDataBz, err := s.cfg.Codec.MarshalJSON(&wasmData)
	s.Require().NoError(err)
	genesisState[types.ModuleName] = wasmDataBz
	s.cfg.GenesisState = genesisState

	s.network = network.New(s.T(), s.cfg)
	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	// deploy contract
	s.codeID = s.deployContract()

	s.verifier = s.network.Validators[0].Address.String()
	s.beneficiary = keeper.RandomAccountAddress(s.T())
	params := fmt.Sprintf("{\"verifier\": \"%s\", \"beneficiary\": \"%s\"}", s.verifier, s.beneficiary)
	s.contractAddress = s.instantiate(s.codeID, params)

	s.nonExistValidAddress = keeper.RandomAccountAddress(s.T()).String()

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
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=%d", flags.FlagGas, 1500000),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.StoreCodeCmd(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())

	// parse codeID
	for _, v := range res.Events {
		if v.Type == types.EventTypeStoreCode {
			return string(v.Attributes[0].Value)
		}
	}

	return ""
}

func (s *IntegrationTestSuite) instantiate(codeID, params string) string {
	val := s.network.Validators[0]
	owner := val.Address.String()

	args := append([]string{
		codeID,
		params,
		fmt.Sprintf("--label=%s", "TestContract"),
		fmt.Sprintf("--admin=%s", owner),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.InstantiateContractCmd(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())

	// parse contractAddress
	for _, v := range res.Events {
		if v.Type == types.EventTypeInstantiate {
			return string(v.Attributes[0].Value)
		}
	}

	return ""
}
