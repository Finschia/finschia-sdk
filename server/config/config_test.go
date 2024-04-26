package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.True(t, cfg.GetMinGasPrices().IsZero())
}

func TestGetAndSetMinimumGas(t *testing.T) {
	cfg := DefaultConfig()

	input := sdk.DecCoins{sdk.NewInt64DecCoin("foo", 5)}
	cfg.SetMinGasPrices(input)
	require.Equal(t, "5.000000000000000000foo", cfg.MinGasPrices)
	require.EqualValues(t, cfg.GetMinGasPrices(), input)

	input = sdk.DecCoins{sdk.NewInt64DecCoin("bar", 1), sdk.NewInt64DecCoin("foo", 5)}
	cfg.SetMinGasPrices(input)
	require.Equal(t, "1.000000000000000000bar,5.000000000000000000foo", cfg.MinGasPrices)
	require.EqualValues(t, cfg.GetMinGasPrices(), input)
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

func TestSetGRPCMsgSize(t *testing.T) {
	cfg := DefaultConfig()
	require.Equal(t, DefaultGRPCMaxRecvMsgSize, cfg.GRPC.MaxRecvMsgSize)
	require.Equal(t, DefaultGRPCMaxSendMsgSize, cfg.GRPC.MaxSendMsgSize)
}
