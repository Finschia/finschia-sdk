package rest

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/link-chain/link/client"

	"github.com/cosmos/cosmos-sdk/client/context"
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

type MintReq struct {
	BaseReq rest.BaseReq   `json:"base_req" yaml:"base_req"`
	To      sdk.AccAddress `json:"to"`
	Amount  int64          `json:"amount"`
}

type BurnReq struct {
	BaseReq rest.BaseReq   `json:"base_req" yaml:"base_req"`
	From    sdk.AccAddress `json:"from"`
	Amount  int64          `json:"amount"`
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

func MintRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MintReq
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "cannot read request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "cannot read request")
			return
		}

		_, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		amount := strconv.FormatInt(req.Amount, 10)
		coins, err := sdk.ParseCoins(strings.Join([]string{amount, symbol}, ""))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgMint(req.To, coins)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func BurnRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BurnReq
		vars := mux.Vars(r)
		symbol := vars["symbol"]

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

		amount := strconv.FormatInt(req.Amount, 10)
		coins, err := sdk.ParseCoins(strings.Join([]string{amount, symbol}, ""))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgBurn(fromAddr, coins)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
