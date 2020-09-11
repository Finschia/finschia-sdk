package rpc

import (
	"github.com/gorilla/mux"
	"github.com/line/link-modules/client/rpc/link/genesis"
	"github.com/line/link-modules/client/rpc/link/mempool"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/rpc"
)

// Register REST endpoints
func RegisterRPCRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/node_info", rpc.NodeInfoRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/syncing", rpc.NodeSyncingRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/blocks/latest", rpc.LatestBlockRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/blocks/{height}", rpc.BlockRequestHandlerFn(cliCtx)).Methods("GET")
	// r.HandleFunc("/blocks_with_tx_results/{from_height}", block.WithTxResultRequestHandlerFn(block.NewBlockUtil(cliCtx))).Methods("GET")
	r.HandleFunc("/validatorsets/latest", rpc.LatestValidatorSetRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/validatorsets/{height}", rpc.ValidatorSetRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/genesis/genutil/gentxs", genesis.QueryGenesisTxRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/genesis/app_state/accounts", genesis.QueryGenesisAccountRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/num_unconfirmed_txs", mempool.NumUnconfirmedTxsRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/unconfirmed_txs", mempool.UnconfirmedTxsRequestHandlerFn(cliCtx)).Methods("GET")
}
