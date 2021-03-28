package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/line/lbm-sdk/v2/x/account/client/utils"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	genutilrest "github.com/cosmos/cosmos-sdk/x/genutil/client/rest"
)

// *****
// Original code: `github.com/cosmos/cosmos-sdk/x/auth/client/rest/query.go`
// Difference: referring import path of `utils`
// *****

// QueryTxsHandlerFn implements a REST handler that searches for transactions.
// Genesis transactions are returned if the height parameter is set to zero,
// otherwise the transactions are searched for by events.
func QueryTxsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			rest.WriteErrorResponse(
				w, http.StatusBadRequest,
				fmt.Sprintf("failed to parse query parameters: %s", err),
			)
			return
		}

		// if the height query param is set to zero, query for genesis transactions
		heightStr := r.FormValue("height")
		if heightStr != "" {
			if height, err := strconv.ParseInt(heightStr, 10, 64); err == nil && height == 0 {
				genutilrest.QueryGenesisTxs(cliCtx, w)
				return
			}
		}

		var (
			events      []string
			txs         []sdk.TxResponse
			page, limit int
		)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		if len(r.Form) == 0 {
			rest.PostProcessResponseBare(w, cliCtx, txs)
			return
		}

		events, page, limit, err = utils.ParseHTTPArgs(r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		searchResult, err := utils.QueryTxsByEvents(cliCtx, events, page, limit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, searchResult)
	}
}

// QueryTxRequestHandlerFn implements a REST handler that queries a transaction
// by hash in a committed block.
func QueryTxRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hashHexStr := vars["hash"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		output, err := utils.QueryTx(cliCtx, hashHexStr)
		if err != nil {
			if strings.Contains(err.Error(), hashHexStr) {
				rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if output.Empty() {
			rest.WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("no transaction found with hash %s", hashHexStr))
		}

		rest.PostProcessResponseBare(w, cliCtx, output)
	}
}

func QueryBlockWithTxsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		fromBlockHeight, err := strconv.ParseInt(vars["from_height"], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("couldn't parse block height. Assumed format is '/blocks_with_tx_results/{from_height}'. because of %s", err.Error()))
			return
		}
		fetchSize, err := strconv.ParseInt(r.URL.Query().Get("fetchsize"), 10, 8)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("couldn't parse fetchsize. because of %s", err.Error()))
			return
		}

		latestBlockHeight, err := utils.LatestBlockHeight(cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't get latestBlockHeight. because of %s", err.Error()))
			return
		}

		if fromBlockHeight >= latestBlockHeight {
			rest.WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("the block height does not exist. Requested: %d, Latest: %d", fromBlockHeight, latestBlockHeight))
			return
		}

		output, err := utils.BlockWithTxResponses(cliCtx, latestBlockHeight, fromBlockHeight, fetchSize)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't process request. because of %s", err.Error()))
			return
		}
		rest.PostProcessResponseBare(w, cliCtx, output)
	}
}
