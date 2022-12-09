package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	storetypes "github.com/line/lbm-sdk/store/types"
	sdk "github.com/line/lbm-sdk/types"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.True(t, cfg.GetMinGasPrices().IsZero())
}

func TestSetMinimumFees(t *testing.T) {
	cfg := DefaultConfig()
	cfg.SetMinGasPrices(sdk.DecCoins{sdk.NewInt64DecCoin("foo", 5)})
	require.Equal(t, "5.000000000000000000foo", cfg.MinGasPrices)
}

func TestValidateBasic(t *testing.T) {
	cfg := DefaultConfig()
	cfg.SetMinGasPrices(sdk.DecCoins{sdk.NewInt64DecCoin("foo", 5)})
	err := cfg.ValidateBasic()
	require.NoError(t, err)

	cfg.Pruning = storetypes.PruningOptionEverything
	cfg.StateSync.SnapshotInterval = 5
	err = cfg.ValidateBasic()
	require.Error(t, err)
}
