package rest

import (
	"github.com/line/lbm-sdk/v2/client"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
	registerNewTxRoutes(cliCtx, r)
}
