package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/line/lbm-sdk/v2/baseapp"
	"github.com/line/lbm-sdk/v2/client"
	"github.com/line/lbm-sdk/v2/client/debug"
	"github.com/line/lbm-sdk/v2/client/flags"
	"github.com/line/lbm-sdk/v2/client/keys"
	"github.com/line/lbm-sdk/v2/client/rpc"
	"github.com/line/lbm-sdk/v2/server"
	servertypes "github.com/line/lbm-sdk/v2/server/types"
	"github.com/line/lbm-sdk/v2/snapshots"
	"github.com/line/lbm-sdk/v2/store"
	sdk "github.com/line/lbm-sdk/v2/types"
	authclient "github.com/line/lbm-sdk/v2/x/auth/client"
	authcmd "github.com/line/lbm-sdk/v2/x/auth/client/cli"
	authtypes "github.com/line/lbm-sdk/v2/x/auth/types"
	vestingcli "github.com/line/lbm-sdk/v2/x/auth/vesting/client/cli"
	banktypes "github.com/line/lbm-sdk/v2/x/bank/types"
	"github.com/line/lbm-sdk/v2/x/crisis"
	genutilcli "github.com/line/lbm-sdk/v2/x/genutil/client/cli"
	"github.com/line/lbm-sdk/v2/x/wasm"
	ostcli "github.com/line/ostracon/libs/cli"
	"github.com/line/ostracon/libs/log"
	dbm "github.com/line/tm-db/v2"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/line/lbm-sdk/v2/x/wasm/linkwasmd/app"
	"github.com/line/lbm-sdk/v2/x/wasm/linkwasmd/app/params"
	lbmtypes "github.com/line/lbm-sdk/v2/x/wasm/linkwasmd/types"
)

const (
	flagTestnet = "testnet"
)

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	encodingConfig := app.MakeEncodingConfig()

	initClientCtx := client.Context{}.
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(app.DefaultNodeHome)

	rootCmd := &cobra.Command{
		Use:   "linkwasmd",
		Short: "Wasm Daemon (server)",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(cmd)
		},
	}

	initRootCmd(rootCmd, encodingConfig)

	return rootCmd, encodingConfig
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {
	authclient.Codec = encodingConfig.Marshaler

	rootCmd.AddCommand(
		genutilcli.InitCmd(app.ModuleBasics, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		genutilcli.GenTxCmd(app.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		AddGenesisAccountCmd(app.DefaultNodeHome),
		ostcli.NewCompletionCmd(rootCmd, true),
		testnetCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}),
		debug.Cmd(),
	)

	server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, createLinkAppAndExport, addModuleInitFlags)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		queryCommand(),
		txCommand(),
		keys.Commands(app.DefaultNodeHome),
	)
}
func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	app.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		flags.LineBreak,
		vestingcli.GetTxCmd(),
	)

	app.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

// newApp is an AppCreator
func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	var cache sdk.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	snapshotDir := filepath.Join(cast.ToString(appOpts.Get(flags.FlagHome)), "data", "snapshots")
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		panic(err)
	}
	snapshotDB, err := sdk.NewLevelDB("metadata", snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}
	var emptyWasmOpts []wasm.Option
	return app.NewLinkApp(
		logger, db, traceStore, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		app.MakeEncodingConfig(), // Ideally, we would reuse the one created by NewRootCmd.
		wasm.EnableAllProposals,
		appOpts,
		emptyWasmOpts,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(cast.ToString(appOpts.Get(server.FlagMinGasPrices))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(server.FlagIndexEvents))),
		baseapp.SetSnapshotStore(snapshotStore),
		baseapp.SetSnapshotInterval(cast.ToUint64(appOpts.Get(server.FlagStateSyncSnapshotInterval))),
		baseapp.SetSnapshotKeepRecent(cast.ToUint32(appOpts.Get(server.FlagStateSyncSnapshotKeepRecent))),
	)
}

func createLinkAppAndExport(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailAllowedAddrs []string,
	appOpts servertypes.AppOptions) (servertypes.ExportedApp, error) {

	var linkApp *app.LinkApp
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}
	var emptyWasmOpts []wasm.Option
	if height != -1 {
		linkApp = app.NewLinkApp(logger, db, traceStore, false, map[int64]bool{}, homePath, uint(1), app.MakeEncodingConfig(), app.GetEnabledProposals(), appOpts, emptyWasmOpts)

		if err := linkApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		linkApp = app.NewLinkApp(logger, db, traceStore, true, map[int64]bool{}, homePath, uint(1), app.MakeEncodingConfig(), app.GetEnabledProposals(), appOpts, emptyWasmOpts)
	}

	return linkApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs)
}

func lbmPreRunE(cmd *cobra.Command) (err error) {
	err = server.InterceptConfigsPreRunHandler(cmd)

	testnet := viper.GetBool(flagTestnet)
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(lbmtypes.Bech32PrefixAcc(testnet), lbmtypes.Bech32PrefixAccPub(testnet))
	config.SetBech32PrefixForConsensusNode(lbmtypes.Bech32PrefixConsAddr(testnet), lbmtypes.Bech32PrefixConsPub(testnet))
	config.SetBech32PrefixForValidator(lbmtypes.Bech32PrefixValAddr(testnet), lbmtypes.Bech32PrefixValPub(testnet))
	config.SetCoinType(lbmtypes.CoinType)
	config.SetFullFundraiserPath(lbmtypes.FullFundraiserPath)
	config.Seal()

	ctx := server.GetServerContextFromCmd(cmd)
	if cmd.Name() == server.StartCmd(nil, "").Name() {
		var networkMode string
		if testnet {
			networkMode = "testnet"
		} else {
			networkMode = "mainnet"
		}
		ctx.Logger.Info(fmt.Sprintf("Network mode is %s", networkMode))
	}
	return
}
