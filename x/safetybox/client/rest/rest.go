package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

// RegisterRoutes register distribution REST routes.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc("/safetybox/{id}", SafetyBoxQueryHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/safetybox/{id}/role/{role}/{address}", SafetyBoxRoleQueryHandlerFn(cliCtx)).Methods("GET")
}
