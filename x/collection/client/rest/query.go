package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/line/link-modules/x/collection/client/internal/types"
	"github.com/line/link-modules/x/collection/internal/types"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/types/rest"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/collection/{contract_id}/fts/{token_id}/supply", QueryTokenTotalRequestHandlerFn(cliCtx, types.QuerySupply)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/fts/{token_id}/mint", QueryTokenTotalRequestHandlerFn(cliCtx, types.QueryMint)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/fts/{token_id}/burn", QueryTokenTotalRequestHandlerFn(cliCtx, types.QueryBurn)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/fts/{token_id}", QueryTokenRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/nfts/{token_id}/parent", QueryParentRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/nfts/{token_id}/root", QueryRootRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/nfts/{token_id}/children", QueryChildrenRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/nfts/{token_id}", QueryTokenRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/tokens", QueryTokensRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/tokentypes/{token_type}/tokens", QueryTokensWithTokenTypeRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/tokentypes/{token_type}/count", QueryCountRequestHandlerFn(cliCtx, types.QueryNFTCount)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/tokentypes/{token_type}/mint", QueryCountRequestHandlerFn(cliCtx, types.QueryNFTMint)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/tokentypes/{token_type}/burn", QueryCountRequestHandlerFn(cliCtx, types.QueryNFTBurn)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/tokentypes/{token_type}", QueryTokenTypeRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/tokentypes", QueryTokenTypesRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/accounts/{address}/permissions", QueryPermRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/accounts/{address}/proxies/{approver}", QueryIsApprovedRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/accounts/{address}/balances", QueryBalancesRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/accounts/{address}/approvers", QueryApproversRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/accounts/{address}/balances/{token_id}", QueryBalanceRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/collection/{contract_id}/collection", QueryCollectionRequestHandlerFn(cliCtx)).Methods("GET")
}

func QueryBalancesRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]
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

		coins, height, err := retriever.GetAccountBalances(cliCtx, contractID, addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, coins)
	}
}

func QueryBalanceRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]
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

		supply, height, err := retriever.GetAccountBalance(cliCtx, contractID, tokenID, addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, supply)
	}
}

func QueryTokenTypeRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]
		tokenTypeID := vars["token_type"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		tokenType, height, err := retriever.GetTokenType(cliCtx, contractID, tokenTypeID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, tokenType)
	}
}

func QueryTokenTypesRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		tokenTypes, height, err := retriever.GetTokenTypes(cliCtx, contractID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, tokenTypes)
	}
}

func QueryTokenRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]
		tokenID := vars["token_id"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		token, height, err := retriever.GetToken(cliCtx, contractID, tokenID)
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
		vars := mux.Vars(r)
		contractID := vars["contract_id"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		tokens, height, err := retriever.GetTokens(cliCtx, contractID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, tokens)
	}
}
func QueryCollectionRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		collection, height, err := retriever.GetCollection(cliCtx, contractID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, collection)
	}
}

func QueryTokenTotalRequestHandlerFn(cliCtx context.CLIContext, target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]
		tokenID := vars["token_id"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		supply, height, err := retriever.GetTotal(cliCtx, contractID, tokenID, target)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, supply)
	}
}

func QueryTokensWithTokenTypeRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]
		tokenType := vars["token_type"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)
		tokens, height, err := retriever.GetTokensWithTokenType(cliCtx, contractID, tokenType)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, tokens)
	}
}

func QueryCountRequestHandlerFn(cliCtx context.CLIContext, target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]
		tokenID := vars["token_type"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		nftcount, height, err := retriever.GetCollectionNFTCount(cliCtx, contractID, tokenID, target)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, nftcount)
	}
}

func QueryPermRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		addr, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("address cannot parsed: %s", err))
			return
		}
		contractID := vars["contract_id"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		pms, height, err := retriever.GetAccountPermission(cliCtx, contractID, addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, pms)
	}
}

// nolint:dupl
func QueryParentRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]
		tokenID := vars["token_id"]

		if len(contractID) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "contract_id absent")
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

		if err := tokenGetter.EnsureExists(cliCtx, contractID, tokenID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		token, height, err := tokenGetter.GetParent(cliCtx, contractID, tokenID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, token)
	}
}

// nolint:dupl
func QueryRootRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]
		tokenID := vars["token_id"]

		if len(contractID) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "contract_id absent")
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

		if err := tokenGetter.EnsureExists(cliCtx, contractID, tokenID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		token, height, err := tokenGetter.GetRoot(cliCtx, contractID, tokenID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, token)
	}
}

// nolint:dupl
func QueryChildrenRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]
		tokenID := vars["token_id"]

		if len(contractID) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "contract_id absent")
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

		if err := tokenGetter.EnsureExists(cliCtx, contractID, tokenID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tokens, height, err := tokenGetter.GetChildren(cliCtx, contractID, tokenID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, tokens)
	}
}

func QueryApproversRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]

		proxy, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("approver[%s] cannot parsed: %s", proxy.String(), err))
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		approvers, height, err := retriever.GetApprovers(cliCtx, contractID, proxy)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, approvers)
	}
}

func QueryIsApprovedRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)

		proxy, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("proxy[%s] cannot parsed: %s", proxy.String(), err))
			return
		}

		approver, err := sdk.AccAddressFromBech32(vars["approver"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("approver[%s] cannot parsed: %s", approver.String(), err))
			return
		}

		contractID := vars["contract_id"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		approved, height, err := retriever.IsApproved(cliCtx, contractID, proxy, approver)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, approved)
	}
}
