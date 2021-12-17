package rest

import (
	"net/http"

	"github.com/line/lbm-sdk/client/tx"

	"github.com/gorilla/mux"

	govrest "github.com/line/lbm-sdk/x/gov/client/rest"

	"github.com/line/lbm-sdk/client"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/rest"
	"github.com/line/lbm-sdk/x/consortium/types"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

func registerTxHandlers(
	clientCtx client.Context,
	r *mux.Router) {
	r.HandleFunc("/consortium/params", newPostUpdateConsortiumParamsHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/consortium/validators", newPostUpdateValidatorAuthsHandler(clientCtx)).Methods("POST")
}

// UpdateConsortiumParamsRequest defines a proposal to update parameters of consortium.
type UpdateConsortiumParamsRequest struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	Deposit     sdk.Coins    `json:"deposit" yaml:"deposit"`

	Params *types.Params `json:"params" yaml:"params"`
}

// UpdateValidatorAuthsRequest defines a proposal to update validator authorizations.
type UpdateValidatorAuthsRequest struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	Deposit     sdk.Coins    `json:"deposit" yaml:"deposit"`

	Auths []*types.ValidatorAuth `json:"auths" yaml:"auths"`
}

func ProposalUpdateConsortiumParamsRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "consortium",
		Handler:  newPostUpdateConsortiumParamsHandler(clientCtx),
	}
}

func ProposalUpdateValidatorAuthsRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "consortium",
		Handler:  newPostUpdateValidatorAuthsHandler(clientCtx),
	}
}

func newPostUpdateConsortiumParamsHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateConsortiumParamsRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr := sdk.AccAddress(req.BaseReq.From)

		content := types.NewUpdateConsortiumParamsProposal(req.Title, req.Description, req.Params)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func newPostUpdateValidatorAuthsHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateValidatorAuthsRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr := sdk.AccAddress(req.BaseReq.From)

		content := types.NewUpdateValidatorAuthsProposal(req.Title, req.Description, req.Auths)

		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
