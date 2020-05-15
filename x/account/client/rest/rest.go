package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	cosmosrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
)

// RegisterTxRoutes registers all transaction routes on the provided router.
func RegisterTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/txs/{hash}", QueryTxRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/txs", QueryTxsRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/txs", cosmosrest.BroadcastTxRequest(cliCtx)).Methods("POST")
	r.HandleFunc("/txs/encode", cosmosrest.EncodeTxRequestHandlerFn(cliCtx)).Methods("POST")
}
