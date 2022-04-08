package testutil

import (
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/testutil/network"
	"github.com/line/lbm-sdk/x/foundation"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	genesisState := s.cfg.GenesisState

	var foundationData foundation.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[foundation.ModuleName], &foundationData))

	// enable foundation
	params := &foundation.Params{
		Enabled: true,
	}
	foundationData.Params = params

	foundationDataBz, err := s.cfg.Codec.MarshalJSON(&foundationData)
	s.Require().NoError(err)
	genesisState[foundation.ModuleName] = foundationDataBz
	s.cfg.GenesisState = genesisState

	s.network = network.New(s.T(), s.cfg)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}
