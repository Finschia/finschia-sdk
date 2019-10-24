package rest

import (
	"github.com/gorilla/mux"
	"github.com/link-chain/link/client"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx client.CLIContext, r *mux.Router) {
	r.HandleFunc("/token/tokens/{symbol}", QueryTokensRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/tokens", QueryAllTokensRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/tokens/{symbol}/publish", PublishRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/token/tokens/{symbol}/mint", MintRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/token/tokens/{symbol}/burn", BurnRequestHandlerFn(cliCtx)).Methods("POST")
}
