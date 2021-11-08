package simapp

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/line/ostracon/libs/log"
	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/line/tm-db/v2/memdb"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/baseapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/module"
	abci "github.com/line/ostracon/abci/types"

	"github.com/line/lbm-sdk/x/auth"
	"github.com/line/lbm-sdk/x/auth/vesting"
	banktypes "github.com/line/lbm-sdk/x/bank/types"
	"github.com/line/lbm-sdk/x/capability"
	"github.com/line/lbm-sdk/x/crisis"
	"github.com/line/lbm-sdk/x/distribution"
	"github.com/line/lbm-sdk/x/evidence"
	"github.com/line/lbm-sdk/x/genutil"
	"github.com/line/lbm-sdk/x/gov"
	transfer "github.com/line/lbm-sdk/x/ibc/applications/transfer"
	ibc "github.com/line/lbm-sdk/x/ibc/core"
	"github.com/line/lbm-sdk/x/mint"
	"github.com/line/lbm-sdk/x/params"
	"github.com/line/lbm-sdk/x/slashing"
	"github.com/line/lbm-sdk/x/staking"
	"github.com/line/lbm-sdk/x/upgrade"
)

func TestSimAppExportAndBlockedAddrs(t *testing.T) {
	encCfg := MakeTestEncodingConfig()
	db := memdb.NewDB()
	app := NewSimApp(log.NewOCLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, DefaultNodeHome, 0, encCfg, EmptyAppOptions{})

	for acc := range maccPerms {
		require.True(
			t,
			app.BankKeeper.BlockedAddr(app.AccountKeeper.GetModuleAddress(acc)),
			"ensure that blocked addresses are properly set in bank keeper",
		)
	}

	genesisState := NewDefaultGenesisState(encCfg.Marshaler)
	stateBytes, err := json.MarshalIndent(genesisState, "", "  ")
	require.NoError(t, err)

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	app.Commit()

	// Making a new app object with the db, so that initchain hasn't been called
	app2 := NewSimApp(log.NewOCLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, DefaultNodeHome, 0, encCfg, EmptyAppOptions{})
	_, err = app2.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

func TestGetMaccPerms(t *testing.T) {
	dup := GetMaccPerms()
	require.Equal(t, maccPerms, dup, "duplicated module account permissions differed from actual module account permissions")
}

func TestRunMigrations(t *testing.T) {
	db := memdb.NewDB()
	encCfg := MakeTestEncodingConfig()
	logger := log.NewOCLogger(log.NewSyncWriter(os.Stdout))
	app := NewSimApp(logger, db, nil, true, map[int64]bool{}, DefaultNodeHome, 0, encCfg, EmptyAppOptions{})

	// Create a new baseapp and configurator for the purpose of this test.
	bApp := baseapp.NewBaseApp(appName, logger, db, encCfg.TxConfig.TxDecoder())
	bApp.SetCommitMultiStoreTracer(nil)
	bApp.SetInterfaceRegistry(encCfg.InterfaceRegistry)
	app.BaseApp = bApp
	app.configurator = module.NewConfigurator(app.MsgServiceRouter(), app.GRPCQueryRouter())

	// We register all modules on the Configurator, except x/bank. x/bank will
	// serve as the test subject on which we run the migration tests.
	//
	// The loop below is the same as calling `RegisterServices` on
	// ModuleManager, except that we skip x/bank.
	for _, module := range app.mm.Modules {
		if module.Name() == banktypes.ModuleName {
			continue
		}

		module.RegisterServices(app.configurator)
	}

	// Initialize the chain
	app.InitChain(abci.RequestInitChain{})
	app.Commit()

	testCases := []struct {
		name         string
		moduleName   string
		forVersion   uint64
		expRegErr    bool // errors while registering migration
		expRegErrMsg string
		expRunErr    bool // errors while running migration
		expRunErrMsg string
		expCalled    int
	}{
		{
			"cannot register migration for version 0",
			"bank", 0,
			true, "module migration versions should start at 1: invalid version", false, "", 0,
		},
		// {
		// 	"throws error on RunMigrations if no migration registered for bank",
		// 	"", 1,
		// 	false, "", true, "no migrations found for module bank: not found", 0,
		// },
		// {
		// 	"can register and run migration handler for x/bank",
		// 	"bank", 1,
		// 	false, "", false, "", 1,
		// },
		// {
		// 	"cannot register migration handler for same module & forVersion",
		// 	"bank", 1,
		// 	true, "another migration for module bank and version 1 already exists: internal logic error", false, "", 0,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error

			// Since it's very hard to test actual in-place store migrations in
			// tests (due to the difficulty of maintaing multiple versions of a
			// module), we're just testing here that the migration logic is
			// called.
			called := 0

			if tc.moduleName != "" {
				// Register migration for module from version `forVersion` to `forVersion+1`.
				err = app.configurator.RegisterMigration(tc.moduleName, tc.forVersion, func(sdk.Context) error {
					called++

					return nil
				})

				if tc.expRegErr {
					require.EqualError(t, err, tc.expRegErrMsg)

					return
				}
			}
			require.NoError(t, err)

			// Run migrations only for bank. That's why we put the initial
			// version for bank as 1, and for all other modules, we put as
			// their latest ConsensusVersion.
			err = app.RunMigrations(
				app.NewContext(true, ocproto.Header{Height: app.LastBlockHeight()}),
				module.MigrationMap{
					"bank":         1,
					"auth":         auth.AppModule{}.ConsensusVersion(),
					"staking":      staking.AppModule{}.ConsensusVersion(),
					"mint":         mint.AppModule{}.ConsensusVersion(),
					"distribution": distribution.AppModule{}.ConsensusVersion(),
					"slashing":     slashing.AppModule{}.ConsensusVersion(),
					"gov":          gov.AppModule{}.ConsensusVersion(),
					"params":       params.AppModule{}.ConsensusVersion(),
					"ibc":          ibc.AppModule{}.ConsensusVersion(),
					"upgrade":      upgrade.AppModule{}.ConsensusVersion(),
					"vesting":      vesting.AppModule{}.ConsensusVersion(),
					"transfer":     transfer.AppModule{}.ConsensusVersion(),
					"evidence":     evidence.AppModule{}.ConsensusVersion(),
					"crisis":       crisis.AppModule{}.ConsensusVersion(),
					"genutil":      genutil.AppModule{}.ConsensusVersion(),
					"capability":   capability.AppModule{}.ConsensusVersion(),
				},
			)
			if tc.expRunErr {
				require.EqualError(t, err, tc.expRunErrMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expCalled, called)
			}
		})
	}
}
