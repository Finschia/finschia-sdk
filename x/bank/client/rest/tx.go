package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/line/lfb-sdk/client"
	"github.com/line/lfb-sdk/client/tx"
	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/types/rest"
	"github.com/line/lfb-sdk/x/bank/types"
)

// SendReq defines the properties of a send request's body.
type SendReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  sdk.Coins    `json:"amount" yaml:"amount"`
}

// NewSendRequestHandlerFn returns an HTTP REST handler for creating a MsgSend
// transaction.
func NewSendRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		toAddr, err := sdk.AccAddressFromBech32(bech32Addr)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		var req SendReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		msg := types.NewMsgSend(fromAddr, toAddr, req.Amount)
		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
