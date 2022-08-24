package testutil

import (
	"fmt"
	"io/ioutil"
	"testing"

	ostcli "github.com/line/ostracon/libs/cli"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/client/flags"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	"github.com/line/lbm-sdk/testutil/network"
	sdk "github.com/line/lbm-sdk/types"
	wasmcli "github.com/line/lbm-sdk/x/wasm/client/cli"
	"github.com/line/lbm-sdk/x/wasm/keeper"
	wasmtypes "github.com/line/lbm-sdk/x/wasm/types"
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

	// TODO: should add below comment after adding InitGenesis of lbm's wasm
	// wasmkeeper.BuildContractAddress(1, 100)
	//s.inactiveContractAddress = "link1mujpjkwhut9yjw4xueyugc02evfv46y0dtmnz4lh8xxkkdapym9skz93hr"
	//
	//// add inactive contract to genesis
	//var wasmData lbmwasmtypes.GenesisState
	//genesisState := s.cfg.GenesisState
	//s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[wasmtypes.ModuleName], &wasmData))
	//wasmData.InactiveContractAddresses = []string{s.inactiveContractAddress}
	//
	//wasmDataBz, err := s.cfg.Codec.MarshalJSON(&wasmData)
	//s.Require().NoError(err)
	//genesisState[wasmtypes.ModuleName] = wasmDataBz
	//s.cfg.GenesisState = genesisState

	s.network = network.New(s.T(), s.cfg)
	_, err := s.network.WaitForHeight(1)
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

	wasmPath := "../../../keeper/testdata/hackatom.wasm"
	_, err := ioutil.ReadFile(wasmPath)
	s.Require().NoError(err)

	args := append([]string{
		wasmPath,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=%d", flags.FlagGas, 1500000),
	}, commonArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, wasmcli.StoreCodeCmd(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())

	// parse codeID
	for _, v := range res.Events {
		if v.Type == wasmtypes.EventTypeStoreCode {
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

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, wasmcli.InstantiateContractCmd(), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().EqualValues(0, res.Code, out.String())

	// parse contractAddress
	for _, v := range res.Events {
		if v.Type == wasmtypes.EventTypeInstantiate {
			return string(v.Attributes[0].Value)
		}
	}

	return ""
}
