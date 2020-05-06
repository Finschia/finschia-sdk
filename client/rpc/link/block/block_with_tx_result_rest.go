package block

import (
	"fmt"
	"net/http"
	"strconv"

	cdc "github.com/line/link/client/rpc/link/block/codec"

	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func WithTxResultRequestHandlerFn(util *Util) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode, output, hasError := processReq(r, util)
		if hasError {
			rest.WriteErrorResponse(w, statusCode, output.(string))
		} else {
			rest.PostProcessResponseBare(w, util.lcliCtx.CosmosCliCtx(), output.(*cdc.HasMoreResponseWrapper))
		}
	}
}

func processReq(r *http.Request, util *Util) (int, interface{}, bool) {
	vars := mux.Vars(r)

	height, err := strconv.ParseInt(vars["from_height"], 10, 64)
	if err != nil {
		return http.StatusBadRequest, fmt.Sprintf("couldn't parse block height. Assumed format is '/blocks_with_tx_results/{from_height}'. because of %s", err.Error()), true
	}
	fetchsize, err := strconv.ParseInt(r.URL.Query().Get("fetchsize"), 10, 8)
	if err != nil {
		return http.StatusBadRequest, fmt.Sprintf("couldn't parse fetchsize. because of %s",
			err.Error()), true
	}

	latestBlockHeight, err := util.LatestBlockHeight()
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("couldn't get latestBlockHeight. because of %s",
			err.Error()), true
	}

	if height > latestBlockHeight {
		return http.StatusNotFound, fmt.Sprintf("the block height does not exist. Requested: %d, Latest: %d", height, latestBlockHeight), true
	}

	fetchSizeInt8 := int8(fetchsize)

	output, err := util.fetchByBlockHeights(&latestBlockHeight, &height, &fetchSizeInt8)
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("couldn't process request. because of %s",
			err.Error()), true
	}
	return http.StatusOK, output, false
}
