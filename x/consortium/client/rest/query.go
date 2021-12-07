package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/types/rest"
	"github.com/line/lbm-sdk/x/consortium/types"
)

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		"/consortium/enabled", getEnabledHandler(clientCtx),
	).Methods("GET")
	r.HandleFunc(
		"/consortium/allowed_operator", getAllowedOperatorHandler(clientCtx),
	).Methods("GET")
	r.HandleFunc(
		"/consortium/allowed_operators", getAllowedOperatorsHandler(clientCtx),
	).Methods("GET")
}

func getEnabledHandler(clientCtx client.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		// ignore height for now
		res, _, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s", types.QuerierKey, types.QueryEnabled))
		if rest.CheckInternalServerError(w, err) {
			return
		}
		if len(res) == 0 {
			http.NotFound(w, request)
			return
		}

		var enabled bool
		err = clientCtx.LegacyAmino.UnmarshalJSON(res, &enabled)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, clientCtx, enabled)
	}
}

func getAllowedOperatorsHandler(clientCtx client.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s", types.QuerierKey, types.QueryAllowedValidators))
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if len(res) == 0 {
			http.NotFound(w, r)
			return
		}

		allowedOperators := []string{}
		err = clientCtx.LegacyAmino.UnmarshalJSON(res, &allowedOperators)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, clientCtx, allowedOperators)
	}
}

func getAllowedOperatorHandler(clientCtx client.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s", types.QuerierKey, types.QueryAllowedValidator))
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if len(res) == 0 {
			http.NotFound(w, r)
			return
		}

		allowedOperators := []string{}
		err = clientCtx.LegacyAmino.UnmarshalJSON(res, &allowedOperators)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, clientCtx, allowedOperators)
	}
}
