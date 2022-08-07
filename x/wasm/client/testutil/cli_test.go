package testutil

import (
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/line/lbm-sdk/testutil/network"
)

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = 1
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
