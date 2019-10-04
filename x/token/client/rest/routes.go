package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/token/tokens/{symbol}", QueryTokensRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/tokens/{symbol}/publish", PublishRequestHandlerFn(cliCtx)).Methods("POST")
}

/*
TODO
GET  /token/tokens
GET  /token/tokens/{symbol}
POST /token/tokens/{symbol}/publish
POST /token/tokens/{symbol}/transfer
POST /token/tokens/{symbol}/mint
POST /token/tokens/{symbol}/burn
*/
