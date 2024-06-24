package server

// DONTCOVER

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/Finschia/ostracon/abci/server"
	ostcmd "github.com/Finschia/ostracon/cmd/ostracon/commands"
	tmcfg "github.com/Finschia/ostracon/config"
	tmlog "github.com/Finschia/ostracon/libs/log"
	"github.com/Finschia/ostracon/node"
	"github.com/Finschia/ostracon/p2p"
	pvm "github.com/Finschia/ostracon/privval"
	"github.com/Finschia/ostracon/proxy"
	rpchttp "github.com/Finschia/ostracon/rpc/client/http"
	"github.com/Finschia/ostracon/rpc/client/local"
	tmtypes "github.com/Finschia/ostracon/types"
	"github.com/hashicorp/go-metrics"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/server/api"
	serverconfig "github.com/Finschia/finschia-sdk/server/config"
	servergrpc "github.com/Finschia/finschia-sdk/server/grpc"
	"github.com/Finschia/finschia-sdk/server/types"
	"github.com/Finschia/finschia-sdk/store/cache"
	"github.com/Finschia/finschia-sdk/store/iavl"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	"github.com/Finschia/finschia-sdk/telemetry"
	"github.com/Finschia/finschia-sdk/version"
)

// Ostracon full-node start flags
const (
	flagWithOstracon        = "with-ostracon"
	flagAddress             = "address"
	flagTransport           = "transport"
	flagTraceStore          = "trace-store"
	flagCPUProfile          = "cpu-profile"
	FlagMinGasPrices        = "minimum-gas-prices"
	FlagHaltHeight          = "halt-height"
	FlagHaltTime            = "halt-time"
	FlagInterBlockCache     = "inter-block-cache"
	FlagInterBlockCacheSize = "inter-block-cache-size"
	FlagUnsafeSkipUpgrades  = "unsafe-skip-upgrades"
	FlagTrace               = "trace"
	FlagInvCheckPeriod      = "inv-check-period"
	FlagPrometheus          = "prometheus"
	FlagChanCheckTxSize     = "chan-check-tx-size"

	FlagPruning           = "pruning"
	FlagPruningKeepRecent = "pruning-keep-recent"
	FlagPruningKeepEvery  = "pruning-keep-every"
	FlagPruningInterval   = "pruning-interval"
	FlagIndexEvents       = "index-events"
	FlagMinRetainBlocks   = "min-retain-blocks"
	FlagIAVLCacheSize     = "iavl-cache-size"
	FlagIAVLFastNode      = "iavl-disable-fastnode"

	// state sync-related flags
	FlagStateSyncSnapshotInterval   = "state-sync.snapshot-interval"
	FlagStateSyncSnapshotKeepRecent = "state-sync.snapshot-keep-recent"

	// gRPC-related flags
	flagGRPCOnly    = "grpc-only"
	flagGRPCEnable  = "grpc.enable"
	flagGRPCAddress = "grpc.address"
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
everything: all saved states will be deleted, storing only the current and previous state; pruning at 10 block intervals
custom: allow pruning options to be manually specified through 'pruning-keep-recent', 'pruning-keep-every', and 'pruning-interval'

Node halting configurations exist in the form of two flags: '--halt-height' and '--halt-time'. During
the ABCI Commit phase, the node will check if the current block height is greater than or equal to
the halt-height or if the current block time is greater than or equal to the halt-time. If so, the
node will attempt to gracefully shutdown and the block will not be committed. In addition, the node
will not be able to commit subsequent blocks.

For profiling and benchmarking purposes, CPU profiling can be enabled via the '--cpu-profile' flag
which accepts a path for the resulting pprof file.

The node may be started in a 'query only' mode where only the gRPC and JSON HTTP
API services are enabled via the 'grpc-only' flag. In this mode, Tendermint is
bypassed and can be used when legacy queries are needed after an on-chain upgrade
is performed. Note, when enabled, gRPC will also be automatically enabled.
`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := GetServerContextFromCmd(cmd)

			// Bind flags to the Context's Viper so the app construction can set
			// options accordingly.
			if err := serverCtx.Viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			_, err := GetPruningOptionsFromFlags(serverCtx.Viper)
			return err
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			svrCtx := GetServerContextFromCmd(cmd)
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			svrCfg, err := getAndValidateConfig(svrCtx)
			if err != nil {
				return err
			}

			tmetrics, err := startTelemetry(svrCfg)
			if err != nil {
				return err
			}

			emitServerInfoMetrics()

			db, err := openDB(svrCtx.Config.RootDir)
			if err != nil {
				return err
			}

			traceWriter, traceCleanupFn, err := SetupTraceWriter(svrCtx.Logger, svrCtx.Viper.GetString(flagTraceStore))
			if err != nil {
				return err
			}
			defer traceCleanupFn()

			app := appCreator(svrCtx.Logger, db, traceWriter, svrCtx.Viper)

			withTM, _ := cmd.Flags().GetBool(flagWithOstracon)
			if !withTM {
				svrCtx.Logger.Info("starting ABCI without Tendermint")

				return wrapCPUProfile(svrCtx, func() error {
					return startStandAlone(svrCtx, svrCfg, clientCtx, app, tmetrics)
				})
			}

			return wrapCPUProfile(svrCtx, func() error {
				return startInProcess(svrCtx, svrCfg, clientCtx, app, tmetrics)
			})
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
	cmd.Flags().Int(FlagInterBlockCacheSize, cache.DefaultCommitKVStoreCacheSize, "The maximum bytes size of the inter-block cache")
	cmd.Flags().Int(FlagIAVLCacheSize, iavl.DefaultIAVLCacheSize, "The maximum units size of the iavl node cache (1 unit is 128 bytes).")
	cmd.Flags().String(flagCPUProfile, "", "Enable CPU profiling and write to the provided file")
	cmd.Flags().Bool(FlagTrace, false, "Provide full stack traces for errors in ABCI Log")
	cmd.Flags().String(FlagPruning, storetypes.PruningOptionDefault, "Pruning strategy (default|nothing|everything|custom)")
	cmd.Flags().Uint64(FlagPruningKeepRecent, 0, "Number of recent heights to keep on disk (ignored if pruning is not 'custom')")
	cmd.Flags().Uint64(FlagPruningKeepEvery, 0, "Offset heights to keep on disk after 'keep-every' (ignored if pruning is not 'custom')")
	cmd.Flags().Uint64(FlagPruningInterval, 0, "Height interval at which pruned heights are removed from disk (ignored if pruning is not 'custom')")
	cmd.Flags().Uint(FlagInvCheckPeriod, 0, "Assert registered invariants every N blocks")
	cmd.Flags().Uint64(FlagMinRetainBlocks, 0, "Minimum block height offset during ABCI commit to prune Ostracon blocks")

	cmd.Flags().Bool(flagGRPCOnly, false, "Start the node in gRPC query only mode (no Tendermint process is started)")
	cmd.Flags().Bool(flagGRPCEnable, true, "Define if the gRPC server should be enabled")
	cmd.Flags().String(flagGRPCAddress, serverconfig.DefaultGRPCAddress, "the gRPC server address to listen on")

	cmd.Flags().Uint64(FlagStateSyncSnapshotInterval, 0, "State sync snapshot interval")
	cmd.Flags().Uint32(FlagStateSyncSnapshotKeepRecent, 2, "State sync snapshot to keep")

	cmd.Flags().Bool(FlagIAVLFastNode, true, "Enable fast node for IAVL tree")

	cmd.Flags().Bool(FlagPrometheus, false, "Enable prometheus metric for app")

	cmd.Flags().Uint(FlagChanCheckTxSize, serverconfig.DefaultChanCheckTxSize, "The size of the channel check tx")

	// add support for all Ostracon-specific command line options
	ostcmd.AddNodeFlags(cmd)
	return cmd
}

func startStandAlone(svrCtx *Context, svrCfg serverconfig.Config, clientCtx client.Context, app types.Application, tmetrics *telemetry.Metrics) error {
	svr, err := server.NewServer(svrCtx.Viper.GetString(flagAddress), svrCtx.Viper.GetString(flagTransport), app)
	if err != nil {
		return fmt.Errorf("error creating listener: %w", err)
	}

	svr.SetLogger(svrCtx.Logger.With("module", "abci-server"))

	g, ctx := getCtx(svrCtx, false)

	// Add the tx service to the gRPC router. We only need to register this
	// service if API or gRPC is enabled, and avoid doing so in the general
	// case, because it spawns a new local CometBFT RPC client.
	if svrCfg.API.Enable || svrCfg.GRPC.Enable {
		// create tendermint client
		// assumes the rpc listen address is where tendermint has its rpc server
		rpcclient, err := rpchttp.New(svrCtx.Config.RPC.ListenAddress, "/websocket")
		if err != nil {
			return err
		}
		// re-assign for making the client available below
		// do not use := to avoid shadowing clientCtx
		clientCtx = clientCtx.WithClient(rpcclient)

		// use the provided clientCtx to register the services
		app.RegisterTxService(clientCtx)
		app.RegisterTendermintService(clientCtx)
		if a, ok := app.(types.ApplicationQueryService); ok {
			a.RegisterNodeService(clientCtx)
		}
	}

	clientCtx, err = startGrpcServer(ctx, g, svrCfg.GRPC, clientCtx, svrCtx, app)
	if err != nil {
		return err
	}

	err = startAPIServer(ctx, g, clientCtx, svrCfg, svrCtx, app, svrCtx.Config.RootDir, tmetrics)
	if err != nil {
		return err
	}

	g.Go(func() error {
		if err := svr.Start(); err != nil {
			svrCtx.Logger.Error("failed to start out-of-process ABCI server", "err", err)
			return err
		}

		// Wait for the calling process to be canceled or close the provided context,
		// so we can gracefully stop the ABCI server.
		<-ctx.Done()
		svrCtx.Logger.Info("stopping the ABCI server...")
		return svr.Stop()
	})

	return g.Wait()
}

func startInProcess(svrCtx *Context, svrCfg serverconfig.Config, clientCtx client.Context, app types.Application,
	tmetrics *telemetry.Metrics,
) error {
	tmCfg := svrCtx.Config
	gRPCOnly := svrCtx.Viper.GetBool(flagGRPCOnly)

	g, ctx := getCtx(svrCtx, true)

	if gRPCOnly {
		// TODO: Generalize logic so that gRPC only is really in startStandAlone
		svrCtx.Logger.Info("starting node in gRPC only mode; Tendermint is disabled")
		svrCfg.GRPC.Enable = true
	} else {
		svrCtx.Logger.Info("starting node with ABCI Tendermint in-process")
		tmNode, cleanupFn, err := startTmNode(ctx, tmCfg, app, svrCtx)
		if err != nil {
			return err
		}
		defer cleanupFn()

		// Add the tx service to the gRPC router. We only need to register this
		// service if API or gRPC is enabled, and avoid doing so in the general
		// case, because it spawns a new local tendermint RPC client.
		if svrCfg.API.Enable || svrCfg.GRPC.Enable {
			// Re-assign for making the client available below do not use := to avoid
			// shadowing the clientCtx variable.
			clientCtx = clientCtx.WithClient(local.New(tmNode))

			app.RegisterTxService(clientCtx)
			app.RegisterTendermintService(clientCtx)

			if a, ok := app.(types.ApplicationQueryService); ok {
				a.RegisterNodeService(clientCtx)
			}
		}
	}

	clientCtx, err := startGrpcServer(ctx, g, svrCfg.GRPC, clientCtx, svrCtx, app)
	if err != nil {
		return err
	}

	err = startAPIServer(ctx, g, clientCtx, svrCfg, svrCtx, app, svrCtx.Config.RootDir, tmetrics)
	if err != nil {
		return err
	}

	// wait for signal capture and gracefully return
	// we are guaranteed to be waiting for the "ListenForQuitSignals" goroutine.
	return g.Wait()
}

func genPvFileOnlyWhenKmsAddressEmpty(cfg *tmcfg.Config) *pvm.FilePV {
	if len(strings.TrimSpace(cfg.PrivValidatorListenAddr)) == 0 {
		return pvm.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
	}
	return nil
}

func getAndValidateConfig(svrCtx *Context) (serverconfig.Config, error) {
	svrcfg, err := serverconfig.GetConfig(svrCtx.Viper)
	if err != nil {
		return svrcfg, err
	}

	if err := svrcfg.ValidateBasic(); err != nil {
		return svrcfg, err
	}
	return svrcfg, nil
}

// returns a function which returns the genesis doc from the genesis file.
func getGenDocProvider(cfg *tmcfg.Config) func() (*tmtypes.GenesisDoc, error) {
	return node.DefaultGenesisDocProviderFunc(cfg)
}

// SetupTraceWriter sets up the trace writer and returns a cleanup function.
func SetupTraceWriter(logger tmlog.Logger, traceWriterFile string) (traceWriter io.WriteCloser, cleanup func(), err error) {
	// clean up the traceWriter when the server is shutting down
	cleanup = func() {}

	traceWriter, err = openTraceWriter(traceWriterFile)
	if err != nil {
		return traceWriter, cleanup, err
	}

	// if flagTraceStore is not used then traceWriter is nil
	if traceWriter != nil {
		cleanup = func() {
			if err = traceWriter.Close(); err != nil {
				logger.Error("failed to close trace writer", "err", err)
			}
		}
	}

	return traceWriter, cleanup, nil
}

// TODO: Move nodeKey into being created within the function.
func startTmNode(
	_ context.Context,
	cfg *tmcfg.Config,
	app types.Application,
	svrCtx *Context,
) (tmNode *node.Node, cleanupFn func(), err error) {
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, cleanupFn, err
	}

	tmNode, err = node.NewNode(
		cfg,
		genPvFileOnlyWhenKmsAddressEmpty(cfg),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		getGenDocProvider(cfg),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		svrCtx.Logger,
	)
	if err != nil {
		return tmNode, cleanupFn, err
	}

	svrCtx.Logger.Debug("initialization: tmNode created")
	if err := tmNode.Start(); err != nil {
		return tmNode, cleanupFn, err
	}
	svrCtx.Logger.Debug("initialization: tmNode started")

	cleanupFn = func() {
		if tmNode != nil && tmNode.IsRunning() {
			_ = tmNode.Stop()
		}
	}

	return tmNode, cleanupFn, nil
}

func startGrpcServer(
	ctx context.Context,
	g *errgroup.Group,
	config serverconfig.GRPCConfig,
	clientCtx client.Context,
	svrCtx *Context,
	app types.Application,
) (client.Context, error) {
	if !config.Enable {
		// return grpcServer as nil if gRPC is disabled
		return clientCtx, nil
	}
	_, _, err := net.SplitHostPort(config.Address)
	if err != nil {
		return clientCtx, err
	}

	maxSendMsgSize := config.MaxSendMsgSize
	if maxSendMsgSize == 0 {
		maxSendMsgSize = serverconfig.DefaultGRPCMaxSendMsgSize
	}

	maxRecvMsgSize := config.MaxRecvMsgSize
	if maxRecvMsgSize == 0 {
		maxRecvMsgSize = serverconfig.DefaultGRPCMaxRecvMsgSize
	}

	// if gRPC is enabled, configure gRPC client for gRPC gateway
	grpcClient, err := grpc.NewClient(
		config.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.ForceCodec(codec.NewProtoCodec(clientCtx.InterfaceRegistry).GRPCCodec()),
			grpc.MaxCallRecvMsgSize(maxRecvMsgSize),
			grpc.MaxCallSendMsgSize(maxSendMsgSize),
		),
	)
	if err != nil {
		return clientCtx, err
	}

	clientCtx = clientCtx.WithGRPCClient(grpcClient)
	svrCtx.Logger.Debug("gRPC client assigned to client context", "target", config.Address)

	grpcSrv, err := servergrpc.NewGRPCServer(clientCtx, app, config)
	if err != nil {
		return clientCtx, err
	}

	// Start the gRPC server in a goroutine. Note, the provided ctx will ensure
	// that the server is gracefully shut down.
	g.Go(func() error {
		return servergrpc.StartGRPCServer(ctx, svrCtx.Logger.With("module", "grpc-server"), config, grpcSrv)
	})
	return clientCtx, nil
}

func startAPIServer(
	ctx context.Context,
	g *errgroup.Group,
	clientCtx client.Context,
	svrCfg serverconfig.Config,
	svrCtx *Context,
	app types.Application,
	home string,
	metrics *telemetry.Metrics,
) error {
	if !svrCfg.API.Enable {
		return nil
	}

	clientCtx = clientCtx.WithHomeDir(home)

	apiSrv := api.New(clientCtx, svrCtx.Logger.With("module", "api-server"))
	app.RegisterAPIRoutes(apiSrv, svrCfg.API)

	if svrCfg.Telemetry.Enabled {
		apiSrv.SetTelemetry(metrics)
	}

	g.Go(func() error {
		return apiSrv.Start(ctx, svrCfg)
	})
	return nil
}

func startTelemetry(cfg serverconfig.Config) (*telemetry.Metrics, error) {
	if !cfg.Telemetry.Enabled {
		return nil, nil
	}
	return telemetry.New(cfg.Telemetry)
}

// wrapCPUProfile starts CPU profiling, if enabled, and executes the provided
// callbackFn in a separate goroutine, then will wait for that callback to
// return.
//
// NOTE: We expect the caller to handle graceful shutdown and signal handling.
func wrapCPUProfile(svrCtx *Context, callbackFn func() error) error {
	if cpuProfile := svrCtx.Viper.GetString(flagCPUProfile); cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			return err
		}

		svrCtx.Logger.Info("starting CPU profiler", "profile", cpuProfile)

		if err := pprof.StartCPUProfile(f); err != nil {
			return err
		}

		defer func() {
			svrCtx.Logger.Info("stopping CPU profiler", "profile", cpuProfile)
			pprof.StopCPUProfile()

			if err := f.Close(); err != nil {
				svrCtx.Logger.Info("failed to close cpu-profile file", "profile", cpuProfile, "err", err.Error())
			}
		}()
	}

	errCh := make(chan error)
	go func() {
		errCh <- callbackFn()
	}()

	return <-errCh
}

// emitServerInfoMetrics emits server info related metrics using application telemetry.
func emitServerInfoMetrics() {
	var ls []metrics.Label

	versionInfo := version.NewInfo()
	if len(versionInfo.GoVersion) > 0 {
		ls = append(ls, telemetry.NewLabel("go", versionInfo.GoVersion))
	}
	if len(versionInfo.LbmSdkVersion) > 0 {
		ls = append(ls, telemetry.NewLabel("version", versionInfo.LbmSdkVersion))
	}

	if len(ls) == 0 {
		return
	}

	telemetry.SetGaugeWithLabels([]string{"server", "info"}, 1, ls)
}

func getCtx(svrCtx *Context, block bool) (*errgroup.Group, context.Context) {
	ctx, cancelFn := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)
	// listen for quit signals so the calling parent process can gracefully exit
	ListenForQuitSignals(g, block, cancelFn, svrCtx.Logger)
	return g, ctx
}
