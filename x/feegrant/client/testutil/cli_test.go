//go:build norace
// +build norace

package testutil

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Finschia/finschia-rdk/testutil/network"
)

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = 3
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
