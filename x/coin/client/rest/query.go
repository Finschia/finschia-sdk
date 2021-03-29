package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/line/lbm-sdk/v2/x/coin/internal/types"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// query accountREST Handler
func QueryBalancesRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		bech32addr := vars["address"]
		denom := vars["denom"] // "" if the key is not exists

		arrBech32addrs := strings.Split(bech32addr, ",")
		if len(arrBech32addrs) > 1 {
			addrs := make([]sdk.AccAddress, len(arrBech32addrs))

			for i, bech32addr := range arrBech32addrs {
				addr, err := sdk.AccAddressFromBech32(bech32addr)
				if err != nil {
					rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
					return
				}
				addrs[i] = addr
			}

			params := types.NewQueryBulkBalanceParams(addrs)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			queryBalancesRequest(w, r, cliCtx, bz, "bulk_balances")
		} else {
			addr, err := sdk.AccAddressFromBech32(bech32addr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			params := types.NewQueryBalanceParams(addr, denom)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			queryBalancesRequest(w, r, cliCtx, bz, "balances")
		}
	}
}

func queryBalancesRequest(w http.ResponseWriter, r *http.Request, cliCtx context.CLIContext, bz []byte, path string) {
	cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
	if !ok {
		return
	}

	res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/coin/%s", path), bz)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	cliCtx = cliCtx.WithHeight(height)

	// the query will return empty if there is no data for this account
	if len(res) == 0 {
		rest.PostProcessResponse(w, cliCtx, sdk.Coins{})
		return
	}

	rest.PostProcessResponse(w, cliCtx, res)
}
