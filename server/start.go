package server

// DONTCOVER

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/line/ostracon/abci/server"
	ostcmd "github.com/line/ostracon/cmd/ostracon/commands"
	ostos "github.com/line/ostracon/libs/os"
	"github.com/line/ostracon/node"
	"github.com/line/ostracon/p2p"
	pvm "github.com/line/ostracon/privval"
	"github.com/line/ostracon/proxy"
	"github.com/line/ostracon/rpc/client/local"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/line/lbm-sdk/v2/client"
	"github.com/line/lbm-sdk/v2/client/flags"
	"github.com/line/lbm-sdk/v2/server/api"
	"github.com/line/lbm-sdk/v2/server/config"
	servergrpc "github.com/line/lbm-sdk/v2/server/grpc"
	"github.com/line/lbm-sdk/v2/server/types"
	storetypes "github.com/line/lbm-sdk/v2/store/types"
)

// Ostracon full-node start flags
const (
	flagWithOstracon       = "with-ostracon"
	flagAddress            = "address"
	flagTransport          = "transport"
	flagTraceStore         = "trace-store"
	flagCPUProfile         = "cpu-profile"
	FlagMinGasPrices       = "minimum-gas-prices"
	FlagHaltHeight         = "halt-height"
	FlagHaltTime           = "halt-time"
	FlagInterBlockCache    = "inter-block-cache"
	FlagUnsafeSkipUpgrades = "unsafe-skip-upgrades"
	FlagTrace              = "trace"
	FlagInvCheckPeriod     = "inv-check-period"

	FlagPruning           = "pruning"
	FlagPruningKeepRecent = "pruning-keep-recent"
	FlagPruningKeepEvery  = "pruning-keep-every"
	FlagPruningInterval   = "pruning-interval"
	FlagIndexEvents       = "index-events"
	FlagMinRetainBlocks   = "min-retain-blocks"
)

// GRPC-related flags.
const (
	flagGRPCEnable  = "grpc.enable"
	flagGRPCAddress = "grpc.address"
)

// State sync-related flags.
const (
	FlagStateSyncSnapshotInterval   = "state-sync.snapshot-interval"
	FlagStateSyncSnapshotKeepRecent = "state-sync.snapshot-keep-recent"
)

