package testutil

import (
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/testutil/network"
	"github.com/line/lbm-sdk/x/consortium"
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

	var consortiumData consortium.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[consortium.ModuleName], &consortiumData))

	// enable consortium
	params := &consortium.Params{
		Enabled: true,
	}
	consortiumData.Params = params

	consortiumDataBz, err := s.cfg.Codec.MarshalJSON(&consortiumData)
	s.Require().NoError(err)
	genesisState[consortium.ModuleName] = consortiumDataBz
	s.cfg.GenesisState = genesisState

	s.network = network.New(s.T(), s.cfg)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}
