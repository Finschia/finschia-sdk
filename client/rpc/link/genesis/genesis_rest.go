package genesis

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	genutilrest "github.com/cosmos/cosmos-sdk/x/genutil/client/rest"
	"github.com/line/link/x/account/client/utils"
)

// QueryGenesisTxRequestHandlerFn implements a REST handler to get the genesis
func QueryGenesisTxRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		genutilrest.QueryGenesisTxs(cliCtx, w)
	}
}

// QueryGenesisAccountRequestHandlerFn implements a REST handler to get the genesis accounts
func QueryGenesisAccountRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgs(r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		output, err := utils.QueryGenesisAccount(cliCtx, page, limit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, output)
	}
}
