package rest

import (
	"net/http"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/tx"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/rest"
	govrest "github.com/line/lbm-sdk/x/gov/client/rest"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
	lbmwasmtypes "github.com/line/lbm-sdk/x/wasm/lbm/types"
)

type InactiveContractProposalJSONReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string    `json:"title" yaml:"title"`
	Description string    `json:"description" yaml:"description"`
	Proposer    string    `json:"proposer" yaml:"proposer"`
	Deposit     sdk.Coins `json:"deposit" yaml:"deposit"`

	Contract string `json:"contract" yaml:"contract"`
}

func (s InactiveContractProposalJSONReq) Content() govtypes.Content {
	return &lbmwasmtypes.DeactivateContractProposal{
		Title:       s.Title,
		Description: s.Description,
		Contract:    s.Contract,
	}
}

func (s InactiveContractProposalJSONReq) GetProposer() string {
	return s.Proposer
}
func (s InactiveContractProposalJSONReq) GetDeposit() sdk.Coins {
	return s.Deposit
}
func (s InactiveContractProposalJSONReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}
func DeactivateContractProposalHandler(cliCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "deactivate_contract",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req InactiveContractProposalJSONReq
			if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

func ActivateContractProposalHandler(cliCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "activate_contract",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req InactiveContractProposalJSONReq
			if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type wasmProposalData interface {
	Content() govtypes.Content
	GetProposer() string
	GetDeposit() sdk.Coins
	GetBaseReq() rest.BaseReq
}

func toStdTxResponse(cliCtx client.Context, w http.ResponseWriter, data wasmProposalData) {
	proposerAddr, err := sdk.AccAddressFromBech32(data.GetProposer())
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	msg, err := govtypes.NewMsgSubmitProposal(data.Content(), data.GetDeposit(), proposerAddr)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := msg.ValidateBasic(); err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	baseReq := data.GetBaseReq().Sanitize()
	if !baseReq.ValidateBasic(w) {
		return
	}
	tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
}