// StartCmd runs the service passed in, either stand-alone or in-process with
// Ostracon.
func StartCmd(appCreator types.AppCreator, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Run the full node",
		Long: `Run the full node application with Ostracon in or out of process. By
default, the application will run with Ostracon in process.

Pruning options can be provided via the '--pruning' flag or alternatively with '--pruning-keep-recent',
'pruning-keep-every', and 'pruning-interval' together.

For '--pruning' the options are as follows:

default: the last 100 states are kept in addition to every 500th state; pruning at 10 block intervals
nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
everything: all saved states will be deleted, storing only the current state; pruning at 10 block intervals
custom: allow pruning options to be manually specified through 'pruning-keep-recent', 'pruning-keep-every', and 'pruning-interval'

Node halting configurations exist in the form of two flags: '--halt-height' and '--halt-time'. During
the ABCI Commit phase, the node will check if the current block height is greater than or equal to
the halt-height or if the current block time is greater than or equal to the halt-time. If so, the
node will attempt to gracefully shutdown and the block will not be committed. In addition, the node
will not be able to commit subsequent blocks.

For profiling and benchmarking purposes, CPU profiling can be enabled via the '--cpu-profile' flag
which accepts a path for the resulting pprof file.
`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := GetServerContextFromCmd(cmd)

			// Bind flags to the Context's Viper so the app construction can set
			// options accordingly.
			serverCtx.Viper.BindPFlags(cmd.Flags())

			_, err := GetPruningOptionsFromFlags(serverCtx.Viper)
			return err
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := GetServerContextFromCmd(cmd)
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			withOST, _ := cmd.Flags().GetBool(flagWithOstracon)
			if !withOST {
				serverCtx.Logger.Info("starting ABCI without Ostracon")
				return startStandAlone(serverCtx, appCreator)
			}

			serverCtx.Logger.Info("starting ABCI with Ostracon")

			// amino is needed here for backwards compatibility of REST routes
			err = startInProcess(serverCtx, clientCtx, appCreator)
			errCode, ok := err.(ErrorCode)
			if !ok {
				return err
			}

			serverCtx.Logger.Debug(fmt.Sprintf("received quit signal: %d", errCode.Code))
			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().Bool(flagWithOstracon, true, "Run abci app embedded in-process with ostracon")
	cmd.Flags().String(flagAddress, "tcp://0.0.0.0:26658", "Listen address")
	cmd.Flags().String(flagTransport, "socket", "Transport protocol: socket, grpc")
	cmd.Flags().String(flagTraceStore, "", "Enable KVStore tracing to an output file")
	cmd.Flags().String(FlagMinGasPrices, "", "Minimum gas prices to accept for transactions; Any fee in a tx must meet this minimum (e.g. 0.01photino;0.0001stake)")
	cmd.Flags().IntSlice(FlagUnsafeSkipUpgrades, []int{}, "Skip a set of upgrade heights to continue the old binary")
	cmd.Flags().Uint64(FlagHaltHeight, 0, "Block height at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().Uint64(FlagHaltTime, 0, "Minimum block time (in Unix seconds) at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().Bool(FlagInterBlockCache, true, "Enable inter-block caching")
	cmd.Flags().String(flagCPUProfile, "", "Enable CPU profiling and write to the provided file")
	cmd.Flags().Bool(FlagTrace, false, "Provide full stack traces for errors in ABCI Log")
	cmd.Flags().String(FlagPruning, storetypes.PruningOptionDefault, "Pruning strategy (default|nothing|everything|custom)")
	cmd.Flags().Uint64(FlagPruningKeepRecent, 0, "Number of recent heights to keep on disk (ignored if pruning is not 'custom')")
	cmd.Flags().Uint64(FlagPruningKeepEvery, 0, "Offset heights to keep on disk after 'keep-every' (ignored if pruning is not 'custom')")
	cmd.Flags().Uint64(FlagPruningInterval, 0, "Height interval at which pruned heights are removed from disk (ignored if pruning is not 'custom')")
	cmd.Flags().Uint(FlagInvCheckPeriod, 0, "Assert registered invariants every N blocks")
	cmd.Flags().Uint64(FlagMinRetainBlocks, 0, "Minimum block height offset during ABCI commit to prune Ostracon blocks")

	cmd.Flags().Bool(flagGRPCEnable, true, "Define if the gRPC server should be enabled")
	cmd.Flags().String(flagGRPCAddress, config.DefaultGRPCAddress, "the gRPC server address to listen on")

	cmd.Flags().Uint64(FlagStateSyncSnapshotInterval, 0, "State sync snapshot interval")
	cmd.Flags().Uint32(FlagStateSyncSnapshotKeepRecent, 2, "State sync snapshot to keep")

	// add support for all Ostracon-specific command line options
	ostcmd.AddNodeFlags(cmd)
	return cmd
}

func startStandAlone(ctx *Context, appCreator types.AppCreator) error {
	addr := ctx.Viper.GetString(flagAddress)
	transport := ctx.Viper.GetString(flagTransport)
	home := ctx.Viper.GetString(flags.FlagHome)

	db, err := openDB(home)
	if err != nil {
		return err
	}

	traceWriterFile := ctx.Viper.GetString(flagTraceStore)
	traceWriter, err := openTraceWriter(traceWriterFile)
	if err != nil {
		return err
	}

	app := appCreator(ctx.Logger, db, traceWriter, ctx.Viper)

	svr, err := server.NewServer(addr, transport, app)
	if err != nil {
		return fmt.Errorf("error creating listener: %v", err)
	}

	svr.SetLogger(ctx.Logger.With("module", "abci-server"))

	err = svr.Start()
	if err != nil {
		ostos.Exit(err.Error())
	}

	defer func() {
		if err = svr.Stop(); err != nil {
			ostos.Exit(err.Error())
		}
	}()

	// Wait for SIGINT or SIGTERM signal
	return WaitForQuitSignals()
}

// legacyAminoCdc is used for the legacy REST API
func startInProcess(ctx *Context, clientCtx client.Context, appCreator types.AppCreator) error {
	cfg := ctx.Config
	home := cfg.RootDir
	var cpuProfileCleanup func()

	if cpuProfile := ctx.Viper.GetString(flagCPUProfile); cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			return err
		}

		ctx.Logger.Info("starting CPU profiler", "profile", cpuProfile)
		if err := pprof.StartCPUProfile(f); err != nil {
			return err
		}

		cpuProfileCleanup = func() {
			ctx.Logger.Info("stopping CPU profiler", "profile", cpuProfile)
			pprof.StopCPUProfile()
			f.Close()
		}
	}

	traceWriterFile := ctx.Viper.GetString(flagTraceStore)
	db, err := openDB(home)
	if err != nil {
		return err
	}

	traceWriter, err := openTraceWriter(traceWriterFile)
	if err != nil {
		return err
	}

	app := appCreator(ctx.Logger, db, traceWriter, ctx.Viper)

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return err
	}

	genDocProvider := node.DefaultGenesisDocProviderFunc(cfg)
	tmNode, err := node.NewNode(
		cfg,
		pvm.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		genDocProvider,
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		ctx.Logger,
	)
	if err != nil {
		return err
	}

	ctx.Logger.Debug("initialization: tmNode created")
	if err := tmNode.Start(); err != nil {
		return err
	}
	ctx.Logger.Debug("initialization: tmNode started")

	config := config.GetConfig(ctx.Viper)

	// Add the tx service to the gRPC router. We only need to register this
	// service if API or gRPC is enabled, and avoid doing so in the general
	// case, because it spawns a new local ostracon RPC client.
	if config.API.Enable || config.GRPC.Enable {
		clientCtx = clientCtx.WithClient(local.New(tmNode))

		app.RegisterTxService(clientCtx)
		app.RegisterTendermintService(clientCtx)
	}

	var apiSrv *api.Server

	if config.API.Enable {
		genDoc, err := genDocProvider()
		if err != nil {
			return err
		}

		clientCtx := clientCtx.
			WithHomeDir(home).
			WithChainID(genDoc.ChainID)

		apiSrv = api.New(clientCtx, ctx.Logger.With("module", "api-server"))
		app.RegisterAPIRoutes(apiSrv, config.API)
		errCh := make(chan error)

		go func() {
			if err := apiSrv.Start(config); err != nil {
				errCh <- err
			}
		}()

		select {
		case err := <-errCh:
			return err
		case <-time.After(5 * time.Second): // assume server started successfully
		}
	}

	var grpcSrv *grpc.Server
	if config.GRPC.Enable {
		grpcSrv, err = servergrpc.StartGRPCServer(clientCtx, app, config.GRPC.Address)
		if err != nil {
			return err
		}
	}

	defer func() {
		if tmNode.IsRunning() {
			_ = tmNode.Stop()
		}

		if cpuProfileCleanup != nil {
			cpuProfileCleanup()
		}

		if apiSrv != nil {
			_ = apiSrv.Close()
		}

		if grpcSrv != nil {
			grpcSrv.Stop()
		}

		ctx.Logger.Info("exiting...")
	}()

	// Wait for SIGINT or SIGTERM signal
	return WaitForQuitSignals()
}
