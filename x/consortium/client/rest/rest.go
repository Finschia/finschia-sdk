package rest

import (
	"github.com/gorilla/mux"

	"github.com/line/lbm-sdk/client/rest"

	"github.com/line/lbm-sdk/client"
)

// RegisterRoutes registers REST routes for the consortium module under the path specified by routeName.
func RegisterRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := rest.WithHTTPDeprecationHeaders(rtr)
	registerQueryRoutes(clientCtx, r)
}
