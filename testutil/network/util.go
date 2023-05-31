package network

import (
	"context"
	"encoding/json"
	"github.com/Finschia/finschia-sdk/compat"
	"github.com/tendermint/tendermint/privval"
	"path/filepath"
	"time"

	"github.com/Finschia/finschia-sdk/server/api"
	servergrpc "github.com/Finschia/finschia-sdk/server/grpc"
	srvtypes "github.com/Finschia/finschia-sdk/server/types"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	banktypes "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/genutil"
	genutiltypes "github.com/Finschia/finschia-sdk/x/genutil/types"
	ostos "github.com/Finschia/ostracon/libs/os"
	"github.com/Finschia/ostracon/node"
	"github.com/Finschia/ostracon/types"
	osttime "github.com/Finschia/ostracon/types/time"
	"github.com/tendermint/tendermint/p2p"

	rollconf "github.com/Finschia/ramus/config"
	rollconv "github.com/Finschia/ramus/conv"
	rollnode "github.com/Finschia/ramus/node"
	rollrpc "github.com/Finschia/ramus/rpc"
)

func startInProcess(cfg Config, val *Validator) error {
	logger := val.Ctx.Logger
	tmCfg := val.Ctx.Config
	tmCfg.Instrumentation.Prometheus = false

	if err := val.AppConfig.ValidateBasic(); err != nil {
		return err
	}

	nodeKey, err := p2p.LoadOrGenNodeKey(tmCfg.NodeKeyFile())
	if err != nil {
		return err
	}
	pval := privval.LoadOrGenFilePV(tmCfg.PrivValidatorKeyFile(), tmCfg.PrivValidatorStateFile())
	// keys in Rollkit format
	p2pKey, err := rollconv.GetNodeKey(nodeKey)
	if err != nil {
		return err
	}
	signingKey, err := rollconv.GetNodeKey(&p2p.NodeKey{PrivKey: pval.Key.PrivKey})
	if err != nil {
		return err
	}

	app := cfg.AppConstructor(*val)
	genDocProvider := node.DefaultGenesisDocProviderFunc(tmCfg)
	genDoc, err := genDocProvider()
	if err != nil {
		return err
	}

	nodeConfig := rollconf.NodeConfig{}
	err = nodeConfig.GetViperConfig(val.Ctx.Viper)
	nodeConfig.Aggregator = true
	nodeConfig.DALayer = "mock"
	if err != nil {
		return err
	}
	rollconv.GetNodeConfig(&nodeConfig, tmCfg)
	err = rollconv.TranslateAddresses(&nodeConfig)
	if err != nil {
		return err
	}
	val.tmNode, err = rollnode.NewNode(
		context.Background(),
		nodeConfig,
		p2pKey,
		signingKey,
		compat.NewTMClientCreator(app),
		compat.NewTMGenesisDoc(genDoc),
		compat.NewTMLogger(logger.With("module", val.Moniker)),
	)
	if err != nil {
		return err
	}

	if err := val.tmNode.Start(); err != nil {
		return err
	}

	if val.RPCAddress != "" {
		server := rollrpc.NewServer(val.tmNode, compat.NewTMRPCConfig(tmCfg.RPC), compat.NewTMLogger(logger))
		err = server.Start()
		if err != nil {
			return err
		}
		val.RPCClient = server.Client()
	}

	// We'll need a RPC client if the validator exposes a gRPC or REST endpoint.
	if val.APIAddress != "" || val.AppConfig.GRPC.Enable {
		val.ClientCtx = val.ClientCtx.
			WithClient(val.RPCClient)

		app.RegisterTxService(val.ClientCtx)
		app.RegisterTendermintService(val.ClientCtx)

		if a, ok := app.(srvtypes.ApplicationQueryService); ok {
			a.RegisterNodeService(val.ClientCtx)
		}
	}

	if val.APIAddress != "" {
		apiSrv := api.New(val.ClientCtx, logger.With("module", "api-server"))
		app.RegisterAPIRoutes(apiSrv, val.AppConfig.API)

		errCh := make(chan error)

		go func() {
			if err := apiSrv.Start(*val.AppConfig); err != nil {
				errCh <- err
			}
		}()

		select {
		case err := <-errCh:
			return err
		case <-time.After(srvtypes.ServerStartTime): // assume server started successfully
		}

		val.api = apiSrv
	}

	if val.AppConfig.GRPC.Enable {
		grpcSrv, err := servergrpc.StartGRPCServer(val.ClientCtx, app, val.AppConfig.GRPC.Address)
		if err != nil {
			return err
		}

		val.grpc = grpcSrv

		if val.AppConfig.GRPCWeb.Enable {
			val.grpcWeb, err = servergrpc.StartGRPCWeb(grpcSrv, *val.AppConfig)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func collectGenFiles(cfg Config, vals []*Validator, outputDir string) error {
	genTime := osttime.Now()

	for i := 0; i < cfg.NumValidators; i++ {
		tmCfg := vals[i].Ctx.Config

		nodeDir := filepath.Join(outputDir, vals[i].Moniker, "simd")
		gentxsDir := filepath.Join(outputDir, "gentxs")

		tmCfg.Moniker = vals[i].Moniker
		tmCfg.SetRoot(nodeDir)

		initCfg := genutiltypes.NewInitConfig(cfg.ChainID, gentxsDir, vals[i].NodeID, vals[i].PubKey)

		genFile := tmCfg.GenesisFile()
		genDoc, err := types.GenesisDocFromFile(genFile)
		if err != nil {
			return err
		}

		appState, err := genutil.GenAppStateFromConfig(cfg.Codec, cfg.TxConfig,
			tmCfg, initCfg, *genDoc, banktypes.GenesisBalancesIterator{})
		if err != nil {
			return err
		}

		// overwrite each validator's genesis file to have a canonical genesis time
		if err := genutil.ExportGenesisFileWithTime(genFile, cfg.ChainID, nil, appState, genTime); err != nil {
			return err
		}
	}

	return nil
}

func initGenFiles(cfg Config, genAccounts []authtypes.GenesisAccount, genBalances []banktypes.Balance, genFiles []string) error {
	// set the accounts in the genesis state
	var authGenState authtypes.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[authtypes.ModuleName], &authGenState)

	accounts, err := authtypes.PackAccounts(genAccounts)
	if err != nil {
		return err
	}

	authGenState.Accounts = append(authGenState.Accounts, accounts...)
	cfg.GenesisState[authtypes.ModuleName] = cfg.Codec.MustMarshalJSON(&authGenState)

	// set the balances in the genesis state
	var bankGenState banktypes.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[banktypes.ModuleName], &bankGenState)

	bankGenState.Balances = append(bankGenState.Balances, genBalances...)
	cfg.GenesisState[banktypes.ModuleName] = cfg.Codec.MustMarshalJSON(&bankGenState)

	appGenStateJSON, err := json.MarshalIndent(cfg.GenesisState, "", "  ")
	if err != nil {
		return err
	}

	genDoc := types.GenesisDoc{
		ChainID:    cfg.ChainID,
		AppState:   appGenStateJSON,
		Validators: nil,
	}

	// generate empty genesis files for each validator and save
	for i := 0; i < cfg.NumValidators; i++ {
		if err := genDoc.SaveAs(genFiles[i]); err != nil {
			return err
		}
	}

	return nil
}

func writeFile(name string, dir string, contents []byte) error {
	writePath := filepath.Join(dir) //nolint:gocritic
	file := filepath.Join(writePath, name)

	err := ostos.EnsureDir(writePath, 0755)
	if err != nil {
		return err
	}

	err = ostos.WriteFile(file, contents, 0644)
	if err != nil {
		return err
	}

	return nil
}
