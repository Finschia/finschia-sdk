package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/line/lbm-sdk/v2/x/token/client/internal/types"
	"github.com/line/lbm-sdk/v2/x/token/internal/types"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/types/rest"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/token/{contract_id}/supply", QueryTotalRequestHandlerFn(cliCtx, types.QuerySupply)).Methods("GET")
	r.HandleFunc("/token/{contract_id}/mint", QueryTotalRequestHandlerFn(cliCtx, types.QueryMint)).Methods("GET")
	r.HandleFunc("/token/{contract_id}/burn", QueryTotalRequestHandlerFn(cliCtx, types.QueryBurn)).Methods("GET")
	r.HandleFunc("/token/{contract_id}/accounts/{address}/balance", QueryBalanceRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/{contract_id}/accounts/{address}/permissions", QueryPermRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/{contract_id}/token", QueryTokenRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/{contract_id}/accounts/{address}/proxies/{approver}", QueryIsApprovedRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/token/{contract_id}/accounts/{address}/approvers", QueryApproversRequestHandlerFn(cliCtx)).Methods("GET")
}

func QueryTokenRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		token, height, err := retriever.GetToken(cliCtx, contractID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, token)
	}
}

func QueryBalanceRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]
		addr, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		supply, height, err := retriever.GetAccountBalance(cliCtx, contractID, addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, supply)
	}
}

func QueryTotalRequestHandlerFn(cliCtx context.CLIContext, target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		retriever := clienttypes.NewRetriever(cliCtx)

		total, height, err := retriever.GetTotal(cliCtx, contractID, target)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, total)
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

func QueryApproversRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		contractID := vars["contract_id"]

		proxy, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("proxy[%s] cannot parsed: %s", proxy.String(), err))
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
