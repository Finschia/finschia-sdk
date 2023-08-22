package cmd

import (
	"context"
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	"github.com/Finschia/ostracon/libs/log"

	"github.com/Finschia/finschia-rdk/client"
	"github.com/Finschia/finschia-rdk/client/flags"
	"github.com/Finschia/finschia-rdk/server"
	"github.com/Finschia/finschia-rdk/simapp"
	banktypes "github.com/Finschia/finschia-rdk/x/bank/types"
	genutiltest "github.com/Finschia/finschia-rdk/x/genutil/client/testutil"
	genutiltypes "github.com/Finschia/finschia-rdk/x/genutil/types"
)

func Test_TestnetCmd(t *testing.T) {
	home := t.TempDir()
	encodingConfig := simapp.MakeTestEncodingConfig()
	logger := log.NewNopLogger()
	cfg, err := genutiltest.CreateDefaultTendermintConfig(home)
	require.NoError(t, err)

	err = genutiltest.ExecInitCmd(simapp.ModuleBasics, home, encodingConfig.Marshaler)
	require.NoError(t, err)

	serverCtx := server.NewContext(viper.New(), cfg, logger)
	clientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithHomeDir(home).
		WithTxConfig(encodingConfig.TxConfig)

	ctx := context.Background()
	ctx = context.WithValue(ctx, server.ServerContextKey, serverCtx)
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)
	cmd := testnetCmd(simapp.ModuleBasics, banktypes.GenesisBalancesIterator{})
	cmd.SetArgs([]string{fmt.Sprintf("--%s=test", flags.FlagKeyringBackend), fmt.Sprintf("--output-dir=%s", home)})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	genFile := cfg.GenesisFile()
	appState, _, err := genutiltypes.GenesisStateFromGenFile(genFile)
	require.NoError(t, err)

	bankGenState := banktypes.GetGenesisStateFromAppState(encodingConfig.Marshaler, appState)
	require.NotEmpty(t, bankGenState.Supply.String())
}
