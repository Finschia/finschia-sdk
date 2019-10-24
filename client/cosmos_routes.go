package client

import (
	"github.com/gorilla/mux"
)

// Register routes
func RegisterRoutes(cliCtx CLIContext, r *mux.Router) {
	RegisterRPCRoutes(cliCtx, r)
}
