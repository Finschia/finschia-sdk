package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

// RegisterRoutes register distribution REST routes.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc("/proxy/allowance/{proxy}/{on_behalf_of}/{denom}", ProxyQueryAllowanceHandlerFn(cliCtx)).Methods("GET")
}
