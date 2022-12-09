package rest

import (
	"github.com/gorilla/mux"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/rest"
)

func RegisterHandlers(clientCtx client.Context, rtr *mux.Router) {
	r := rest.WithHTTPDeprecationHeaders(rtr)

	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)
}
