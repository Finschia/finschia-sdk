package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/line/link/x/bank/internal/types"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// Query the balance of a single denom
	r.HandleFunc("/bank/balances/{address}/{denom}", QueryBalanceOfDenomRequestHandlerFn(cliCtx)).Methods("GET")
}

// query accountREST Handler
func QueryBalanceOfDenomRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		bech32addr := vars["address"]
		denom := mux.Vars(r)["denom"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryBalanceOfParams(addr, denom)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/bank/balance_of", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		// the query will return empty if there is no data for this account
		if len(res) == 0 {
			rest.PostProcessResponse(w, cliCtx, 0)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
