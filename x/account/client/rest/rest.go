package rest

import (
	"github.com/gorilla/mux"
	"github.com/line/lbm-sdk/client/context"
	cosmosrest "github.com/line/lbm-sdk/x/auth/client/rest"
)

// RegisterTxRoutes registers all transaction routes on the provided router.
func RegisterTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/txs/{hash}", QueryTxRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/txs", QueryTxsRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/txs", cosmosrest.BroadcastTxRequest(cliCtx)).Methods("POST")
	r.HandleFunc("/txs/encode", cosmosrest.EncodeTxRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/txs/simulate", SimulateTxRequest(cliCtx)).Methods("POST")
	r.HandleFunc("/blocks_with_tx_results/{from_height}", QueryBlockWithTxsRequestHandlerFn(cliCtx)).Methods("GET")
}
