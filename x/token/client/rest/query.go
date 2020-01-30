package rest

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/line/link/x/token/client/internal/types"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/line/link/client"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func RegisterRoutes(cliCtx client.CLIContext, r *mux.Router) {
	r.HandleFunc("/token/tokens/{symbol}/supply", QuerySupplyRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/tokens/{symbol}", QueryTokenRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/tokens", QueryTokensRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/collections/{symbol}/tokens/{token_id}/supply", QueryCollectionTokenSupplyRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/collections/{symbol}/tokens/{token_type}/count", QueryCountRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/collections/{symbol}/tokens/{token_id}", QueryCollectionTokenRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/collections/{symbol}/tokens", QueryCollectionTokensRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/collections/{symbol}", QueryCollectionRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/collections", QuerCollectionsRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/permissions/{address}", QueryPermRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/parent/{symbol}/{token_id}", QueryParentRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/root/{symbol}/{token_id}", QueryRootRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/children/{symbol}/{token_id}", QueryChildrenRequestHandlerFn(cliCtx)).Methods("GET")
}
func QueryTokenRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		token, height, err := retriever.GetToken(cliCtx, symbol, "")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, token)
	}
}

func QueryCollectionTokenRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		symbol := vars["symbol"]
		tokenID := vars["token_id"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		token, height, err := retriever.GetToken(cliCtx, symbol, tokenID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, token)
	}
}

func QueryTokensRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		tokens, height, err := retriever.GetTokens(cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, tokens)
	}
}

func QueryCollectionTokensRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		collection, height, err := retriever.GetCollection(cliCtx, symbol)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, collection.Tokens)
	}
}
func QueryCollectionRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		collection, height, err := retriever.GetCollection(cliCtx, symbol)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, collection)
	}
}

func QuerCollectionsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		collectionGetter := clienttypes.NewRetriever(cliCtx)

		collections, height, err := collectionGetter.GetCollections(cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, collections)
	}
}

func QuerySupplyRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		supply, height, err := retriever.GetSupply(cliCtx, symbol, "")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, supply)
	}
}
func QueryCollectionTokenSupplyRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		symbol := vars["symbol"]
		tokenID := vars["token_id"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		supply, height, err := retriever.GetSupply(cliCtx, symbol, tokenID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, supply)
	}
}

func QueryCountRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		symbol := vars["symbol"]
		tokenID := vars["token_type"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		nftcount, height, err := retriever.GetCollectionNFTCount(cliCtx, symbol, tokenID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, nftcount)
	}
}

func QueryPermRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		addr, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("address cannot parsed: %s", err))
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		nftcount, height, err := retriever.GetAccountPermission(cliCtx, addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, nftcount)
	}
}

func QueryParentRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		symbol := vars["symbol"]
		tokenID := vars["token_id"]

		if len(symbol) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "symbol absent")
			return
		}

		if len(tokenID) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "token_id absent")
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

		token, height, err := tokenGetter.GetParent(cliCtx, symbol, tokenID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, token)
	}
}

func QueryRootRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		symbol := vars["symbol"]
		tokenID := vars["token_id"]

		if len(symbol) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "symbol absent")
			return
		}

		if len(tokenID) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "token_id absent")
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

		token, height, err := tokenGetter.GetRoot(cliCtx, symbol, tokenID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, token)
	}
}

func QueryChildrenRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		symbol := vars["symbol"]
		tokenID := vars["token_id"]

		if len(symbol) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "symbol absent")
			return
		}

		if len(tokenID) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "token_id absent")
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

		tokens, height, err := tokenGetter.GetChildren(cliCtx, symbol, tokenID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, tokens)
	}
}
