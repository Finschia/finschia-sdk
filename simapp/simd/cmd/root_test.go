package cmd

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"

	"github.com/Finschia/ostracon/libs/log"
	octypes "github.com/Finschia/ostracon/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Finschia/finschia-rdk/client/flags"
	"github.com/Finschia/finschia-rdk/server"
	"github.com/Finschia/finschia-rdk/simapp"
	"github.com/Finschia/finschia-rdk/store/types"
)

func TestNewApp(t *testing.T) {
	encodingConfig := simapp.MakeTestEncodingConfig()
	a := appCreator{encodingConfig}
	db := dbm.NewMemDB()
	tempDir := t.TempDir()
	ctx := server.NewDefaultContext()
	ctx.Viper.Set(flags.FlagHome, tempDir)
	ctx.Viper.Set(server.FlagPruning, types.PruningOptionNothing)
	app := a.newApp(log.NewOCLogger(log.NewSyncWriter(os.Stdout)), db, nil, ctx.Viper)
	require.NotNil(t, app)
}

func TestAppExport(t *testing.T) {
	encodingConfig := simapp.MakeTestEncodingConfig()
	logger := log.NewOCLogger(log.NewSyncWriter(os.Stdout))
	a := appCreator{encodingConfig}
	db := dbm.NewMemDB()
	tempDir := t.TempDir()
	ctx := server.NewDefaultContext()
	ctx.Viper.Set(flags.FlagHome, tempDir)
	ctx.Viper.Set(server.FlagPruning, types.PruningOptionNothing)

	// create default genesis data and save to store
	app := simapp.NewSimApp(logger, db, nil, true, map[int64]bool{}, tempDir, 0, encodingConfig, simapp.EmptyAppOptions{})
	genesisState := simapp.NewDefaultGenesisState(encodingConfig.Marshaler)
	stateBytes, err := json.MarshalIndent(genesisState, "", "  ")
	require.NoError(t, err)
	genDoc := &octypes.GenesisDoc{}
	genDoc.ChainID = "theChainId"
	genDoc.Validators = nil
	genDoc.AppState = stateBytes
	app.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simapp.DefaultConsensusParams,
			AppStateBytes:   genDoc.AppState,
		},
	)
	app.Commit()

	tests := []struct {
		name      string
		height    int64
		expectErr bool
	}{
		{
			"height error",
			3,
			true,
		},
		{
			"valid export",
			-1,
			false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			exported, err := a.appExport(logger, db, nil, tc.height, false, []string{}, ctx.Viper)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.NotNil(t, exported)
		})
	}

}
