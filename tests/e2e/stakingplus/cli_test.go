//go:build e2e
// +build e2e

package stakingplus

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/cosmos/go-bip39"

	"github.com/Finschia/finschia-sdk/simapp"
)

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig(simapp.NewTestNetworkFixture)

	cfg.NumValidators = 1

	entropySeed, err := bip39.NewEntropy(256)
	require.NoError(t, err)
	mnemonic, err := bip39.NewMnemonic(entropySeed)
	require.NoError(t, err)
	cfg.Mnemonics = []string{mnemonic}

	suite.Run(t, NewE2ETestSuite(cfg))
}
