package rest

import (
	"net/http"

	"github.com/link-chain/link/client"
	"github.com/link-chain/link/x/token/types"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func QueryTokensRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		if len(symbol) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "symbol length > 2")
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		tokenGetter := types.NewTokenRetriever(cliCtx)

		if err := tokenGetter.EnsureExists(symbol); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		token, height, err := tokenGetter.GetTokenWithHeight(symbol)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, token)
	}
}

func QueryAllTokensRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		tokenGetter := types.NewTokenRetriever(cliCtx)

		tokens, height, err := tokenGetter.GetAllTokensWithHeight()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, tokens)
	}
}
