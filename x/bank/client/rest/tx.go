package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/tx"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/rest"
	"github.com/line/lbm-sdk/x/bank/types"
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

		toAddr := sdk.AccAddress(bech32Addr)

		var req SendReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr := sdk.AccAddress(req.BaseReq.From)

		msg := types.NewMsgSend(fromAddr, toAddr, req.Amount)
		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
