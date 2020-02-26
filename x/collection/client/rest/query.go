package rest

import (
	"fmt"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
	clienttypes "github.com/line/link/x/collection/client/internal/types"

	"github.com/gorilla/mux"
	"github.com/line/link/client"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func RegisterRoutes(cliCtx client.CLIContext, r *mux.Router) {
	r.HandleFunc("/collection/collections/{symbol}/tokens/{token_id}/supply", QueryTokenSupplyRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/collections/{symbol}/tokens/{token_type}/count", QueryCountRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/collections/{symbol}/tokens/{token_id}", QueryTokenRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/collections/{symbol}/tokens", QueryTokensRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/collections/{symbol}", QueryCollectionRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/collections", QuerCollectionsRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/parent/{symbol}/{token_id}", QueryParentRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/root/{symbol}/{token_id}", QueryRootRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/children/{symbol}/{token_id}", QueryChildrenRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/approved/{proxy}/{approver}/{symbol}", QueryIsApprovedRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/collections/{symbol}/balances/{address}/{token_id}", QueryBalanceRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/permissions/{address}", QueryPermRequestHandlerFn(cliCtx)).Methods("GET")
}

func QueryBalanceRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		symbol := vars["symbol"]
		tokenID := vars["token_id"]
		addr, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("addr[%s] cannot parsed: %s", vars["address"], err))
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		supply, height, err := retriever.GetAccountBalance(cliCtx, symbol, tokenID, addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, supply)
	}
}

func QueryTokenRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
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

func QueryTokensRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		tokens, height, err := retriever.GetTokens(cliCtx, symbol)
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

func QueryTokenSupplyRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
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

//nolint:dupl
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

		tokenGetter := clienttypes.NewRetriever(cliCtx)

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

//nolint:dupl
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

		tokenGetter := clienttypes.NewRetriever(cliCtx)

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

//nolint:dupl
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

		tokenGetter := clienttypes.NewRetriever(cliCtx)

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

func QueryIsApprovedRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)

		proxy, err := sdk.AccAddressFromBech32(vars["proxy"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("proxy[%s] cannot parsed: %s", proxy.String(), err))
			return
		}

		approver, err := sdk.AccAddressFromBech32(vars["approver"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("approver[%s] cannot parsed: %s", approver.String(), err))
			return
		}

		symbol := vars["symbol"]
		if err := types.ValidateSymbolUserDefined(symbol); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid symbol[%s]: %s", symbol, err))
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		approved, height, err := retriever.IsApproved(cliCtx, proxy, approver, symbol)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, approved)
	}
}
