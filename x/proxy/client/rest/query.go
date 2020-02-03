package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/line/link/x/proxy/types"
)

func ProxyQueryAllowanceHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		proxyStr, onBehalfOfStr, denom := vars["proxy"], vars["on_behalf_of"], vars["denom"]
		if len(proxyStr) == 0 || len(onBehalfOfStr) == 0 || len(denom) == 0 {
			rest.WriteErrorResponse(
				w,
				http.StatusBadRequest,
				"All proxy, on_behalf_of addresses and denom is required",
			)
			return
		}

		proxy, err := sdk.AccAddressFromBech32(proxyStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		onBehalfOf, err := sdk.AccAddressFromBech32(onBehalfOfStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		proxyAllowanceGetter := types.NewProxyAllowanceRetriever(cliCtx)

		allowance, height, err := proxyAllowanceGetter.GetProxyAllowance(proxy, onBehalfOf, denom)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), allowance)
	}
}
