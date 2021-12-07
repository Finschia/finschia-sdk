package rest

import (
	"net/http"

	"github.com/line/lbm-sdk/client/tx"

	"github.com/gorilla/mux"

	govrest "github.com/line/lbm-sdk/x/gov/client/rest"

	"github.com/line/lbm-sdk/client"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/rest"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
	"github.com/line/lbm-sdk/x/consortium/types"
)

func registerTxHandlers(
	clientCtx client.Context,
	r *mux.Router) {
	r.HandleFunc("/consortium/disable_consortium", newPostDisableConsortiumHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/consortium/edit_allowed_validators", newEditAllowedValidatorsHandler(clientCtx)).Methods("POST")
}

// DisableConsortiumRequest defines a proposal to disable consortium.
type DisableConsortiumRequest struct {
	BaseReq       rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title         string       `json:"title" yaml:"title"`
	Description   string       `json:"description" yaml:"description"`
	Deposit       sdk.Coins    `json:"deposit" yaml:"deposit"`
}

// EditAllowedValidatorsRequest defines a proposal to edit allowed validators.
type EditAllowedValidatorsRequest struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	Deposit     sdk.Coins    `json:"deposit" yaml:"deposit"`

	AddingAddresses   []string `json:"adding_addresses" yaml:"adding_addresses"`
	RemovingAddresses []string `json:"removing_addresses" yaml:"removing_addresses"`
}

func ProposalDisableConsortiumRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "consortium",
		Handler:  newPostDisableConsortiumHandler(clientCtx),
	}
}

func ProposalEditAllowedValidatorsRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "consortium",
		Handler:  newEditAllowedValidatorsHandler(clientCtx),
	}
}

func newPostDisableConsortiumHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DisableConsortiumRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr := sdk.AccAddress(req.BaseReq.From)

		content := types.NewDisableConsortiumProposal(req.Title, req.Description)
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

func newEditAllowedValidatorsHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EditAllowedValidatorsRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr := sdk.AccAddress(req.BaseReq.From)

		content := types.NewEditAllowedValidatorsProposal(req.Title, req.Description, req.AddingAddresses, req.RemovingAddresses)

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
