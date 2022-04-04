package rest

import (
	"github.com/gorilla/mux"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/rest"
)

// RegisterRoutes registers REST routes for the upgrade module under the path specified by routeName.
func RegisterRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := rest.WithHTTPDeprecationHeaders(rtr)
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)
}
