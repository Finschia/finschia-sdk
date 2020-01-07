package block

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

// REST handler to get a block
func RequestHandlerFn(util *Util) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode, output, hasErr := processBlockFetchReq(r, util)
		if hasErr {
			rest.WriteErrorResponse(w, statusCode, output.(string))
		} else {
			rest.PostProcessResponseBare(w, util.lcliCtx.CosmosCliCtx(), output)
		}
	}
}

func processBlockFetchReq(r *http.Request, util *Util) (int, interface{}, bool) {
	vars := mux.Vars(r)

	height, err := strconv.ParseInt(vars["height"], 10, 64)
	if err != nil {
		return http.StatusBadRequest, fmt.Sprintf("couldn't parse block height. Assumed format is '/block/{height}'.  because of %s", err.Error()), true
	}

	chainHeight, err := util.LatestBlockHeight()
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("failed to parse chain height. because of %s", err.Error()), true
	}

	if height > chainHeight {
		return http.StatusNotFound, fmt.Sprintf("requested block height is bigger then the chain length"), true
	}

	output, err := GetBlock(util, &height, parseExtended(r))
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("temporary error. because of %s", err.Error()), true
	}

	return http.StatusOK, output, false
}

// REST handler to get the latest block
func LatestBlockRequestHandlerFn(util *Util) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output, err := GetBlock(util, nil, parseExtended(r))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponseBare(w, util.lcliCtx.CosmosCliCtx(), output)
	}
}

func parseExtended(r *http.Request) bool {
	return r.URL.Query().Get("extended") == "true"
}
