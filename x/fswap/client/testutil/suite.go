package testutil

import (
	"github.com/stretchr/testify/suite"

	"github.com/Finschia/finschia-sdk/testutil/network"
	sdk "github.com/Finschia/finschia-sdk/types"
	banktypes "github.com/Finschia/finschia-sdk/x/bank/types"
	fswaptypes "github.com/Finschia/finschia-sdk/x/fswap/types"
)

const jsonOutputFormat string = "JSON"

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	authority sdk.AccAddress
	toDenom   banktypes.Metadata
	dummySwap fswaptypes.Swap
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	genesisState := s.cfg.GenesisState
	var bankGenesis banktypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[banktypes.ModuleName], &bankGenesis))
	s.toDenom = banktypes.Metadata{
		Name:        "Dummy Coin",
		Symbol:      "Dummy",
		Description: "The native token of Dummy chain",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "dummy",
				Exponent: 0,
				Aliases:  []string{},
			},
		},
		Base:    "dummy",
		Display: "dummy",
	}
	bankGenesis.DenomMetadata = []banktypes.Metadata{s.toDenom}
	bankDataBz, err := s.cfg.Codec.MarshalJSON(&bankGenesis)
	s.Require().NoError(err)
	genesisState[banktypes.ModuleName] = bankDataBz

	bondDenom := s.cfg.BondDenom
	toDenom := s.toDenom.Base

	var fswapData fswaptypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[fswaptypes.ModuleName], &fswapData))
	fswapData.SwapStats.SwapCount = 1
	s.dummySwap = fswaptypes.Swap{
		FromDenom:           bondDenom,
		ToDenom:             toDenom,
		AmountCapForToDenom: sdk.NewInt(12340000000000),
		SwapRate:            sdk.MustNewDecFromStr("1234"),
	}

	fswapData.Swaps = []fswaptypes.Swap{s.dummySwap}
	fswapData.Swappeds = []fswaptypes.Swapped{
		{
			FromCoinAmount: sdk.Coin{
				Denom:  bondDenom,
				Amount: sdk.ZeroInt(),
			},
			ToCoinAmount: sdk.Coin{
				Denom:  toDenom,
				Amount: sdk.ZeroInt(),
			},
		},
	}

	fswapDataBz, err := s.cfg.Codec.MarshalJSON(&fswapData)
	s.Require().NoError(err)
	genesisState[fswaptypes.ModuleName] = fswapDataBz
	s.cfg.GenesisState = genesisState

	s.network = network.New(s.T(), s.cfg)
	s.authority = fswaptypes.DefaultAuthority()

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}
