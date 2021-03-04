package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	"github.com/line/lbm-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestLinkdExport(t *testing.T) {
	db := db.NewMemDB()
	gapp := NewLinkApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, 0)
	err := setGenesis(gapp)
	require.NoError(t, err)

	// Making a new app object with the db, so that initchain hasn't been called
	newGapp := NewLinkApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, 0)
	_, _, err = newGapp.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

// ensure that black listed addresses are properly set in bank keeper
func TestBlacklistedAddrs(t *testing.T) {
	db := db.NewMemDB()
	app := NewLinkApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, 0)

	for acc := range maccPerms {
		require.True(t, app.bankKeeper.BlacklistedAddr(app.supplyKeeper.GetModuleAddress(acc)))
	}
}

func setGenesis(gapp *LinkApp) error {
	genesisState := NewDefaultGenesisState()
	stateBytes, err := codec.MarshalJSONIndent(gapp.cdc, genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	gapp.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	gapp.Commit()
	return nil
}
