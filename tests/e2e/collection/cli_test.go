//go:build e2e
// +build e2e

package collection

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/stretchr/testify/suite"

	"github.com/Finschia/finschia-sdk/simapp"
)

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig(simapp.NewTestNetworkFixture)
	cfg.NumValidators = 1
	suite.Run(t, NewE2ETestSuite(cfg))
}
