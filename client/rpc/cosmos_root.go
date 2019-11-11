package rpc

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/rpc"
)

// Register REST endpoints
func RegisterRPCRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/node_info", rpc.NodeInfoRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/syncing", rpc.NodeSyncingRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/blocks/latest", LatestBlockRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/blocks/{height}", BlockRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/validatorsets/latest", rpc.LatestValidatorSetRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/validatorsets/{height}", rpc.ValidatorSetRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/genesis/genutil/gentxs", QueryGenesisTxRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/genesis/app_state/accounts", QueryGenesisAccountRequestHandlerFn(cliCtx)).Methods("GET")
}
