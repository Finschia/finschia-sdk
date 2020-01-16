package rest

import (
	clienttypes "github.com/line/link/x/token/client/internal/types"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/line/link/client"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func RegisterRoutes(cliCtx client.CLIContext, r *mux.Router) {
	r.HandleFunc("/token/tokens/{symbol}/{token_id}", QueryTokenRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/tokens/{symbol}", QueryTokenRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/tokens", QueryAllTokensRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/supply/{symbol}", QuerySupplysRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/collections/{symbol}", QueryCollectionRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/collections", QueryAllCollectionsRequestHandlerFn(cliCtx)).Methods("GET")
}

func QueryTokenRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		symbol := vars["symbol"]
		tokenID := ""
		if tokenIDVal, ok := vars["token_id"]; ok {
			tokenID = tokenIDVal
		}

		if len(symbol) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "symbol length > 2")
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		tokenGetter := clienttypes.NewTokenRetriever(cliCtx)

		if err := tokenGetter.EnsureExists(cliCtx, symbol, tokenID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		token, height, err := tokenGetter.GetTokenWithHeight(cliCtx, symbol, tokenID)
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

		tokenGetter := clienttypes.NewTokenRetriever(cliCtx)

		tokens, height, err := tokenGetter.GetAllTokensWithHeight(cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, tokens)
	}
}

func QueryCollectionRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
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

		collectionGetter := clienttypes.NewCollectionRetriever(cliCtx)

		if err := collectionGetter.EnsureExists(cliCtx, symbol); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		collection, height, err := collectionGetter.GetCollectionWithHeight(cliCtx, symbol)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, collection)
	}
}

func QueryAllCollectionsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		collectionGetter := clienttypes.NewCollectionRetriever(cliCtx)

		collections, height, err := collectionGetter.GetAllCollectionsWithHeight(cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, collections)
	}
}

func QuerySupplysRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		if len(symbol) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "symbol length should be greater than 2")
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		supplyGetter := clienttypes.NewSupplyRetriever(cliCtx)

		if err := supplyGetter.EnsureExists(cliCtx, symbol); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		supply, height, err := supplyGetter.GetSupplyWithHeight(cliCtx, symbol)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, supply)
	}
}
