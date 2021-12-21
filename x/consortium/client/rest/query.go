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
		"/consortium/enabled", getParamsHandler(clientCtx),
	).Methods("GET")
	r.HandleFunc(
		"/consortium/validator", getValidatorAuthHandler(clientCtx),
	).Methods("GET")
	r.HandleFunc(
		"/consortium/validators", getValidatorAuthsHandler(clientCtx),
	).Methods("GET")
}

func getParamsHandler(clientCtx client.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		// ignore height for now
		res, _, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s", types.QuerierKey, types.QueryParams))
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

func getValidatorAuthsHandler(clientCtx client.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s", types.QuerierKey, types.QueryValidatorAuths))
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if len(res) == 0 {
			http.NotFound(w, r)
			return
		}

		var allowedOperators []string
		err = clientCtx.LegacyAmino.UnmarshalJSON(res, &allowedOperators)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, clientCtx, allowedOperators)
	}
}

func getValidatorAuthHandler(clientCtx client.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s", types.QuerierKey, types.QueryValidatorAuth))
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if len(res) == 0 {
			http.NotFound(w, r)
			return
		}

		var allowedOperators []string
		err = clientCtx.LegacyAmino.UnmarshalJSON(res, &allowedOperators)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, clientCtx, allowedOperators)
	}
}
