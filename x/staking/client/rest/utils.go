package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/line/lbm-sdk/client"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/rest"
	authclient "github.com/line/lbm-sdk/x/auth/client"
	"github.com/line/lbm-sdk/x/staking/types"
)

// contains checks if the a given query contains one of the tx types
func contains(stringSlice []string, txType string) bool {
	for _, word := range stringSlice {
		if word == txType {
			return true
		}
	}

	return false
}

// queries staking txs
func queryTxs(clientCtx client.Context, action string, delegatorAddr string) (*sdk.SearchTxsResult, error) {
	page := 1
	limit := 100
	events := []string{
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, action),
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeySender, delegatorAddr),
	}

	return authclient.QueryTxsByEvents(clientCtx, events, page, limit, "")
}

func queryBonds(clientCtx client.Context, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32delegator := vars["delegatorAddr"]
		bech32validator := vars["validatorAddr"]

		err := sdk.ValidateAccAddress(bech32delegator)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		delegatorAddr := sdk.AccAddress(bech32delegator)

		err = sdk.ValidateValAddress(bech32validator)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		validatorAddr := sdk.ValAddress(bech32validator)

		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		params := types.QueryDelegatorValidatorRequest{DelegatorAddr: delegatorAddr.String(), ValidatorAddr: validatorAddr.String()}

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, height, err := clientCtx.QueryWithData(endpoint, bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryDelegator(clientCtx client.Context, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32delegator := vars["delegatorAddr"]

		err := sdk.ValidateAccAddress(bech32delegator)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		delegatorAddr := sdk.AccAddress(bech32delegator)

		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryDelegatorParams(delegatorAddr)

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, height, err := clientCtx.QueryWithData(endpoint, bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryValidator(clientCtx client.Context, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32validatorAddr := vars["validatorAddr"]

		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		err = sdk.ValidateValAddress(bech32validatorAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		validatorAddr := sdk.ValAddress(bech32validatorAddr)

		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryValidatorParams(validatorAddr, page, limit)

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, height, err := clientCtx.QueryWithData(endpoint, bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}
