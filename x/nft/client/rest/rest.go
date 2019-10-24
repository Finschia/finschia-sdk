package rest

import (
	"github.com/gorilla/mux"
	"github.com/link-chain/link/client"

	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterRoutes register distribution REST routes.
func RegisterRoutes(cliCtx client.CLIContext, r *mux.Router, cdc *codec.Codec, queryRoute string) {
	registerQueryRoutes(cliCtx, r, cdc, queryRoute)
	registerTxRoutes(cliCtx, r, cdc, queryRoute)
}
