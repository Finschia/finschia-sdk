package rest

import (
	"github.com/gorilla/mux"
	"github.com/link-chain/link/client"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx client.CLIContext, r *mux.Router) {
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
