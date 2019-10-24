package rest

import (
	"github.com/gorilla/mux"
	"github.com/link-chain/link/client"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/link-chain/link/x/auth/client/utils"
	"github.com/link-chain/link/x/token/types"
)

// SendReq defines the properties of a send request's body.
type PublishReq struct {
	BaseReq  rest.BaseReq   `json:"base_req" yaml:"base_req"`
	Name     string         `json:"name"`
	Owner    sdk.AccAddress `json:"owner"`
	Amount   int64          `json:"amount"`
	Mintable bool           `json:"mintable"`
}

// SendRequestHandlerFn - http request handler to send coins to a address.
func PublishRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		//TODO: symbol validate

		var req PublishReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "cannot read request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "cannot read request")
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !fromAddr.Equals(req.Owner) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "owner is different from from address")
			return
		}

		msg := types.NewMsgPublishToken(req.Name, symbol, sdk.NewInt(req.Amount), req.Owner, req.Mintable)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
