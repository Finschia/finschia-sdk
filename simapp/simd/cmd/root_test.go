package cmd

import (
	"os"
	"testing"

	"github.com/line/ostracon/libs/log"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"

	"github.com/line/lbm-sdk/client/flags"
	"github.com/line/lbm-sdk/server"
	"github.com/line/lbm-sdk/simapp"
	"github.com/line/lbm-sdk/store/types"
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
	a := appCreator{encodingConfig}
	db := dbm.NewMemDB()
	tempDir := t.TempDir()
	ctx := server.NewDefaultContext()
	ctx.Viper.Set(flags.FlagHome, tempDir)
	ctx.Viper.Set(server.FlagPruning, types.PruningOptionNothing)
	exported, err := a.appExport(log.NewOCLogger(log.NewSyncWriter(os.Stdout)), db, nil, 3, false, []string{}, ctx.Viper)
	require.Error(t, err)
	require.NotNil(t, exported)
}
