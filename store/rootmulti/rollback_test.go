package rootmulti_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"

	abci "github.com/line/ostracon/abci/types"
	"github.com/line/ostracon/libs/log"
	ocproto "github.com/line/ostracon/proto/ostracon/types"

	"github.com/line/lbm-sdk/simapp"
)

func setup(withGenesis bool, invCheckPeriod uint, db dbm.DB) (*simapp.SimApp, simapp.GenesisState) {
	encCdc := simapp.MakeTestEncodingConfig()
	app := simapp.NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, simapp.DefaultNodeHome, invCheckPeriod, encCdc, simapp.EmptyAppOptions{}, nil)
	if withGenesis {
		return app, simapp.NewDefaultGenesisState(encCdc.Marshaler)
	}
	return app, simapp.GenesisState{}
}

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func SetupWithDB(isCheckTx bool, db dbm.DB) *simapp.SimApp {
	app, genesisState := setup(!isCheckTx, 5, db)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

func TestRollback(t *testing.T) {
	db := dbm.NewMemDB()
	app := SetupWithDB(false, db)
	app.Commit()
	ver0 := app.LastBlockHeight()
	// commit 10 blocks
	for i := int64(1); i <= 10; i++ {
		header := ocproto.Header{
			Height:  ver0 + i,
			AppHash: app.LastCommitID().Hash,
		}
		app.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx := app.NewContext(false, header)
		store := ctx.KVStore(app.GetKey("bank"))
		store.Set([]byte("key"), []byte(fmt.Sprintf("value%d", i)))
		app.Commit()
	}

	require.Equal(t, ver0+10, app.LastBlockHeight())
	store := app.NewContext(true, ocproto.Header{}).KVStore(app.GetKey("bank"))
	require.Equal(t, []byte("value10"), store.Get([]byte("key")))

	// rollback 5 blocks
	target := ver0 + 5
	require.NoError(t, app.CommitMultiStore().RollbackToVersion(target))
	require.Equal(t, target, app.LastBlockHeight())

	// recreate app to have clean check state
	encCdc := simapp.MakeTestEncodingConfig()
	app = simapp.NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, simapp.DefaultNodeHome, 5, encCdc, simapp.EmptyAppOptions{}, nil)
	store = app.NewContext(true, ocproto.Header{}).KVStore(app.GetKey("bank"))
	require.Equal(t, []byte("value5"), store.Get([]byte("key")))

	// commit another 5 blocks with different values
	for i := int64(6); i <= 10; i++ {
		header := ocproto.Header{
			Height:  ver0 + i,
			AppHash: app.LastCommitID().Hash,
		}
		app.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx := app.NewContext(false, header)
		store := ctx.KVStore(app.GetKey("bank"))
		store.Set([]byte("key"), []byte(fmt.Sprintf("VALUE%d", i)))
		app.Commit()
	}

	require.Equal(t, ver0+10, app.LastBlockHeight())
	store = app.NewContext(true, ocproto.Header{}).KVStore(app.GetKey("bank"))
	require.Equal(t, []byte("VALUE10"), store.Get([]byte("key")))
}
