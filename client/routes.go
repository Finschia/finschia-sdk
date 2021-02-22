package client

import (
	"github.com/gorilla/mux"

	"github.com/line/lbm-sdk/client/context"
	"github.com/line/lbm-sdk/client/rpc"
)

// Register routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	rpc.RegisterRPCRoutes(cliCtx, r)
}
