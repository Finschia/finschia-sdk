package cmd

import (
	"os"
	"testing"

	"github.com/line/lbm-sdk/v2/client/flags"
	"github.com/line/lbm-sdk/v2/server"
	"github.com/line/lbm-sdk/v2/simapp"
	"github.com/line/lbm-sdk/v2/store/types"
	"github.com/line/ostracon/libs/log"
	"github.com/line/tm-db/v2/memdb"
	"github.com/stretchr/testify/require"
)

func TestNewApp(t *testing.T) {
	encodingConfig := simapp.MakeTestEncodingConfig()
	a := appCreator{encodingConfig}
	db := memdb.NewDB()
	tempDir := t.TempDir()
	ctx := server.NewDefaultContext()
	ctx.Viper.Set(flags.FlagHome, tempDir)
	ctx.Viper.Set(server.FlagPruning, types.PruningOptionNothing)
	app := a.newApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, ctx.Viper)
	require.NotNil(t, app)
}
